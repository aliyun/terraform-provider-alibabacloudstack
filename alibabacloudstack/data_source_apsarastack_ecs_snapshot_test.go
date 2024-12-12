package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackEcsSnapshotsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsSnapshotsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_snapshots.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsSnapshotsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_snapshots.default.id}_fake"]`,
		}),
	}

	encryptedConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsSnapshotsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_snapshots.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsSnapshotsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_snapshots.default.id}_fake"]`,
		}),
	}

	resource_group_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsSnapshotsSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_ecs_snapshots.default.id}"]`,
			"resource_group_id": `"${alibabacloudstack_ecs_snapshots.default.ResourceGroupId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsSnapshotsSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_ecs_snapshots.default.id}_fake"]`,
			"resource_group_id": `"${alibabacloudstack_ecs_snapshots.default.ResourceGroupId}_fake"`,
		}),
	}

	snapshot_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsSnapshotsSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_ecs_snapshots.default.id}"]`,
			"snapshot_id": `"${alibabacloudstack_ecs_snapshots.default.SnapshotId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsSnapshotsSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_ecs_snapshots.default.id}_fake"]`,
			"snapshot_id": `"${alibabacloudstack_ecs_snapshots.default.SnapshotId}_fake"`,
		}),
	}

	snapshot_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsSnapshotsSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_ecs_snapshots.default.id}"]`,
			"snapshot_name": `"${alibabacloudstack_ecs_snapshots.default.SnapshotName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsSnapshotsSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_ecs_snapshots.default.id}_fake"]`,
			"snapshot_name": `"${alibabacloudstack_ecs_snapshots.default.SnapshotName}_fake"`,
		}),
	}

	snapshot_typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsSnapshotsSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_ecs_snapshots.default.id}"]`,
			"snapshot_type": `"${alibabacloudstack_ecs_snapshots.default.SnapshotType}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsSnapshotsSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_ecs_snapshots.default.id}_fake"]`,
			"snapshot_type": `"${alibabacloudstack_ecs_snapshots.default.SnapshotType}_fake"`,
		}),
	}

	source_disk_typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsSnapshotsSourceConfig(rand, map[string]string{
			"ids":              `["${alibabacloudstack_ecs_snapshots.default.id}"]`,
			"source_disk_type": `"${alibabacloudstack_ecs_snapshots.default.SourceDiskType}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsSnapshotsSourceConfig(rand, map[string]string{
			"ids":              `["${alibabacloudstack_ecs_snapshots.default.id}_fake"]`,
			"source_disk_type": `"${alibabacloudstack_ecs_snapshots.default.SourceDiskType}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsSnapshotsSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_ecs_snapshots.default.id}"]`,
			"status": `"${alibabacloudstack_ecs_snapshots.default.Status}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsSnapshotsSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_ecs_snapshots.default.id}_fake"]`,
			"status": `"${alibabacloudstack_ecs_snapshots.default.Status}_fake"`,
		}),
	}

	usageConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsSnapshotsSourceConfig(rand, map[string]string{
			"ids":   `["${alibabacloudstack_ecs_snapshots.default.id}"]`,
			"usage": `"${alibabacloudstack_ecs_snapshots.default.Usage}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsSnapshotsSourceConfig(rand, map[string]string{
			"ids":   `["${alibabacloudstack_ecs_snapshots.default.id}_fake"]`,
			"usage": `"${alibabacloudstack_ecs_snapshots.default.Usage}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsSnapshotsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_snapshots.default.id}"]`,

			"resource_group_id": `"${alibabacloudstack_ecs_snapshots.default.ResourceGroupId}"`,
			"snapshot_id":       `"${alibabacloudstack_ecs_snapshots.default.SnapshotId}"`,
			"snapshot_name":     `"${alibabacloudstack_ecs_snapshots.default.SnapshotName}"`,
			"snapshot_type":     `"${alibabacloudstack_ecs_snapshots.default.SnapshotType}"`,
			"source_disk_type":  `"${alibabacloudstack_ecs_snapshots.default.SourceDiskType}"`,
			"status":            `"${alibabacloudstack_ecs_snapshots.default.Status}"`,
			"usage":             `"${alibabacloudstack_ecs_snapshots.default.Usage}"`}),
		fakeConfig: testAccCheckAlibabacloudstackEcsSnapshotsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_snapshots.default.id}_fake"]`,

			"resource_group_id": `"${alibabacloudstack_ecs_snapshots.default.ResourceGroupId}_fake"`,
			"snapshot_id":       `"${alibabacloudstack_ecs_snapshots.default.SnapshotId}_fake"`,
			"snapshot_name":     `"${alibabacloudstack_ecs_snapshots.default.SnapshotName}_fake"`,
			"snapshot_type":     `"${alibabacloudstack_ecs_snapshots.default.SnapshotType}_fake"`,
			"source_disk_type":  `"${alibabacloudstack_ecs_snapshots.default.SourceDiskType}_fake"`,
			"status":            `"${alibabacloudstack_ecs_snapshots.default.Status}_fake"`,
			"usage":             `"${alibabacloudstack_ecs_snapshots.default.Usage}_fake"`}),
	}

	AlibabacloudstackEcsSnapshotsCheckInfo.dataSourceTestCheck(t, rand, idsConf, encryptedConf, resource_group_idConf, snapshot_idConf, snapshot_nameConf, snapshot_typeConf, source_disk_typeConf, statusConf, usageConf, allConf)
}

var existAlibabacloudstackEcsSnapshotsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"snapshots.#":    "1",
		"snapshots.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackEcsSnapshotsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"snapshots.#": "0",
	}
}

var AlibabacloudstackEcsSnapshotsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_ecs_snapshots.default",
	existMapFunc: existAlibabacloudstackEcsSnapshotsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackEcsSnapshotsMapFunc,
}

func testAccCheckAlibabacloudstackEcsSnapshotsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackEcsSnapshots%d"
}






data "alibabacloudstack_ecs_snapshots" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
