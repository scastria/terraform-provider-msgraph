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
* `is_public` - **(Optional, Boolean)** Whether this group can be joined without owner's approval. Default: `false`.
* `primary_owner_id` - **(Optional, ForceNew, String)** An owner can be attached to the group at time of creation.  This property is only used at creation time so any changes forces a new resource.
## Attribute Reference
* `id` - **(String)** Guid
* `mail` - **(String)** The email address of the group if `is_unified` is true
* `mail_enabled` - **(Boolean)** Same as `is_unified`
## Import
Groups can be imported using a proper value of `id` as described above
