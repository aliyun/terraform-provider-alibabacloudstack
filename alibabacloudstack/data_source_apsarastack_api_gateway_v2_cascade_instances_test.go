package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackAPIGatewayV2CascadeInstancesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	testAcc := dataSourceAttr{
		resourceId: "data.alibabacloudstack_api_gateway_v2_cascade_instances.default",
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                           "1",
				"names.#":                         "1",
				"instances.#":                     "1",
				"instances.0.instance_name":       CHECKSET,
				"instances.0.instance_type":       CHECKSET,
				"instances.0.cascade_instance_id": CHECKSET,
				"instances.0.status":              CHECKSET,
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":       "0",
				"names.#":     "0",
				"instances.#": "0",
			}
		},
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: APIGatewayV2CascadeInstancesDependenceNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_api_gateway_v2_cascade_instance.default.instance_name}"`,
		}),
		fakeConfig: APIGatewayV2CascadeInstancesDependenceNew(rand, map[string]string{
			"name_regex": `"^fake-name.*"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: APIGatewayV2CascadeInstancesDependenceNew(rand, map[string]string{
			"ids": `["${alibabacloudstack_api_gateway_v2_cascade_instance.default.id}"]`,
		}),
		fakeConfig: APIGatewayV2CascadeInstancesDependenceNew(rand, map[string]string{
			"ids": `["fake-id"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: APIGatewayV2CascadeInstancesDependenceNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_api_gateway_v2_cascade_instance.default.instance_name}"`,
			"ids":        `["${alibabacloudstack_api_gateway_v2_cascade_instance.default.id}"]`,
		}),
		fakeConfig: APIGatewayV2CascadeInstancesDependenceNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_api_gateway_v2_cascade_instance.default.instance_name}"`,
			"ids":        `["fake-id"]`,
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

// Dependence template generation method
func APIGatewayV2CascadeInstancesDependenceNew(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccApiGwV2%d"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name      = "${var.name}source"
  node_number        = 1
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

resource "alibabacloudstack_api_gateway_v2_cascade_instance" "default" {
  instance_type       = "0"
  instance_name       = "${var.name}"
  cascade_instance_id = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
}

data "alibabacloudstack_api_gateway_v2_cascade_instances" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
}
