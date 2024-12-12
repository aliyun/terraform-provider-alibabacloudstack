package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackDmsenterpriseUsersDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackDmsenterpriseUsersSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_dmsenterprise_users.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackDmsenterpriseUsersSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_dmsenterprise_users.default.id}_fake"]`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackDmsenterpriseUsersSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_dmsenterprise_users.default.id}"]`,
			"status": `"${alibabacloudstack_dmsenterprise_users.default.Status}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackDmsenterpriseUsersSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_dmsenterprise_users.default.id}_fake"]`,
			"status": `"${alibabacloudstack_dmsenterprise_users.default.Status}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackDmsenterpriseUsersSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_dmsenterprise_users.default.id}"]`,

			"status": `"${alibabacloudstack_dmsenterprise_users.default.Status}"`}),
		fakeConfig: testAccCheckAlibabacloudstackDmsenterpriseUsersSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_dmsenterprise_users.default.id}_fake"]`,

			"status": `"${alibabacloudstack_dmsenterprise_users.default.Status}_fake"`}),
	}

	AlibabacloudstackDmsenterpriseUsersCheckInfo.dataSourceTestCheck(t, rand, idsConf, statusConf, allConf)
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

func testAccCheckAlibabacloudstackDmsenterpriseUsersSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackDmsenterpriseUsers%d"
}






data "alibabacloudstack_dmsenterprise_users" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
