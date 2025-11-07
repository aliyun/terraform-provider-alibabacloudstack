package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackApiGatewayV2SignaturesDataSource(t *testing.T) {
	testAcc := dataSourceAttr{
		resourceId: "data.alibabacloudstack_api_gateway_v2_signatures.default",
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

	testAcc.dataSourceTestCheck(t, -1,
		testAccCheckAlibabacloudStackAPIGatewayV2SignaturesDataSourceConfigBasic,
		testAccCheckAlibabacloudStackAPIGatewayV2SignaturesDataSourceConfigByNameRegex,
		testAccCheckAlibabacloudStackAPIGatewayV2SignaturesDataSourceConfigByIds,
	)
}

func resourceApiGatewayV2SignatureDependenceNew(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
variable "signature_name" {
  default = "test-%d"
}

variable "signature_algorithm" {
  default = "HmacSM3"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name      = "testtf-apigw-instance-%d"
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

data "alibabacloudstack_api_gateway_v2_signatures" "default" {
  %s
}
`, rand, rand, strings.Join(pairs, "\n  "))
}

var testAccCheckAlibabacloudStackAPIGatewayV2SignaturesDataSourceConfigBasic = dataSourceTestAccConfig{
	existConfig: resourceApiGatewayV2SignatureDependenceNew(-1, map[string]string{
		"gw_instance_id": "alibabacloudstack_api_gateway_v2_signature.default.gw_instance_id",
	}),
	fakeConfig: resourceApiGatewayV2SignatureDependenceNew(-1, map[string]string{
		"gw_instance_id": "\"fake-id\"",
	}),
}

var testAccCheckAlibabacloudStackAPIGatewayV2SignaturesDataSourceConfigByNameRegex = dataSourceTestAccConfig{
	existConfig: resourceApiGatewayV2SignatureDependenceNew(-1, map[string]string{
		"gw_instance_id": "alibabacloudstack_api_gateway_v2_signature.default.gw_instance_id",
		"name_regex":     "\"test-.*\"",
	}),
	fakeConfig: resourceApiGatewayV2SignatureDependenceNew(-1, map[string]string{
		"gw_instance_id": "alibabacloudstack_api_gateway_v2_signature.default.gw_instance_id",
		"name_regex":     "\"nonexistent-.*\"",
	}),
}

var testAccCheckAlibabacloudStackAPIGatewayV2SignaturesDataSourceConfigByIds = dataSourceTestAccConfig{
	existConfig: resourceApiGatewayV2SignatureDependenceNew(-1, map[string]string{
		"gw_instance_id": "alibabacloudstack_api_gateway_v2_signature.default.gw_instance_id",
		"ids":            "[alibabacloudstack_api_gateway_v2_signature.default.id]",
	}),
	fakeConfig: resourceApiGatewayV2SignatureDependenceNew(-1, map[string]string{
		"gw_instance_id": "alibabacloudstack_api_gateway_v2_signature.default.gw_instance_id",
		"ids":            "[\"fake-id\"]",
	}),
}
