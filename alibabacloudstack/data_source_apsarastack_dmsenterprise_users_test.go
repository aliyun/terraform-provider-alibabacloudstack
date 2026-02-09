package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackDmsenterpriseUsersDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)
	resourceId := AlibabacloudstackDmsenterpriseUsersCheckInfo.resourceId
	name := fmt.Sprintf("tftestdomain%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlibabacloudstackDmsenterpriseUsersSourceConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_dmsenterprise_user.default.uid}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_dmsenterprise_user.default.uid}-fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_dmsenterprise_user.default.user_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_dmsenterprise_user.default.user_name}-fake",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alibabacloudstack_dmsenterprise_user.default.uid}"},
			"status": "NORMAL",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alibabacloudstack_dmsenterprise_user.default.uid}"},
			"status": "DISABLE",
		}),
	}

	roleConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"${alibabacloudstack_dmsenterprise_user.default.uid}"},
			"role": "USER",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"${alibabacloudstack_dmsenterprise_user.default.uid}"},
			"role": "DBA",
		}),
	}

	AlibabacloudstackDmsenterpriseUsersCheckInfo.dataSourceTestCheck(t, rand, idsConf, statusConf, nameRegexConf, roleConf)
}

var existAlibabacloudstackDmsenterpriseUsersMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"users.#":    "1",
		"users.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackDmsenterpriseUsersMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"users.#": "0",
	}
}

var AlibabacloudstackDmsenterpriseUsersCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_dmsenterprise_users.default",
	existMapFunc: existAlibabacloudstackDmsenterpriseUsersMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackDmsenterpriseUsersMapFunc,
}

func testAccCheckAlibabacloudstackDmsenterpriseUsersSourceConfig(name string) string {
	return AlibabacloudTestAccDmsenterpriseUserBasicdependence0(name) + `
resource "alibabacloudstack_dmsenterprise_user" "default" {
	uid=               "${alibabacloudstack_ascm_user.user.user_uid}"
	user_name=         "${alibabacloudstack_ascm_user.user.login_name}"
	max_execute_count= 10
	max_result_count=  10
	role_names=        ["USER"]
}
`
}
