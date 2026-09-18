package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackApiGatewayV2SignaturesDataSource(t *testing.T) {
	rand := -1
	resourceId := "data.alibabacloudstack_api_gateway_v2_signatures.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("test-%d", rand),
		dataSourceApiGatewayV2SignaturesConfigDependence)

	testAcc := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                        "1",
				"names.#":                      "1",
				"signatures.#":                 "1",
				"signatures.0.sig_scheme_name": fmt.Sprintf("test-%d", rand),
				"signatures.0.sig_alg":         "HmacSM3",
				"signatures.0.status":          "0",
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":        "0",
				"names.#":      "0",
				"signatures.#": "0",
			}
		},
	}

	configBasic := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"gw_instance_id": "${alibabacloudstack_api_gateway_v2_signature.default.gw_instance_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"gw_instance_id": "fake-id",
		}),
	}

	configByNameRegex := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"gw_instance_id": "${alibabacloudstack_api_gateway_v2_signature.default.gw_instance_id}",
			"name_regex":     "test-.*",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"gw_instance_id": "${alibabacloudstack_api_gateway_v2_signature.default.gw_instance_id}",
			"name_regex":     "nonexistent-.*",
		}),
	}

	configByIds := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"gw_instance_id": "${alibabacloudstack_api_gateway_v2_signature.default.gw_instance_id}",
			"ids":            []string{"${alibabacloudstack_api_gateway_v2_signature.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"gw_instance_id": "${alibabacloudstack_api_gateway_v2_signature.default.gw_instance_id}",
			"ids":            []string{"fake-id"},
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, configBasic, configByNameRegex, configByIds)
}

func dataSourceApiGatewayV2SignaturesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "signature_name" {
  default = "%s"
}

variable "signature_algorithm" {
  default = "HmacSM3"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name      = "testtf-apigw-instance-%s"
  node_number        = "1"
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

data "alibabacloudstack_api_gateway_v2_instances" "default" {
  name_regex = alibabacloudstack_api_gateway_v2_instance.default.instance_name
}

resource "alibabacloudstack_api_gateway_v2_signature" "default" {
  sig_scheme_name = var.signature_name
  sig_alg         = var.signature_algorithm
  gw_instance_id  = alibabacloudstack_api_gateway_v2_instance.default.id
}
`, name, name)
}
