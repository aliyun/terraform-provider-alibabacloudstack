package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackNasMountTargetDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_nas_mount_targets.default"
	name := fmt.Sprintf("tf-testnasfs%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlibabacloudStackMountTargetDataSourceConfig)

	fileSystemIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"file_system_id": "${alibabacloudstack_nas_mount_target.default.file_system_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"file_system_id":    "${alibabacloudstack_nas_mount_target.default.file_system_id}_fake",
		}),
	}
	accessGroupNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"file_system_id":    "${alibabacloudstack_nas_mount_target.default.file_system_id}",
			"access_group_name": "${alibabacloudstack_nas_access_group.default.access_group_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"file_system_id":    "${alibabacloudstack_nas_mount_target.default.file_system_id}",
			"access_group_name": "${alibabacloudstack_nas_access_group.default.access_group_name}_fake",
		}),
	}
	typeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"file_system_id": "${alibabacloudstack_nas_mount_target.default.file_system_id}",
			"type":           "${alibabacloudstack_nas_access_group.default.access_group_type}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"file_system_id": "${alibabacloudstack_nas_mount_target.default.file_system_id}",
			"type":           "${alibabacloudstack_nas_access_group.default.access_group_type}_fake",
		}),
	}
	netWorkTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"file_system_id": "${alibabacloudstack_nas_mount_target.default.file_system_id}",
			"network_type":   "${alibabacloudstack_nas_access_group.default.access_group_type}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"file_system_id": "${alibabacloudstack_nas_mount_target.default.file_system_id}",
			"network_type":   "${alibabacloudstack_nas_access_group.default.access_group_type}_fake",
		}),
	}
	mountTargetDomainConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"file_system_id":      "${alibabacloudstack_nas_mount_target.default.file_system_id}",
			"mount_target_domain": "${alibabacloudstack_nas_mount_target.default.mount_target_domain}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"file_system_id":      "${alibabacloudstack_nas_mount_target.default.file_system_id}",
			"mount_target_domain": "fake",
		}),
	}
	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"file_system_id": "${alibabacloudstack_nas_mount_target.default.file_system_id}",
			"vpc_id":         "${alibabacloudstack_vpc_vpc.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"file_system_id": "${alibabacloudstack_nas_mount_target.default.file_system_id}",
			"vpc_id":         "${alibabacloudstack_vpc_vpc.default.id}_fake",
		}),
	}
	vswitchIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"file_system_id": "${alibabacloudstack_nas_mount_target.default.file_system_id}",
			"vswitch_id":     "${alibabacloudstack_nas_mount_target.default.vswitch_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"file_system_id": "${alibabacloudstack_nas_mount_target.default.file_system_id}",
			"vswitch_id":     "fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"file_system_id": "${alibabacloudstack_nas_mount_target.default.file_system_id}",
			"ids":            []string{"${alibabacloudstack_nas_mount_target.default.mount_target_domain}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"file_system_id": "${alibabacloudstack_nas_mount_target.default.file_system_id}",
			"ids":            []string{"fake"},
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"file_system_id": "${alibabacloudstack_nas_mount_target.default.file_system_id}",
			"status":         "Active",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"file_system_id": "${alibabacloudstack_nas_mount_target.default.file_system_id}",
			"status":         "Inactive",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"file_system_id":      "${alibabacloudstack_nas_mount_target.default.file_system_id}",
			"access_group_name":   "${alibabacloudstack_nas_mount_target.default.access_group_name}",
			"vswitch_id":          "${alibabacloudstack_nas_mount_target.default.vswitch_id}",
			"type":                "${alibabacloudstack_nas_access_group.default.access_group_type}",
			"network_type":        "${alibabacloudstack_nas_access_group.default.access_group_type}",
			"vpc_id":              "${alibabacloudstack_vpc_vpc.default.id}",
			"mount_target_domain": "${alibabacloudstack_nas_mount_target.default.mount_target_domain}",
			"status":              "Active",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"file_system_id":      "${alibabacloudstack_nas_mount_target.default.file_system_id}",
			"access_group_name":   "${alibabacloudstack_nas_mount_target.default.access_group_name}",
			"vswitch_id":          "${alibabacloudstack_nas_mount_target.default.vswitch_id}_fake",
			"type":                "${alibabacloudstack_nas_access_group.default.access_group_type}_fake",
			"network_type":        "${alibabacloudstack_nas_access_group.default.access_group_type}_fake}",
			"vpc_id":              "${alibabacloudstack_vpc_vpc.default.id}",
			"mount_target_domain": "fake",
			"status":              "Inactive",
		}),
	}

	var existMountTargetMapCheck = func(rand int) map[string]string {
		return map[string]string{
			"targets.0.type":                "Vpc",
			"targets.0.network_type":        "Vpc",
			"targets.0.status":              "Active",
			"targets.0.vpc_id":              CHECKSET,
			"targets.0.mount_target_domain": CHECKSET,
			"targets.0.vswitch_id":          CHECKSET,
			"targets.0.access_group_name":   name,
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
		resourceId:   resourceId,
		existMapFunc: existMountTargetMapCheck,
		fakeMapFunc:  fakeMountTargetMapCheck,
	}

	mountTargetCheckInfo.dataSourceTestCheck(t, rand, fileSystemIdConf, accessGroupNameConf, typeConf, netWorkTypeConf, mountTargetDomainConf, vpcIdConf, vswitchIdConf, idsConf, statusConf, allConf)
}

func testAccCheckAlibabacloudStackMountTargetDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
}



%s

%s

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
`, name, VSwitchCommonTestCase, NasCommonTestCase)
}
