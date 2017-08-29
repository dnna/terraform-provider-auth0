package main

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"domain": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MANAGEMENT_API_DOMAIN", nil),
			},
			"access_token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MANAGEMENT_API_ACCESS_TOKEN", nil),
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

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Println("[INFO] Initializing Auth0 client")
	return Config{domain: d.Get("domain").(string), accessToken: d.Get("access_token").(string)}, nil
}
