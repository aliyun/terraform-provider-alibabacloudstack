package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"
)

func TestAccAlibabacloudStackEbsDiskReplicaPairsDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_ebs_diskreplicapairs.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAcc%sEbsDiskReplicaPairsDataSource-%d", defaultRegionToTest, rand),
		dataSourceEbsDiskReplicaPairsDependence)

	descriptionRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_ebs_diskreplicapair.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_ebs_diskreplicapair.default.description}-fakeTestAcccc",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ebs_diskreplicapair.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ebs_diskreplicapair.default.id}-fakeTestAcccc"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":        "${alibabacloudstack_ebs_diskreplicapair.default.description}",
			"ids":               []string{"${alibabacloudstack_ebs_diskreplicapair.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":        "${alibabacloudstack_ebs_diskreplicapair.default.description}-fakeTestAcccc",
			"ids":               []string{"${alibabacloudstack_ebs_diskreplicapair.default.id}-fakeTestAcccc"},
		}),
	}

	var existEbsDiskReplicaPairsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"ids.0":                            CHECKSET,
			"disk_replica_pairs.#":             "1",
			"disk_replica_pairs.0.description": fmt.Sprintf("tf-testAcc%sEbsDiskReplicaPairsDataSource-%d", defaultRegionToTest, rand),
			"disk_replica_pairs.0.destination_disk_id":   CHECKSET,
			"disk_replica_pairs.0.destination_region_id": CHECKSET,
			"disk_replica_pairs.0.destination_zone_id":   CHECKSET,
			"disk_replica_pairs.0.id":                    CHECKSET,
			"disk_replica_pairs.0.source_disk_id":        CHECKSET,
			"disk_replica_pairs.0.source_region_id":      CHECKSET,
			"disk_replica_pairs.0.source_zone_id":        CHECKSET,
		}
	}

	var fakeEbsDiskReplicaPairsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                "0",
			"disk_replica_pairs.#": "0",
		}
	}

	var EbsDiskReplicaPairsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEbsDiskReplicaPairsMapFunc,
		fakeMapFunc:  fakeEbsDiskReplicaPairsMapFunc,
	}

	EbsDiskReplicaPairsCheckInfo.dataSourceTestCheck(t, rand, descriptionRegexConf, idsConf, allConf)
}

func dataSourceEbsDiskReplicaPairsDependence(name string) string {
	region := os.Getenv("ALIBABACLOUDSTACK_REGION")
	return fmt.Sprintf(`

variable "name" {
	default = "%s"
}

variable "region" {
  default = "%s"
}

// resource "alibabacloudstack_ecs_disk" "disk1" {
// 	availability_zone = "${data.alibabacloudstack_zones.default.zones[0].id}"
// 	size = "20"
// 	name = "${var.name}"
// 	category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
// }

// resource "alibabacloudstack_ecs_disk" "disk2" {
// 	availability_zone = "${data.alibabacloudstack_zones.default.zones[1].id}"
// 	size = "20"
// 	name = "${var.name}"
// 	category = "${data.alibabacloudstack_zones.default.zones.1.available_disk_categories.0}"
// }
	
resource "alibabacloudstack_ebs_diskreplicapair" "default" {
	disk_replica_pair_name = "${var.name}"
	description =            "${var.name}"
	source_zone_id =         "cn-wulan-env82-amtest83002-b"
	source_region_id =       "${var.region}"
	source_disk_id =         "d-9rt00vs0qsrd2502tx5u"
	destination_zone_id =    "cn-wulan-env82-amtest82001-a"
	destination_region_id =  "${var.region}"
	destination_disk_id =    "d-9rt00vs0qsrd2502tx5p"
	rpo =                    300
}
 `, name, region)
}
