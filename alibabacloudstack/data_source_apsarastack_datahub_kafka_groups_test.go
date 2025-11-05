package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackDatahubKafkaGroupsDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_datahub_kafka_groups.default"
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf_testacc_group%d", rand)

	testAcc := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                           "1",
				"kafka_groups.#":                  "1",
				"kafka_groups.0.group_name":       name,
				"kafka_groups.0.comment":          "test group",
				"kafka_groups.0.creator":          CHECKSET,
				"kafka_groups.0.create_time":      CHECKSET,
				"kafka_groups.0.last_modify_time": CHECKSET,
				"kafka_groups.0.topic_list.#":     "2",
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":          "0",
				"kafka_groups.#": "0",
			}
		},
	}

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceDatahubGroupConfigDependence)

	projectConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"project_name": "${alibabacloudstack_datahub_kafka_group.default.project_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"project_name": "${alibabacloudstack_datahub_project.default.name}_fake",
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"project_name": "${alibabacloudstack_datahub_kafka_group.default.project_name}",
			"name_regex":   "^${alibabacloudstack_datahub_kafka_group.default.group_name}$",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"project_name": "${alibabacloudstack_datahub_kafka_group.default.project_name}",
			"name_regex":   "^fake.*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"project_name": "${alibabacloudstack_datahub_kafka_group.default.project_name}",
			"ids":          []string{"${alibabacloudstack_datahub_kafka_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"project_name": "${alibabacloudstack_datahub_kafka_group.default.project_name}",
			"ids":          []string{"fake_id"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"project_name": "${alibabacloudstack_datahub_kafka_group.default.project_name}",
			"name_regex":   "^${alibabacloudstack_datahub_kafka_group.default.group_name}$",
			"ids":          []string{"${alibabacloudstack_datahub_kafka_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"project_name": "${alibabacloudstack_datahub_kafka_group.default.project_name}",
			"name_regex":   "^fake.*",
			"ids":          []string{"fake_id"},
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, projectConf, nameRegexConf, idsConf, allConf)
}

func dataSourceDatahubGroupConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}
	resource "alibabacloudstack_datahub_project" "default" {
	  comment = "test"
	  name    = var.name
	}

	resource "alibabacloudstack_datahub_topic" "default" {
	  count        =2
	  name         = "${var.name}_${count.index}"
	  comment      = "test"
	  record_type  = "BLOB"
	  project_name = alibabacloudstack_datahub_project.default.name
	}

	resource "alibabacloudstack_datahub_kafka_group" "default" {
	  project_name = alibabacloudstack_datahub_project.default.name
	  comment      = "test group"
	  group_name   = var.name
	  topic_list = [
	    "${alibabacloudstack_datahub_topic.default.0.name}",
		"${alibabacloudstack_datahub_topic.default.1.name}"
	  ]
	}
`, name)
}
