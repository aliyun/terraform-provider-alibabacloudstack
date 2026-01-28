package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackAscmUsersDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_ascm_users.default"
	name := fmt.Sprintf("tf-testacc-ascmuser%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceAscmUsersConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_ascm_user.default.login_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_ascm_user.default.login_name}_fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ascm_user.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ascm_user.default.id}_fake"},
		}),
	}
	loginNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"login_name": "${alibabacloudstack_ascm_user.default.login_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"login_name": "${alibabacloudstack_ascm_user.default.login_name}_fake",
		}),
	}
	loginPolicyIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"login_policy_id": "${alibabacloudstack_ascm_user.default.login_policy_id}",
			"ids":             []string{"${alibabacloudstack_ascm_user.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"login_policy_id": "2",
			"ids":             []string{"${alibabacloudstack_ascm_user.default.id}"},
		}),
	}
	organizationIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"organization_id": "${alibabacloudstack_ascm_user.default.organization_id}",
			"ids":             []string{"${alibabacloudstack_ascm_user.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"organization_id": "5",
			"ids":             []string{"${alibabacloudstack_ascm_user.default.id}"},
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"status": "ACTIVE",
			"ids":    []string{"${alibabacloudstack_ascm_user.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"status": "INACTIVE",
			"ids":    []string{"${alibabacloudstack_ascm_user.default.id}"},
		}),
	}

	var existAscmUsersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"users.#":                 "1",
			"users.0.id":              CHECKSET,
			"users.0.login_name":      name,
			"users.0.organization_id": CHECKSET,
			"users.0.cellphone_num":   CHECKSET,
			"users.0.display_name":    CHECKSET,
			"users.0.email":           CHECKSET,
		}
	}

	var fakeAscmUsersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"users.#": "0",
		}
	}

	var ascmUsersCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAscmUsersMapFunc,
		fakeMapFunc:  fakeAscmUsersMapFunc,
	}

	ascmUsersCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, loginNameConf, loginPolicyIdConf, organizationIdConf, statusConf)
}

func dataSourceAscmUsersConfigDependence(name string) string {
	return fmt.Sprintf(`
variable name{
 default = "%s"
}

resource "alibabacloudstack_ascm_organization" "default" {
  name = "${var.name}"
  parent_id = "1"
} 

resource "alibabacloudstack_ascm_user_group" "default" {
 group_name =      "${var.name}"
 organization_id = "${alibabacloudstack_ascm_organization.default.id}"
}

resource "alibabacloudstack_ascm_user" "default" {
 cellphone_number = "13900000000"
 email = "test@gmail.com"
 display_name = "${var.name}"
 organization_id = "${alibabacloudstack_ascm_organization.default.id}"
 mobile_nation_code = "86"
 login_name = "${var.name}"
 login_policy_id = 1
}
`, name)
}
