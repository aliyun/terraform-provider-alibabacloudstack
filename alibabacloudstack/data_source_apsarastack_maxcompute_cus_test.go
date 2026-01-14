package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackMaxcomputeCu_DataSource(t *testing.T) {
	rand := getAccTestRandInt(1000, 9999)
	resourceId := "data.alibabacloudstack_maxcompute_cus.default"
	name := fmt.Sprintf("tf_testAcc%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceMaxcomputeCusConfigDependence)

	// Test with name_regex filter (should match the created CU)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"${alibabacloudstack_maxcompute_cu.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"fake_cu_id"},
		}),
	}
	
	// Test with name_regex filter (should match the created CU)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":  "^${alibabacloudstack_maxcompute_cu.default.cu_name}$",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake-nonexistent-cu",
		}),
	}

	// Test with cluster_name filter (using the cluster from data source)
	clusterConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cluster_name": "${data.alibabacloudstack_maxcompute_clusters.default.clusters.0.cluster}",
			"name_regex":  "^${alibabacloudstack_maxcompute_cu.default.cu_name}$",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cluster_name": "fake-cluster-name",
			"name_regex":   "${alibabacloudstack_maxcompute_cu.default.cu_name}",
		}),
	}

	var existMaxcomputeCusMapFunc = func(rand int) map[string]string {
		cuName := fmt.Sprintf("tf_testAcc%d", rand)
		return map[string]string{
			"ids.#":                    "1",
			"cus.#":                    "1",
			"cus.0.id":                 CHECKSET,
			"cus.0.cu_name":            cuName,
			"cus.0.cu_num":             "1",
			"cus.0.cluster_name":       CHECKSET,
		}
	}

	var fakeMaxcomputeCusMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
			"cus.#": "0",
		}
	}

	var maxcomputeCusCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMaxcomputeCusMapFunc,
		fakeMapFunc:  fakeMaxcomputeCusMapFunc,
	}
	maxcomputeCusCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, clusterConf)
}

func dataSourceMaxcomputeCusConfigDependence(name string) string {
	return fmt.Sprintf(`
data "alibabacloudstack_maxcompute_clusters" "default" {
  name_regex = "HYBRIDODPSCLUSTER-.*"
}

resource "alibabacloudstack_maxcompute_cu" "default" {
  cu_name      = "%s"
  cu_num       = 1
  cluster_name = data.alibabacloudstack_maxcompute_clusters.default.clusters.0.cluster
}

`, name)
}
