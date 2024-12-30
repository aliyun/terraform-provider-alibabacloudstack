package alibabacloudstack

import (
	"fmt"

	"testing"
)

func TestAccAlibabacloudStackSnapshotsDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testaccSnapshotDataSourceBasic%d", rand)
	resourceId := "data.alibabacloudstack_snapshots.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceSnapshotsConfigDependence)

	idsConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_snapshot.default.name}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_snapshot.default.name}_fake"},
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
			"disk_id": "${alibabacloudstack_ecs_instance.default.system_disk_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"disk_id": "${alibabacloudstack_ecs_instance.default.system_disk_id}_fake",
		}),
	}

	nameRegexConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": name + "_fake",
		}),
	}

	//	statusConfig := dataSourceTestAccConfig{
	//		existConfig: testAccConfig(map[string]interface{}{
	//			"ids":    []string{"${alibabacloudstack_snapshot.default.id}"},
	//			"status": "accomplished",
	//		}),
	//		fakeConfig: testAccConfig(map[string]interface{}{
	//			"ids":    []string{"${alibabacloudstack_snapshot.default.id}"},
	//			"status": "failed",
	//		}),
	//	}

	// typeConfig := dataSourceTestAccConfig{
	// 	existConfig: testAccConfig(map[string]interface{}{
	// 		"ids":  []string{"${alibabacloudstack_snapshot.default.id}"},
	// 		"type": "user",
	// 	}),
	// 	fakeConfig: testAccConfig(map[string]interface{}{
	// 		"ids":  []string{"${alibabacloudstack_snapshot.default.id}"},
	// 		"type": "auto",
	// 	}),
	// }

	//	sourceDiskTypeConfig := dataSourceTestAccConfig{
	//		existConfig: testAccConfig(map[string]interface{}{
	//			"ids":              []string{"${alibabacloudstack_snapshot.default.id}"},
	//			"source_disk_type": "Data",
	//		}),
	//		fakeConfig: testAccConfig(map[string]interface{}{
	//			"ids":              []string{"${alibabacloudstack_snapshot.default.id}"},
	//			"source_disk_type": "System",
	//		}),
	//	}

	//	usageConfig := dataSourceTestAccConfig{
	//		existConfig: testAccConfig(map[string]interface{}{
	//			"ids":   []string{"${alibabacloudstack_snapshot.default.id}"},
	//			"usage": "none",
	//		}),
	//		fakeConfig: testAccConfig(map[string]interface{}{
	//			"ids":   []string{"${alibabacloudstack_snapshot.default.id}"},
	//			"usage": "image",
	//		}),
	//	}
	//
	//	tagsConfig := dataSourceTestAccConfig{
	//		existConfig: testAccConfig(map[string]interface{}{
	//			"ids": []string{"${alibabacloudstack_snapshot.default.id}"},
	//			"tags": map[string]interface{}{
	//				"version": "1.0",
	//			},
	//		}),
	//		fakeConfig: testAccConfig(map[string]interface{}{
	//			"ids": []string{"${alibabacloudstack_snapshot.default.id}"},
	//			"tags": map[string]interface{}{
	//				"version": "1.0_fake",
	//			},
	//		}),
	//	}

	//	allConfig := dataSourceTestAccConfig{
	//		existConfig: testAccConfig(map[string]interface{}{
	//			"ids":              []string{"${alibabacloudstack_snapshot.default.id}"},
	//			"instance_id":      "${alibabacloudstack_instance.default.id}",
	//			"disk_id":          "${alibabacloudstack_disk_attachment.default.disk_id}",
	//			"name_regex":       name,
	//			"status":           "accomplished",
	//			"type":             "user",
	//			"source_disk_type": "Data",
	//			"usage":            "none",
	//			"tags": map[string]interface{}{
	//				"version": "1.0",
	//			},
	//		}),
	//		fakeConfig: testAccConfig(map[string]interface{}{
	//			"ids":              []string{"${alibabacloudstack_snapshot.default.id}"},
	//			"instance_id":      "${alibabacloudstack_instance.default.id}",
	//			"disk_id":          "${alibabacloudstack_disk_attachment.default.disk_id}",
	//			"name_regex":       name,
	//			"status":           "accomplished",
	//			"type":             "user",
	//			"source_disk_type": "Data",
	//			"usage":            "none",
	//			"tags": map[string]interface{}{
	//				"version": "1.0_fake",
	//			},
	//		}),
	//	}

	var existSnapshotsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"test": NOSET,
			// "ids.#": "1",
			// "names.#":                      "1",
			// "snapshots.#":                  "1",
			// "snapshots.0.id": CHECKSET,
			// "snapshots.0.name":             name,
			// "snapshots.0.description":      name,
			// "snapshots.0.progress":         CHECKSET,
			// "snapshots.0.source_disk_id":   CHECKSET,
			// "snapshots.0.source_disk_size": "20",
			// "snapshots.0.source_disk_type": CHECKSET,
			// "snapshots.0.product_code":     "",
			// "snapshots.0.remain_time":      CHECKSET,
			// "snapshots.0.creation_time":    CHECKSET,
			// "snapshots.0.status":           "accomplished",
			// "snapshots.0.usage":            "none",
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

	snapshotsCheckInfo.dataSourceTestCheck(t, rand, idsConfig, instanceIdConfig, diskIdConfig, nameRegexConfig) // statusConfig,sourceDiskTypeConfig,usageConfig, tagsConfig,allConfig
}

func dataSourceSnapshotsConfigDependence(name string) string {
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
