package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackAscmMaxcomputeClusterDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_maxcompute_clusters.default"
	name := "maxcompute"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceMaxcomputeClustersConfigDependence)

	// Basic configuration without filters
	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{}),
		// No fake config since this data source queries system-wide clusters
		// and doesn't support filtering that results in empty lists with fake patterns
	}

	// Test with name_regex filter
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${data.alibabacloudstack_maxcompute_clusters.anyone.clusters.0.cluster}$",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake-nonexistent-cluster",
		}),
	}

	// Test with ids filter (using a placeholder that should exist)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${data.alibabacloudstack_maxcompute_clusters.anyone.clusters.0.cluster}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"fake-cluster-id-12345"},
		}),
	}

	var existMaxcomputeClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       CHECKSET, // Should contain at least one cluster
			"clusters.#":  CHECKSET, // Should contain at least one cluster
			"clusters.0.cluster":   CHECKSET,
			"clusters.0.core_arch": CHECKSET,
			"clusters.0.project":   CHECKSET,
			"clusters.0.region":    CHECKSET,
		}
	}

	// When no clusters match the filter, expect empty results
	var fakeMaxcomputeClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"clusters.#": "0",
		}
	}

	var maxcomputeClustersCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMaxcomputeClustersMapFunc,
		fakeMapFunc:  fakeMaxcomputeClustersMapFunc,
	}
	maxcomputeClustersCheckInfo.dataSourceTestCheck(t, 0, basicConf, nameRegexConf, idsConf)
}

func dataSourceMaxcomputeClustersConfigDependence(name string) string {
	return `
	data "alibabacloudstack_maxcompute_clusters" "anyone" {
		
	}
`
}
