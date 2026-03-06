package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackMqttTopicsDataSource_basic(t *testing.T) {
	resourceId := "data.alibabacloudstack_mqtt_topics.default"
	rand := getAccTestRandInt(1000, 9999)
	name := fmt.Sprintf("tf_testAcck%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, resourceMqttTopicsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"store_instance_id": "${alibabacloudstack_mqtt_topic.default.store_instance_id}",
			"name_regex":        "${alibabacloudstack_mqtt_topic.default.topic}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"store_instance_id": "${alibabacloudstack_mqtt_topic.default.store_instance_id}",
			"name_regex":        "^fake-topic-[a-z0-9]{8}$",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"store_instance_id": "${alibabacloudstack_mqtt_topic.default.store_instance_id}",
			"ids": []string{
				"${alibabacloudstack_mqtt_topic.default.store_instance_id}:${alibabacloudstack_mqtt_topic.default.topic}",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"store_instance_id": "${alibabacloudstack_mqtt_topic.default.store_instance_id}",
			"ids": []string{
				"${alibabacloudstack_mqtt_topic.default.store_instance_id}:fake-topic-id",
			},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"store_instance_id": "${alibabacloudstack_mqtt_topic.default.store_instance_id}",
			"name_regex":        "${alibabacloudstack_mqtt_topic.default.topic}",
			"ids": []string{
				"${alibabacloudstack_mqtt_topic.default.store_instance_id}:${alibabacloudstack_mqtt_topic.default.topic}",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"store_instance_id": "${alibabacloudstack_mqtt_topic.default.store_instance_id}",
			"name_regex":        "${alibabacloudstack_mqtt_topic.default.topic}",
			"ids": []string{
				"${alibabacloudstack_ons_instance.default.id}:fake-topic-id",
			},
		}),
	}

	var existMqttTopicsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"topics.#":                    "1",
			"topics.0.store_instance_id":  CHECKSET,
			"topics.0.topic":              CHECKSET,
			"topics.0.remark":             "test",
			"topics.0.order_type":         "1",
			"topics.0.independent_naming": "true",
			"topics.0.update_time":        CHECKSET,
			"topics.0.relation":           "1",
			"topics.0.relation_name":      CHECKSET,
			"topics.0.create_time":        CHECKSET,
			"topics.0.namespace_id":       CHECKSET,
			"topics.0.unit_flag":          "false",
			"topics.0.status_name":        CHECKSET,
			"topics.0.channel_name":       "ALIYUN",
			"topics.0.channel_id":         "0",
			"topics.0.status":             "0",
			"ids.#":                       "1",
			"ids.0":                       CHECKSET,
		}
	}

	var fakeMqttTopicsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"topics.#": "0",
			"ids.#":    "0",
		}
	}

	var alibabacloudstackMqttTopicsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMqttTopicsMapFunc,
		fakeMapFunc:  fakeMqttTopicsMapFunc,
	}

	alibabacloudstackMqttTopicsCheckInfo.dataSourceTestCheck(t, -1, nameRegexConf, idsConf, allConf)
}

func resourceMqttTopicsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s

resource "alibabacloudstack_mqtt_topic" "default" {
  topic = "${var.name}"
  order_type = 1
  remark = "test"
  store_instance_id = "${alibabacloudstack_mqtt_instance.default.store_instance_id}"
}
`, name, MqttCommonTestCase)
}
