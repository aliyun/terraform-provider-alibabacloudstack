package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackAPIGatewayV2InstanceTypesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_api_gateway_v2_instance_types.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccAPIGatewayV2InstanceTypesDataSource_%d", rand),
		dataSourceAPIGatewayV2InstanceTypesConfigDependence)

	cpuConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cpu": "2",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cpu": "200",
		}),
	}
	memoryConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"memory": "2",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"memory": "3200",
		}),
	}

	var existAPIGatewayV2InstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           CHECKSET,
			"ids.0":                           CHECKSET,
			"instance_types.#":                CHECKSET,
			"instance_types.0.id":             CHECKSET,
			"instance_types.0.cpu":            CHECKSET,
			"instance_types.0.memory":         CHECKSET,
			"instance_types.0.instance_class": CHECKSET,
		}
	}

	var fakeAPIGatewayV2InstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":            "0",
			"instance_types.#": "0",
		}
	}

	var APIGatewayV2InstanceTypesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAPIGatewayV2InstanceTypesMapFunc,
		fakeMapFunc:  fakeAPIGatewayV2InstanceTypesMapFunc,
	}

	APIGatewayV2InstanceTypesCheckInfo.dataSourceTestCheck(t, rand, cpuConf, memoryConf)
}

func dataSourceAPIGatewayV2InstanceTypesConfigDependence(name string) string {
	return ""
}
