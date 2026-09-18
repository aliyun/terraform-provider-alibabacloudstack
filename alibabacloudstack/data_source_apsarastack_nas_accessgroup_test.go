package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackNasAccessGroupsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_nas_access_groups.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAlibabacloudstackNasAccessGroups%d", rand),
		dataSourceNasAccessGroupsConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_nas_access_groups.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_nas_access_groups.default.id}_fake"},
		}),
	}

	access_group_nameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alibabacloudstack_nas_access_groups.default.id}"},
			"access_group_name": "${alibabacloudstack_nas_access_groups.default.AccessGroupName}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alibabacloudstack_nas_access_groups.default.id}_fake"},
			"access_group_name": "${alibabacloudstack_nas_access_groups.default.AccessGroupName}_fake",
		}),
	}

	file_system_typeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alibabacloudstack_nas_access_groups.default.id}"},
			"file_system_type": "${alibabacloudstack_nas_access_groups.default.FileSystemType}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alibabacloudstack_nas_access_groups.default.id}_fake"},
			"file_system_type": "${alibabacloudstack_nas_access_groups.default.FileSystemType}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alibabacloudstack_nas_access_groups.default.id}"},
			"access_group_name": "${alibabacloudstack_nas_access_groups.default.AccessGroupName}",
			"file_system_type":  "${alibabacloudstack_nas_access_groups.default.FileSystemType}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alibabacloudstack_nas_access_groups.default.id}_fake"},
			"access_group_name": "${alibabacloudstack_nas_access_groups.default.AccessGroupName}_fake",
			"file_system_type":  "${alibabacloudstack_nas_access_groups.default.FileSystemType}_fake",
		}),
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

func dataSourceNasAccessGroupsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
`, name)
}
