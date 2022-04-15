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
}
```
## Argument Reference
* `app_id` - **(Required, String, ForceNew)** The `app_id` of the corresponding app registration
## Attribute Reference
* `id` - **(String)** Guid
* `display_name` - **(String)** The `display_name` of the corresponding app registration
## Import
Enterprise apps can be imported using a proper value of `id` as described above
