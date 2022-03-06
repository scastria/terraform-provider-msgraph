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
* `wait_until_exists` - **(Optional, Boolean)** Whether to wait and keep checking for existence of the user instead of immediately returning an error.  Default: `false`
* `wait_timeout` - **(Optional, Integer)** How many total seconds to wait for existence until giving up.  Default: `60`
* `wait_polling_interval` - **(Optional, Integer)** How many seconds to wait between existence checks.  Default: `10`
## Attribute Reference
* `id` - **(String)** Guid
