package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackSecurityGroupsDataSourceBasic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackSecurityGroupsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_security_groups.default"),
					resource.TestCheckResourceAttr("data.apsarastack_security_groups.default", "groups.#", "1"),
					resource.TestCheckResourceAttrSet("data.apsarastack_security_groups.default", "ids.#"),
				),
			},
		},
	})
}

const testAccCheckApsaraStackSecurityGroupsDataSourceConfig = `

variable "name" {
  default = "tf-securityGroupdatasource"
}
data "apsarastack_zones" "default" {
	available_resource_creation = "VSwitch"
}
resource "apsarastack_vpc" "vpc" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}
resource "apsarastack_vswitch" "vswitch" {
  vpc_id            = apsarastack_vpc.vpc.id
  cidr_block        = "172.16.0.0/24"
  availability_zone =  data.apsarastack_zones.default.zones.0.id
  name              = "test45"
}
resource "apsarastack_security_group" "group" {
  name        = var.name
  description = "foo"
  vpc_id      = apsarastack_vpc.vpc.id
}
data "apsarastack_security_groups" "default" {
  ids = [apsarastack_security_group.group.id]
}
`
