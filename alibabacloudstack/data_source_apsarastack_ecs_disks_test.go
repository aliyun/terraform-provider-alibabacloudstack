package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackDisksDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_ecs_disks.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccCheckAlibabacloudStackDisksDataSource_ids-%d", rand),
		dataSourceEcsDisksConfigDependence)

	existDisksMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"disks.0.id":                CHECKSET,
			"disks.0.name":              fmt.Sprintf("tf-testAccCheckAlibabacloudStackDisksDataSource_ids-%d", rand),
			"disks.0.description":       fmt.Sprintf("tf-testAccCheckAlibabacloudStackDisksDataSource_ids-%d_description", rand),
			"disks.0.region_id":         CHECKSET,
			"disks.0.availability_zone": CHECKSET,
			"disks.0.status":            "Available",
			"disks.0.type":              "data",
			"disks.0.category":          CHECKSET,
			"disks.0.size":              "20",
			"disks.0.image_id":          "",
			"disks.0.snapshot_id":       "",
			"disks.0.instance_id":       "",
			"disks.0.creation_time":     CHECKSET,
			"disks.0.attached_time":     "",
			"disks.0.detached_time":     "",
			"disks.0.tags.%":            "2",
		}
	}
	fakeDisksMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"disks.#": "0",
		}
	}

	checkInfo := dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDisksMapFunc,
		fakeMapFunc:  fakeDisksMapFunc,
	}

	checkInfo.dataSourceTestCheck(t, rand,
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"${alibabacloudstack_ecs_disk.default.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"${alibabacloudstack_ecs_disk.default.id}_fake"},
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_ecs_disk.default.name}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_ecs_disk.default.name}_fake",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_ecs_disk.default.name}",
				"type":       "data",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_ecs_disk.default.name}",
				"type":       "system",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"instance_id": "${alibabacloudstack_ecs_diskattachment.default.instance_id}",
				"name_regex":  "${alibabacloudstack_ecs_disk.default.name}",
			}),
			existChangMap: map[string]string{
				"disks.0.instance_id":   CHECKSET,
				"disks.0.attached_time": CHECKSET,
				"disks.0.status":        "In_use",
			},
			fakeConfig: testAccConfig(map[string]interface{}{
				"instance_id": "${alibabacloudstack_ecs_diskattachment.default.instance_id}_fake",
				"name_regex":  "${alibabacloudstack_ecs_disk.default.name}",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids":         []string{"${alibabacloudstack_ecs_disk.default.id}"},
				"name_regex":  "${alibabacloudstack_ecs_disk.default.name}",
				"type":        "data",
				"category":    "cloud_efficiency",
				"tags":        map[string]interface{}{"Name": "TerraformTest", "Name1": "TerraformTest"},
				"instance_id": "${alibabacloudstack_ecs_diskattachment.default.instance_id}",
			}),
			existChangMap: map[string]string{
				"disks.0.instance_id":   CHECKSET,
				"disks.0.attached_time": CHECKSET,
				"disks.0.status":        "In_use",
			},
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids":         []string{"${alibabacloudstack_ecs_disk.default.id}"},
				"name_regex":  "${alibabacloudstack_ecs_disk.default.name}",
				"type":        "data",
				"category":    "cloud_efficiency",
				"tags":        map[string]interface{}{"Name": "TerraformTest", "Name1": "TerraformTest"},
				"instance_id": "${alibabacloudstack_ecs_diskattachment.default.instance_id}_fake",
			}),
		},
	)
}

func dataSourceEcsDisksConfigDependence(name string) string {
	return fmt.Sprintf(`
%s

variable "name" {
	default = "%s"
}

resource "alibabacloudstack_ecs_disk" "default" {
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
	name = "${var.name}"
	description = "${var.name}_description"
	tags = {
		Name = "TerraformTest"
		Name1 = "TerraformTest"
	}
	size = "20"
}

resource "alibabacloudstack_ecs_diskattachment" "default" {
	disk_id = "${alibabacloudstack_ecs_disk.default.id}"
	instance_id = "${alibabacloudstack_ecs_instance.default.id}"
}
`, ECSInstanceCommonTestCase, name)
}
