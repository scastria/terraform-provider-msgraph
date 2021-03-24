---
subcategory: "Application"
---
# Data Source: msgraph_app_registration
Represents an application registration
## Example usage
```hcl
data "msgraph_app_registration" "example" {
  search_display_name = "MyApp"
//  display_name = "Company MyApp"
}
```
## Argument Reference
* `search_display_name` - **(Optional, String)** The search string to apply to the display name of the app registration. Uses contains.
* `display_name` - **(Optional, String)** The filter string to apply to the display name of the app registration. Uses equality.
* `app_id` - **(Optional, String)** The filter string to apply to the app id of the app registration. Uses equality.
## Attribute Reference
* `id` - **(String)** Guid
