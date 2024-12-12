package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackNasAccessGroupsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackNasAccessGroupsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_nas_access_groups.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackNasAccessGroupsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_nas_access_groups.default.id}_fake"]`,
		}),
	}

	access_group_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackNasAccessGroupsSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_nas_access_groups.default.id}"]`,
			"access_group_name": `"${alibabacloudstack_nas_access_groups.default.AccessGroupName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackNasAccessGroupsSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_nas_access_groups.default.id}_fake"]`,
			"access_group_name": `"${alibabacloudstack_nas_access_groups.default.AccessGroupName}_fake"`,
		}),
	}

	file_system_typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackNasAccessGroupsSourceConfig(rand, map[string]string{
			"ids":              `["${alibabacloudstack_nas_access_groups.default.id}"]`,
			"file_system_type": `"${alibabacloudstack_nas_access_groups.default.FileSystemType}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackNasAccessGroupsSourceConfig(rand, map[string]string{
			"ids":              `["${alibabacloudstack_nas_access_groups.default.id}_fake"]`,
			"file_system_type": `"${alibabacloudstack_nas_access_groups.default.FileSystemType}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackNasAccessGroupsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_nas_access_groups.default.id}"]`,

			"access_group_name": `"${alibabacloudstack_nas_access_groups.default.AccessGroupName}"`,
			"file_system_type":  `"${alibabacloudstack_nas_access_groups.default.FileSystemType}"`}),
		fakeConfig: testAccCheckAlibabacloudstackNasAccessGroupsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_nas_access_groups.default.id}_fake"]`,

			"access_group_name": `"${alibabacloudstack_nas_access_groups.default.AccessGroupName}_fake"`,
			"file_system_type":  `"${alibabacloudstack_nas_access_groups.default.FileSystemType}_fake"`}),
	}

	AlibabacloudstackNasAccessGroupsCheckInfo.dataSourceTestCheck(t, rand, idsConf, access_group_nameConf, file_system_typeConf, allConf)
}

var existAlibabacloudstackNasAccessGroupsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#":    "1",
		"groups.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackNasAccessGroupsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#": "0",
	}
}

var AlibabacloudstackNasAccessGroupsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_nas_access_groups.default",
	existMapFunc: existAlibabacloudstackNasAccessGroupsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackNasAccessGroupsMapFunc,
}

func testAccCheckAlibabacloudstackNasAccessGroupsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackNasAccessGroups%d"
}






data "alibabacloudstack_nas_access_groups" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
