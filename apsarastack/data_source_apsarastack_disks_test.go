package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"strings"
	"testing"

	"fmt"
)

func TestAccApsaraStackDisksDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)

	idsConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackDisksDataSourceConfig(rand, map[string]string{
			"ids": `[ "${apsarastack_disk.default.id}" ]`,
		}),
		fakeConfig: testAccCheckApsaraStackDisksDataSourceConfig(rand, map[string]string{
			"ids": `[ "${apsarastack_disk.default.id}_fake" ]`,
		}),
	}

	nameRegexConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackDisksDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_disk.default.name}"`,
		}),
		fakeConfig: testAccCheckApsaraStackDisksDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_disk.default.name}_fake"`,
		}),
	}

	typeConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackDisksDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_disk.default.name}"`,
			"type":       `"data"`,
		}),
		fakeConfig: testAccCheckApsaraStackDisksDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_disk.default.name}"`,
			"type":       `"system"`,
		}),
	}

	//	categoryConfig := dataSourceTestAccConfig{
	//		existConfig: testAccCheckApsaraStackDisksDataSourceConfig(rand, map[string]string{
	//			"name_regex": `"${apsarastack_disk.default.name}"`,
	//			"category":   `"cloud_efficiency"`,
	//		}),
	//		fakeConfig: testAccCheckApsaraStackDisksDataSourceConfig(rand, map[string]string{
	//			"name_regex": `"${apsarastack_disk.default.name}"`,
	//			"category":   `"cloud"`,
	//		}),
	//	}

	//	tagsConfig := dataSourceTestAccConfig{
	//		existConfig: testAccCheckApsaraStackDisksDataSourceConfig(rand, map[string]string{
	//			"name_regex": `"${apsarastack_disk.default.name}"`,
	//			"tags":       `"${apsarastack_disk.default.tags}"`,
	//		}),
	//		fakeConfig: testAccCheckApsaraStackDisksDataSourceConfig(rand, map[string]string{
	//			"name_regex": `"${apsarastack_disk.default.name}"`,
	//			"tags": `{
	//                           Name = "TerraformTest_fake"
	//                           Name1 = "TerraformTest_fake"
	//                        }`,
	//		}),
	//	}

	instanceIdConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackDisksDataSourceConfigWithCommon(rand, map[string]string{
			"instance_id": `"${apsarastack_disk_attachment.default.instance_id}"`,
			"name_regex":  `"${apsarastack_disk.default.name}"`,
		}),
		existChangMap: map[string]string{
			"disks.0.instance_id":   CHECKSET,
			"disks.0.attached_time": CHECKSET,
			"disks.0.status":        "In_use",
		},
		fakeConfig: testAccCheckApsaraStackDisksDataSourceConfigWithCommon(rand, map[string]string{
			"instance_id": `"${apsarastack_disk_attachment.default.instance_id}_fake"`,
			"name_regex":  `"${apsarastack_disk.default.name}"`,
		}),
	}

	allConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackDisksDataSourceConfigWithCommon(rand, map[string]string{
			"ids":         `[ "${apsarastack_disk.default.id}" ]`,
			"name_regex":  `"${apsarastack_disk.default.name}"`,
			"type":        `"data"`,
			"category":    `"cloud_efficiency"`,
			"tags":        `"${apsarastack_disk.default.tags}"`,
			"instance_id": `"${apsarastack_disk_attachment.default.instance_id}"`,
		}),
		existChangMap: map[string]string{
			"disks.0.instance_id":   CHECKSET,
			"disks.0.attached_time": CHECKSET,
			"disks.0.status":        "In_use",
		},
		fakeConfig: testAccCheckApsaraStackDisksDataSourceConfigWithCommon(rand, map[string]string{
			"ids":         `[ "${apsarastack_disk.default.id}" ]`,
			"name_regex":  `"${apsarastack_disk.default.name}"`,
			"type":        `"data"`,
			"category":    `"cloud_efficiency"`,
			"tags":        `"${apsarastack_disk.default.tags}"`,
			"instance_id": `"${apsarastack_disk_attachment.default.instance_id}_fake"`,
		}),
	}

	disksCheckInfo.dataSourceTestCheck(t, rand, idsConfig, nameRegexConfig, typeConfig,
		instanceIdConfig, allConfig)
}

func testAccCheckApsaraStackDisksDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testAccCheckApsaraStackDisksDataSource_ids-%d"
}

data "apsarastack_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "apsarastack_disk" "default" {
	availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
	description = "${var.name}_description"
	size = "20"
	tags = {
		Name = "TerraformTest"
		Name1 = "TerraformTest"
	}
}

data "apsarastack_disks" "default" {
	%s
}
	`, rand, strings.Join(pairs, "\n	"))
	return config
}

func testAccCheckApsaraStackDisksDataSourceConfigWithCommon(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
%s

variable "name" {
	default = "tf-testAccCheckApsaraStackDisksDataSource_ids-%d"
}
resource "apsarastack_disk" "default" {
	availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
	description = "${var.name}_description"
	tags = {
		Name = "TerraformTest"
		Name1 = "TerraformTest"
	}
	size = "20"
}

resource "apsarastack_instance" "default" {
	vswitch_id = "${apsarastack_vswitch.default.id}"
	private_ip = "172.16.0.10"
	image_id = "${data.apsarastack_images.default.images.0.id}"
	instance_type = "${local.instance_type_id}"
	instance_name = "${var.name}"
	system_disk_category = "cloud_efficiency"
	security_groups = ["${apsarastack_security_group.default.id}"]
}

resource "apsarastack_disk_attachment" "default" {
	disk_id = "${apsarastack_disk.default.id}"
	instance_id = "${apsarastack_instance.default.id}"
}

data "apsarastack_disks" "default" {
	%s
}
`, EcsInstanceCommonTestCase, rand, strings.Join(pairs, "\n	"))
	return config
}

var existDisksMapFunc = func(rand int) map[string]string {
	return map[string]string{
		//"disks.#":                   "1",
		"disks.0.id":                CHECKSET,
		"disks.0.name":              fmt.Sprintf("tf-testAccCheckApsaraStackDisksDataSource_ids-%d", rand),
		"disks.0.description":       fmt.Sprintf("tf-testAccCheckApsaraStackDisksDataSource_ids-%d_description", rand),
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
	resourceId:   "data.apsarastack_disks.default",
	existMapFunc: existDisksMapFunc,
	fakeMapFunc:  fakeDisksMapFunc,
}
