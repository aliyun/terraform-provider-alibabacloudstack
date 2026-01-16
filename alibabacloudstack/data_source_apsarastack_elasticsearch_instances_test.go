package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackElasticsearchInstancesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_elasticsearch_instances.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccES%d", rand),
		dataSourceElasticsearchConfigDependence)

	descriptionRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_elasticsearch_instance.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_elasticsearch_instance.default.description}-F",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_elasticsearch_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_elasticsearch_instance.default.id}-F"},
		}),
	}

	versionConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"version": "7.10.0_ali1.6.0",
			"ids":     []string{"${alibabacloudstack_elasticsearch_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"version": "7.10.0_ali1.6.0-F",
			"ids":     []string{"${alibabacloudstack_elasticsearch_instance.default.id}"},
		}),
	}

	vpcConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"vpc_id": "${alibabacloudstack_vpc_vpc.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"vpc_id": "${alibabacloudstack_vpc_vpc.default.id}-F",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_elasticsearch_instance.default.description}",
			"ids":               []string{"${alibabacloudstack_elasticsearch_instance.default.id}"},
			"version":           "7.10.0_ali1.6.0",
			"vpc_id":            "${alibabacloudstack_vpc_vpc.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_elasticsearch_instance.default.description}-F",
			"ids":               []string{"${alibabacloudstack_elasticsearch_instance.default.id}-F"},
			"version":           "7.10.0_ali1.6.0-F",
			"vpc_id":            "${alibabacloudstack_vpc_vpc.default.id}-F",
		}),
	}

	var elasticsearchCheckInfo = dataSourceAttr{
		resourceId:        resourceId,
		existMapFunc:      existElasticsearchMapFunc,
		fakeMapFunc:       fakeElasticsearchMapFunc,
		ExternalProviders: testAccExternalProviders,
	}
	elasticsearchCheckInfo.dataSourceTestCheck(t, rand, descriptionRegexConf, idsConf, versionConf, vpcConf, allConf)
}

var existElasticsearchMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                        "1",
		"ids.0":                        CHECKSET,
		"descriptions.#":               "1",
		"descriptions.0":               fmt.Sprintf("tf-testAccES%d", rand),
		"instances.#":                  "1",
		"instances.0.id":               CHECKSET,
		"instances.0.description":      fmt.Sprintf("tf-testAccES%d", rand),
		"instances.0.data_node_amount": "3",
		"instances.0.data_node_spec":   "1C 2Gi",
		"instances.0.status":           "active",
		"instances.0.version":          "7.10.0_ali1.6.0",
		"instances.0.vswitch_id":       CHECKSET,
	}
}

var fakeElasticsearchMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instances.#":    "0",
		"ids.#":          "0",
		"descriptions.#": "0",
	}
}

var esTestPassword = getAccTestPassword(12)

func dataSourceElasticsearchConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

%s

resource "alibabacloudstack_elasticsearch_instance" "default" {
  data_node_disk_type = "yoda-lvm"
  data_node_spec = "1C 2Gi"
  password = "${random_password.password.0.result}"
  monitor_password = "${random_password.password.0.result}"
  data_node_amount = "3"
  cpu_type = "Intel"
  scene = "normal"
  description = "${var.name}"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
  version = "7.10.0_ali1.6.0"
  data_node_disk_size = "500"
  kibana_node_spec=      "1C 2Gi"
  kibana_password=       "${random_password.password.0.result}"
  master_node_amount=    3
  master_node_spec=      "1C 2Gi"
  master_node_disk_size= "100"
  master_node_disk_type= "yoda-lvm"
  client_node_amount=    2
  client_node_spec=      "1C 2Gi"
}
`, name, RandomPasswordTestCase(12, 1), VSwitchCommonTestCase)
}
