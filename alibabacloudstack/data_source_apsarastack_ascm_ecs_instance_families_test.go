package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackAscmEcsInstanceFamiliesDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_ascm_ecs_instance_families.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceAscmEcsInstanceFamiliesConfigDependence)

	// Test with ids parameter (though typically not used with status, but supported by schema)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${data.alibabacloudstack_ascm_ecs_instance_families.anyone.ids.0}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"fake-family-id-12345"},
		}),
	}

	var existAscmEcsInstanceFamiliesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      CHECKSET, // Should be set with at least one ID
			"families.#": CHECKSET, // Should contain at least one family
			"families.0.instance_type_family_id": CHECKSET,
			// Individual family attributes are computed but we don't know exact values
			// so we only verify the list structure exists
		}
	}

	var fakeAscmEcsInstanceFamiliesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"families.#": "0",
		}
	}

	var ascmEcsInstanceFamiliesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAscmEcsInstanceFamiliesMapFunc,
		fakeMapFunc:  fakeAscmEcsInstanceFamiliesMapFunc,
	}
	ascmEcsInstanceFamiliesCheckInfo.dataSourceTestCheck(t, 0, idsConf)
}

func dataSourceAscmEcsInstanceFamiliesConfigDependence(name string) string {
	return `
	data "alibabacloudstack_ascm_ecs_instance_families" "anyone" {
	}
`
}
