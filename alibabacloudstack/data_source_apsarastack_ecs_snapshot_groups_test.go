package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackEcsSnapshotGroupsDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_ecs_snapshot_groups.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAcc%sEcsSnapshotGroupsDataSource-%d", defaultRegionToTest, rand),
		dataSourceEcsSnapshotGroupsDependence)

	descriptionRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_ecs_snapshot_group.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_ecs_snapshot_group.default.description}-fakeTestAcccc",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ecs_snapshot_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ecs_snapshot_group.default.id}-fakeTestAcccc"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_ecs_snapshot_group.default.description}",
			"ids":        []string{"${alibabacloudstack_ecs_snapshot_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_ecs_snapshot_group.default.description}-fakeTestAcccc",
			"ids":        []string{"${alibabacloudstack_ecs_snapshot_group.default.id}-fakeTestAcccc"},
		}),
	}

	var existEcsSnapshotGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"ids.0":                            CHECKSET,
			"snapshot_groups.#":             "1",
			"snapshot_groups.0.description": fmt.Sprintf("tf-testAcc%sEcsSnapshotGroupsDataSource-%d", defaultRegionToTest, rand),
			"snapshot_groups.0.create_time":   CHECKSET,
			"snapshot_groups.0.instance_id": CHECKSET,
			"snapshot_groups.0.disk_ids":   CHECKSET,
			"snapshot_groups.0.instant_access":                    CHECKSET,
			"snapshot_groups.0.instant_access_retention_days":        CHECKSET,
			"snapshot_groups.0.snapshot_group_id":      CHECKSET,
			"snapshot_groups.0.snapshot_group_name":        CHECKSET,
		}
	}

	var fakeEcsSnapshotGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                "0",
			"snapshot_groups.#": "0",
		}
	}

	var EcsSnapshotGroupsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEcsSnapshotGroupsMapFunc,
		fakeMapFunc:  fakeEcsSnapshotGroupsMapFunc,
	}

	EcsSnapshotGroupsCheckInfo.dataSourceTestCheck(t, rand, descriptionRegexConf, idsConf, allConf)
}

func dataSourceEcsSnapshotGroupsDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}

%s

%s

%s

resource "alibabacloudstack_ecs_instance" "default" {
  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type        = "${local.default_instance_type_id}"
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  system_disk_size     = 20
  system_disk_name     = "test_sys_disk"
  security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
  instance_name        = "${var.name}_ecs"
  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
  zone_id    		   = data.alibabacloudstack_zones.default.zones.0.id
  is_outdated          = false
  data_disks {
      name                 = "disk1"
      category             = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
      size                 = 20
      delete_with_instance = true
    }
  data_disks {
      name                 = "disk2"
      category             = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
      size                 = 20
	  delete_with_instance = true
    }
  lifecycle {
    ignore_changes = [
      instance_type
    ]
  }
}

data "alibabacloudstack_ecs_disks" "disks" {
	instance_id = "${alibabacloudstack_ecs_instance.default.id}"
}

resource "alibabacloudstack_ecs_snapshot_group" "default" {
  	description =                   "${var.name}"
	instance_id =                   "${alibabacloudstack_ecs_instance.default.id}"
	instant_access =                "true"
	instant_access_retention_days = "7"
	snapshot_group_name =           "${var.name}"
	disk_ids = [
		"${data.alibabacloudstack_ecs_disks.disks.disks.0.id}",
		"${data.alibabacloudstack_ecs_disks.disks.disks.1.id}",
		"${data.alibabacloudstack_ecs_disks.disks.disks.2.id}",
	]
}

 `, name, SecurityGroupCommonTestCase, DataAlibabacloudstackImages, DataAlibabacloudstackInstanceTypes)
}
