package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackAdbDbClustersDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackAdbDbClustersSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_adb_db_clusters.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackAdbDbClustersSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_adb_db_clusters.default.id}_fake"]`,
		}),
	}

	db_cluster_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackAdbDbClustersSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_adb_db_clusters.default.id}"]`,
			"db_cluster_id": `"${alibabacloudstack_adb_db_clusters.default.DBClusterId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackAdbDbClustersSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_adb_db_clusters.default.id}_fake"]`,
			"db_cluster_id": `"${alibabacloudstack_adb_db_clusters.default.DBClusterId}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackAdbDbClustersSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_adb_db_clusters.default.id}"]`,
			"status": `"${alibabacloudstack_adb_db_clusters.default.Status}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackAdbDbClustersSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_adb_db_clusters.default.id}_fake"]`,
			"status": `"${alibabacloudstack_adb_db_clusters.default.Status}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackAdbDbClustersSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_adb_db_clusters.default.id}"]`,

			"db_cluster_id": `"${alibabacloudstack_adb_db_clusters.default.DBClusterId}"`,
			"status":        `"${alibabacloudstack_adb_db_clusters.default.Status}"`}),
		fakeConfig: testAccCheckAlibabacloudstackAdbDbClustersSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_adb_db_clusters.default.id}_fake"]`,

			"db_cluster_id": `"${alibabacloudstack_adb_db_clusters.default.DBClusterId}_fake"`,
			"status":        `"${alibabacloudstack_adb_db_clusters.default.Status}_fake"`}),
	}

	AlibabacloudstackAdbDbClustersCheckInfo.dataSourceTestCheck(t, rand, idsConf, db_cluster_idConf, statusConf, allConf)
}

var existAlibabacloudstackAdbDbClustersMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"clusters.#":    "1",
		"clusters.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackAdbDbClustersMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"clusters.#": "0",
	}
}

var AlibabacloudstackAdbDbClustersCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_adb_db_clusters.default",
	existMapFunc: existAlibabacloudstackAdbDbClustersMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackAdbDbClustersMapFunc,
}

func testAccCheckAlibabacloudstackAdbDbClustersSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackAdbDbClusters%d"
}






data "alibabacloudstack_adb_db_clusters" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
