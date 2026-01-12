package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackAscmInstanceFamiliesDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_ascm_instance_families.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceAscmInstanceFamiliesConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${data.alibabacloudstack_ascm_instance_families.anyone.ids.0}"},
			"resource_type": "DRDS",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"fake-id-12345"},
			"resource_type": "DRDS",
		}),
	}

	// Test with name_regex parameter
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"resource_type": "DRDS",
			"name_regex":    ".* Edition$", // Match all families
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"resource_type": "DRDS",
			"name_regex":    ".* Fake$", // Match all families
		}),
	}

	var existAscmInstanceFamiliesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":        CHECKSET, // Should be set with at least one ID
			"families.#":   CHECKSET, // Should contain at least one family
			"resource_type": "DRDS",
			// Individual family attributes are computed but we don't know exact values
			// so we only verify the list structure exists
		}
	}

	// Since this data source queries system-wide instance families and doesn't support
	// filtering that results in empty lists (invalid parameters typically cause errors),
	// we provide a minimal fake function for framework compatibility
	var fakeAscmInstanceFamiliesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"families.#": "0",
		}
	}

	var ascmInstanceFamiliesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAscmInstanceFamiliesMapFunc,
		fakeMapFunc:  fakeAscmInstanceFamiliesMapFunc,
	}
	ascmInstanceFamiliesCheckInfo.dataSourceTestCheck(t, 0, idsConf, nameRegexConf)
}

func dataSourceAscmInstanceFamiliesConfigDependence(name string) string {
	return `
	data "alibabacloudstack_ascm_instance_families" "anyone" {
	resource_type="DRDS"
	}
`
}
