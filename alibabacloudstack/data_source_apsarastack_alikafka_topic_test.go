package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackAlikafkaTopicsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackAlikafkaTopicsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_alikafka_topics.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackAlikafkaTopicsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_alikafka_topics.default.id}_fake"]`,
		}),
	}

	instance_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackAlikafkaTopicsSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_alikafka_topics.default.id}"]`,
			"instance_id": `"${alibabacloudstack_alikafka_topics.default.InstanceId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackAlikafkaTopicsSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_alikafka_topics.default.id}_fake"]`,
			"instance_id": `"${alibabacloudstack_alikafka_topics.default.InstanceId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackAlikafkaTopicsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_alikafka_topics.default.id}"]`,

			"instance_id": `"${alibabacloudstack_alikafka_topics.default.InstanceId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackAlikafkaTopicsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_alikafka_topics.default.id}_fake"]`,

			"instance_id": `"${alibabacloudstack_alikafka_topics.default.InstanceId}_fake"`}),
	}

	AlibabacloudstackAlikafkaTopicsCheckInfo.dataSourceTestCheck(t, rand, idsConf, instance_idConf, allConf)
}

var existAlibabacloudstackAlikafkaTopicsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"topics.#":    "1",
		"topics.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackAlikafkaTopicsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"topics.#": "0",
	}
}

var AlibabacloudstackAlikafkaTopicsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_alikafka_topics.default",
	existMapFunc: existAlibabacloudstackAlikafkaTopicsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackAlikafkaTopicsMapFunc,
}

func testAccCheckAlibabacloudstackAlikafkaTopicsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackAlikafkaTopics%d"
}






data "alibabacloudstack_alikafka_topics" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
