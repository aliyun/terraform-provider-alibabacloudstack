package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackCSAutoscalingConfigsDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_cs_autoscaling_configs.default"
	name := fmt.Sprintf("tf-testAccAutoscalingConfigs-%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceCSAutoscalingConfigsDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${local.k8s_cluster_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${local.k8s_cluster_id}_fake"},
		}),
	}

	testAcc := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#": CHECKSET,
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":     "0",
				"configs.#": "0",
			}
		},
		ExternalProviders: testAccExternalProviders,
	}

	testAcc.dataSourceTestCheck(t, rand, idsConf)
}

func dataSourceCSAutoscalingConfigsDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s

`, name, AckK8sCommonTestCase())
}
