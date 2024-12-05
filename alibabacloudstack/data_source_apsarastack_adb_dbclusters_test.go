package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackAdbDbClustersDataSource(t *testing.T) {
	rand := acctest.RandInt()
	nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackAdbDbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alibabacloudstack_adb_db_cluster.default.description}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackAdbDbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"^test1234"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackAdbDbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alibabacloudstack_adb_db_cluster.default.description}"`,
			"status":            `"Running"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackAdbDbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alibabacloudstack_adb_db_cluster.default.description}"`,
			"status":            `"Creating"`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackAdbDbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alibabacloudstack_adb_db_cluster.default.description}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackAdbDbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alibabacloudstack_adb_db_cluster.default.description}"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackAdbDbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alibabacloudstack_adb_db_cluster.default.description}"`,
			"status":            `"Running"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackAdbDbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"^test1234"`,
			"status":            `"Creating"`,
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
			"clusters.0.db_node_count": "2",
			//"clusters.0.db_node_storage":    "300",
			//"clusters.0.compute_resource": "8Core50GB",
			//"clusters.0.elastic_io_resource": "0",
			//"clusters.0.zone_id":             CHECKSET,
			//"clusters.0.db_cluster_category": "Cluster",
			//"clusters.0.maintain_time":       "23:00Z-00:00Z",
			"clusters.0.security_ips.#": "2",
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
		resourceId:   "data.alibabacloudstack_adb_db_clusters.default",
		existMapFunc: existAdbClusterMapFunc,
		fakeMapFunc:  fakeAdbClusterMapFunc,
	}

	AdbClusterCheckInfo.dataSourceTestCheck(t, rand, nameConf, statusConf, tagsConf, allConf)
}

func testAccCheckAlibabacloudStackAdbDbClusterDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
	%s
variable "creation" {	
	default = "ADB"
}

variable "name" {
	default = "tf-testAccADBConfig_%d"
}

resource "alibabacloudstack_adb_db_cluster" "default" {
	db_cluster_category = "Basic"
	db_cluster_class = "C8"
	db_node_storage = "200"
	db_cluster_version = "3.0"
	db_node_count = "2"
	vswitch_id              = "${alibabacloudstack_vswitch.default.id}"
	description             = "${var.name}"
	mode					= "reserver"
	cluster_type =        "analyticdb"
	cpu_type =            "intel"
	security_ips      = ["10.168.1.12", "10.168.1.11"]
}

data "alibabacloudstack_adb_db_clusters" "default" {	
	enable_details = true
	%s
}
`, AdbCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}
