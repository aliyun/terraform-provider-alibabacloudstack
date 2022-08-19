package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackSlbsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackSlbsDataSource,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_slbs.default"),
					resource.TestCheckResourceAttr("data.apsarastack_slbs.default", "slbs.#", "1"),
					resource.TestCheckResourceAttrSet("data.apsarastack_slbs.default", "ids.#"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})

}

const testAccCheckApsaraStackSlbsDataSource = `
variable "name" {
	default = "tf-SlbDataSourceSlbs"
}
data "apsarastack_zones" "default" {
	available_resource_creation= "VSwitch"
}
resource "apsarastack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}
resource "apsarastack_vswitch" "default" {
  vpc_id = "${apsarastack_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name = "${var.name}"
}
resource "apsarastack_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${apsarastack_vswitch.default.id}"
}
data "apsarastack_slbs" "default" {
 ids = ["${apsarastack_slb.default.id}"]
}
`
