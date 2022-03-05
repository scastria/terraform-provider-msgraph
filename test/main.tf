terraform {
  required_providers {
    msgraph = {
      version = "0.1"
      source = "github.com/scastria/msgraph"
    }
  }
}

provider "msgraph" {
  tenant_id = "WWWW"
  client_id = "XXXX"
  client_secret = "YYYY"
}

data "msgraph_group" "MyGroup" {
  display_name = "ShawnTest"
#  wait_until_exists = true
#  wait_timeout = 55
#  wait_polling_interval = 3
}

output "MyOutput" {
  value = data.msgraph_group.MyGroup.id
}
