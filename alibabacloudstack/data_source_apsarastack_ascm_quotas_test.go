package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackAscmQuotasDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_ascm_quotas.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceAscmQuotasConfigDependence)

	// Test with valid required parameters
	validConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"quota_type":    "organization",
			"quota_type_id": "1",
			"product_name":  "SLB",
		}),
		// No fake config since invalid parameters may cause API errors rather than empty results
		// We'll test different valid parameter combinations instead
	}

	// Test with target_type parameter
	targetTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"quota_type":    "organization",
			"quota_type_id": "1",
			"product_name":  "SLB",
			"target_type":   "organization",
		}),
	}

	var existAscmQuotasMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      CHECKSET, // Should be set with at least one ID
			"quotas.#":   CHECKSET, // Should contain at least one quota entry
			"quota_type": "organization",
			"product_name": "SLB",
			"quota_type_id": "1",
			// Individual quota attributes are computed but we don't know exact values
			// so we only verify the list structure exists
		}
	}

	// Since this data source requires valid parameters and queries system quotas,
	// invalid parameters typically cause API errors rather than empty results.
	// We provide a minimal fake function for framework compatibility.
	var fakeAscmQuotasMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"quotas.#": "0",
		}
	}

	var ascmQuotasCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAscmQuotasMapFunc,
		fakeMapFunc:  fakeAscmQuotasMapFunc,
	}
	ascmQuotasCheckInfo.dataSourceTestCheck(t, 0, validConf, targetTypeConf)
}

func dataSourceAscmQuotasConfigDependence(name string) string {
	return `
`
}
