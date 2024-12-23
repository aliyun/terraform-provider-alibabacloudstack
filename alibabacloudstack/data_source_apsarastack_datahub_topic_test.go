package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackDatahubTopicsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackDatahubTopicsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_datahub_topics.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackDatahubTopicsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_datahub_topics.default.id}_fake"]`,
		}),
	}

	project_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackDatahubTopicsSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_datahub_topics.default.id}"]`,
			"project_name": `"${alibabacloudstack_datahub_topics.default.ProjectName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackDatahubTopicsSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_datahub_topics.default.id}_fake"]`,
			"project_name": `"${alibabacloudstack_datahub_topics.default.ProjectName}_fake"`,
		}),
	}

	topic_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackDatahubTopicsSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_datahub_topics.default.id}"]`,
			"topic_name": `"${alibabacloudstack_datahub_topics.default.TopicName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackDatahubTopicsSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_datahub_topics.default.id}_fake"]`,
			"topic_name": `"${alibabacloudstack_datahub_topics.default.TopicName}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackDatahubTopicsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_datahub_topics.default.id}"]`,

			"project_name": `"${alibabacloudstack_datahub_topics.default.ProjectName}"`,
			"topic_name":   `"${alibabacloudstack_datahub_topics.default.TopicName}"`}),
		fakeConfig: testAccCheckAlibabacloudstackDatahubTopicsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_datahub_topics.default.id}_fake"]`,

			"project_name": `"${alibabacloudstack_datahub_topics.default.ProjectName}_fake"`,
			"topic_name":   `"${alibabacloudstack_datahub_topics.default.TopicName}_fake"`}),
	}

	AlibabacloudstackDatahubTopicsCheckInfo.dataSourceTestCheck(t, rand, idsConf, project_nameConf, topic_nameConf, allConf)
}

var existAlibabacloudstackDatahubTopicsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"topics.#":    "1",
		"topics.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackDatahubTopicsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"topics.#": "0",
	}
}

var AlibabacloudstackDatahubTopicsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_datahub_topics.default",
	existMapFunc: existAlibabacloudstackDatahubTopicsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackDatahubTopicsMapFunc,
}

func testAccCheckAlibabacloudstackDatahubTopicsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackDatahubTopics%d"
}






data "alibabacloudstack_datahub_topics" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
