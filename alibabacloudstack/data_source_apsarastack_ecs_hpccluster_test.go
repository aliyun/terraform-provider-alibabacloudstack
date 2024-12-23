package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackEcsHpcClustersDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsHpcClustersSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_hpc_clusters.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsHpcClustersSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_hpc_clusters.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsHpcClustersSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_hpc_clusters.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsHpcClustersSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_hpc_clusters.default.id}_fake"]`,
		}),
	}

	AlibabacloudstackEcsHpcClustersCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}

var existAlibabacloudstackEcsHpcClustersMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"clusters.#":    "1",
		"clusters.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackEcsHpcClustersMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"clusters.#": "0",
	}
}

var AlibabacloudstackEcsHpcClustersCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_ecs_hpc_clusters.default",
	existMapFunc: existAlibabacloudstackEcsHpcClustersMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackEcsHpcClustersMapFunc,
}

func testAccCheckAlibabacloudstackEcsHpcClustersSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackEcsHpcClusters%d"
}






data "alibabacloudstack_ecs_hpc_clusters" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
