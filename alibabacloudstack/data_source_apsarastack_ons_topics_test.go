package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackOnsTopicsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_ons_topics.default"
	name := fmt.Sprintf("tf-onstopicdata%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceOnsTopicsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_ons_instance.default.id}",
			"name_regex":  "${alibabacloudstack_ons_topic.default.topic}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_ons_instance.default.id}",
			"name_regex":  "fake_tf-testacc*",
		}),
	}

	var existOnsTopicsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"topics.#":                    "1",
			"topics.0.topic":              name,
			"topics.0.message_type":       "0",
			"topics.0.independent_naming": "true",
			"topics.0.remark":             "alibabacloudstack_ons_topic_remark",
		}
	}

	var fakeOnsTopicsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"topics.#": "0",
		}
	}

	var onsTopicsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOnsTopicsMapFunc,
		fakeMapFunc:  fakeOnsTopicsMapFunc,
	}

	onsTopicsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf)
}

func dataSourceOnsTopicsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
 default = "%v"
}

%s

resource "alibabacloudstack_ons_topic" "default" {
  instance_id = "${alibabacloudstack_ons_instance.default.id}"
  topic = "${var.name}"
  message_type = "0"
  remark = "alibabacloudstack_ons_topic_remark"
}
`, name, OnsCommonTestCase)
}
