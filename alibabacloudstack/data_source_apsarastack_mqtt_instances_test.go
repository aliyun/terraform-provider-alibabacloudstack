package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackMqttInstancesDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_mqtt_instances.default"
	rand := getAccTestRandInt(1000, 9999)
	name := fmt.Sprintf("tf_instncedata%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccMqttInstancesConfigDependence)
	testAcc := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                            "1",
				"names.#":                          "1",
				"instances.#":                      "1",
				"instances.0.instance_id":          CHECKSET,
				"instances.0.instance_name":        CHECKSET,
				"instances.0.max_conn":             CHECKSET,
				"instances.0.max_sub":              CHECKSET,
				"instances.0.max_up_tps":           CHECKSET,
				"instances.0.max_down_tps":         CHECKSET,
				"instances.0.independent_naming":   CHECKSET,
				"instances.0.store_instance_id":    CHECKSET,
				"instances.0.store_type":           CHECKSET,
				"instances.0.create_time":          CHECKSET,
				"instances.0.instance_status":      CHECKSET,
				"instances.0.instance_type":        CHECKSET,
				"instances.0.namespace_rules_type": CHECKSET,
				"instances.0.order_id":             CHECKSET,
				"instances.0.sp_instance_type":     CHECKSET,
				"instances.0.max_tps":              CHECKSET,
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
			"name_regex": "${alibabacloudstack_mqtt_instance.default.instance_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^test-fake.*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_mqtt_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"MQTT-fake-id"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_mqtt_instance.default.instance_name}",
			"ids":        []string{"${alibabacloudstack_mqtt_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_mqtt_instance.default.instance_name}",
			"ids":        []string{"MQTT-fake-id"},
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccMqttInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

%s

`, name, MqttCommonTestCase)
}
