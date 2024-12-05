package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackSlbsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackSlbsDataSource,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_slbs.default"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_slbs.default", "slbs.#", "1"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_slbs.default", "ids.#"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})

}

const testAccCheckAlibabacloudStackSlbsDataSource = `
variable "name" {
	default = "tf-SlbDataSourceSlbs"
}
data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}
resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}
resource "alibabacloudstack_vswitch" "default" {
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name = "${var.name}"
}
resource "alibabacloudstack_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}
data "alibabacloudstack_slbs" "default" {
 ids = ["${alibabacloudstack_slb.default.id}"]
}
`
