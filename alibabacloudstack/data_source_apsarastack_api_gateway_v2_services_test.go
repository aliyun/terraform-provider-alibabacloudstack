package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackApiGatewayV2ServicesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_api_gateway_v2_services.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("testtf-apigw-%d", rand),
		dataSourceApiGatewayV2ServicesConfigDependence)

	testAcc := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                    "1",
				"names.#":                  "1",
				"services.#":               "1",
				"services.0.name":          fmt.Sprintf("testtf-apigw-%d", rand),
				"services.0.description":   fmt.Sprintf("testtf-apigw-%d", rand),
				"services.0.protocol":      "HTTP",
				"services.0.upstream_type": "1",
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":      "0",
				"names.#":    "0",
				"services.#": "0",
			}
		},
		ExternalProviders: testAccExternalProviders,
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"gw_instance_id": "${alibabacloudstack_api_gateway_v2_service.default.gw_instance_id}",
			"name_regex":     "testtf-apigw-*",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"gw_instance_id": "${alibabacloudstack_api_gateway_v2_service.default.gw_instance_id}",
			"name_regex":     "tf-apigw-fakeservice*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"gw_instance_id": "${alibabacloudstack_api_gateway_v2_service.default.gw_instance_id}",
			"ids":            []string{"${alibabacloudstack_api_gateway_v2_service.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"gw_instance_id": "${alibabacloudstack_api_gateway_v2_service.default.gw_instance_id}",
			"ids":            []string{"fake-id"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"gw_instance_id": "${alibabacloudstack_api_gateway_v2_service.default.gw_instance_id}",
			"name_regex":     "testtf-apigw-*",
			"ids":            []string{"${alibabacloudstack_api_gateway_v2_service.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"gw_instance_id": "${alibabacloudstack_api_gateway_v2_service.default.gw_instance_id}",
			"name_regex":     "tf-apigw-fakeservice*",
			"ids":            []string{"fake-id"},
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceApiGatewayV2ServicesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

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
  health_check_struct {
	type = "1"
	health_path = "/check"
	http_statuses = "200"
	timeout = "20000"
	health_interval = "30"
	un_health_interval = "30"
	http_successes = "1"
	http_failures = "0"
  }
}
`, name, ApiGatwayV2K8sInstanceTestCase("SCG", "apig_k8s"))
}
