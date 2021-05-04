---
subcategory: "Group"
---
# Data Source: msgraph_group
Represents a group
## Example usage
```hcl
data "msgraph_group" "example" {
  search_display_name = "Scientists"
//  display_name = "Data Scientists"
}
```
## Argument Reference
* `search_display_name` - **(Optional, String)** The search string to apply to the display name of the group. Uses contains.
* `display_name` - **(Optional, String)** The filter string to apply to the display name of the group. Uses equality.
* `mail_nickname` - **(Optional, String)** The filter string to apply to the mail nickname of the group. Uses equality.
* `mail` - **(Optional, String)** The filter string to apply to the mail of the group. Uses equality.
## Attribute Reference
* `id` - **(String)** Guid
* `description` - **(String)** The description of the group.
* `security_enabled` - **(Boolean)** Whether this group can be used for security purposes.
* `is_unified` - **(Boolean)** Whether this group is the new `Unified` group type.
* `is_public` - **(Boolean)** Whether this group can be joined without owner's approval.
* `mail_enabled` - **(Boolean)** Same as `is_unified`
