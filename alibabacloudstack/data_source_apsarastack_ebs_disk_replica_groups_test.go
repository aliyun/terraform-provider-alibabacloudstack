package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"
)

func TestAccAlibabacloudStackEbsDiskReplicaGroupsDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_ebs_diskreplicagroups.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAcc%sEbsDiskReplicaGroupsDataSource-%d", defaultRegionToTest, rand),
		dataSourceEbsDiskReplicaGroupsDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_ebs_diskreplicagroup.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_ebs_diskreplicagroup.default.description}-fakeTestAcccc",
		}),
	}

	descriptionRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_ebs_diskreplicagroup.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_ebs_diskreplicagroup.default.description}-fakeTestAcccc",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ebs_diskreplicagroup.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ebs_diskreplicagroup.default.id}-fakeTestAcccc"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":        "${alibabacloudstack_ebs_diskreplicagroup.default.description}",
			"description_regex": "${alibabacloudstack_ebs_diskreplicagroup.default.description}",
			"ids":               []string{"${alibabacloudstack_ebs_diskreplicagroup.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":        "${alibabacloudstack_ebs_diskreplicagroup.default.description}-fakeTestAcccc",
			"description_regex": "${alibabacloudstack_ebs_diskreplicagroup.default.description}-fakeTestAcccc",
			"ids":               []string{"${alibabacloudstack_ebs_diskreplicagroup.default.id}-fakeTestAcccc"},
		}),
	}

	var existEbsDiskReplicaGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"ids.0":                             CHECKSET,
			"disk_replica_groups.#":             "1",
			"disk_replica_groups.0.description": fmt.Sprintf("tf-testAcc%sEbsDiskReplicaGroupsDataSource-%d", defaultRegionToTest, rand),
			"disk_replica_groups.0.destination_region_id": CHECKSET,
			"disk_replica_groups.0.destination_zone_id":   CHECKSET,
			"disk_replica_groups.0.id":                    CHECKSET,
			"disk_replica_groups.0.source_region_id":      CHECKSET,
			"disk_replica_groups.0.source_zone_id":        CHECKSET,
		}
	}

	var fakeEbsDiskReplicaGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                 "0",
			"disk_replica_groups.#": "0",
		}
	}

	var EbsDiskReplicaGroupsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEbsDiskReplicaGroupsMapFunc,
		fakeMapFunc:  fakeEbsDiskReplicaGroupsMapFunc,
	}

	EbsDiskReplicaGroupsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, descriptionRegexConf, idsConf, allConf)
}

func dataSourceEbsDiskReplicaGroupsDependence(name string) string {
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
	
resource "alibabacloudstack_ebs_diskreplicagroup" "default" {
    disk_replica_group_name = "${var.name}"
    description = "${var.name}"
    destination_region_id = "cn-wulan-env82-d01"
    destination_zone_id ="cn-wulan-env82-amtest83002-b"
    site = "production"
    source_region_id = "cn-wulan-env82-d01"
    source_zone_id = "cn-wulan-env82-amtest82001-a"
}
 `, name, region)
}
