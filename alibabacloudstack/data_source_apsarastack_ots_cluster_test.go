package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackOtsClustersDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_ots_clusters.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceOtsClustersConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${data.alibabacloudstack_ots_clusters.anyone.clusters.0.cluster_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_*",
		}),
	}

	namesConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"names": []string{"${data.alibabacloudstack_ots_clusters.anyone.clusters.0.cluster_name}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"names": []string{"${data.alibabacloudstack_ots_clusters.anyone.clusters.0.cluster_name}_fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"names":               []string{"${data.alibabacloudstack_ots_clusters.anyone.clusters.0.cluster_name}"},
			"name_regex":          "${data.alibabacloudstack_ots_clusters.anyone.clusters.0.cluster_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"names":               []string{"${data.alibabacloudstack_ots_clusters.anyone.clusters.0.cluster_name}_fake"},
			"name_regex":          "${data.alibabacloudstack_ots_clusters.anyone.clusters.0.cluster_name}_fake",
		}),
	}

	var existOtsClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                    "1",
			"clusters.#":                 "1",
			"clusters.0.cluster_name":    CHECKSET,
			"clusters.0.cluster_type":    CHECKSET,
			"clusters.0.support_replica": CHECKSET,
		}
	}

	var fakeOtsClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":    "0",
			"clusters.#": "0",
		}
	}

	var otsClustersCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOtsClustersMapFunc,
		fakeMapFunc:  fakeOtsClustersMapFunc,
	}

	otsClustersCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, namesConf, allConf)
}

func dataSourceOtsClustersConfigDependence(name string) string {
	return `
data "alibabacloudstack_ots_clusters" "anyone" {
}
`
}
