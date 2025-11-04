package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackApigatewayv2K8sClustersDataSource(t *testing.T) {
	randInt := getAccTestRandInt(10000, 20000)
	testAcc := dataSourceAttr{
		resourceId: "data.alibabacloudstack_api_gateway_v2_k8s_clusters.default",
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                       "1",
				"clusters.#":                  "1",
				"clusters.0.id":               CHECKSET,
				"clusters.0.k8s_cluster_name": CHECKSET,
				"clusters.0.cluster_type":     CHECKSET,
				"clusters.0.cs_cluster_id":    CHECKSET,
				"clusters.0.cs_cluster_name":  CHECKSET,
				"clusters.0.slb_type":         CHECKSET,
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":      "0",
				"clusters.#": "0",
			}
		},
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: APIGateWayV2InstanceClusterDependenceNew(randInt, map[string]string{
			"name_regex": `"^${alibabacloudstack_api_gateway_v2_k8s_cluster.default.k8s_cluster_name}$"`,
		}),
		fakeConfig: APIGateWayV2InstanceClusterDependenceNew(randInt, map[string]string{
			"name_regex": `"^fake.*"`,
		}),
	}

	k8sClusterNameConf := dataSourceTestAccConfig{
		existConfig: APIGateWayV2InstanceClusterDependenceNew(randInt, map[string]string{
			"k8s_cluster_name": `"${alibabacloudstack_api_gateway_v2_k8s_cluster.default.k8s_cluster_name}"`,
		}),
		fakeConfig: APIGateWayV2InstanceClusterDependenceNew(randInt, map[string]string{
			"k8s_cluster_name": `"fake-name"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: APIGateWayV2InstanceClusterDependenceNew(randInt, map[string]string{
			"ids": `["${alibabacloudstack_api_gateway_v2_k8s_cluster.default.id}"]`,
		}),
		fakeConfig: APIGateWayV2InstanceClusterDependenceNew(randInt, map[string]string{
			"ids": `["fake-id"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: APIGateWayV2InstanceClusterDependenceNew(randInt, map[string]string{
			"name_regex":       `"^${alibabacloudstack_api_gateway_v2_k8s_cluster.default.k8s_cluster_name}$"`,
			"k8s_cluster_name": `"${alibabacloudstack_api_gateway_v2_k8s_cluster.default.k8s_cluster_name}"`,
			"ids":              `["${alibabacloudstack_api_gateway_v2_k8s_cluster.default.id}"]`,
		}),
		fakeConfig: APIGateWayV2InstanceClusterDependenceNew(randInt, map[string]string{
			"name_regex":       `"^${alibabacloudstack_api_gateway_v2_k8s_cluster.default.k8s_cluster_name}$"`,
			"k8s_cluster_name": `"${alibabacloudstack_api_gateway_v2_k8s_cluster.default.k8s_cluster_name}"`,
			"ids":              `["fake-id"]`,
		}),
	}

	testAcc.dataSourceTestCheck(t, randInt, idsConf, nameRegexConf, k8sClusterNameConf, allConf)
}

func APIGateWayV2InstanceClusterDependenceNew(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
variable "name" {
  default = "tf_testAccApiGatewayv2cluster_%d"
}

%s

resource "alibabacloudstack_api_gateway_v2_k8s_cluster" "default" {
	cs_cluster_id =   "${local.k8s_cluster_id}"
	k8s_cluster_name = "${var.name}"
}

data "alibabacloudstack_api_gateway_v2_k8s_clusters" "default" {
    %s
}
`, rand, AckK8sCommonTestCase(), strings.Join(pairs, "\n   "))
}
