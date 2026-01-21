package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackAdbDbClustersDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_adb_db_clusters.default"
	name := fmt.Sprintf("tf-testAccADBConfig%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlibabacloudStackAdbDbClusterDataSourceConfig)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_adb_db_cluster.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"^test1234"},
		}),
	}
		nameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_adb_db_cluster.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "^test1234",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_adb_db_cluster.default.description}",
			"status":            "Running",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_adb_db_cluster.default.description}",
			"status":            "Creating",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_adb_db_cluster.default.description}",
			"status":            "Running",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "^test1234",
			"status":            "Creating",
		}),
	}

	var existAdbClusterMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                  "1",
			"descriptions.#":         "1",
			"clusters.#":             "1",
			"clusters.0.id":          CHECKSET,
			"clusters.0.description": CHECKSET,
			//"clusters.0.charge_type":        "PostPaid",
			"clusters.0.region_id": CHECKSET,
			//"clusters.0.expired":            "false",
			"clusters.0.create_time":        CHECKSET,
			"clusters.0.db_cluster_version": "3.0",
			//"clusters.0.db_node_class":      "C8",
			//"clusters.0.db_node_storage":    "300",
			//"clusters.0.compute_resource": "8Core50GB",
			//"clusters.0.elastic_io_resource": "0",
			//"clusters.0.zone_id":             CHECKSET,
			//"clusters.0.db_cluster_category": "Cluster",
			//"clusters.0.maintain_time":       "23:00Z-00:00Z",
		}
	}

	var fakeAdbClusterMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"clusters.#":     CHECKSET,
			"ids.#":          CHECKSET,
			"descriptions.#": CHECKSET,
		}
	}

	var AdbClusterCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAdbClusterMapFunc,
		fakeMapFunc:  fakeAdbClusterMapFunc,
	}
	
	AdbClusterCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameConf, statusConf, allConf)
}

func testAccCheckAlibabacloudStackAdbDbClusterDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s

data "alibabacloudstack_zones" "adb" {
  available_resource_creation = "ADB"
}

data "alibabacloudstack_adb_cluster_types" "default" {
  status = "Available"
  sorted_by = "CPU"
}

resource "alibabacloudstack_adb_db_cluster" "default" {
  db_cluster_category         = "${local.adb_instance_types.0.cluster_category}"
  db_cluster_version          = "3.0"
  db_node_class               = "${local.adb_instance_types.0.id}"
  description                 = name
  db_node_count               = "${local.adb_instance_types.0.node_min}"
  db_node_storage             = "${local.adb_instance_types.0.storage_min}"
  mode                        = "${local.adb_instance_types.0.mode}"
  vswitch_id                  = "${alibabacloudstack_vpc_vswitch.default.id}"
  cluster_type                = "${local.adb_instance_types.0.cluster_type}"
  cpu_type                    = "${local.adb_instance_types.0.cpu_type}"
}
`, name, VSwitchCommonTestCase)
}
