package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackDatahubTopicsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	namesConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackDatahubTopicsSourceConfig(rand, map[string]string{
			"names": `["${alibabacloudstack_datahub_topic.default.name}"]`,
			"project_name": `"${alibabacloudstack_datahub_topic.default.project_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackDatahubTopicsSourceConfig(rand, map[string]string{
			"names": `["${alibabacloudstack_datahub_topic.default.name}_fake"]`,
			"project_name": `"${alibabacloudstack_datahub_topic.default.project_name}"`,
		}),
	}

	projectNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackDatahubTopicsSourceConfig(rand, map[string]string{
			"project_name": `"${alibabacloudstack_datahub_topic.default.project_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackDatahubTopicsSourceConfig(rand, map[string]string{
			"project_name": `"${alibabacloudstack_datahub_topic.default.project_name}_fake"`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackDatahubTopicsSourceConfig(rand, map[string]string{
			"name_regex": `"^${alibabacloudstack_datahub_topic.default.name}$"`,
			"project_name": `"${alibabacloudstack_datahub_topic.default.project_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackDatahubTopicsSourceConfig(rand, map[string]string{
			"name_regex": `"fake_topic_name"`,
			"project_name": `"${alibabacloudstack_datahub_topic.default.project_name}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackDatahubTopicsSourceConfig(rand, map[string]string{
			"names":        `["${alibabacloudstack_datahub_topic.default.name}"]`,
			"project_name": `"${alibabacloudstack_datahub_topic.default.project_name}"`,
			"name_regex":   `"^${alibabacloudstack_datahub_topic.default.name}$"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackDatahubTopicsSourceConfig(rand, map[string]string{
			"names":        `["${alibabacloudstack_datahub_topic.default.name}_fake"]`,
			"project_name": `"${alibabacloudstack_datahub_topic.default.project_name}_fake"`,
			"name_regex":   `"fake_topic_name"`,
		}),
	}

	AlibabacloudstackDatahubTopicsCheckInfo.dataSourceTestCheck(t, rand, projectNameConf, namesConf, nameRegexConf, allConf)
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
	default = "tf_testdatahubtopics%d"
}
resource "alibabacloudstack_datahub_project" "default" {
  comment = "test"
  name    = var.name
}

resource "alibabacloudstack_datahub_topic" "default" {
  name = var.name
  record_schemas {
    type       = "STRING"
    allow_null = false
    comment    = "test comment 1"
    name       = "test1"
  }

  comment      = "test"
  record_type  = "TUPLE"
  project_name = alibabacloudstack_datahub_project.default.name
}

data "alibabacloudstack_datahub_topics" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
