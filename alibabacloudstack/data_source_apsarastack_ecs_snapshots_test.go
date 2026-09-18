package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackEcsSnapshotsDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testaccSnapshotDataSourceBasic%d", rand)
	resourceId := "data.alibabacloudstack_snapshots.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEcsSnapshotsConfigDependence)

	idsConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_snapshot.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_snapshot.default.id}_fake"},
		}),
	}

	instanceIdConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_ecs_instance.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_ecs_instance.default.id}_fake",
		}),
	}

	diskIdConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"disk_id": "${alibabacloudstack_snapshot.default.disk_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"disk_id": "${alibabacloudstack_snapshot.default.disk_id}_fake",
		}),
	}

	nameRegexConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_snapshot.default.snapshot_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_snapshot.default.snapshot_name}_fake",
		}),
	}

	statusConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alibabacloudstack_snapshot.default.id}"},
			"status": "accomplished",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alibabacloudstack_snapshot.default.id}"},
			"status": "failed",
		}),
	}

	typeConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"${alibabacloudstack_snapshot.default.id}"},
			"type": "user",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"${alibabacloudstack_snapshot.default.id}"},
			"type": "auto",
		}),
	}

	sourceDiskTypeConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alibabacloudstack_snapshot.default.id}"},
			"source_disk_type": "System",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alibabacloudstack_snapshot.default.id}"},
			"source_disk_type": "Data",
		}),
	}

	usageConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":   []string{"${alibabacloudstack_snapshot.default.id}"},
			"usage": "none",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":   []string{"${alibabacloudstack_snapshot.default.id}"},
			"usage": "image",
		}),
	}

	allConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alibabacloudstack_snapshot.default.id}"},
			"instance_id":      "${alibabacloudstack_ecs_instance.default.id}",
			"disk_id":          "${alibabacloudstack_snapshot.default.disk_id}",
			"name_regex":       "${alibabacloudstack_snapshot.default.snapshot_name}",
			"status":           "accomplished",
			"type":             "user",
			"source_disk_type": "System",
			"usage":            "none",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alibabacloudstack_snapshot.default.id}_fake"},
			"instance_id":      "${alibabacloudstack_ecs_instance.default.id}_fake",
			"disk_id":          "${alibabacloudstack_snapshot.default.disk_id}_fake",
			"name_regex":       "${alibabacloudstack_snapshot.default.snapshot_name}_fake",
			"status":           "failed",
			"type":             "auto",
			"source_disk_type": "Data",
			"usage":            "image",
		}),
	}

	var existSnapshotsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"names.#":                      "1",
			"snapshots.#":                  "1",
			"snapshots.0.id":               CHECKSET,
			"snapshots.0.name":             name,
			"snapshots.0.description":      name,
			"snapshots.0.progress":         CHECKSET,
			"snapshots.0.source_disk_id":   CHECKSET,
			"snapshots.0.source_disk_size": "20",
			"snapshots.0.source_disk_type": CHECKSET,
			"snapshots.0.product_code":     "",
			"snapshots.0.remain_time":      CHECKSET,
			"snapshots.0.creation_time":    CHECKSET,
			"snapshots.0.status":           "accomplished",
			"snapshots.0.usage":            "none",
		}
	}

	var fakeSnapshotsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"snapshots.#": "0",
		}
	}

	var snapshotsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existSnapshotsMapFunc,
		fakeMapFunc:  fakeSnapshotsMapFunc,
	}

	snapshotsCheckInfo.dataSourceTestCheck(t, rand, allConfig, sourceDiskTypeConfig, nameRegexConfig, statusConfig, typeConfig, usageConfig, idsConfig, instanceIdConfig, diskIdConfig)
}

func dataSourceEcsSnapshotsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s

resource "alibabacloudstack_snapshot" "default" {
	disk_id = "${alibabacloudstack_ecs_instance.default.system_disk_id}"
	name = "${var.name}"
	description = "${var.name}"
}
`, name, ECSInstanceCommonTestCase)
}
