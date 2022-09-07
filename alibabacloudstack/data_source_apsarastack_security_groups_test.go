package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackSecurityGroupsDataSourceBasic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackSecurityGroupsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_security_groups.default"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_groups.default", "groups.#", "1"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_security_groups.default", "ids.#"),
				),
			},
		},
	})
}

const testAccCheckAlibabacloudStackSecurityGroupsDataSourceConfig = `

variable "name" {
  default = "tf-securityGroupdatasource"
}
data "alibabacloudstack_zones" "default" {
	available_resource_creation = "VSwitch"
}
resource "alibabacloudstack_vpc" "vpc" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vswitch" "vswitch" {
  vpc_id            = alibabacloudstack_vpc.vpc.id
  cidr_block        = "172.16.0.0/24"
  availability_zone =  data.alibabacloudstack_zones.default.zones.0.id
  name              = "test45"
}
resource "alibabacloudstack_security_group" "group" {
  name        = var.name
  description = "foo"
  vpc_id      = alibabacloudstack_vpc.vpc.id
}
data "alibabacloudstack_security_groups" "default" {
  ids = [alibabacloudstack_security_group.group.id]
}
`
