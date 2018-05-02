package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ClientRequest struct {
	GrantType    string `json:"grant_type"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience"`
}

func TestProviderConfigRawSad(t *testing.T) {
	assert := assert.New(t)
	clientSecret := "cauliflower"
	clientId := "joebang"

	testServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Header().Set("content-type", "application/json")
		fmt.Fprintf(w, `{"error":"access_denied","error_description":"Service not enabled within domain: https://dshbreak.auth0.com/api/v2/"}`)
	}))
	defer testServer.Close()

	testDomain := testServer.URL[8:]
	result, _ := providerConfigureRaw(testServer.Client(), testDomain, clientId, clientSecret, "")
	assert.Equal(Config{}, result)

}

func TestProviderConfigWithToken(t *testing.T) {
	assert := assert.New(t)

	result, err := providerConfigureRaw(http.DefaultClient, "contoso.auth0.com", "", "", "not_a_real_token")
	assert.Equal("contoso.auth0.com", result.(Config).domain)
	assert.Equal("not_a_real_token", result.(Config).accessToken)
	assert.Equal(err, nil)
}

func TestProviderConfigRaw(t *testing.T) {
	assert := assert.New(t)

	times := 0
	clientSecret := "cauliflower"
	clientId := "joebang"
	token := "wubbalubbadubdub"

	testServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		times++
		assert.Equal("POST", r.Method)
		body, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			t.Fatalf("Failed to read request: %s", readErr)
		}

		var clientRequest ClientRequest
		unmarshalErr := json.Unmarshal(body, &clientRequest)
		if unmarshalErr != nil {
			t.Fatalf("Failed to parse request: %s", unmarshalErr)
		}

		assert.Equal(clientSecret, clientRequest.ClientSecret)
		assert.Equal(clientId, clientRequest.ClientId)
		assert.Equal("client_credentials", clientRequest.GrantType)

		clientResponse := &Auth0Token{
			AccessToken: token,
			ExpiresIn:   86400,
			Scope:       "superman:all",
			TokenType:   "type",
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(clientResponse)
	}))
	defer testServer.Close()

	testDomain := testServer.URL[8:]
	result, _ := providerConfigureRaw(testServer.Client(), testDomain, clientId, clientSecret, "")

	assert.Equal(1, times)
	assert.Equal(testDomain, result.(Config).domain)
	assert.Equal(token, result.(Config).accessToken)
}
