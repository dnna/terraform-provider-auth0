# terraform-provider-auth0

A custom provider for terraform.

## Installation

1. Download the latest [release](github.com/dnna/terraform-provider-auth0/releases) for your platform
2. rename the file to `terraform-provider-auth0`
3. Copy the file to the same directory as terraform `dirname $(which terraform)` is installed

## Usage

### auth0_client

#### Example Usage

Provides an auth0 client

```hcl
provider "auth0" {
    domain = "abc.eu.auth0.com"
    client_id = "<CLIENT_ID>"
    client_secret = "<CLIENT_SECRET>"
}

resource "auth0_client" "test-client" {
    name = "p-test-dev2"
    is_token_endpoint_ip_header_trusted = true
    is_first_party = false
    description = ""
    cross_origin_auth = false
    sso = true
    token_endpoint_auth_method = "client_secret_post"
    grant_types = [ "authorization_code", "implicit", "refresh_token", "client_credentials" ]
    app_type = "non_interactive"
    custom_login_page_on = true
}
```

The output will look like this:

```sh
resource.auth0_client.test-client.client_id = "generated_client_id"
resource.auth0_client.test-client.client_secret = "generated_client_secret"
```

#### Argument Reference

Arguments have the same names as provided in Auth0 Management API documentation (https://auth0.com/docs/api/management/v2#!/Clients).

The provider itself requires 3 parameters:


- `domain` - The provided domain for the Auth0 account
- `client_id` - The client ID for the Application
- `client_secret` - The client secret for the Applicaton.

You can consult [the Auth0 Documentation](https://auth0.com/docs/api/management/v2/tokens#1-create-and-authorize-an-application) for steps on creating a Machine-to-Machine Application in your Auth0 tenant with access to its Auth0 Management API.

Optionally, you may also manually specify a token using the `access_token` parameter:

```
provider "auth0" {
    domain = "abc.eu.auth0.com"
    access_token = "<ACCESS_TOKEN>"
}
```

When you do this, the provider will not request a token on your behalf.

#### Attributes Reference

- `client_id` - The client ID of the new client
- `client_secret` - The client secret of the new client

## Develop

```sh
go get github.com/dnna/terraform-provider-auth0
cd $GOPATH/src/github.com/dnna/terraform-provider-auth0
go get ./...
$EDITOR .
```
