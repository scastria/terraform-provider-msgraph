# Microsoft Graph Provider
The Microsoft Graph provider is used to interact Azure Active Directory via the Microsoft Graph REST API.  The provider
needs to be configured with the proper credentials before it can be used.
## Example Usage
Terraform 0.13 and later:
```hcl
terraform {
  required_providers {
    msgraph = {
      source  = "scastria/msgraph"
      version = "~> 0.1.0"
    }
  }
}

# Configure the MSGraph Provider
provider "msgraph" {
  tenant_id = "WWWW"
  client_id = "XXXX"
  client_secret = "YYYY"
//  access_token = "Use access token instead of client_id/client_secret"
}

# Create a Group
resource "msgraph_group" "example" {
  display_name = "TestGroup"
  description = "TestGroup"
  mail_nickname = "TestGroup"
  is_unified = true
  security_enabled = true
}
```
Terraform 0.12 and earlier:
```hcl
# Configure the Apigee Provider
provider "msgraph" {
  version = "~> 0.1.0"
  tenant_id = "WWWW"
  client_id = "XXXX"
  client_secret = "YYYY"
  //  access_token = "Use access token instead of client_id/client_secret"
}

# Create a Group
resource "msgraph_group" "example" {
  display_name = "TestGroup"
  description = "TestGroup"
  mail_nickname = "TestGroup"
  is_unified = true
  security_enabled = true
}
```
## Argument Reference
* `tenant_id` - **(Required, String)** The tenant id for your Azure Active Directory. Can be specified via env variable `MSGRAPH_TENANT_ID`.
* `access_token` - **(Optional, String)** The access token obtained via authentication that can be used instead of `client_id` and `client_secret`. Token Authentication. Can be specified via env variable `MSGRAPH_ACCESS_TOKEN`.
* `client_id` - **(Optional, String)** The client_id that will invoke all MS Graph API commands. Client Credentials Authentication. Can be specified via env variable `MSGRAPH_CLIENT_ID`.
* `client_secret` - **(Optional, String)** The client_secret for the client_id. Client Credentials Authentication. Can be specified via env variable `MSGRAPH_CLIENT_SECRET`.
