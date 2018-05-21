package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceClient() *schema.Resource {
	return &schema.Resource{
		Create: resourceClientCreate,
		Read:   resourceClientRead,
		Update: resourceClientUpdate,
		Delete: resourceClientDelete,

		Schema: map[string]*schema.Schema{
			"client_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_secret": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Sensitive: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_token_endpoint_ip_header_trusted": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_first_party": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cross_origin_auth": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"sso": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"token_endpoint_auth_method": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"grant_types": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"custom_login_page_on": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"app_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"callbacks": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

type Client struct {
	Client_id                           string        `json:"client_id,omitempty"`
	Client_secret                       string        `json:"client_secret,omitempty"`
	Name                                string        `json:"name"`
	Is_token_endpoint_ip_header_trusted bool          `json:"is_token_endpoint_ip_header_trusted"`
	Is_first_party                      bool          `json:"is_first_party"`
	Description                         string        `json:"description"`
	Cross_origin_auth                   bool          `json:"cross_origin_auth"`
	Sso                                 bool          `json:"sso"`
	Token_endpoint_auth_method          string        `json:"token_endpoint_auth_method,omitempty"`
	Grant_types                         []interface{} `json:"grant_types"`
	App_type                            string        `json:"app_type,omitempty"`
	Custom_login_page_on                bool          `json:"custom_login_page_on"`
}

func resourceClientCreate(d *schema.ResourceData, m interface{}) error {
	reqClient := Client{
		Name: d.Get("name").(string),
		Is_token_endpoint_ip_header_trusted: d.Get("is_token_endpoint_ip_header_trusted").(bool),
		Is_first_party:                      d.Get("is_first_party").(bool),
		Description:                         d.Get("description").(string),
		Cross_origin_auth:                   d.Get("cross_origin_auth").(bool),
		Sso:                                 d.Get("sso").(bool),
		Token_endpoint_auth_method: d.Get("token_endpoint_auth_method").(string),
		Grant_types:                d.Get("grant_types").([]interface{}),
		App_type:                   d.Get("app_type").(string),
		Custom_login_page_on:       d.Get("custom_login_page_on").(bool),
	}

	jsonValue, err := json.Marshal(reqClient)
	if err != nil {
		return err
	}
	log.Println("[DEBUG] Request JSON: " + string(jsonValue))

	config := m.(Config)
	client := http.Client{}
	req, err := http.NewRequest("POST", "https://"+config.domain+"/api/v2/clients", bytes.NewBuffer(jsonValue))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+config.accessToken)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println("[DEBUG] Response JSON:" + string(data))

	if resp.StatusCode != 201 {
		return errors.New("Error: Invalid status code during create: " + string(data))
	}

	var respClient Client
	err = json.Unmarshal(data, &respClient)
	if err != nil {
		return err
	}

	log.Println("New client_id:" + respClient.Client_id)
	d.Set("client_id", respClient.Client_id)
	d.Set("client_secret", respClient.Client_secret)
	d.SetId(respClient.Client_id)
	return nil
}

func resourceClientRead(d *schema.ResourceData, m interface{}) error {
	config := m.(Config)
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://"+config.domain+"/api/v2/clients/"+d.Id(), nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+config.accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println("Response JSON:" + string(data))

	if resp.StatusCode != 200 {
		if resp.StatusCode == 404 {
			// Client was deleted
			d.SetId("")
			return nil
		} else {
			return errors.New("Error: Invalid status code during read: " + string(data))
		}
	}

	var respClient Client
	err = json.Unmarshal(data, &respClient)
	if err != nil {
		return err
	}

	return nil
}

func resourceClientUpdate(d *schema.ResourceData, m interface{}) error {
	reqClient := Client{
		Name: d.Get("name").(string),
		Is_token_endpoint_ip_header_trusted: d.Get("is_token_endpoint_ip_header_trusted").(bool),
		Is_first_party:                      d.Get("is_first_party").(bool),
		Description:                         d.Get("description").(string),
		Cross_origin_auth:                   d.Get("cross_origin_auth").(bool),
		Sso:                                 d.Get("sso").(bool),
		Token_endpoint_auth_method: d.Get("token_endpoint_auth_method").(string),
		Grant_types:                d.Get("grant_types").([]interface{}),
		App_type:                   d.Get("app_type").(string),
		Custom_login_page_on:       d.Get("custom_login_page_on").(bool),
	}

	jsonValue, err := json.Marshal(reqClient)
	if err != nil {
		return err
	}
	log.Println("Request JSON: " + string(jsonValue))

	config := m.(Config)
	client := http.Client{}
	req, err := http.NewRequest("PATCH", "https://"+config.domain+"/api/v2/clients/"+d.Id(), bytes.NewBuffer(jsonValue))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+config.accessToken)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println("Response JSON:" + string(data))

	if resp.StatusCode != 200 {
		return errors.New("Error: Invalid status code during update: " + string(data))
	}

	var respClient Client
	err = json.Unmarshal(data, &respClient)
	if err != nil {
		return err
	}

	return nil
}

func resourceClientDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(Config)
	client := http.Client{}
	req, err := http.NewRequest("DELETE", "https://"+config.domain+"/api/v2/clients/"+d.Id(), nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+config.accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	log.Println("Response Status Code:" + string(resp.StatusCode))

	if resp.StatusCode != 204 {
		return errors.New("Error: Invalid status code during delete")
	}

	return nil
}
