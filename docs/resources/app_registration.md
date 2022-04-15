---
subcategory: "Application"
---
# Resource: msgraph_app_registration
Represents an app registration
## Example usage
```hcl
resource "msgraph_app_registration" "example" {
  display_name = "TestAppRegistration"
}
```
## Argument Reference
* `display_name` - **(Required, String)** The display name of the app registration.
## Attribute Reference
* `id` - **(String)** Guid
* `app_id` - **(String)** The app id (client id) of the app registration
## Import
App registrations can be imported using a proper value of `id` as described above
