package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackSlbRule_basic(t *testing.T) {
	var v *slb.DescribeRuleAttributeResponse
	resourceId := "alibabacloudstack_slb_rule.default"
	ra := resourceAttrInit(resourceId, ruleMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testAccSlbRuleBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbRuleBasicDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                         name,
					"delete_protection_validation": true,
					"load_balancer_id":             "${alibabacloudstack_slb_server_group.default.load_balancer_id}",
					"frontend_port":                "${alibabacloudstack_slb_listener.default.frontend_port}",
					"domain":                       "*.aliyun.com",
					"url":                          "/image",
					"server_group_id":              "${alibabacloudstack_slb_server_group.default.id}",
					"listener_sync":                "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":             name,
						"load_balancer_id": CHECKSET,
						"frontend_port":    CHECKSET,
						"domain":           "*.aliyun.com",
						"url":              "/image",
						"server_group_id":  CHECKSET,
						"listener_sync":    "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cookie":                    "23ffsa",
					"health_check_http_code":    "http_2xx",
					"health_check_interval":     "10",
					"health_check_uri":          "/test",
					"health_check_connect_port": "80",
					"health_check_timeout":      "10",
					"healthy_threshold":         "3",
					"unhealthy_threshold":       "3",
					"sticky_session":            "on",
					"sticky_session_type":       "server",
					"listener_sync":             "off",
					"scheduler":                 "rr",
					"health_check_domain":       "test.com",
					"health_check":              "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cookie":                    "23ffsa",
						"health_check_http_code":    "http_2xx",
						"health_check_interval":     "10",
						"health_check_uri":          "/test",
						"health_check_connect_port": "80",
						"health_check_timeout":      "10",
						"healthy_threshold":         "3",
						"unhealthy_threshold":       "3",
						"sticky_session":            "on",
						"sticky_session_type":       "server",
						"listener_sync":             "off",
						"scheduler":                 "rr",
						"health_check_domain":       "test.com",
						"health_check":              "on",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"cookie_timeout":      "100",
					"sticky_session":      "on",
					"sticky_session_type": "insert",
					"listener_sync":       "off",
					"scheduler":           "wrr",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"cookie_timeout":      "100",
						"cookie":              "",
						"sticky_session":      "on",
						"sticky_session_type": "insert",
						"listener_sync":       "off",
						"scheduler":           "wrr",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// delete_protection_validation is a local attribute and cannot be loaded from the remote
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
		},
	})
}

func resourceSlbRuleBasicDependence(name string) string {
	return fmt.Sprintf(`
%s

%s

%s

variable "name" {
  default = "%s"
}

resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
  name = "${var.name}"
}
resource "alibabacloudstack_security_group" "default" {
  name = "${var.name}"
  vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_instance" "default" {
  image_id = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type = "${local.default_instance_type_id}"
  security_groups = "${alibabacloudstack_security_group.default.*.id}"
  internet_max_bandwidth_out = "10"
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
  system_disk_category = "cloud_sperf"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
  instance_name = "${var.name}"
}

resource "alibabacloudstack_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_slb_listener" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  backend_port = 22
  frontend_port = 22
  protocol = "http"
  bandwidth = 5
  health_check_connect_port = "20"
  health_check = "on"
  sticky_session = "off"
}

resource "alibabacloudstack_slb_server_group" "default" {
  load_balancer_id = "${alibabacloudstack_slb_listener.default.load_balancer_id}"
  name = "${var.name}"
  servers {
      server_ids = "${alibabacloudstack_instance.default.*.id}"
      port = 80
      weight = 100
    }
}
`, DataAlibabacloudstackVswitchZones, DataAlibabacloudstackInstanceTypes, DataAlibabacloudstackImages, name)
}

var ruleMap = map[string]string{
	"load_balancer_id": CHECKSET,
	"frontend_port":    "22",
	"domain":           "*.aliyun.com",
	"url":              "/image",
	"server_group_id":  CHECKSET,
}
