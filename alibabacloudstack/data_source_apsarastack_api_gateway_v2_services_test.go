package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackAPIGatewayV2ServicesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	testAcc := dataSourceAttr{
		resourceId: "data.alibabacloudstack_api_gateway_v2_services.default",
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
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: ApiGatewayV2ServicesDataSourceNew(rand, map[string]string{
			"name_regex": `"testtf-apigw-*"`,
		}),
		fakeConfig: ApiGatewayV2ServicesDataSourceNew(rand, map[string]string{
			"name_regex": `"tf-apigw-fakeservice*"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: ApiGatewayV2ServicesDataSourceNew(rand, map[string]string{
			"ids": `["${alibabacloudstack_api_gateway_v2_service.default.id}"]`,
		}),
		fakeConfig: ApiGatewayV2ServicesDataSourceNew(rand, map[string]string{
			"ids": `["fake-id"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: ApiGatewayV2ServicesDataSourceNew(rand, map[string]string{
			"name_regex": `"testtf-apigw-*"`,
			"ids":        `["${alibabacloudstack_api_gateway_v2_service.default.id}"]`,
		}),
		fakeConfig: ApiGatewayV2ServicesDataSourceNew(rand, map[string]string{
			"name_regex": `"tf-apigw-fakeservice*"`,
			"ids":        `["fake-id"]`,
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func ApiGatewayV2ServiceCommonTestCaseNew(rand int) string {
	return fmt.Sprintf(`

variable "name" {
  default = "testtf-apigw-%d"
}

data "alibabacloudstack_api_gateway_v2_instance_types" "default" {
	sorted_by = "CPU"
}

resource "alibabacloudstack_api_gateway_v2_k8s_cluster" "default" {
	cs_cluster_id =   "${local.k8s_cluster_id}"
	k8s_cluster_name = "${var.name}"
}

%s

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  	instance_name = "${var.name}"
	node_number = "1"
	instance_class = "mini"
	broker_engine_type = "SCG"
	deploy_mode = "k8s"
	deploy_cluster_code = "${alibabacloudstack_api_gateway_v2_k8s_cluster.default.id}"
	deploy_cluster_namespace = "${var.name}-namespace"
	ingress_class_name = "${var.name}-class"
	sls_enabled = true
	prometheus_enabled = true
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
 `, rand, AckK8sCommonTestCase())
}
func ApiGatewayV2ServicesDataSourceNew(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, "  "+k+" = "+v)
	}
	config := fmt.Sprintf(`
%s

data "alibabacloudstack_api_gateway_v2_services" "default" {
	gw_instance_id = alibabacloudstack_api_gateway_v2_service.default.gw_instance_id
%s
}
`, ApiGatewayV2ServiceCommonTestCaseNew(rand), strings.Join(pairs, "\n"))

	return config
}
