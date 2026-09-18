package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackApiGatewayV2ConsumersDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_api_gateway_v2_consumers.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc%d", rand),
		dataSourceApiGatewayV2ConsumersConfigDependence)

	existMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"consumers.#":                CHECKSET,
			"consumers.0.app_name":       fmt.Sprintf("tf-testacc%d", rand),
			"consumers.0.description":    fmt.Sprintf("tf-testacc%d", rand),
			"consumers.0.auth_type":      "1",
			"consumers.0.auth_type_name": "BASIC",
			"consumers.0.app_id":         CHECKSET,
			"consumers.0.use_white_list": "false",
			"consumers.0.enable":         "true",
			"consumers.0.groups.#":       "1",
			"consumers.0.groups.0":       "test",
			"ids.#":                      CHECKSET,
			"ids.0":                      CHECKSET,
		}
	}

	fakeMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"consumers.#": "0",
			"ids.#":       "0",
		}
	}

	dsa := dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}

	// Test case 1: Basic test with gw_instance_id
	config1 := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"gw_instance_id": "${alibabacloudstack_api_gateway_v2_consumer.default.gw_instance_id}",
			"appid":          "${alibabacloudstack_api_gateway_v2_consumer.default.app_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"gw_instance_id": "${alibabacloudstack_api_gateway_v2_consumer.default.gw_instance_id}",
			"appid":          "fake-gw-instance-appid",
		}),
	}

	// Test case 2: Filter by name_regex
	config2 := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"gw_instance_id": "${alibabacloudstack_api_gateway_v2_consumer.default.gw_instance_id}",
			"name_regex":     fmt.Sprintf("tf-testacc%d", rand),
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"gw_instance_id": "${alibabacloudstack_api_gateway_v2_consumer.default.gw_instance_id}",
			"name_regex":     "nonexistent-consumer",
		}),
	}

	// Test case 3: Filter by ids
	config3 := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"gw_instance_id": "${alibabacloudstack_api_gateway_v2_consumer.default.gw_instance_id}",
			"ids":            []string{"${alibabacloudstack_api_gateway_v2_consumer.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"gw_instance_id": "${alibabacloudstack_api_gateway_v2_consumer.default.gw_instance_id}",
			"ids":            []string{"fake-consumer-id"},
		}),
	}

	dsa.dataSourceTestCheck(t, rand, config1, config2, config3)
}

func dataSourceApiGatewayV2ConsumersConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alibabacloudstack_api_gateway_v2_instance_types" "default" {
	sorted_by = "CPU"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  	instance_name = "${var.name}"
	node_number = "1"
	instance_class = "${data.alibabacloudstack_api_gateway_v2_instance_types.default.instance_types[0].id}"
	broker_engine_type = "SCG"
	deploy_mode = "custom"
}

resource "alibabacloudstack_api_gateway_v2_consumer" "default" {
	auth_type = "1"
	app_name = "${var.name}"
	description = "${var.name}"
	key = "root"
	password = "admin"
	gw_instance_id = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
	groups = ["test"]
}
`, name)
}
