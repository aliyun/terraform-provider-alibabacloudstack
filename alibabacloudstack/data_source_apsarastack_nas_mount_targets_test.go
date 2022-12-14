package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackNasMountTargetDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100000, 999999)
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
			"type":           `"${alibabacloudstack_nas_access_group.default.type}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alibabacloudstack_nas_mount_target.default.file_system_id}"`,
			"type":           `"${alibabacloudstack_nas_access_group.default.type}_fake"`,
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
			"vpc_id":         `"${alibabacloudstack_vpc.default.id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alibabacloudstack_nas_mount_target.default.file_system_id}"`,
			"vpc_id":         `"${alibabacloudstack_vpc.default.id}"`,
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
			"type":                `"${alibabacloudstack_nas_access_group.default.type}"`,
			"network_type":        `"${alibabacloudstack_nas_access_group.default.access_group_type}"`,
			"vpc_id":              `"${alibabacloudstack_vpc.default.id}"`,
			"mount_target_domain": `split(":",alibabacloudstack_nas_mount_target.default.id)[1]`,
			"status":              `"Active"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackMountTargetDataSourceConfig(rand, map[string]string{
			"file_system_id":      `"${alibabacloudstack_nas_mount_target.default.file_system_id}"`,
			"access_group_name":   `"${alibabacloudstack_nas_mount_target.default.access_group_name}"`,
			"vswitch_id":          `"${alibabacloudstack_nas_mount_target.default.vswitch_id}_fake"`,
			"type":                `"${alibabacloudstack_nas_access_group.default.type}_fake"`,
			"network_type":        `"${alibabacloudstack_nas_access_group.default.access_group_type}_fake}"`,
			"vpc_id":              `"${alibabacloudstack_vpc.default.id}"`,
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
			default = "tf-testAccVswitch"
}
variable "description" {
  default = "tf-testAccCheckAlibabacloudStackFileSystemsDataSource"
}
//data "alibabacloudstack_vpcs" "default" {
//			name_regex = "lz_vpc3"
//}

resource "alibabacloudstack_vpc" "default" {
  cidr_block = "172.16.0.0/16"
  name = "${var.name}"
}
data "alibabacloudstack_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name = "${var.name}"
}
variable "storage_type" {
  default = "Capacity"
}
data "alibabacloudstack_nas_protocols" "default" {
        type = "${var.storage_type}"
}
resource "alibabacloudstack_nas_file_system" "default" {
  description = "${var.description}"
  storage_type = "${var.storage_type}"
  protocol_type = "${data.alibabacloudstack_nas_protocols.default.protocols.0}"
}
resource "alibabacloudstack_nas_access_group" "default" {
			access_group_name = "tf-testAccNasConfig-%d"
			access_group_type = "Vpc"
			description = "tf-testAccNasConfig"
}
resource "alibabacloudstack_nas_mount_target" "default" {
			file_system_id = "${alibabacloudstack_nas_file_system.default.id}"
			access_group_name = "${alibabacloudstack_nas_access_group.default.access_group_name}"
			vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}
data "alibabacloudstack_nas_mount_targets" "default" {
		%s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existMountTargetMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"targets.0.type":                "Vpc",
		"targets.0.network_type":        "Vpc",
		"targets.0.status":              "Active",
		"targets.0.vpc_id":              CHECKSET,
		"targets.0.mount_target_domain": CHECKSET,
		"targets.0.vswitch_id":          CHECKSET,
		"targets.0.access_group_name":   fmt.Sprintf("tf-testAccNasConfig-%d", rand),
		"ids.#":                         "1",
		"ids.0":                         CHECKSET,
	}
}

var fakeMountTargetMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"targets.#": "0",
		"ids.#":     "0",
	}
}

var mountTargetCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_nas_mount_targets.default",
	existMapFunc: existMountTargetMapCheck,
	fakeMapFunc:  fakeMountTargetMapCheck,
}
