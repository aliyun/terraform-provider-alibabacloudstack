package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackSlbRulesDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_slb_rules.default"
	name := fmt.Sprintf("tf-testaccslbrulesdatasourcebasic%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceSlbRulesConfigDependence)

	frontendPortConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_rule.default.load_balancer_id}",
			"frontend_port":    "80",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_rule.default.load_balancer_id}",
			"frontend_port":    "999",
		}),
	}

	name_regexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_rule.default.load_balancer_id}",
			"frontend_port":    "80",
			"name_regex":       "${alibabacloudstack_slb_rule.default.rule_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_rule.default.load_balancer_id}",
			"frontend_port":    "80",
			"name_regex":       "${alibabacloudstack_slb_rule.default.rule_name}_fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_rule.default.load_balancer_id}",
			"frontend_port":    "80",
			"ids":              []string{"${alibabacloudstack_slb_rule.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"load_balancer_id": "${alibabacloudstack_slb_rule.default.load_balancer_id}",
			"frontend_port":    "80",
			"ids":              []string{"${alibabacloudstack_slb_rule.default.id}_fake"},
		}),
	}

	var existSlbRulesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"slb_rules.#":                 "1",
			"slb_rules.0.id":              CHECKSET,
			"slb_rules.0.name":            name,
			"slb_rules.0.domain":          "*.aliyun.com",
			"slb_rules.0.url":             "/image",
			"slb_rules.0.server_group_id": CHECKSET,
		}
	}

	var fakeSlbRulesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"slb_rules.#": "0",
		}
	}

	var slbRulesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existSlbRulesMapFunc,
		fakeMapFunc:  fakeSlbRulesMapFunc,
	}
	slbRulesCheckInfo.dataSourceTestCheck(t, rand, frontendPortConf, idsConf, name_regexConf)
}

func dataSourceSlbRulesConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	%s


	resource "alibabacloudstack_slb" "default" {
	  name = "${var.name}"
	  vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
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
	}

	resource "alibabacloudstack_slb_server_group" "default" {
	  load_balancer_id = "${alibabacloudstack_slb.default.id}"
    name = "${var.name}"
	  servers {
	      server_ids = ["${alibabacloudstack_ecs_instance.default.id}"]
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
	`, name, ECSInstanceCommonTestCase)
}
