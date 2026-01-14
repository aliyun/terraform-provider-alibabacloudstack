package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackAscmMaxcomputeClusterQuotasDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_maxcompute_cluster_quotas.default"
	name := "maxcompute_quotas"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceMaxcomputeClusterQuotasConfigDependence)

	// Basic configuration that depends on maxcompute clusters data source
	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cluster" :  "${data.alibabacloudstack_maxcompute_clusters.default.clusters.0.cluster}",
		}),
		// No fake config since this depends on real cluster data
	}

	var existMaxcomputeClusterQuotasMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"cluster":        CHECKSET,
			"cu_total":       CHECKSET,
			"disk_available": CHECKSET,
			"cu_available":   CHECKSET,
			"disk_total":     CHECKSET,
		}
	}

	// Since this data source requires a valid cluster, fake scenarios aren't applicable
	// But we provide a minimal fake function for framework compatibility
	var fakeMaxcomputeClusterQuotasMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"cluster":        "",
			"cu_total":       "",
			"disk_available": "",
			"cu_available":   "",
			"disk_total":     "",
		}
	}

	var maxcomputeClusterQuotasCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMaxcomputeClusterQuotasMapFunc,
		fakeMapFunc:  fakeMaxcomputeClusterQuotasMapFunc,
	}
	maxcomputeClusterQuotasCheckInfo.dataSourceTestCheck(t, 0, basicConf)
}

func dataSourceMaxcomputeClusterQuotasConfigDependence(name string) string {
	return `
data "alibabacloudstack_maxcompute_clusters" "default" {
}
`
}
