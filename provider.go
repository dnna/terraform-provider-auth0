package main

import (
	"log"

	"encoding/json"
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"domain": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MANAGEMENT_API_DOMAIN", nil),
			},
			"client_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MANAGEMENT_API_CLIENT_ID", nil),
			},
			"client_secret": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MANAGEMENT_API_CLIENT_SECRET", nil),
			},
			"access_token":  &schema.Schema{
				Type:		 schema.TypeString,
				Optional:	 true,
				DefaultFunc: schema.EnvDefaultFunc("MANAGEMENT_ACCESS_TOKEN", nil), 
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"auth0_client": resourceClient(),
			"auth0_grant":  resourceGrant(),
		},

		ConfigureFunc: providerConfigure,
	}
}

type Config struct {
	domain      string
	accessToken string
}

type Auth0Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Println("[INFO] Initializing Auth0 client")
	domain := d.Get("domain").(string)
	clientSecret := d.Get("client_secret").(string)
	clientId := d.Get("client_id").(string)
	accessToken := d.Get("access_token").(string)
	log.Printf("[DEBUG] d %s, cs %s, ci %s, at %s", domain, clientSecret, clientId, accessToken)

	return providerConfigureRaw(http.DefaultClient, domain, clientId, clientSecret, accessToken)
}

func providerConfigureRaw(client *http.Client, domain string, clientId string, clientSecret string, accessToken string) (interface{}, error) {
	if (accessToken == "") && (clientId == "" && clientSecret == "") {
		message := "Must set either access_token or client_id AND client_secret."
		return nil, errors.New(message)
	}

	// for the sake of compatibility, make sure you can still use the token if you've got it.
	if accessToken != "" {
		log.Printf("[INFO] Skipping token request, token already present.")
		return Config{domain: domain, accessToken: accessToken}, nil
	}

	url := "https://" + domain + "/oauth/token"

	payload := strings.NewReader(`{ "grant_type": "client_credentials", ` +
		`"client_id": "` + clientId + `", ` +
		`"client_secret": "` + clientSecret + `", ` +
		`"audience": "https://` + domain + `/api/v2/"}`)

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("content-type", "application/json")

	res, postErr := client.Do(req)
	if postErr != nil {
		log.Printf("[ERROR] Failed to contact Auth0 API: %s", postErr)
		return Config{}, postErr
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read Response: %s", err)
		return Config{}, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		log.Printf("[ERROR] Received HTTP Response Code %d from %s: %s", res.StatusCode, url, body)
		errorText := url + ` responded with code ` + strconv.Itoa(res.StatusCode) + `: ` + string(body)
		return Config{}, errors.New(errorText)
	}

	log.Printf("[DEBUG] Response was: %s", body)

	var auth0Token Auth0Token
	unmarshalErr := json.Unmarshal(body, &auth0Token)
	if unmarshalErr != nil {
		log.Printf("[ERROR] Failed to Unmarshal Auth0 Response: %s", unmarshalErr)
		return Config{}, unmarshalErr
	}

	log.Printf("[DEBUG] Domain is %s, token is %s", domain, auth0Token.AccessToken)

	return Config{domain: domain, accessToken: auth0Token.AccessToken}, nil
}
