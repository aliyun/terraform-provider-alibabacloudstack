package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackDatahubTopicsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_datahub_topics.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testdatahubtopics%d", rand),
		dataSourceDatahubTopicsConfigDependence)

	namesConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"project_name": "${alibabacloudstack_datahub_topic.default.project_name}",
			"names":        []string{"${alibabacloudstack_datahub_topic.default.name}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"project_name": "${alibabacloudstack_datahub_topic.default.project_name}",
			"names":        []string{"${alibabacloudstack_datahub_topic.default.name}_fake"},
		}),
	}

	projectNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"project_name": "${alibabacloudstack_datahub_topic.default.project_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"project_name": "${alibabacloudstack_datahub_topic.default.project_name}_fake",
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"project_name": "${alibabacloudstack_datahub_topic.default.project_name}",
			"name_regex":   "^${alibabacloudstack_datahub_topic.default.name}$",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"project_name": "${alibabacloudstack_datahub_topic.default.project_name}",
			"name_regex":   "fake_topic_name",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"project_name": "${alibabacloudstack_datahub_topic.default.project_name}",
			"names":        []string{"${alibabacloudstack_datahub_topic.default.name}"},
			"name_regex":   "^${alibabacloudstack_datahub_topic.default.name}$",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"project_name": "${alibabacloudstack_datahub_topic.default.project_name}_fake",
			"names":        []string{"${alibabacloudstack_datahub_topic.default.name}_fake"},
			"name_regex":   "fake_topic_name",
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

func dataSourceDatahubTopicsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
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
`, name)
}
