---
subcategory: "Group"
---
# Resource: msgraph_group_member
Represents a member of a group
## Example usage
```hcl
resource "msgraph_group" "TestGroup" {
  display_name = "TestGroup"
  description = "TestGroup"
  mail_nickname = "TestGroup"
  is_unified = true
  security_enabled = true
}
data "msgraph_user" "TestUser" {
  user_principal_name = "John.Smith@company.com"
}
resource "msgraph_group_member" "example" {
  group_id = msgraph_group.TestGroup.id
  member_id = data.msgraph_user.TestUser.id
}
```
## Argument Reference
* `group_id` - **(Required, ForceNew, String)** The id of the group.
* `member_id` - **(Required, ForceNew, String)** The id of the user member.
## Attribute Reference
* `id` - **(String)** Same as `group_id`:`member_id`
## Import
Group members can be imported using a proper value of `id` as described above
