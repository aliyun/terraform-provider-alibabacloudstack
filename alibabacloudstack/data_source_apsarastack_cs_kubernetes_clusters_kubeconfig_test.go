package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackCSKubernetesClustersKubeConfigDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacckubernetesconfig-%d", rand)
	resourceId := "data.alibabacloudstack_cs_kubernetes_clusters_kubeconfig.k8s_clusters_kubeconfig"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceCSKubernetesClustersKubeConfigDependence)

	// Basic configuration - get kubeconfig for an existing cluster
	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cluster_id":      `${local.k8s_cluster_id}`,
			"private_address": "true",
		}),
	}

	var existCSKubeConfigMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"kubeconfig": CHECKSET,
		}
	}

	var fakeCSKubeConfigMapFunc = func(rand int) map[string]string {
		return map[string]string{}
	}

	var csKubeConfigCheckInfo = dataSourceAttr{
		resourceId:        resourceId,
		existMapFunc:      existCSKubeConfigMapFunc,
		fakeMapFunc:       fakeCSKubeConfigMapFunc,
		ExternalProviders: testAccExternalProviders,
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
			testAccPreCheckApiAvailable(t, "GET", "CS", "2015-12-15", "DescribeClusterUserKubeconfig", "")
		},
	}
	csKubeConfigCheckInfo.dataSourceTestCheck(t, 0, basicConf)
}

func dataSourceCSKubernetesClustersKubeConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s
`, name, AckK8sCommonTestCase())
}
