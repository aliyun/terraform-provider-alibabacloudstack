package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackAscmRamServiceRoles_DataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_ascm_ram_service_roles.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceAscmRamServiceRolesConfigDependence)

	productConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"product": "${data.alibabacloudstack_ascm_ram_service_roles.anyone.roles.0.product}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"product": "fake-product",
		}),
	}

	// Since this data source returns roles based on product filter,
	// and we cannot create/delete service roles in tests,
	// we assume that "ecs" product will return at least one role in real environment.
	// For fake product, it should return empty list.

	var existAscmRamServiceRolesMapFunc = func(rand int) map[string]string {
		// We expect at least one role for "ecs" product
		return map[string]string{
			"roles.#":                   CHECKSET, // At least one role expected
			"roles.0.id":                CHECKSET,
			"roles.0.name":              CHECKSET,
			"roles.0.description":       CHECKSET,
			"roles.0.role_type":         CHECKSET,
			"roles.0.product":           CHECKSET,
			"roles.0.organization_name": CHECKSET,
			"roles.0.aliyun_user_id":    CHECKSET,
		}
	}

	var fakeAscmRamServiceRolesMapFunc = func(rand int) map[string]string {
		// Fake product should return no roles
		return map[string]string{
			"roles.#": "0",
		}
	}

	var ascmRamServiceRolesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAscmRamServiceRolesMapFunc,
		fakeMapFunc:  fakeAscmRamServiceRolesMapFunc,
	}
	ascmRamServiceRolesCheckInfo.dataSourceTestCheck(t, 0, productConf)
}

func dataSourceAscmRamServiceRolesConfigDependence(name string) string {
	return `
	data "alibabacloudstack_ascm_ram_service_roles" "anyone" {
	}
`
}
