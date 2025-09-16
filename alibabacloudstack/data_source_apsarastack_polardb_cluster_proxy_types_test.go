package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackPolardbClusterProxyTypesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_polardb_cluster_proxy_types.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccPolardbClusterProxyTypesDataSource_%d", rand),
		dataSourcePolardbClusterProxyTypesPresetDependence)

	baseConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_type":    "MySQL",
			"db_version": "5.7",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${data.alibabacloudstack_polardb_cluster_proxy_types.preset.proxy_classes.0.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"xxxxx"},
		}),
	}
	cpuConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"core_count": "4",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"core_count": "800",
		}),
	}

	var existPolardbClusterProxyTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      CHECKSET,
			"ids.0":                      CHECKSET,
			"proxy_classes.#":            CHECKSET,
			"proxy_classes.0.id":         CHECKSET,
			"proxy_classes.0.core_count": CHECKSET,
		}
	}

	var fakePolardbClusterProxyTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":           "0",
			"proxy_classes.#": "0",
		}
	}

	var PolardbClusterProxyTypesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existPolardbClusterProxyTypesMapFunc,
		fakeMapFunc:  fakePolardbClusterProxyTypesMapFunc,
	}

	PolardbClusterProxyTypesCheckInfo.dataSourceTestCheck(t, rand, baseConf, idsConf, cpuConf)
}

func dataSourcePolardbClusterProxyTypesPresetDependence(name string) string {
	return `
	data "alibabacloudstack_polardb_cluster_proxy_types" "preset" {
	}
`
}
