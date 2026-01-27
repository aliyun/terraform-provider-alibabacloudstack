package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackAscmUserGroupsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_ascm_user_groups.default"
	name := fmt.Sprintf("tf-testacc-ascmusergroup-%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceAscmUserGroupsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_ascm_user_group.demo.group_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake-group-name-12345",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ascm_user_group.demo.user_group_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"fake-group-id"},
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_ascm_user_group.demo.group_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake-group-name-12345",
		}),
	}

	var existAscmUserGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#":                 "1",
			"groups.0.id":              CHECKSET,
			"groups.0.group_name":      name,
			"groups.0.organization_id": CHECKSET,
			"groups.0.user_group_id":   CHECKSET,
			"groups.0.role_ids.#":     CHECKSET,
			"groups.0.users.#":        CHECKSET,
		}
	}

	var fakeAscmUserGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#": "0",
		}
	}

	var ascmUserGroupsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAscmUserGroupsMapFunc,
		fakeMapFunc:  fakeAscmUserGroupsMapFunc,
	}

	ascmUserGroupsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}

func dataSourceAscmUserGroupsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_ascm_user_group" "demo" {
  group_name = var.name
  role_ids   = ["2"]
}
`, name)
}
