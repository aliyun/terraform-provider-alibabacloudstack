package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackInstancesDataSourceBasic(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackInstancesDataSource,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_instances.default"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_instances.default", "instances.#", "1"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_instances.default", "ids.#"),
				),
			},
		},
	})
}

const testAccCheckAlibabacloudStackInstancesDataSource = DataAlibabacloudstackVswitchZones + DataAlibabacloudstackInstanceTypes + DataAlibabacloudstackImages + `
variable "name" {
  default = "Tf-EcsInstanceDataSource"
}

resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name = "${var.name}"
}
resource "alibabacloudstack_security_group" "default" {
  name = "${var.name}"
  vpc_id = "${alibabacloudstack_vpc.default.id}"
}
resource "alibabacloudstack_instance" "default" {
  image_id = data.alibabacloudstack_images.default.images.0.id
  instance_type = "ecs.e4.customize.undjfvanfg"
  instance_name = "${var.name}"
  internet_max_bandwidth_out = "10"
  security_groups = "${alibabacloudstack_security_group.default.*.id}"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  system_disk_category = "cloud_ssd"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}
data "alibabacloudstack_instances" "default" {
  ids = ["${alibabacloudstack_instance.default.id}"]
}
`
