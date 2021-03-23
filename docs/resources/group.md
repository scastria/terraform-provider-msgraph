---
subcategory: "Group"
---
# Resource: msgraph_group
Represents a group
## Example usage
```hcl
resource "msgraph_group" "example" {
  display_name = "TestGroup"
  description = "TestGroup"
  mail_nickname = "TestGroup"
  is_unified = true
  security_enabled = true
}
```
## Argument Reference
* `display_name` - **(Required, String)** The display name of the group.
* `mail_nickname` - **(Required, String)** The mail nickname of the group.
* `description` - **(Required, String)** The description of the group.  Cannot be empty string.
* `security_enabled` - **(Optional, Boolean)** Whether this group can be used for security purposes.
* `is_unified` - **(Optional, ForceNew, Boolean)** Whether this group is the new `Unified` group type.
## Attribute Reference
* `id` - **(String)** Guid
* `mail` - **(String)** The email address of the group if `is_unified` is true
* `mail_enabled` - **(Boolean)** Same as `is_unified`
## Import
Users can be imported using a proper value of `id` as described above
