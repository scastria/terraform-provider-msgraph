---
subcategory: "Group"
---
# Resource: msgraph_group_role_assignment
Represents an assignment of an enterprise app role to a group
## Example usage
```hcl
data "msgraph_enerprise_app" "TestEnterpriseApp" {
  search_display_name = "TestEnterpriseApp"
  //  display_name = "Company TestEnterpriseApp"
}
data "msgraph_enerprise_app_role" "TestEnterpriseAppUserRole" {
  enterprise_app_id = msgraph_enerprise_app.TestEnterpriseApp.id
  search_display_name = "User"
  //  display_name = "User"
}
resource "msgraph_group" "TestGroup" {
  display_name = "TestGroup"
  description = "TestGroup"
  mail_nickname = "TestGroup"
  is_unified = true
  security_enabled = true
}
resource "msgraph_group_role_assignment" "example" {
  group_id = msgraph_group.TestGroup.id
  enterprise_app_role_id = data.msgraph_enerprise_app_role.TestEnterpriseAppUserRole.id
}
```
## Argument Reference
* `group_id` - **(Required, ForceNew, String)** The id of the group.
* `enterprise_app_role_id` - **(Required, ForceNew, String)** The id of the enterprise app role.
## Attribute Reference
* `id` - **(String)** Same as `group_id`:`Guid of the assignment`
## Import
Group role assignments can be imported using a proper value of `id` as described above
