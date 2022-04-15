---
subcategory: "Application"
---
# Resource: msgraph_enterprise_app
Represents an enterprise app
## Example usage
```hcl
resource "msgraph_app_registration" "TestAppReg" {
  display_name = "TestAppRegistration"
}
resource "msgraph_enterprise_app" "example" {
  app_id = msgraph_app_registration.TestAppReg.app_id
  app_role {
    display_name = "TestRole"
    description = "TestRole"
    allowed_member_types = ["User"]
  }
}
```
## Argument Reference
* `app_id` - **(Required, String, ForceNew)** The `app_id` of the corresponding app registration
* `app_role` - **(Optional, set{object})** Configuration block for an app role.  Can be specified multiple times for each app role.  Each block supports the fields documented below.
### app_role
* `allowed_member_types` - **(Required, List of String)** Types of members allowed for this role.  Allowed values: `User` and `Application`
* `display_name` - **(Required, String)** Display name of role
* `description` - **(Required, String)** Description of role
* `is_enabled` - **(Optional, Boolean)** Whether role is enabled or not.  Default: `true`
* `id` - **(Computed, String)** Generated GUID identifier for role
## Attribute Reference
* `id` - **(String)** Guid
* `display_name` - **(String)** The `display_name` of the corresponding app registration
## Import
Enterprise apps can be imported using a proper value of `id` as described above
