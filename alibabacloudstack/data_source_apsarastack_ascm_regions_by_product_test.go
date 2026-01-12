package alibabacloudstack

import (
	"os"
	"testing"
)

func TestAccAlibabacloudStackAscmRegionsByProductDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_ascm_regions_by_product.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceAscmRegionsByProductConfigDependence)

	// Test with valid product_name
	productConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"product_name": "ecs",
		}),
	}
	
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{os.Getenv("ALIBABACLOUDSTACK_REGION")},
			"product_name":  "ecs",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"fake_region"},
			"product_name":  "ecs",
		}),
	}

	var existAscmRegionsByProductMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":            CHECKSET, // Should be set with at least one ID
			"region_list.#":    CHECKSET, // Should contain at least one region
			"product_name":     "ecs",
			// region_list attributes are computed but we don't know exact values
			// so we only verify the list structure exists
		}
	}

	// Since this data source requires a valid product_name and doesn't support
	// filtering that results in empty lists (invalid products typically cause errors),
	// we provide a minimal fake function for framework compatibility
	var fakeAscmRegionsByProductMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":         "0",
			"region_list.#": "0",
		}
	}

	var ascmRegionsByProductCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAscmRegionsByProductMapFunc,
		fakeMapFunc:  fakeAscmRegionsByProductMapFunc,
	}
	ascmRegionsByProductCheckInfo.dataSourceTestCheck(t, 0, idsConf, productConf)
}

func dataSourceAscmRegionsByProductConfigDependence(name string) string {
	return `
`
}
