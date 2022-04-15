---
subcategory: "Application"
---
# Resource: msgraph_enterprise_app_owner
Represents an owner of an enterprise app
## Example usage
```hcl
resource "msgraph_app_registration" "TestAppReg" {
  display_name = "TestAppRegistration"
}
resource "msgraph_enterprise_app" "TestEntApp" {
  app_id = msgraph_app_registration.TestAppReg.app_id
  app_role {
    display_name = "TestRole"
    description = "TestRole"
    allowed_member_types = ["User"]
  }
}
data "msgraph_user" "TestUser" {
  user_principal_name = "John.Smith@company.com"
}
resource "msgraph_enterprise_app_owner" "example" {
  enterprise_app_id = msgraph_enterprise_app.TestEntApp.id
  owner_id = data.msgraph_user.TestUser.id
}
```
## Argument Reference
* `enterprise_app_id` - **(Required, ForceNew, String)** The id of the enterprise app.
* `owner_id` - **(Required, ForceNew, String)** The id of the user owner.
## Attribute Reference
* `id` - **(String)** Same as `enterprise_app_id`:`owner_id`
## Import
Enterprise app owners can be imported using a proper value of `id` as described above
