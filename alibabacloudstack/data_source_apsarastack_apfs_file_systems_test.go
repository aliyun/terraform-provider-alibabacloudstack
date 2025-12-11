package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackApfsFileSystemsDataSource_basic(t *testing.T) {
	resourceId := "data.alibabacloudstack_apfs_file_systems.default"
	rand := getAccTestRandInt(100000, 999999)
	name := fmt.Sprintf("tf-testacc%v", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, ApfsFileSystemDependenceNew)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_apfs_file_system.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^test-fake-name$",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_apfs_file_system.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"fake-id"},
		}),
	}

	typeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"file_system_type": "efs",
			"ids":              []string{"${alibabacloudstack_apfs_file_system.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"file_system_type": "standard",
			"ids":              []string{"${alibabacloudstack_apfs_file_system.default.id}"},
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"status": "Running",
			"ids":    []string{"${alibabacloudstack_apfs_file_system.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"status": "Stopped",
			"ids":    []string{"${alibabacloudstack_apfs_file_system.default.id}"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":       "${alibabacloudstack_apfs_file_system.default.description}",
			"ids":              []string{"${alibabacloudstack_apfs_file_system.default.id}"},
			"file_system_type": "efs",
			"status":           "Running",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":       "${alibabacloudstack_apfs_file_system.default.description}",
			"ids":              []string{"${alibabacloudstack_apfs_file_system.default.id}"},
			"file_system_type": "efs",
			"status":           "Stopped",
		}),
	}

	var existApfsFileSystemsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"file_systems.#":               "1",
			"file_systems.0.id":            CHECKSET,
			"file_systems.0.zone_id":       CHECKSET,
			"file_systems.0.status":        "Running",
			"file_systems.0.protocol_type": "EFS",
		}
	}

	var fakeApfsFileSystemsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"file_systems.#": "0",
		}
	}

	var apfsFileSystemsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existApfsFileSystemsMapFunc,
		fakeMapFunc:  fakeApfsFileSystemsMapFunc,
	}

	apfsFileSystemsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, typeConf, statusConf, allConf)
}

func ApfsFileSystemDependenceNew(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alibabacloudstack_zones" "default" {
  enable_details = true
}

resource "alibabacloudstack_apfs_file_system" "default" {
  zone_id      = "${data.alibabacloudstack_zones.default.zones.0.id}"
  cluster_id   = "EfsStorageCluster-A-20251125-0131"
  storage_type = "parastor_advance_300"
  volume_size  = 1024
  description  = "${var.name}"
}`, name)
}
