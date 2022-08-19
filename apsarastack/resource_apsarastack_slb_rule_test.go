package apsarastack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackSlbRuleCreate(t *testing.T) {
	var v *slb.DescribeRuleAttributeResponse
	resourceId := "apsarastack_slb_rule.default"
	ra := resourceAttrInit(resourceId, ruleMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccSlbRuleBasic")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbRuleBasicDependence)
	resource.Test(t, resource.TestCase{
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
					"name":                      "${var.name}",
					"load_balancer_id":          "${apsarastack_slb.default.id}",
					"frontend_port":             "${apsarastack_slb_listener.default.frontend_port}",
					"domain":                    "*.aliyun.com",
					"url":                       "/image",
					"server_group_id":           "${apsarastack_slb_server_group.default.id}",
					"cookie":                    "23ffsa",
					"cookie_timeout":            "100",
					"health_check_http_code":    "http_2xx",
					"health_check_interval":     "10",
					"health_check_uri":          "/test",
					"health_check_connect_port": "80",
					"health_check_timeout":      "10",
					"healthy_threshold":         "3",
					"unhealthy_threshold":       "3",
					"sticky_session":            "on",
					"sticky_session_type":       "server",
					"listener_sync":             "on",
					"scheduler":                 "rr",
					"health_check_domain":       "test",
					"health_check":              "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "tf-testAccSlbRuleBasic_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccSlbRuleBasic_change",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func resourceSlbRuleBasicDependence(name string) string {
	return fmt.Sprintf(`
%s

%s

%s
provider "apsarastack" {
	assume_role {}
}
variable "name" {
  default = "%s"
}

resource "apsarastack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "apsarastack_vswitch" "default" {
  vpc_id = "${apsarastack_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = data.apsarastack_zones.default.zones.0.id
  name = "${var.name}"
}
resource "apsarastack_security_group" "default" {
  name = "${var.name}"
  vpc_id = "${apsarastack_vpc.default.id}"
}

resource "apsarastack_instance" "default" {
  image_id = "${data.apsarastack_images.default.images.0.id}"
  instance_type = "${local.instance_type_id}"
  security_groups = "${apsarastack_security_group.default.*.id}"
  internet_max_bandwidth_out = "10"
  availability_zone = data.apsarastack_zones.default.zones.0.id
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${apsarastack_vswitch.default.id}"
  instance_name = "${var.name}"
}

resource "apsarastack_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${apsarastack_vswitch.default.id}"
}

resource "apsarastack_slb_listener" "default" {
  load_balancer_id = "${apsarastack_slb.default.id}"
  backend_port = 22
  frontend_port = 22
  protocol = "http"
  bandwidth = 5
  health_check_connect_port = "20"
  health_check = "on"
  sticky_session = "off"
}

resource "apsarastack_slb_server_group" "default" {
  load_balancer_id = "${apsarastack_slb.default.id}"
  servers {
      server_ids = "${apsarastack_instance.default.*.id}"
      port = 80
      weight = 100
    }
}
`, DataApsarastackVswitchZones, DataApsarastackInstanceTypes, DataApsarastackImages, name)
}

var ruleMap = map[string]string{
	"load_balancer_id": CHECKSET,
	"frontend_port":    "22",
	"domain":           "*.aliyun.com",
	"url":              "/image",
	"server_group_id":  CHECKSET,
}
