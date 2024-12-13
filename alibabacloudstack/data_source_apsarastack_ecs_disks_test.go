package alibabacloudstack

import (
	"strings"
	"testing"
	"fmt"
)

func TestAccAlibabacloudStackDisksDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)

	idsConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackDisksDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alibabacloudstack_ecs_disk.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackDisksDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alibabacloudstack_ecs_disk.default.id}_fake" ]`,
		}),
	}

	nameRegexConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackDisksDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_ecs_disk.default.name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackDisksDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_ecs_disk.default.name}_fake"`,
		}),
	}

	typeConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackDisksDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_ecs_disk.default.name}"`,
			"type":       `"data"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackDisksDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_ecs_disk.default.name}"`,
			"type":       `"system"`,
		}),
	}

	//	categoryConfig := dataSourceTestAccConfig{
	//		existConfig: testAccCheckAlibabacloudStackDisksDataSourceConfig(rand, map[string]string{
	//			"name_regex": `"${alibabacloudstack_ecs_disk.default.name}"`,
	//			"category":   `"cloud_efficiency"`,
	//		}),
	//		fakeConfig: testAccCheckAlibabacloudStackDisksDataSourceConfig(rand, map[string]string{
	//			"name_regex": `"${alibabacloudstack_ecs_disk.default.name}"`,
	//			"category":   `"cloud"`,
	//		}),
	//	}

	//	tagsConfig := dataSourceTestAccConfig{
	//		existConfig: testAccCheckAlibabacloudStackDisksDataSourceConfig(rand, map[string]string{
	//			"name_regex": `"${alibabacloudstack_ecs_disk.default.name}"`,
	//			"tags":       `"${alibabacloudstack_ecs_disk.default.tags}"`,
	//		}),
	//		fakeConfig: testAccCheckAlibabacloudStackDisksDataSourceConfig(rand, map[string]string{
	//			"name_regex": `"${alibabacloudstack_ecs_disk.default.name}"`,
	//			"tags": `{
	//                           Name = "TerraformTest_fake"
	//                           Name1 = "TerraformTest_fake"
	//                        }`,
	//		}),
	//	}

	instanceIdConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackDisksDataSourceConfigWithCommon(rand, map[string]string{
			"instance_id": `"${alibabacloudstack_ecs_diskattachment.default.instance_id}"`,
			"name_regex":  `"${alibabacloudstack_ecs_disk.default.name}"`,
		}),
		existChangMap: map[string]string{
			"disks.0.instance_id":   CHECKSET,
			"disks.0.attached_time": CHECKSET,
			"disks.0.status":        "In_use",
		},
		fakeConfig: testAccCheckAlibabacloudStackDisksDataSourceConfigWithCommon(rand, map[string]string{
			"instance_id": `"${alibabacloudstack_ecs_diskattachment.default.instance_id}_fake"`,
			"name_regex":  `"${alibabacloudstack_ecs_disk.default.name}"`,
		}),
	}

	allConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackDisksDataSourceConfigWithCommon(rand, map[string]string{
			"ids":         `[ "${alibabacloudstack_ecs_disk.default.id}" ]`,
			"name_regex":  `"${alibabacloudstack_ecs_disk.default.name}"`,
			"type":        `"data"`,
			"category":    `"cloud_efficiency"`,
			"tags":        `"${alibabacloudstack_ecs_disk.default.tags}"`,
			"instance_id": `"${alibabacloudstack_ecs_diskattachment.default.instance_id}"`,
		}),
		existChangMap: map[string]string{
			"disks.0.instance_id":   CHECKSET,
			"disks.0.attached_time": CHECKSET,
			"disks.0.status":        "In_use",
		},
		fakeConfig: testAccCheckAlibabacloudStackDisksDataSourceConfigWithCommon(rand, map[string]string{
			"ids":         `[ "${alibabacloudstack_ecs_disk.default.id}" ]`,
			"name_regex":  `"${alibabacloudstack_ecs_disk.default.name}"`,
			"type":        `"data"`,
			"category":    `"cloud_efficiency"`,
			"tags":        `"${alibabacloudstack_ecs_disk.default.tags}"`,
			"instance_id": `"${alibabacloudstack_ecs_diskattachment.default.instance_id}_fake"`,
		}),
	}

	disksCheckInfo.dataSourceTestCheck(t, rand, idsConfig, nameRegexConfig, typeConfig,
		instanceIdConfig, allConfig)
}

func testAccCheckAlibabacloudStackDisksDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
%s

variable "name" {
	default = "tf-testAccCheckAlibabacloudStackDisksDataSource_ids-%d"
}

resource "alibabacloudstack_ecs_disk" "default" {
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
	name = "${var.name}"
	description = "${var.name}_description"
	size = "20"
	tags = {
		Name = "TerraformTest"
		Name1 = "TerraformTest"
	}
}

data "alibabacloudstack_ecs_disks" "default" {
	%s
}
	`, DataAlibabacloudstackVswitchZones, rand, strings.Join(pairs, "\n	"))
	return config
}

func testAccCheckAlibabacloudStackDisksDataSourceConfigWithCommon(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
%s

variable "name" {
	default = "tf-testAccCheckAlibabacloudStackDisksDataSource_ids-%d"
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

data "alibabacloudstack_ecs_disks" "default" {
	%s
}
`, ECSInstanceCommonTestCase, rand, strings.Join(pairs, "\n	"))
	return config
}

var existDisksMapFunc = func(rand int) map[string]string {
	return map[string]string{
		//"disks.#":                   "1",
		"disks.0.id":                CHECKSET,
		"disks.0.name":              fmt.Sprintf("tf-testAccCheckAlibabacloudStackDisksDataSource_ids-%d", rand),
		"disks.0.description":       fmt.Sprintf("tf-testAccCheckAlibabacloudStackDisksDataSource_ids-%d_description", rand),
		"disks.0.region_id":         CHECKSET,
		"disks.0.availability_zone": CHECKSET,
		"disks.0.status":            "Available",
		"disks.0.type":              "data",
		"disks.0.category":          "cloud_efficiency",
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

var fakeDisksMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"disks.#": "0",
	}
}

var disksCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_ecs_disks.default",
	existMapFunc: existDisksMapFunc,
	fakeMapFunc:  fakeDisksMapFunc,
}
