package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackEcsDedicatedHostsClusterDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_ecs_dedicated_host_cluster.default"
	name := fmt.Sprintf("tf_testAccEcsDedicatedHostsClusterDataSource_%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEcsDedicatedHostsClusterConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{

			"ids": []string{"${alibabacloudstack_ecs_dedicated_host_cluster.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{

			"ids": []string{"${alibabacloudstack_ecs_dedicated_host_cluster.default.id}-fake"},
		}),
	}

	nameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{

			"dedicated_host_cluster_name": name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{

			"dedicated_host_cluster_name": name + "fake",
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{

			"dedicated_host_cluster_name_regex": name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"dedicated_host_cluster_name_regex": name + "fake",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":                         []string{"${alibabacloudstack_ecs_dedicated_host_cluster.default.id}"},
			"dedicated_host_cluster_name": name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":                         []string{"${alibabacloudstack_ecs_dedicated_host_cluster.default.id}"},
			"dedicated_host_cluster_name": name + "_fake",
		}),
	}

	var existKmsSecretVersionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"ids.0":                     CHECKSET,
			"dedicated_host_clusters.#": "1",
			"dedicated_host_clusters.0.dedicated_host_cluster_name": name,
			"dedicated_host_clusters.0.zone_id":                     CHECKSET,
		}
	}

	var fakeKmsSecretVersionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "0",
			"dedicated_host_cluster_name.#": "0",
		}
	}

	var ecsDedicatedHostsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existKmsSecretVersionsMapFunc,
		fakeMapFunc:  fakeKmsSecretVersionsMapFunc,
	}

	ecsDedicatedHostsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameConf, nameRegexConf, allConf)
}

func dataSourceEcsDedicatedHostsClusterConfigDependence(name string) string {
	return fmt.Sprintf(`
		%s
		resource "alibabacloudstack_ecs_dedicated_host_cluster" "default" {
		  dedicated_host_cluster_name = "%s"
          zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
		}
	`, DataZoneCommonTestCase, name)
}
