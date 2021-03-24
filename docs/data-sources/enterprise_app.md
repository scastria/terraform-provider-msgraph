---
subcategory: "Application"
---
# Data Source: msgraph_enterprise_app
Represents an enterprise application
## Example usage
```hcl
data "msgraph_enerprise_app" "example" {
  search_display_name = "MyEnterpriseApp"
//  display_name = "Company MyEnterpriseApp"
}
```
## Argument Reference
* `search_display_name` - **(Optional, String)** The search string to apply to the display name of the enterprise app. Uses contains.
* `display_name` - **(Optional, String)** The filter string to apply to the display name of the enterprise app. Uses equality.
* `app_id` - **(Optional, String)** The filter string to apply to the given app_id of the enterprise app. Uses equality.
## Attribute Reference
* `id` - **(String)** Guid
* `app_display_name` - **(String)** The app display name of the enterprise app.
* `login_url` - **(String)** The login url of the enterprise app.
* `logout_url` - **(String)** The logout url of the enterprise app.
