package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackMqttGroupsDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_mqtt_groups.default"
	rand := getAccTestRandInt(1000, 9999)
	name := fmt.Sprintf("tf_testAcc%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccMqttGroupsConfigDependence)
	testAcc := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                       "1",
				"names.#":                     "1",
				"groups.#":                    "1",
				"groups.0.group_id":           CHECKSET,
				"groups.0.instance_id":        CHECKSET,
				"groups.0.create_time":        CHECKSET,
				"groups.0.update_time":        CHECKSET,
				"groups.0.channel_name":       CHECKSET,
				"groups.0.region_name":        CHECKSET,
				"groups.0.independent_naming": CHECKSET,
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":    "0",
				"names.#":  "0",
				"groups.#": "0",
			}
		},
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_mqtt_instance.default.id}",
			"name_regex":  "^${alibabacloudstack_mqtt_group.default.group_id}$",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_mqtt_instance.default.id}",
			"name_regex":  "^GID_fake.*$",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_mqtt_instance.default.id}",
			"ids":         []string{"${alibabacloudstack_mqtt_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_mqtt_instance.default.id}",
			"ids":         []string{"MQTT_fake:GID_fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_mqtt_instance.default.id}",
			"name_regex":  "^${alibabacloudstack_mqtt_group.default.group_id}$",
			"ids":         []string{"${alibabacloudstack_mqtt_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_mqtt_instance.default.id}",
			"name_regex":  "^GID_fake.*$",
			"ids":         []string{"MQTT_fake:GID_fake"},
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccMqttGroupsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s

resource "alibabacloudstack_mqtt_group" "default" {
  group_id = "GID_${var.name}"
  instance_id = "${alibabacloudstack_mqtt_instance.default.id}"
}
`, name, MqttCommonTestCase)
}
