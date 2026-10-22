package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"
)

func TestAccAlibabacloudStackElasticsearchInstancesDataSource(t *testing.T) {
	testAccPreCheck(t, "data.alibabacloudstack_elasticsearch_instances.default")
	testAccPreCheckApiAvailable(t, "GET", "elasticsearch-k8s", "2017-06-13", "ListInstance", "/openapi/instances")
	testAccPreCheckElasticsearchInstanceExists(t)
	
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_elasticsearch_instances.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccES%d", rand),
		dataSourceElasticsearchConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${local.es_instance_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${local.es_instance_id}-F"},
		}),
	}

	versionConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"version": "${local.es_version}",
			"ids":     []string{"${local.es_instance_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"version": "${local.es_version}-F",
			"ids":     []string{"${local.es_instance_id}"},
		}),
	}

	vpcConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"vpc_id": "${local.es_vpc_id}",
			"ids":    []string{"${local.es_instance_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"vpc_id": "${local.es_vpc_id}-F",
			"ids":    []string{"${local.es_instance_id}"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "^${local.es_description}$",
			"ids":               []string{"${local.es_instance_id}"},
			"version":           "${local.es_version}",
			"vpc_id":            "${local.es_vpc_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "^${local.es_description}-F$",
			"ids":               []string{"${local.es_instance_id}-F"},
			"version":           "${local.es_version}-F",
			"vpc_id":            "${local.es_vpc_id}-F",
		}),
	}

	var elasticsearchCheckInfo = dataSourceAttr{
		resourceId:        resourceId,
		existMapFunc:      existElasticsearchMapFunc,
		fakeMapFunc:       fakeElasticsearchMapFunc,
		ExternalProviders: testAccExternalProviders,
	}
	elasticsearchCheckInfo.dataSourceTestCheck(t, rand, idsConf, versionConf, vpcConf, allConf)
}

var existElasticsearchMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                        "1",
		"ids.0":                        CHECKSET,
		"descriptions.#":               "1",
		"descriptions.0":               CHECKSET,
		"instances.#":                  "1",
		"instances.0.id":               CHECKSET,
		"instances.0.description":      CHECKSET,
		"instances.0.data_node_amount": CHECKSET,
		"instances.0.data_node_spec":   CHECKSET,
		"instances.0.status":           CHECKSET,
		"instances.0.version":          CHECKSET,
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

variable "existed_es_instance" {
  default = "%s"
}

locals {
  create_es_instance_count = var.existed_es_instance == "" ? 1 : 0
}

%s

%s

resource "alibabacloudstack_elasticsearch_instance" "default" {
  count = local.create_es_instance_count
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

data "alibabacloudstack_elasticsearch_instances" "query" {
  ids = [var.existed_es_instance != "" ? var.existed_es_instance : alibabacloudstack_elasticsearch_instance.default.0.id]
}

locals {
  es_instance_id = var.existed_es_instance != "" ? var.existed_es_instance : alibabacloudstack_elasticsearch_instance.default.0.id
  es_description = var.existed_es_instance != "" ? (length(data.alibabacloudstack_elasticsearch_instances.query.instances) > 0 ? data.alibabacloudstack_elasticsearch_instances.query.instances.0.description : var.name) : var.name
  es_version = var.existed_es_instance != "" ? (length(data.alibabacloudstack_elasticsearch_instances.query.instances) > 0 ? data.alibabacloudstack_elasticsearch_instances.query.instances.0.version : "7.10.0_ali1.6.0") : "7.10.0_ali1.6.0"
  es_vpc_id = var.existed_es_instance != "" ? (length(data.alibabacloudstack_elasticsearch_instances.query.instances) > 0 ? data.alibabacloudstack_elasticsearch_instances.query.instances.0.vpc_id : "") : alibabacloudstack_vpc_vpc.default.id
}
`, name, os.Getenv("ALIBABACLOUDSTACK_TEST_EXISTED_ES_ID"), RandomPasswordTestCase(12, 1), VSwitchCommonTestCase)
}
