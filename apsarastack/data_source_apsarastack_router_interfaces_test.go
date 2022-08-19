package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackRouterInterfacesDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackRouterInterfacesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_router_interfaces.default"),
					resource.TestCheckResourceAttr("data.apsarastack_router_interfaces.default", "interfaces.#", "1"),
					resource.TestCheckResourceAttrSet("data.apsarastack_router_interfaces.default", "ids.#"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})

}

const testAccCheckApsaraStackRouterInterfacesDataSourceConfig = `
provider "apsarastack" {
  region = "${var.region}"
}
variable "region" {
  default = "cn-neimeng-env30-d01"
}
variable "name" {
  default = "tf-testAccCheckApsaraStackRouterInterfacesDataSourceConfig"
}
variable cidr_block_list {
	type = "list"
	default = [ "172.16.0.0/12", "192.168.0.0/16" ]
}

resource "apsarastack_vpc" "default" {
  count = 2
  name = "${var.name}"
  cidr_block = "${element(var.cidr_block_list,count.index)}"
}
resource "apsarastack_router_interface" "initiating" {
  opposite_region = "${var.region}"
  router_type = "VRouter"
  router_id = "${apsarastack_vpc.default.0.router_id}"
  role = "InitiatingSide"
  specification = "Large.2"
  name = "${var.name}_initiating"
  description = "${var.name}_decription"

}
resource "apsarastack_router_interface" "opposite" {
  provider = "apsarastack"
  opposite_region = "${var.region}"
  router_type = "VRouter"
  router_id = "${apsarastack_vpc.default.1.router_id}"
  role = "AcceptingSide"
  specification = "Large.1"
  name = "${var.name}_opposite"
  description = "${var.name}_decription"

}

resource "apsarastack_router_interface_connection" "foo" {
  interface_id = "${apsarastack_router_interface.initiating.id}"
  opposite_interface_id = "${apsarastack_router_interface.opposite.id}"
  depends_on = ["apsarastack_router_interface_connection.bar"]
  opposite_interface_owner_id = "1262302482727553"
  opposite_router_id = apsarastack_vpc.default.0.router_id
  opposite_router_type = "VRouter"
}

resource "apsarastack_router_interface_connection" "bar" {
  provider = "apsarastack"
  interface_id = "${apsarastack_router_interface.opposite.id}"
  opposite_interface_id = "${apsarastack_router_interface.initiating.id}"
  opposite_interface_owner_id =  "1262302482727553"
  opposite_router_id = apsarastack_vpc.default.1.router_id
  opposite_router_type = "VRouter"
}
data "apsarastack_router_interfaces" "default" {
ids = ["${apsarastack_router_interface.initiating.id}"]
}`
