package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackLogClustersDataSource_basic(t *testing.T) {

	resourceId := "data.alibabacloudstack_log_clusters.default"
	name := ""

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceLogClustersConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${local.cluster_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${local.cluster_name}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"fake-cluster-id"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${local.cluster_name}"},
			"name_regex": "${local.cluster_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"fake-cluster-id"},
			"name_regex": "fake_*",
		}),
	}

	var existLogClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":             "1",
			"clusters.#":        "1",
			"clusters.0.name":   CHECKSET,
			"clusters.0.status": CHECKSET,
		}
	}

	var fakeLogClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"clusters.#": "0",
		}
	}

	var logClustersCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existLogClustersMapFunc,
		fakeMapFunc:  fakeLogClustersMapFunc,
	}
	logClustersCheckInfo.dataSourceTestCheck(t, 0, nameRegexConf, idsConf, allConf)
}

func dataSourceLogClustersConfigDependence(name string) string {
	return `
data "alibabacloudstack_log_clusters" "example" {
}

locals {
	cluster_name = data.alibabacloudstack_log_clusters.example.clusters.0.name
}
`
}
