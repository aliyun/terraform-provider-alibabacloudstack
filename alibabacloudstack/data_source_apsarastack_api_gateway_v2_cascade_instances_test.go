package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackApiGatewayV2CascadeInstancesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_api_gateway_v2_cascade_instances.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccApiGwV2%d", rand),
		dataSourceApiGatewayV2CascadeInstancesConfigDependence)

	testAcc := dataSourceAttr{
		resourceId: resourceId,
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
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_api_gateway_v2_cascade_instance.default.instance_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^fake-name.*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_api_gateway_v2_cascade_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"fake-id"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_api_gateway_v2_cascade_instance.default.instance_name}",
			"ids":        []string{"${alibabacloudstack_api_gateway_v2_cascade_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_api_gateway_v2_cascade_instance.default.instance_name}",
			"ids":        []string{"fake-id"},
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceApiGatewayV2CascadeInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name      = "${var.name}source"
  node_number        = 1
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

resource "alibabacloudstack_api_gateway_v2_cascade_instance" "default" {
  instance_name       = "${var.name}"
  cascade_instance_id = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
}
`, name)
}
