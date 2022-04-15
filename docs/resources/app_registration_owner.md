---
subcategory: "Application"
---
# Resource: msgraph_app_registration_owner
Represents an owner of an app registration
## Example usage
```hcl
resource "msgraph_app_registration" "TestAppReg" {
  display_name = "TestAppRegistration"
}
data "msgraph_user" "TestUser" {
  user_principal_name = "John.Smith@company.com"
}
resource "msgraph_app_registration_owner" "example" {
  app_registration_id = msgraph_app_registration.TestAppReg.id
  owner_id = data.msgraph_user.TestUser.id
}
```
## Argument Reference
* `app_registration_id` - **(Required, ForceNew, String)** The id of the app registration.
* `owner_id` - **(Required, ForceNew, String)** The id of the user owner.
## Attribute Reference
* `id` - **(String)** Same as `app_registration_id`:`owner_id`
## Import
App registration owners can be imported using a proper value of `id` as described above
