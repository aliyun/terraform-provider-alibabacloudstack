package apsarastack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"testing"
)

func TestAccApsaraStackSnapshotsDataSourceBasic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccSnapshotDataSourceBasic%d", rand)
	resourceId := "data.apsarastack_snapshots.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceSnapshotsConfigDependence)

	idsConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${apsarastack_snapshot.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${apsarastack_snapshot.default.id}_fake"},
		}),
	}

	instanceIdConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${apsarastack_instance.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${apsarastack_instance.default.id}_fake",
		}),
	}

	diskIdConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"disk_id": "${apsarastack_disk_attachment.default.disk_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"disk_id": "${apsarastack_disk_attachment.default.disk_id}_fake",
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
	//			"ids":    []string{"${apsarastack_snapshot.default.id}"},
	//			"status": "accomplished",
	//		}),
	//		fakeConfig: testAccConfig(map[string]interface{}{
	//			"ids":    []string{"${apsarastack_snapshot.default.id}"},
	//			"status": "failed",
	//		}),
	//	}

	typeConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"${apsarastack_snapshot.default.id}"},
			"type": "user",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"${apsarastack_snapshot.default.id}"},
			"type": "auto",
		}),
	}

	//	sourceDiskTypeConfig := dataSourceTestAccConfig{
	//		existConfig: testAccConfig(map[string]interface{}{
	//			"ids":              []string{"${apsarastack_snapshot.default.id}"},
	//			"source_disk_type": "Data",
	//		}),
	//		fakeConfig: testAccConfig(map[string]interface{}{
	//			"ids":              []string{"${apsarastack_snapshot.default.id}"},
	//			"source_disk_type": "System",
	//		}),
	//	}

	//	usageConfig := dataSourceTestAccConfig{
	//		existConfig: testAccConfig(map[string]interface{}{
	//			"ids":   []string{"${apsarastack_snapshot.default.id}"},
	//			"usage": "none",
	//		}),
	//		fakeConfig: testAccConfig(map[string]interface{}{
	//			"ids":   []string{"${apsarastack_snapshot.default.id}"},
	//			"usage": "image",
	//		}),
	//	}
	//
	//	tagsConfig := dataSourceTestAccConfig{
	//		existConfig: testAccConfig(map[string]interface{}{
	//			"ids": []string{"${apsarastack_snapshot.default.id}"},
	//			"tags": map[string]interface{}{
	//				"version": "1.0",
	//			},
	//		}),
	//		fakeConfig: testAccConfig(map[string]interface{}{
	//			"ids": []string{"${apsarastack_snapshot.default.id}"},
	//			"tags": map[string]interface{}{
	//				"version": "1.0_fake",
	//			},
	//		}),
	//	}

	//	allConfig := dataSourceTestAccConfig{
	//		existConfig: testAccConfig(map[string]interface{}{
	//			"ids":              []string{"${apsarastack_snapshot.default.id}"},
	//			"instance_id":      "${apsarastack_instance.default.id}",
	//			"disk_id":          "${apsarastack_disk_attachment.default.disk_id}",
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
	//			"ids":              []string{"${apsarastack_snapshot.default.id}"},
	//			"instance_id":      "${apsarastack_instance.default.id}",
	//			"disk_id":          "${apsarastack_disk_attachment.default.disk_id}",
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

	snapshotsCheckInfo.dataSourceTestCheck(t, rand, idsConfig, instanceIdConfig, diskIdConfig, nameRegexConfig,
		typeConfig) // statusConfig,sourceDiskTypeConfig,usageConfig, tagsConfig,allConfig
}

func dataSourceSnapshotsConfigDependence(name string) string {
	return fmt.Sprintf(`

%s

%s

%s

variable "name" {
  default = "%s"
}

resource "apsarastack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "192.168.0.0/16"
}
resource "apsarastack_vswitch" "default" {
  name = "${var.name}"
  cidr_block = "192.168.0.0/24"
  availability_zone = data.apsarastack_zones.default.zones[0].id
  vpc_id = "${apsarastack_vpc.default.id}"
}
resource "apsarastack_security_group" "default" {
  name        = "${var.name}"
  description = "${var.name}"
  vpc_id = "${apsarastack_vpc.default.id}"
}
resource "apsarastack_disk" "default" {
  availability_zone = data.apsarastack_zones.default.zones[0].id
  category          = "cloud_efficiency"
  size              = "20"
}

resource "apsarastack_instance" "default" {
  instance_name   = "${var.name}"
  image_id        = "${data.apsarastack_images.default.images.0.id}"
  instance_type   = local.instance_type_id
  security_groups = ["${apsarastack_security_group.default.id}"]
  vswitch_id      = "${apsarastack_vswitch.default.id}"
  availability_zone = data.apsarastack_zones.default.zones[0].id
}
resource "apsarastack_disk_attachment" "default" {
  disk_id     = "${apsarastack_disk.default.id}"
  instance_id = "${apsarastack_instance.default.id}"
}
resource "apsarastack_snapshot" "default" {
  disk_id = "${apsarastack_disk_attachment.default.disk_id}"
  name = "${var.name}"
  description = "${var.name}"
  tags = {
    version = "1.0"
  }
}
`, DataApsarastackVswitchZones, DataApsarastackInstanceTypes, DataApsarastackImages, name)
}
