package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackInstancesDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackInstancesDataSource,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_instances.default"),
					resource.TestCheckResourceAttr("data.apsarastack_instances.default", "instances.#", "1"),
					resource.TestCheckResourceAttrSet("data.apsarastack_instances.default", "ids.#"),
				),
			},
		},
	})
}

const testAccCheckApsaraStackInstancesDataSource = DataApsarastackVswitchZones + DataApsarastackInstanceTypes + DataApsarastackImages + `
variable "name" {
  default = "Tf-EcsInstanceDataSource"
}

resource "apsarastack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "apsarastack_vswitch" "default" {
  vpc_id = "${apsarastack_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name = "${var.name}"
}
resource "apsarastack_security_group" "default" {
  name = "${var.name}"
  vpc_id = "${apsarastack_vpc.default.id}"
}
resource "apsarastack_instance" "default" {
  image_id = data.apsarastack_images.default.images.0.id
  instance_type = local.instance_type_id
  instance_name = "${var.name}"
  internet_max_bandwidth_out = "10"
  security_groups = "${apsarastack_security_group.default.*.id}"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  system_disk_category = "cloud_ssd"
  vswitch_id = "${apsarastack_vswitch.default.id}"
}
data "apsarastack_instances" "default" {
  ids = ["${apsarastack_instance.default.id}"]
}
`
