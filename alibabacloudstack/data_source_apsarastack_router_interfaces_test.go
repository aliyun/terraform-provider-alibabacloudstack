package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackRouterInterfacesDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackRouterInterfacesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_router_interfaces.default"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_router_interfaces.default", "interfaces.#", "1"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_router_interfaces.default", "ids.#"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})

}

const testAccCheckAlibabacloudStackRouterInterfacesDataSourceConfig = `
provider "alibabacloudstack" {
  region = "${var.region}"
}
variable "region" {
  default = "cn-wulan-env200-d01"
}
variable "name" {
  default = "tf-testAccCheckAlibabacloudStackRouterInterfacesDataSourceConfig"
}
variable cidr_block_list {
	type = "list"
	default = [ "172.16.0.0/12", "192.168.0.0/16" ]
}

resource "alibabacloudstack_vpc" "default" {
  count = 2
  name = "${var.name}"
  cidr_block = "${element(var.cidr_block_list,count.index)}"
}
resource "alibabacloudstack_router_interface" "initiating" {
  opposite_region = "${var.region}"
  router_type = "VRouter"
  router_id = "${alibabacloudstack_vpc.default.0.router_id}"
  role = "InitiatingSide"
  specification = "Large.2"
  name = "${var.name}_initiating"
  description = "${var.name}_decription"

}
resource "alibabacloudstack_router_interface" "opposite" {
  provider = "alibabacloudstack"
  opposite_region = "${var.region}"
  router_type = "VRouter"
  router_id = "${alibabacloudstack_vpc.default.1.router_id}"
  role = "AcceptingSide"
  specification = "Large.1"
  name = "${var.name}_opposite"
  description = "${var.name}_decription"

}

resource "alibabacloudstack_router_interface_connection" "foo" {
  interface_id = "${alibabacloudstack_router_interface.initiating.id}"
  opposite_interface_id = "${alibabacloudstack_router_interface.opposite.id}"
  depends_on = ["alibabacloudstack_router_interface_connection.bar"]
  opposite_interface_owner_id = "1262302482727553"
  opposite_router_id = alibabacloudstack_vpc.default.0.router_id
  opposite_router_type = "VRouter"
}

resource "alibabacloudstack_router_interface_connection" "bar" {
  provider = "alibabacloudstack"
  interface_id = "${alibabacloudstack_router_interface.opposite.id}"
  opposite_interface_id = "${alibabacloudstack_router_interface.initiating.id}"
  opposite_interface_owner_id =  "1262302482727553"
  opposite_router_id = alibabacloudstack_vpc.default.1.router_id
  opposite_router_type = "VRouter"
}
data "alibabacloudstack_router_interfaces" "default" {
ids = ["${alibabacloudstack_router_interface.initiating.id}"]
}`
