package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackNasMountTargetsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackNasMountTargetsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_nas_mount_targets.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackNasMountTargetsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_nas_mount_targets.default.id}_fake"]`,
		}),
	}

	file_system_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackNasMountTargetsSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_nas_mount_targets.default.id}"]`,
			"file_system_id": `"${alibabacloudstack_nas_mount_targets.default.FileSystemId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackNasMountTargetsSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_nas_mount_targets.default.id}_fake"]`,
			"file_system_id": `"${alibabacloudstack_nas_mount_targets.default.FileSystemId}_fake"`,
		}),
	}

	mount_target_domainConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackNasMountTargetsSourceConfig(rand, map[string]string{
			"ids":                 `["${alibabacloudstack_nas_mount_targets.default.id}"]`,
			"mount_target_domain": `"${alibabacloudstack_nas_mount_targets.default.MountTargetDomain}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackNasMountTargetsSourceConfig(rand, map[string]string{
			"ids":                 `["${alibabacloudstack_nas_mount_targets.default.id}_fake"]`,
			"mount_target_domain": `"${alibabacloudstack_nas_mount_targets.default.MountTargetDomain}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackNasMountTargetsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_nas_mount_targets.default.id}"]`,

			"file_system_id":      `"${alibabacloudstack_nas_mount_targets.default.FileSystemId}"`,
			"mount_target_domain": `"${alibabacloudstack_nas_mount_targets.default.MountTargetDomain}"`}),
		fakeConfig: testAccCheckAlibabacloudstackNasMountTargetsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_nas_mount_targets.default.id}_fake"]`,

			"file_system_id":      `"${alibabacloudstack_nas_mount_targets.default.FileSystemId}_fake"`,
			"mount_target_domain": `"${alibabacloudstack_nas_mount_targets.default.MountTargetDomain}_fake"`}),
	}

	AlibabacloudstackNasMountTargetsCheckInfo.dataSourceTestCheck(t, rand, idsConf, file_system_idConf, mount_target_domainConf, allConf)
}

var existAlibabacloudstackNasMountTargetsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"targets.#":    "1",
		"targets.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackNasMountTargetsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"targets.#": "0",
	}
}

var AlibabacloudstackNasMountTargetsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_nas_mount_targets.default",
	existMapFunc: existAlibabacloudstackNasMountTargetsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackNasMountTargetsMapFunc,
}

func testAccCheckAlibabacloudstackNasMountTargetsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackNasMountTargets%d"
}






data "alibabacloudstack_nas_mount_targets" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
