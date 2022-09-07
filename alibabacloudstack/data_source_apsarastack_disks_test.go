package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"strings"
	"testing"

	"fmt"
)

func TestAccAlibabacloudStackDisksDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)

	idsConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackDisksDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alibabacloudstack_disk.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackDisksDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alibabacloudstack_disk.default.id}_fake" ]`,
		}),
	}

	nameRegexConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackDisksDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_disk.default.name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackDisksDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_disk.default.name}_fake"`,
		}),
	}

	typeConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackDisksDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_disk.default.name}"`,
			"type":       `"data"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackDisksDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_disk.default.name}"`,
			"type":       `"system"`,
		}),
	}

	//	categoryConfig := dataSourceTestAccConfig{
	//		existConfig: testAccCheckAlibabacloudStackDisksDataSourceConfig(rand, map[string]string{
	//			"name_regex": `"${alibabacloudstack_disk.default.name}"`,
	//			"category":   `"cloud_efficiency"`,
	//		}),
	//		fakeConfig: testAccCheckAlibabacloudStackDisksDataSourceConfig(rand, map[string]string{
	//			"name_regex": `"${alibabacloudstack_disk.default.name}"`,
	//			"category":   `"cloud"`,
	//		}),
	//	}

	//	tagsConfig := dataSourceTestAccConfig{
	//		existConfig: testAccCheckAlibabacloudStackDisksDataSourceConfig(rand, map[string]string{
	//			"name_regex": `"${alibabacloudstack_disk.default.name}"`,
	//			"tags":       `"${alibabacloudstack_disk.default.tags}"`,
	//		}),
	//		fakeConfig: testAccCheckAlibabacloudStackDisksDataSourceConfig(rand, map[string]string{
	//			"name_regex": `"${alibabacloudstack_disk.default.name}"`,
	//			"tags": `{
	//                           Name = "TerraformTest_fake"
	//                           Name1 = "TerraformTest_fake"
	//                        }`,
	//		}),
	//	}

	instanceIdConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackDisksDataSourceConfigWithCommon(rand, map[string]string{
			"instance_id": `"${alibabacloudstack_disk_attachment.default.instance_id}"`,
			"name_regex":  `"${alibabacloudstack_disk.default.name}"`,
		}),
		existChangMap: map[string]string{
			"disks.0.instance_id":   CHECKSET,
			"disks.0.attached_time": CHECKSET,
			"disks.0.status":        "In_use",
		},
		fakeConfig: testAccCheckAlibabacloudStackDisksDataSourceConfigWithCommon(rand, map[string]string{
			"instance_id": `"${alibabacloudstack_disk_attachment.default.instance_id}_fake"`,
			"name_regex":  `"${alibabacloudstack_disk.default.name}"`,
		}),
	}

	allConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackDisksDataSourceConfigWithCommon(rand, map[string]string{
			"ids":         `[ "${alibabacloudstack_disk.default.id}" ]`,
			"name_regex":  `"${alibabacloudstack_disk.default.name}"`,
			"type":        `"data"`,
			"category":    `"cloud_efficiency"`,
			"tags":        `"${alibabacloudstack_disk.default.tags}"`,
			"instance_id": `"${alibabacloudstack_disk_attachment.default.instance_id}"`,
		}),
		existChangMap: map[string]string{
			"disks.0.instance_id":   CHECKSET,
			"disks.0.attached_time": CHECKSET,
			"disks.0.status":        "In_use",
		},
		fakeConfig: testAccCheckAlibabacloudStackDisksDataSourceConfigWithCommon(rand, map[string]string{
			"ids":         `[ "${alibabacloudstack_disk.default.id}" ]`,
			"name_regex":  `"${alibabacloudstack_disk.default.name}"`,
			"type":        `"data"`,
			"category":    `"cloud_efficiency"`,
			"tags":        `"${alibabacloudstack_disk.default.tags}"`,
			"instance_id": `"${alibabacloudstack_disk_attachment.default.instance_id}_fake"`,
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

variable "name" {
	default = "tf-testAccCheckAlibabacloudStackDisksDataSource_ids-%d"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alibabacloudstack_disk" "default" {
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
	description = "${var.name}_description"
	size = "20"
	tags = {
		Name = "TerraformTest"
		Name1 = "TerraformTest"
	}
}

data "alibabacloudstack_disks" "default" {
	%s
}
	`, rand, strings.Join(pairs, "\n	"))
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
resource "alibabacloudstack_disk" "default" {
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	category = "cloud_efficiency"
	name = "${var.name}"
	description = "${var.name}_description"
	tags = {
		Name = "TerraformTest"
		Name1 = "TerraformTest"
	}
	size = "20"
}

resource "alibabacloudstack_instance" "default" {
	vswitch_id = "${alibabacloudstack_vswitch.default.id}"
	private_ip = "172.16.0.10"
	image_id = "${data.alibabacloudstack_images.default.images.0.id}"
	instance_type = "${local.instance_type_id}"
	instance_name = "${var.name}"
	system_disk_category = "cloud_efficiency"
	security_groups = ["${alibabacloudstack_security_group.default.id}"]
}

resource "alibabacloudstack_disk_attachment" "default" {
	disk_id = "${alibabacloudstack_disk.default.id}"
	instance_id = "${alibabacloudstack_instance.default.id}"
}

data "alibabacloudstack_disks" "default" {
	%s
}
`, EcsInstanceCommonTestCase, rand, strings.Join(pairs, "\n	"))
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
	resourceId:   "data.alibabacloudstack_disks.default",
	existMapFunc: existDisksMapFunc,
	fakeMapFunc:  fakeDisksMapFunc,
}
