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
			"name_regex": "${alibabacloudstack_ebs_diskreplicagroup.default.description}",
			"ids":        []string{"${alibabacloudstack_ebs_diskreplicagroup.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_ebs_diskreplicagroup.default.description}-fakeTestAcccc",
			"ids":        []string{"${alibabacloudstack_ebs_diskreplicagroup.default.id}-fakeTestAcccc"},
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

	EbsDiskReplicaGroupsCheckInfo.dataSourceTestCheck(t, rand, descriptionRegexConf, idsConf, allConf)
}

func dataSourceEbsDiskReplicaGroupsDependence(name string) string {
	region := os.Getenv("ALIBABACLOUDSTACK_REGION")
	return fmt.Sprintf(`

variable "name" {
	default = "%s"
}

variable "region_id" {
  default = "%s"
}

%s

resource "alibabacloudstack_ebs_diskreplicagroup" "default" {
    disk_replica_group_name = "${var.name}"
    description = "${var.name}"
    destination_region_id = "${var.region_id}"
    destination_zone_id ="${data.alibabacloudstack_zones.default.zones[1].id}"
    site = "production"
    source_region_id = "${var.region_id}"
    source_zone_id = "${data.alibabacloudstack_zones.default.zones[0].id}"
}
 `, name, region, DataAlibabacloudstackVswitchZones)
}
