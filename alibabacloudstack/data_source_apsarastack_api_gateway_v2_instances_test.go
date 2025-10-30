package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackAPIGateWayV2InstancesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_api_gateway_v2_instances.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("testtf-apigw-%d", rand),
		dataSourceAPIGateWayV2InstancesDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_api_gateway_v2_instance.default.instance_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_api_gateway_v2_instance.default.instance_name}-fakeTestAcccc",
		}),
	}

	deployModeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"deploy_mode": "custom",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"deploy_mode": "edas",
		}),
	}

	brokerEngineTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"broker_engine_type": "SCG",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"broker_engine_type": "HIGRESS",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_api_gateway_v2_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_api_gateway_v2_instance.default.id}-fakeTestAcccc"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_api_gateway_v2_instance.default.instance_name}",
			"ids":        []string{"${alibabacloudstack_api_gateway_v2_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_api_gateway_v2_instance.default.instance_name}-fakeTestAcccc",
			"ids":        []string{"${alibabacloudstack_api_gateway_v2_instance.default.id}-fakeTestAcccc"},
		}),
	}

	var existAPIGateWayV2InstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              "1",
			"ids.0":                              CHECKSET,
			"instances.#":                        "1",
			"instances.0.instance_name":          fmt.Sprintf("testtf-apigw-%d", rand),
			"instances.0.broker_engine_type":     "SCG",
			"instances.0.deploy_mode":            "custom",
			"instances.0.instance_class":         "mini",
			"instances.0.custom_deploy_config.%": "4",
		}
	}

	var fakeAPIGateWayV2InstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}

	var APIGateWayV2InstancesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAPIGateWayV2InstancesMapFunc,
		fakeMapFunc:  fakeAPIGateWayV2InstancesMapFunc,
	}

	APIGateWayV2InstancesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, deployModeConf, brokerEngineTypeConf, idsConf, allConf)
}

func dataSourceAPIGateWayV2InstancesDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  	instance_name = "${var.name}"
	node_number = "1"
	instance_class = "mini"
	broker_engine_type = "SCG"
	deploy_mode = "custom"
}

 `, name)
}
