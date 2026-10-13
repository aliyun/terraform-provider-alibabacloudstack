package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackCSKubernetesAddonInstanceDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testaccaddoninstance-%d", rand)
	resourceId := "data.alibabacloudstack_cs_kubernetes_addon_instance.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceCSKubernetesAddonInstanceDependence)

	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cluster_id": `${local.k8s_cluster_id}`,
			"addon_name": "cluster-autoscaler",
		}),
	}

	var existCSAddonInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"cluster_id": CHECKSET,
			"name":       CHECKSET,
			"config":     CHECKSET,
			"state":      CHECKSET,
		}
	}

	var fakeCSAddonInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{}
	}

	var csAddonInstanceCheckInfo = dataSourceAttr{
		resourceId:        resourceId,
		existMapFunc:      existCSAddonInstanceMapFunc,
		fakeMapFunc:       fakeCSAddonInstanceMapFunc,
		ExternalProviders: testAccExternalProviders,
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
			testAccPreCheckApiAvailable(t, "GET", "CS", "2015-12-15", "DescribeClusterAddonInstance", "")
		},
	}
	csAddonInstanceCheckInfo.dataSourceTestCheck(t, 0, basicConf)
}

func dataSourceCSKubernetesAddonInstanceDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s
`, name, AckK8sCommonTestCase())
}
