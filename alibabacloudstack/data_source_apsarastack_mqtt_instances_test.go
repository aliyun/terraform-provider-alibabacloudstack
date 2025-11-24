package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackMqttInstancesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	testAcc := dataSourceAttr{
		resourceId: "data.alibabacloudstack_mqtt_instances.default",
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
		existConfig: testAccMqttInstancesConfigDependence(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_mqtt_instance.default.instance_name}"`,
		}),
		fakeConfig: testAccMqttInstancesConfigDependence(rand, map[string]string{
			"name_regex": `"^test-fake.*"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccMqttInstancesConfigDependence(rand, map[string]string{
			"ids": `["${alibabacloudstack_mqtt_instance.default.id}"]`,
		}),
		fakeConfig: testAccMqttInstancesConfigDependence(rand, map[string]string{
			"ids": `["MQTT-fake-id"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccMqttInstancesConfigDependence(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_mqtt_instance.default.instance_name}"`,
			"ids":        `["${alibabacloudstack_mqtt_instance.default.id}"]`,
		}),
		fakeConfig: testAccMqttInstancesConfigDependence(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_mqtt_instance.default.instance_name}"`,
			"ids":        `["MQTT-fake-id"]`,
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccMqttInstancesConfigDependence(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "%v"
}

resource "alibabacloudstack_ons_instance" "default" {
  tps_receive_max = 500
  tps_send_max = 500
  topic_capacity = 50
  cluster = "cluster1"
  independent_naming = "true"
  name = "${var.name}MQ"
  remark = "Ons_instance"
}

resource "alibabacloudstack_mqtt_instance" "default" {
  instance_name = "${var.name}"
  remark = "Mqtt"
  max_conn = 1000
  max_sub = 1000
  max_up_tps = 1000
  max_down_tps = 1000
  independent_naming = true
  store_instance_id = "${alibabacloudstack_ons_instance.default.id}"
}

data "alibabacloudstack_mqtt_instances" "default" {
	%s
}

`, rand, strings.Join(pairs, "\n   "))
	return config
}
