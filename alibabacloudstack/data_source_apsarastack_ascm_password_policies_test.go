package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackAscmPasswordPoliciesDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_ascm_password_policies.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceAscmPasswordPoliciesConfigDependence)

	// Test basic configuration (no filters)
	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{}),
		// No fake config since this data source returns system password policies
		// and doesn't support filtering that results in empty lists
	}

	var existAscmPasswordPoliciesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      CHECKSET, // Should be set with at least one ID
			"policies.#": CHECKSET, // Should contain at least one policy
			"policies.0.require_numbers": CHECKSET, // Should contain at least one policy
			// Individual policy attributes are computed but we don't know exact values
			// so we only verify the list structure exists
		}
	}

	// Since this data source queries system-wide password policies and doesn't support
	// filtering that results in empty lists (it always returns the current policy),
	// we provide a minimal fake function for framework compatibility
	var fakeAscmPasswordPoliciesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"policies.#": "0",
		}
	}

	var ascmPasswordPoliciesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAscmPasswordPoliciesMapFunc,
		fakeMapFunc:  fakeAscmPasswordPoliciesMapFunc,
	}
	ascmPasswordPoliciesCheckInfo.dataSourceTestCheck(t, 0, basicConf)
}

func dataSourceAscmPasswordPoliciesConfigDependence(name string) string {
	return `
`
}
