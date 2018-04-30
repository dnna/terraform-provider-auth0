package main

import (
	"log";

	"github.com/hashicorp/terraform/helper/schema";
	"io/ioutil";
	"strings";
	"net/http";
	"encoding/json";
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
				Type:		 schema.TypeString,
				Required:	 true,
				DefaultFunc: schema.EnvDefaultFunc("MANAGEMENT_API_CLIENT_ID", nil),
			},
			"client_secret": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MANAGEMENT_API_CLIENT_SECRET", nil),
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
	domain       string
	accessToken  string
}

type Auth0Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn int `json:"expires_in"`
	Scope string `json:"scope"`
	TokenType string `json:"token_type"`
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Println("[INFO] Initializing Auth0 client")
	domain := d.Get("domain").(string)
	client_secret := d.Get("client_secret").(string)
	client_id := d.Get("client_id").(string)
	
	url := "https://" + domain + "/oauth/token"
	
	payload := strings.NewReader(`{ "grant_type": "client_credentials", ` +
		`"client_id": "` + client_id +  `", ` +
		`"client_secret": "` + client_secret + `", ` + 
		`"audience": "https://` + domain + `/api/v2"}`)

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var auth0Token Auth0Token
	err := json.Unmarshal(body, &auth0Token)
	if err != nil {
		log.Println("[ERROR] Failed to Unmarshal Auth0 Response")
	}

	return Config{ domain: domain, accessToken: auth0Token.AccessToken}, nil
}
