package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackOssClustersDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_oss_clusters.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceOssClustersConfigDependence)

	// Basic configuration without filters
	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{}),
		// No fake config since this data source queries system-wide OSS clusters
		// and doesn't support filtering that results in empty lists
	}

	// Test with name_regex filter
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{`${data.alibabacloudstack_oss_clusters.anyone.clusters.0.cluster}`},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"fake-nonexistent-cluster"},
		}),
	}

	// Test with name_regex filter
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": `^${data.alibabacloudstack_oss_clusters.anyone.clusters.0.cluster_name}$`,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake-nonexistent-cluster",
		}),
	}

	var existOssClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   CHECKSET, // Should contain at least one cluster
			"clusters.#":              CHECKSET, // Should contain at least one cluster
			"clusters.0.id":           CHECKSET,
			"clusters.0.cluster":      CHECKSET,
			"clusters.0.cluster_name": CHECKSET,
			"clusters.0.location":     CHECKSET,
			"clusters.0.oss_endpoint": CHECKSET,
			"clusters.0.real_zone":    CHECKSET,
			// Boolean fields should be set to either "true" or "false"
			"clusters.0.ha_apsara_stack":                     CHECKSET,
			"clusters.0.oss_ha_enable_single_cluster_access": CHECKSET,
			"clusters.0.is_master_zone":                      CHECKSET,
		}
	}

	// Since this data source queries system-wide OSS clusters, fake configurations
	// with non-existent filters should return empty results
	var fakeOssClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"clusters.#": "0",
		}
	}

	var ossClustersCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOssClustersMapFunc,
		fakeMapFunc:  fakeOssClustersMapFunc,
		PreCheck:     func() { testAccPreCheckOssEndpointList(t) },
	}
	ossClustersCheckInfo.dataSourceTestCheck(t, 0, basicConf, idsConf, nameRegexConf)
}

func dataSourceOssClustersConfigDependence(name string) string {
	return `
data "alibabacloudstack_oss_clusters" "anyone" {
}
`
}
