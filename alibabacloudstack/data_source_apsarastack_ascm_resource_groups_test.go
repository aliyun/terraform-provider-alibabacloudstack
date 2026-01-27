package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackascmResourceGroupDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_ascm_resource_groups.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testaccresourcegroup_%d", rand),
		dataSourceAscmResourceGroupDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${alibabacloudstack_ascm_resource_group.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${alibabacloudstack_ascm_resource_group.default.name}_fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ascm_resource_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ascm_resource_group.default.id}_fake"},
		}),
	}

	organizationIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"organization_id": "${alibabacloudstack_ascm_organization.default.id}",
			"ids":             []string{"${alibabacloudstack_ascm_resource_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"organization_id": "1",
			"ids":             []string{"${alibabacloudstack_ascm_resource_group.default.id}"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":      "${alibabacloudstack_ascm_resource_group.default.name}",
			"organization_id": "${alibabacloudstack_ascm_organization.default.id}",
			"ids":             []string{"${alibabacloudstack_ascm_resource_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":      "${alibabacloudstack_ascm_resource_group.default.name}_fake",
			"organization_id": "1",
			"ids":             []string{"${alibabacloudstack_ascm_resource_group.default.id}_fake"},
		}),
	}

	var existAscmResourceGroupMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"groups.#":                     "1",
			"groups.0.name":                fmt.Sprintf("tf_testaccresourcegroup_%d", rand),
			"groups.0.rs_id":               CHECKSET,
			"groups.0.resource_group_type": CHECKSET,
		}
	}

	var fakeAscmResourceGroupMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"groups.#": "0",
		}
	}

	var ascmResourceGroupCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAscmResourceGroupMapFunc,
		fakeMapFunc:  fakeAscmResourceGroupMapFunc,
	}

	ascmResourceGroupCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, organizationIdConf, allConf)
}

func dataSourceAscmResourceGroupDependence(name string) string {
	return fmt.Sprintf(`
variable name{
	default = "%s"
}

resource "alibabacloudstack_ascm_organization" "default" {
  name = "${var.name}"
  parent_id = "1"
} 

resource "alibabacloudstack_ascm_resource_group" "default" {
  organization_id = "${alibabacloudstack_ascm_organization.default.id}"
  name = "${var.name}"
}
`, name)
}
