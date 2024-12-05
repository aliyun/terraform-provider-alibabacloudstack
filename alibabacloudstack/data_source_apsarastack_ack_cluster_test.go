package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackAlibabacloudstackAckClustersDataSource(t *testing.T) {

	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackAckClustersSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ack_clusters.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackAckClustersSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ack_clusters.default.id}_fake"]`,
		}),
	}

	cluster_typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackAckClustersSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_ack_clusters.default.id}"]`,
			"cluster_type": `"${alibabacloudstack_ack_clusters.default.ClusterType}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackAckClustersSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_ack_clusters.default.id}_fake"]`,
			"cluster_type": `"${alibabacloudstack_ack_clusters.default.ClusterType}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackAckClustersSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ack_clusters.default.id}"]`,

			"cluster_type": `"${alibabacloudstack_ack_clusters.default.ClusterType}"`}),
		fakeConfig: testAccCheckAlibabacloudstackAckClustersSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ack_clusters.default.id}_fake"]`,

			"cluster_type": `"${alibabacloudstack_ack_clusters.default.ClusterType}_fake"`}),
	}

	AlibabacloudstackAckClustersCheckInfo.dataSourceTestCheck(t, rand, idsConf, cluster_typeConf, allConf)
}

var existAlibabacloudstackAckClustersMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"clusters.#":    "1",
		"clusters.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackAckClustersMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"clusters.#": "0",
	}
}

var AlibabacloudstackAckClustersCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_ack_clusters.default",
	existMapFunc: existAlibabacloudstackAckClustersMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackAckClustersMapFunc,
}

func testAccCheckAlibabacloudstackAckClustersSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackAckClusters%d"
}






data "alibabacloudstack_ack_clusters" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
