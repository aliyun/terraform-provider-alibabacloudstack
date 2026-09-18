package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackApigatewayv2K8sClustersDataSource(t *testing.T) {
	randInt := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_api_gateway_v2_k8s_clusters.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccApiGatewayv2cluster_%d", randInt),
		dataSourceApiGatewayV2K8sClustersConfigDependence)

	testAcc := dataSourceAttr{
		resourceId: resourceId,
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
		ExternalProviders: testAccExternalProviders,
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${alibabacloudstack_api_gateway_v2_k8s_cluster.default.k8s_cluster_name}$",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^fake.*",
		}),
	}

	k8sClusterNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"k8s_cluster_name": "${alibabacloudstack_api_gateway_v2_k8s_cluster.default.k8s_cluster_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"k8s_cluster_name": "fake-name",
		}),
	}

	csClusterIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cs_cluster_id": "${local.k8s_cluster_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cs_cluster_id": "fake-name",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_api_gateway_v2_k8s_cluster.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"fake-id"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":       "^${alibabacloudstack_api_gateway_v2_k8s_cluster.default.k8s_cluster_name}$",
			"k8s_cluster_name": "${alibabacloudstack_api_gateway_v2_k8s_cluster.default.k8s_cluster_name}",
			"ids":              []string{"${alibabacloudstack_api_gateway_v2_k8s_cluster.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":       "^${alibabacloudstack_api_gateway_v2_k8s_cluster.default.k8s_cluster_name}$",
			"k8s_cluster_name": "${alibabacloudstack_api_gateway_v2_k8s_cluster.default.k8s_cluster_name}",
			"ids":              []string{"fake-id"},
		}),
	}

	testAcc.dataSourceTestCheck(t, randInt, idsConf, nameRegexConf, k8sClusterNameConf, csClusterIdConf, allConf)
}

func dataSourceApiGatewayV2K8sClustersConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

resource "alibabacloudstack_api_gateway_v2_k8s_cluster" "default" {
	cs_cluster_id =   "${local.k8s_cluster_id}"
	k8s_cluster_name = "${var.name}"
}
`, name, AckK8sCommonTestCase())
}
