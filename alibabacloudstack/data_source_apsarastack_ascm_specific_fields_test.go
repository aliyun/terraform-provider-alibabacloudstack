package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackAscmServiceFieldsDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_ascm_specific_fields.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceAscmServiceFieldsConfigDependence)

	// Test with valid group_field and resource_type
	validConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"group_filed":    "storageType",
			"resource_type":  "OSS",
		}),
		// No fake config needed since the API always returns results for valid inputs
		// We'll test invalid inputs separately if needed
	}

	// Test with label parameter
	labelConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"group_filed":    "storageType",
			"resource_type":  "OSS",
			"label":          "true",
		}),
	}

	var existAscmServiceFieldsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"specific_fields.#":  CHECKSET, // Should contain at least one field
			"group_filed":        "storageType",
			"resource_type":      "OSS",
		}
	}

	// Since this data source doesn't support filtering by non-existent values
	// (it returns fields based on valid group_field/resource_type combinations),
	// we don't have a traditional "fake" scenario that returns empty results.
	// Instead, we focus on validating successful responses.
	var fakeAscmServiceFieldsMapFunc = func(rand int) map[string]string {
		// This function is required by the framework but won't be used
		// since we don't have fake configs that return empty results
		return map[string]string{
			"specific_fields.#": "0",
		}
	}

	var ascmServiceFieldsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAscmServiceFieldsMapFunc,
		fakeMapFunc:  fakeAscmServiceFieldsMapFunc,
	}
	ascmServiceFieldsCheckInfo.dataSourceTestCheck(t, 0, validConf, labelConf)
}

func dataSourceAscmServiceFieldsConfigDependence(name string) string {
	return `
`
}
