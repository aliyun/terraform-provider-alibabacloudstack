package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackEcsDisksDataSource(t *testing.T) {
	// 根据test_meta自动生成的tasecase

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_disks.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_disks.default.id}_fake"]`,
		}),
	}

	auto_snapshot_policy_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids":                     `["${alibabacloudstack_ecs_disks.default.id}"]`,
			"auto_snapshot_policy_id": `"${alibabacloudstack_ecs_disks.default.AutoSnapshotPolicyId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids":                     `["${alibabacloudstack_ecs_disks.default.id}_fake"]`,
			"auto_snapshot_policy_id": `"${alibabacloudstack_ecs_disks.default.AutoSnapshotPolicyId}_fake"`,
		}),
	}

	categoryConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids":      `["${alibabacloudstack_ecs_disks.default.id}"]`,
			"category": `"${alibabacloudstack_ecs_disks.default.Category}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids":      `["${alibabacloudstack_ecs_disks.default.id}_fake"]`,
			"category": `"${alibabacloudstack_ecs_disks.default.Category}_fake"`,
		}),
	}

	delete_auto_snapshotConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_disks.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_disks.default.id}_fake"]`,
		}),
	}

	delete_with_instanceConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_disks.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_disks.default.id}_fake"]`,
		}),
	}

	disk_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids":       `["${alibabacloudstack_ecs_disks.default.id}"]`,
			"disk_name": `"${alibabacloudstack_ecs_disks.default.DiskName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids":       `["${alibabacloudstack_ecs_disks.default.id}_fake"]`,
			"disk_name": `"${alibabacloudstack_ecs_disks.default.DiskName}_fake"`,
		}),
	}

	enable_auto_snapshotConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_disks.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_disks.default.id}_fake"]`,
		}),
	}

	enable_automated_snapshot_policyConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_disks.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_disks.default.id}_fake"]`,
		}),
	}

	encryptedConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_disks.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_disks.default.id}_fake"]`,
		}),
	}

	instance_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_ecs_disks.default.id}"]`,
			"instance_id": `"${alibabacloudstack_ecs_disks.default.InstanceId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_ecs_disks.default.id}_fake"]`,
			"instance_id": `"${alibabacloudstack_ecs_disks.default.InstanceId}_fake"`,
		}),
	}

	kms_key_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_ecs_disks.default.id}"]`,
			"kms_key_id": `"${alibabacloudstack_ecs_disks.default.KmsKeyId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_ecs_disks.default.id}_fake"]`,
			"kms_key_id": `"${alibabacloudstack_ecs_disks.default.KmsKeyId}_fake"`,
		}),
	}

	multi_attachConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_ecs_disks.default.id}"]`,
			"multi_attach": `"${alibabacloudstack_ecs_disks.default.MultiAttach}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_ecs_disks.default.id}_fake"]`,
			"multi_attach": `"${alibabacloudstack_ecs_disks.default.MultiAttach}_fake"`,
		}),
	}

	payment_typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_ecs_disks.default.id}"]`,
			"payment_type": `"${alibabacloudstack_ecs_disks.default.PaymentType}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_ecs_disks.default.id}_fake"]`,
			"payment_type": `"${alibabacloudstack_ecs_disks.default.PaymentType}_fake"`,
		}),
	}

	portableConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_disks.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_disks.default.id}_fake"]`,
		}),
	}

	resource_group_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_ecs_disks.default.id}"]`,
			"resource_group_id": `"${alibabacloudstack_ecs_disks.default.ResourceGroupId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_ecs_disks.default.id}_fake"]`,
			"resource_group_id": `"${alibabacloudstack_ecs_disks.default.ResourceGroupId}_fake"`,
		}),
	}

	snapshot_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_ecs_disks.default.id}"]`,
			"snapshot_id": `"${alibabacloudstack_ecs_disks.default.SnapshotId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_ecs_disks.default.id}_fake"]`,
			"snapshot_id": `"${alibabacloudstack_ecs_disks.default.SnapshotId}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_ecs_disks.default.id}"]`,
			"status": `"${alibabacloudstack_ecs_disks.default.Status}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_ecs_disks.default.id}_fake"]`,
			"status": `"${alibabacloudstack_ecs_disks.default.Status}_fake"`,
		}),
	}

	zone_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids":     `["${alibabacloudstack_ecs_disks.default.id}"]`,
			"zone_id": `"${alibabacloudstack_ecs_disks.default.ZoneId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids":     `["${alibabacloudstack_ecs_disks.default.id}_fake"]`,
			"zone_id": `"${alibabacloudstack_ecs_disks.default.ZoneId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_disks.default.id}"]`,

			"auto_snapshot_policy_id": `"${alibabacloudstack_ecs_disks.default.AutoSnapshotPolicyId}"`,
			"category":                `"${alibabacloudstack_ecs_disks.default.Category}"`,
			"disk_name":               `"${alibabacloudstack_ecs_disks.default.DiskName}"`,
			"instance_id":             `"${alibabacloudstack_ecs_disks.default.InstanceId}"`,
			"kms_key_id":              `"${alibabacloudstack_ecs_disks.default.KmsKeyId}"`,
			"multi_attach":            `"${alibabacloudstack_ecs_disks.default.MultiAttach}"`,
			"payment_type":            `"${alibabacloudstack_ecs_disks.default.PaymentType}"`,
			"resource_group_id":       `"${alibabacloudstack_ecs_disks.default.ResourceGroupId}"`,
			"snapshot_id":             `"${alibabacloudstack_ecs_disks.default.SnapshotId}"`,
			"status":                  `"${alibabacloudstack_ecs_disks.default.Status}"`,
			"zone_id":                 `"${alibabacloudstack_ecs_disks.default.ZoneId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_disks.default.id}_fake"]`,

			"auto_snapshot_policy_id": `"${alibabacloudstack_ecs_disks.default.AutoSnapshotPolicyId}_fake"`,
			"category":                `"${alibabacloudstack_ecs_disks.default.Category}_fake"`,
			"disk_name":               `"${alibabacloudstack_ecs_disks.default.DiskName}_fake"`,
			"instance_id":             `"${alibabacloudstack_ecs_disks.default.InstanceId}_fake"`,
			"kms_key_id":              `"${alibabacloudstack_ecs_disks.default.KmsKeyId}_fake"`,
			"multi_attach":            `"${alibabacloudstack_ecs_disks.default.MultiAttach}_fake"`,
			"payment_type":            `"${alibabacloudstack_ecs_disks.default.PaymentType}_fake"`,
			"resource_group_id":       `"${alibabacloudstack_ecs_disks.default.ResourceGroupId}_fake"`,
			"snapshot_id":             `"${alibabacloudstack_ecs_disks.default.SnapshotId}_fake"`,
			"status":                  `"${alibabacloudstack_ecs_disks.default.Status}_fake"`,
			"zone_id":                 `"${alibabacloudstack_ecs_disks.default.ZoneId}_fake"`}),
	}

	AlibabacloudstackEcsDisksDataCheckInfo.dataSourceTestCheck(t, rand, idsConf, auto_snapshot_policy_idConf, categoryConf, delete_auto_snapshotConf, delete_with_instanceConf, disk_nameConf, enable_auto_snapshotConf, enable_automated_snapshot_policyConf, encryptedConf, instance_idConf, kms_key_idConf, multi_attachConf, payment_typeConf, portableConf, resource_group_idConf, snapshot_idConf, statusConf, zone_idConf, allConf)
}

var existAlibabacloudstackEcsDisksDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"disks.#":    "1",
		"disks.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackEcsDisksDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"disks.#": "0",
	}
}

var AlibabacloudstackEcsDisksDataCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_ecs_disks.default",
	existMapFunc: existAlibabacloudstackEcsDisksDataMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackEcsDisksDataMapFunc,
}

func testAccCheckAlibabacloudstackEcsDisksDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackEcsDisks%d"
}

data "alibabacloudstack_ecs_disks" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

