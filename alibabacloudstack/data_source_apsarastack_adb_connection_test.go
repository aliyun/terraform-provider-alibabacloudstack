package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackAlibabacloudstackAdbConnectionsDataSource(t *testing.T) {

	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackAdbConnectionsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_adb_connections.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackAdbConnectionsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_adb_connections.default.id}_fake"]`,
		}),
	}

	db_cluster_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackAdbConnectionsSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_adb_connections.default.id}"]`,
			"db_cluster_id": `"${alibabacloudstack_adb_connections.default.DBClusterId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackAdbConnectionsSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_adb_connections.default.id}_fake"]`,
			"db_cluster_id": `"${alibabacloudstack_adb_connections.default.DBClusterId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackAdbConnectionsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_adb_connections.default.id}"]`,

			"db_cluster_id": `"${alibabacloudstack_adb_connections.default.DBClusterId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackAdbConnectionsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_adb_connections.default.id}_fake"]`,

			"db_cluster_id": `"${alibabacloudstack_adb_connections.default.DBClusterId}_fake"`}),
	}

	AlibabacloudstackAdbConnectionsCheckInfo.dataSourceTestCheck(t, rand, idsConf, db_cluster_idConf, allConf)
}

var existAlibabacloudstackAdbConnectionsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"connections.#":    "1",
		"connections.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackAdbConnectionsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"connections.#": "0",
	}
}

var AlibabacloudstackAdbConnectionsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_adb_connections.default",
	existMapFunc: existAlibabacloudstackAdbConnectionsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackAdbConnectionsMapFunc,
}

func testAccCheckAlibabacloudstackAdbConnectionsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackAdbConnections%d"
}






data "alibabacloudstack_adb_connections" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
