package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackNasMountTargetDataSource(t *testing.T) {
	rand := getAccTestRandInt(100000, 999999)
	fileSystemIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alibabacloudstack_nas_mount_target.default.file_system_id}"`,
		}),
	}
	accessGroupNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id":    `"${alibabacloudstack_nas_mount_target.default.file_system_id}"`,
			"access_group_name": `"${alibabacloudstack_nas_access_group.default.access_group_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id":    `"${alibabacloudstack_nas_mount_target.default.file_system_id}"`,
			"access_group_name": `"${alibabacloudstack_nas_access_group.default.access_group_name}_fake"`,
		}),
	}
	typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alibabacloudstack_nas_mount_target.default.file_system_id}"`,
			"type":           `"${alibabacloudstack_nas_access_group.default.access_group_type}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alibabacloudstack_nas_mount_target.default.file_system_id}"`,
			"type":           `"${alibabacloudstack_nas_access_group.default.access_group_type}_fake"`,
		}),
	}
	netWorkTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alibabacloudstack_nas_mount_target.default.file_system_id}"`,
			"network_type":   `"${alibabacloudstack_nas_access_group.default.access_group_type}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alibabacloudstack_nas_mount_target.default.file_system_id}"`,
			"network_type":   `"${alibabacloudstack_nas_access_group.default.access_group_type}_fake"`,
		}),
	}
	mountTargetDomainConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id":      `"${alibabacloudstack_nas_mount_target.default.file_system_id}"`,
			"mount_target_domain": `split(":",alibabacloudstack_nas_mount_target.default.id)[1]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id":      `"${alibabacloudstack_nas_mount_target.default.file_system_id}"`,
			"mount_target_domain": `"fake"`,
		}),
	}
	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alibabacloudstack_nas_mount_target.default.file_system_id}"`,
			"vpc_id":         `"${alibabacloudstack_vpc_vpc.default.id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alibabacloudstack_nas_mount_target.default.file_system_id}"`,
			"vpc_id":         `"${alibabacloudstack_vpc_vpc.default.id}"`,
		}),
	}
	vswitchIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alibabacloudstack_nas_mount_target.default.file_system_id}"`,
			"vswitch_id":     `"${alibabacloudstack_nas_mount_target.default.vswitch_id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alibabacloudstack_nas_mount_target.default.file_system_id}"`,
			"vswitch_id":     `"fake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alibabacloudstack_nas_mount_target.default.file_system_id}"`,
			"ids":            `[split(":",alibabacloudstack_nas_mount_target.default.id)[1]]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alibabacloudstack_nas_mount_target.default.file_system_id}"`,
			"ids":            `["fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alibabacloudstack_nas_mount_target.default.file_system_id}"`,
			"status":         `"Active"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alibabacloudstack_nas_mount_target.default.file_system_id}"`,
			"status":         `"Inactive"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id":      `"${alibabacloudstack_nas_mount_target.default.file_system_id}"`,
			"access_group_name":   `"${alibabacloudstack_nas_mount_target.default.access_group_name}"`,
			"vswitch_id":          `"${alibabacloudstack_nas_mount_target.default.vswitch_id}"`,
			"type":                `"${alibabacloudstack_nas_access_group.default.access_group_type}"`,
			"network_type":        `"${alibabacloudstack_nas_access_group.default.access_group_type}"`,
			"vpc_id":              `"${alibabacloudstack_vpc_vpc.default.id}"`,
			"mount_target_domain": `split(":",alibabacloudstack_nas_mount_target.default.id)[1]`,
			"status":              `"Active"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id":      `"${alibabacloudstack_nas_mount_target.default.file_system_id}"`,
			"access_group_name":   `"${alibabacloudstack_nas_mount_target.default.access_group_name}"`,
			"vswitch_id":          `"${alibabacloudstack_nas_mount_target.default.vswitch_id}_fake"`,
			"type":                `"${alibabacloudstack_nas_access_group.default.access_group_type}_fake"`,
			"network_type":        `"${alibabacloudstack_nas_access_group.default.access_group_type}_fake}"`,
			"vpc_id":              `"${alibabacloudstack_vpc_vpc.default.id}"`,
			"mount_target_domain": `"fake"`,
			"status":              `"Inactive"`,
		}),
	}
	mountTargetCheckInfo.dataSourceTestCheck(t, rand, fileSystemIdConf, accessGroupNameConf, typeConf, netWorkTypeConf, mountTargetDomainConf, vpcIdConf, vswitchIdConf, idsConf, statusConf, allConf)
}

func testAccCheckAlibabacloudStackMountTargetDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
			default = "tf-testAccCheck-nasmount%d"
}

variable "storage_type" {
  default = "Capacity"
}

%s

data "alibabacloudstack_nas_protocols" "default" {
        type = "${var.storage_type}"
}
resource "alibabacloudstack_nas_file_system" "default" {
  description = "${var.name}"
  storage_type = "${var.storage_type}"
  protocol_type = "${data.alibabacloudstack_nas_protocols.default.protocols.0}"
}
resource "alibabacloudstack_nas_access_group" "default" {
			access_group_name = "${var.name}"
			access_group_type = "Vpc"
			description = "${var.name}"
}
resource "alibabacloudstack_nas_mount_target" "default" {
			file_system_id = "${alibabacloudstack_nas_file_system.default.id}"
			access_group_name = "${alibabacloudstack_nas_access_group.default.access_group_name}"
			vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
}
data "alibabacloudstack_nas_mount_targets" "default" {
		%s
}`, rand, VSwitchCommonTestCase, strings.Join(pairs, "\n  "))
	return config
}

var existMountTargetMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"test": NOSET,
		// "targets.0.type":                "Vpc",
		// "targets.0.network_type":        "Vpc",
		// "targets.0.status":              "Active",
		// "targets.0.vpc_id":              CHECKSET,
		// "targets.0.mount_target_domain": CHECKSET,
		// "targets.0.vswitch_id":          CHECKSET,
		// "targets.0.access_group_name":   fmt.Sprintf("tf-testAccNasConfig-%d", rand),
		// "ids.#":                         "1",
		// "ids.0":                         CHECKSET,
	}
}

var fakeMountTargetMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"test": NOSET,
		// "targets.#": "0",
		// "ids.#":     "0",
	}
}

var mountTargetCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_nas_mount_targets.default",
	existMapFunc: existMountTargetMapCheck,
	fakeMapFunc:  fakeMountTargetMapCheck,
}
