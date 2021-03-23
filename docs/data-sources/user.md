---
subcategory: "User"
---
# Data Source: msgraph_user
Represents a user
## Example usage
```hcl
data "msgraph_user" "example" {
  search_display_name = "John"
//  display_name = "John Smith"
}
```
## Argument Reference
* `search_display_name` - **(Optional, String)** The search string to apply to the display name of the user. Uses contains.
* `display_name` - **(Optional, String)** The filter string to apply to the display name of the user. Uses equality.
* `given_name` - **(Optional, String)** The filter string to apply to the given name of the user. Uses equality.
* `surname` - **(Optional, String)** The filter string to apply to the surname of the user. Uses equality.
* `job_title` - **(Optional, String)** The filter string to apply to the job title of the user. Uses equality.
* `mail` - **(Optional, String)** The filter string to apply to the mail of the user. Uses equality.
* `user_principal_name` - **(Optional, String)** The filter string to apply to the user principal name of the user. Uses equality.
## Attribute Reference
* `id` - **(String)** Guid
