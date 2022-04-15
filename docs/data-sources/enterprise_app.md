---
subcategory: "Application"
---
# Data Source: msgraph_enterprise_app
Represents an enterprise application
## Example usage
```hcl
data "msgraph_enerprise_app" "example" {
  search_display_name = "TestEnterpriseApp"
//  display_name = "Company TestEnterpriseApp"
}
```
## Argument Reference
* `search_display_name` - **(Optional, String)** The search string to apply to the display name of the enterprise app. Uses contains.
* `display_name` - **(Optional, String)** The filter string to apply to the display name of the enterprise app. Uses equality.
* `app_id` - **(Optional, String)** The filter string to apply to the app id of the enterprise app. Uses equality.
* `wait_until_exists` - **(Optional, Boolean)** Whether to wait and keep checking for existence of the enterprise app instead of immediately returning an error.  Default: `false`
* `wait_timeout` - **(Optional, Integer)** How many total seconds to wait for existence until giving up.  Default: `60`
* `wait_polling_interval` - **(Optional, Integer)** How many seconds to wait between existence checks.  Default: `10`
## Attribute Reference
* `id` - **(String)** Guid
* `login_url` - **(String)** The login url of the enterprise app.
* `logout_url` - **(String)** The logout url of the enterprise app.
