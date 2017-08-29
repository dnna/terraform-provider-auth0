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

func resourceGrant() *schema.Resource {
	return &schema.Resource{
		Create: resourceGrantCreate,
		Read:   resourceGrantRead,
		Update: resourceGrantUpdate,
		Delete: resourceGrantDelete,

		Schema: map[string]*schema.Schema{
			"grant_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"audience": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"scope": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
}

type Grant struct {
	Id        string        `json:"id,omitempty"`
	Client_id string        `json:"client_id,omitempty"`
	Audience  string        `json:"audience,omitempty"`
	Scope     []interface{} `json:"scope"`
}

func resourceGrantCreate(d *schema.ResourceData, m interface{}) error {
	reqClient := Grant{
		Client_id: d.Get("client_id").(string),
		Audience:  d.Get("audience").(string),
		Scope:     d.Get("scope").([]interface{}),
	}

	jsonValue, err := json.Marshal(reqClient)
	if err != nil {
		return err
	}
	log.Println("Request JSON: " + string(jsonValue))

	config := m.(Config)
	client := http.Client{}
	req, err := http.NewRequest("POST", "https://"+config.domain+"/api/v2/client-grants", bytes.NewBuffer(jsonValue))
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

	if resp.StatusCode != 201 {
		return errors.New("Error: Invalid status code during create: " + string(data))
	}

	var respClient Grant
	err = json.Unmarshal(data, &respClient)
	if err != nil {
		return err
	}

	log.Println("New client_id:" + respClient.Id)
	d.Set("grant_id", respClient.Id)
	d.SetId(respClient.Id)
	return nil
}

func resourceGrantRead(d *schema.ResourceData, m interface{}) error {
	config := m.(Config)
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://"+config.domain+"/api/v2/client-grants?client_id="+d.Get("client_id").(string)+"&audience="+d.Get("audience").(string), nil)
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
		return errors.New("Error: Invalid status code during read: " + string(data))
	}

	var respClient []Grant
	err = json.Unmarshal(data, &respClient)
	if err != nil {
		return err
	}

	if len(respClient) > 1 {
		return errors.New("Error: Multiple grants found for the same client and audience: " + string(data))
	}

	if len(respClient) <= 0 {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceGrantUpdate(d *schema.ResourceData, m interface{}) error {
	reqClient := Grant{
		Client_id: d.Get("client_id").(string),
		Audience:  d.Get("audience").(string),
		Scope:     d.Get("scope").([]interface{}),
	}

	jsonValue, err := json.Marshal(reqClient)
	if err != nil {
		return err
	}
	log.Println("Request JSON: " + string(jsonValue))

	config := m.(Config)
	client := http.Client{}
	req, err := http.NewRequest("PATCH", "https://"+config.domain+"/api/v2/client-grants/"+d.Id(), bytes.NewBuffer(jsonValue))
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

	var respClient Grant
	err = json.Unmarshal(data, &respClient)
	if err != nil {
		return err
	}

	return nil
}

func resourceGrantDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(Config)
	client := http.Client{}
	req, err := http.NewRequest("DELETE", "https://"+config.domain+"/api/v2/client-grants/"+d.Id(), nil)
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
