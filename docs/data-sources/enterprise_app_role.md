---
subcategory: "Application"
---
# Data Source: msgraph_enterprise_app_role
Represents an enterprise application role
## Example usage
```hcl
data "msgraph_enerprise_app" "TestEnterpriseApp" {
  search_display_name = "TestEnterpriseApp"
//  display_name = "Company TestEnterpriseApp"
}
data "msgraph_enerprise_app_role" "example" {
  enterprise_app_id = msgraph_enerprise_app.TestEnterpriseApp.id
  search_display_name = "User"
  //  display_name = "User"
}
```
## Argument Reference
* `enterprise_app_id` - **(Required, String)** The enterprise app id for which to find roles within.
* `search_display_name` - **(Optional, String)** The search string to apply to the display name of the enterprise app role. Uses contains.
* `display_name` - **(Optional, String)** The filter string to apply to the display name of the enterprise app role. Uses equality.
* `description` - **(Optional, String)** The filter string to apply to the description the enterprise app role. Uses equality.
## Attribute Reference
* `id` - **(String)** Same as `enterprise_app_id`:`Guid of the role`
