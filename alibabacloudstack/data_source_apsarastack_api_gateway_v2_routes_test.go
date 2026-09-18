package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackApiGatewayV2RoutesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_api_gateway_v2_routes.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacca-%d", rand),
		dataSourceApiGatewayV2RoutesConfigDependence)

	testAcc := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                  "1",
				"routes.#":               "1",
				"routes.0.route_name":    fmt.Sprintf("tf-testacca-%d", rand),
				"routes.0.group_id":      "DEFAULT",
				"routes.0.enable_status": "true",
				"routes.0.route_path.#":  "2",
				"routes.0.route_path.0":  "/testtc",
				"routes.0.route_path.1":  "/test/aaa",
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":    "0",
				"routes.#": "0",
			}
		},
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"gw_instance_id": "${alibabacloudstack_api_gateway_v2_route.default.gw_instance_id}",
			"name_regex":     fmt.Sprintf("tf-testacca-%d", rand),
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"gw_instance_id": "${alibabacloudstack_api_gateway_v2_route.default.gw_instance_id}",
			"name_regex":     "^fake-name$",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"gw_instance_id": "${alibabacloudstack_api_gateway_v2_route.default.gw_instance_id}",
			"ids":            []string{"${alibabacloudstack_api_gateway_v2_route.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"gw_instance_id": "${alibabacloudstack_api_gateway_v2_route.default.gw_instance_id}",
			"ids":            []string{"fake-route-id"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"gw_instance_id": "${alibabacloudstack_api_gateway_v2_route.default.gw_instance_id}",
			"ids":            []string{"${alibabacloudstack_api_gateway_v2_route.default.id}"},
			"name_regex":     fmt.Sprintf("tf-testacca-%d", rand),
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"gw_instance_id": "${alibabacloudstack_api_gateway_v2_route.default.gw_instance_id}",
			"ids":            []string{"fake-route-id"},
			"name_regex":     "^fake-name$",
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, allConf, nameRegexConf, idsConf)
}

func dataSourceApiGatewayV2RoutesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name      = var.name
  node_number        = "1"
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

resource "alibabacloudstack_api_gateway_v2_service" "default" {
  name = "${var.name}"
  description = "${var.name}"
  protocol = "HTTP"
  upstream_type = "1"
  load_balance_type = "1"
  gw_instance_id = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
  service_nodes {
    ip = "127.0.0.1"
    port = "80"
    weight = "100"
    enable = "true"
  }
}

resource "alibabacloudstack_api_gateway_v2_domain" "default" {
  domain         = "${var.name}.com"
  instance_id    = alibabacloudstack_api_gateway_v2_instance.default.id
  protocol       = "HTTP"
}

resource "alibabacloudstack_api_gateway_v2_route" "default" {
	gw_instance_id = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
	route_name     = "${var.name}"
	strip_prefix   = "2"
	order          = "100"
	route_path     = ["/testtc", "/test/aaa"]
	methods        = ["GET", "POST", "PUT", "DELETE"]
	header {
		key   = "header"
		value = "aaaaa"
	}
	cookie {
		key   = "cookie"
		value = "bbbbb"
	}
	query_param {
		key   = "query"
		value = "ccccc"
	}
	domain_ids = ["${alibabacloudstack_api_gateway_v2_domain.default.domain_id}"]
	service_id = "${alibabacloudstack_api_gateway_v2_service.default.service_id}"
}
`, name)
}
