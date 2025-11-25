package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackMqttGroupsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	testAcc := dataSourceAttr{
		resourceId: "data.alibabacloudstack_mqtt_groups.default",
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
		existConfig: testAccMqttGroupsConfigDependenceNew(rand, map[string]string{
			"name_regex": `"GID_testacc.*$"`,
		}),
		fakeConfig: testAccMqttGroupsConfigDependenceNew(rand, map[string]string{
			"name_regex": `"^GID_fake.*$"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccMqttGroupsConfigDependenceNew(rand, map[string]string{
			"ids": `["${alibabacloudstack_mqtt_group.default.id}"]`,
		}),
		fakeConfig: testAccMqttGroupsConfigDependenceNew(rand, map[string]string{
			"ids": `["MQTT_fake:GID_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccMqttGroupsConfigDependenceNew(rand, map[string]string{
			"name_regex": `"GID_testacc.*$"`,
			"ids":        `["${alibabacloudstack_mqtt_group.default.id}"]`,
		}),
		fakeConfig: testAccMqttGroupsConfigDependenceNew(rand, map[string]string{
			"name_regex": `"^GID_fake.*$"`,
			"ids":        `["MQTT_fake:GID_fake"]`,
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccMqttGroupsConfigDependenceNew(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "testacc-%d"
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

resource "alibabacloudstack_mqtt_group" "default" {
  group_id = "GID_${var.name}"
  instance_id = "${alibabacloudstack_mqtt_instance.default.id}"
}

data "alibabacloudstack_mqtt_groups" "default" {
  instance_id = "${alibabacloudstack_mqtt_group.default.instance_id}"
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
