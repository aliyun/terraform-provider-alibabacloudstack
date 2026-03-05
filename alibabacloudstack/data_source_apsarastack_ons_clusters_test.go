package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackOnsClustersDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_ons_clusters.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceOnsClustersConfigDependence)

	// Basic configuration without filters
	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{}),
		// No fake config since this data source queries system-wide ONS clusters
		// and doesn't support filtering that results in empty lists
	}

	// Test with name_regex filter
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{`${data.alibabacloudstack_ons_clusters.anyone.clusters.0.id}`},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"fake-nonexistent-cluster"},
		}),
	}

	// Test with name_regex filter
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": `^${data.alibabacloudstack_ons_clusters.anyone.clusters.0.name}$`,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake-nonexistent-cluster",
		}),
	}

	var existOnsClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":           CHECKSET, // Should contain at least one cluster
			"clusters.#":      CHECKSET, // Should contain at least one cluster
			"clusters.0.id":   CHECKSET,
			"clusters.0.name": CHECKSET,
		}
	}

	// Since this data source queries system-wide ONS clusters, fake configurations
	// with non-existent filters should return empty results
	var fakeOnsClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"clusters.#": "0",
		}
	}

	var onsClustersCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOnsClustersMapFunc,
		fakeMapFunc:  fakeOnsClustersMapFunc,
	}
	onsClustersCheckInfo.dataSourceTestCheck(t, 0, basicConf, idsConf, nameRegexConf)
}

func dataSourceOnsClustersConfigDependence(name string) string {
	return `
data "alibabacloudstack_ons_clusters" "anyone" {
}
`
}
