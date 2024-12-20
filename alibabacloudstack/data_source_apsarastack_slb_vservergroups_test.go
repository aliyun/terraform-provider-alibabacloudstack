package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackSlbServerGroupsDataSource_basic(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackSlbServerGroupsDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_slb_server_groups.default"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_slb_server_groups.default", "load_balancer_id.#", "0"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

const testAccCheckAlibabacloudStackSlbServerGroupsDataSourceBasic = DataAlibabacloudstackVswitchZones + DataAlibabacloudstackInstanceTypes + DataAlibabacloudstackImages + `
variable "name" {
	default = "tf-testAccslbservergroupsdatasourcebasic"
}

data "alibabacloudstack_zones" "slbverver" {
	available_resource_creation= "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}
resource "alibabacloudstack_vswitch" "default" {
  name = "${var.name}"
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = data.alibabacloudstack_zones.slbverver.zones.0.id
}
resource "alibabacloudstack_security_group" "default" {
	name = "${var.name}"
	vpc_id = "${alibabacloudstack_vpc.default.id}"
}
resource "alibabacloudstack_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}
resource "alibabacloudstack_slb_listener" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "http"
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie = "${var.name}"
  cookie_timeout = 86400
  health_check = "on"
  health_check_uri = "/cons"
  health_check_connect_port = 20
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  bandwidth = 10
  x_forwarded_for  {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
}
resource "alibabacloudstack_instance" "default" {
  image_id = "${data.alibabacloudstack_images.default.images.0.id}"
  availability_zone = data.alibabacloudstack_zones.slbverver.zones.0.id
  instance_type = "${local.default_instance_type_id}"
  system_disk_category = "cloud_efficiency"
  security_groups = ["${alibabacloudstack_security_group.default.id}"]
  instance_name = "${var.name}"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}
resource "alibabacloudstack_slb_server_group" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  name = "${var.name}"
  servers {
      server_ids = ["${alibabacloudstack_instance.default.id}"]
      port = 80
      weight = 100
    }
}
resource "alibabacloudstack_slb_rule" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  frontend_port = "${alibabacloudstack_slb_listener.default.frontend_port}"
  name = "${var.name}"
  domain = "*.aliyun.com"
  url = "/image"
  server_group_id = "${alibabacloudstack_slb_server_group.default.id}"
}
data "alibabacloudstack_slb_server_groups" "default" {
 load_balancer_id = "${alibabacloudstack_slb.default.id}"
  
}`
