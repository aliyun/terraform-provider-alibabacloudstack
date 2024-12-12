package alibabacloudstack

import (
	"fmt"
	"testing"

	
)

func TestAccAlibabacloudStackDmsEnterpriseUsersDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_dms_enterprise_users.default"
	name := fmt.Sprintf("tf_testAccDmsEnterpriseUsersDataSource_%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceDmsEnterpriseUsersConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_dms_enterprise_user.default.uid}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_dms_enterprise_user.default.uid}-fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_dms_enterprise_user.default.user_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_dms_enterprise_user.default.user_name}-fake",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alibabacloudstack_dms_enterprise_user.default.uid}"},
			"status": "NORMAL",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alibabacloudstack_dms_enterprise_user.default.uid}"},
			"status": "DISABLE",
		}),
	}

	roleConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"${alibabacloudstack_dms_enterprise_user.default.uid}"},
			"role": "DBA",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"${alibabacloudstack_dms_enterprise_user.default.uid}"},
			"role": "USER",
		}),
	}

	var existDmsEnterpriseUsersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                "1",
			"ids.0":                CHECKSET,
			"users.#":              "1",
			"users.0.mobile":       "15910799999",
			"users.0.nick_name":    name,
			"users.0.parent_uid":   CHECKSET,
			"users.0.role_ids.#":   "1",
			"users.0.role_names.#": "1",
			"users.0.status":       "NORMAL",
			"users.0.id":           CHECKSET,
			"users.0.user_id":      CHECKSET,
		}
	}

	var fakeDmsEnterpriseUsersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"users.#": "0",
		}
	}

	var kmsKeysCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDmsEnterpriseUsersMapFunc,
		fakeMapFunc:  fakeDmsEnterpriseUsersMapFunc,
	}

	kmsKeysCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, statusConf, roleConf)
}

func dataSourceDmsEnterpriseUsersConfigDependence(name string) string {
	return fmt.Sprintf(`        
		resource "alibabacloudstack_ascm_organization" "default" {
		 name = "Test_binder"
		 parent_id = "1"
		}
		
		resource "alibabacloudstack_ascm_user" "user" {
		 cellphone_number = "13900000000"
		 email = "test@gmail.com"
		 display_name = "C2C-DELTA"
		 organization_id = alibabacloudstack_ascm_organization.default.org_id
		 mobile_nation_code = "91"
		 login_name = "%s"
		 login_policy_id = 1
		}
		
		resource "alibabacloudstack_dms_enterprise_user" "default" {
		  uid = alibabacloudstack_ascm_user.user.user_id
		  user_name = alibabacloudstack_ascm_user.user.login_name
		  mobile = "15910799999"
		  role_names = ["DBA"]
	}`, name)
}
