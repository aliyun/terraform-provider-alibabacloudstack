package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackSlbBackendServersDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackSlbBackendServersDataSource,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_slb_backend_servers.default"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_slb_backend_servers.default", "load_balancer_id.#", "0"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

const testAccCheckAlibabacloudStackSlbBackendServersDataSource = DataAlibabacloudstackVswitchZones + DataAlibabacloudstackInstanceTypes + DataAlibabacloudstackImages + `
variable "name" {
	default = "tf-slbBackendServersdatasourcebasic"
}


resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
  name = "${var.name}"
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
}

resource "alibabacloudstack_security_group" "default" {
	name = "${var.name}"
	vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_instance" "default" {
  image_id = "${data.alibabacloudstack_images.default.images.0.id}"

  instance_type = "${local.instance_type_id}"
  system_disk_category = "cloud_efficiency"

  security_groups = ["${alibabacloudstack_security_group.default.id}"]
  instance_name = "${var.name}"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_slb_backend_server" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"

  backend_servers {
    server_id = "${alibabacloudstack_instance.default.id}"
    weight     = 100
  }
}

data "alibabacloudstack_slb_backend_servers" "default" {
 load_balancer_id = "${alibabacloudstack_slb.default.id}"
}
`
