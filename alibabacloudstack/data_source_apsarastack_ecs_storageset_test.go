package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackEcsEbsStorageSetsDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000, 9999)
	resourceId := "data.alibabacloudstack_ecs_ebs_storage_sets.default"
	name := fmt.Sprintf("tf-testAcc_storage_set%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEcsEbsStorageSetsConfigDependence)

	nameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"storage_set_name": "${alibabacloudstack_ecs_ebs_storage_set.default.storage_set_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"storage_set_name": "fake-nonexistent-storage-set",
		}),
	}

	zoneIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"storage_set_id": "${alibabacloudstack_ecs_ebs_storage_set.default.id}",
			"zone_id": "${alibabacloudstack_ecs_ebs_storage_set.default.zone_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"storage_set_id": "${alibabacloudstack_ecs_ebs_storage_set.default.id}",
			"zone_id": "fake-zone-id",
		}),
	}

	storageSetIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"storage_set_id": "${alibabacloudstack_ecs_ebs_storage_set.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"storage_set_id": "fake-storage-set-id",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"storage_set_name": "${alibabacloudstack_ecs_ebs_storage_set.default.storage_set_name}",
			"zone_id":          "${alibabacloudstack_ecs_ebs_storage_set.default.zone_id}",
			"storage_set_id":   "${alibabacloudstack_ecs_ebs_storage_set.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"storage_set_name": "fake-storage-set-name",
			"zone_id":          "fake-zone-id",
			"storage_set_id":   "fake-storage-set-id",
		}),
	}

	var existEcsEbsStorageSetsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "1",
			"names.#":    "1",
			"storages.#": "1",
			"storages.0.storage_set_name": name,
			// Other storage attributes are computed but we don't know exact values
			// so we only validate the fields that are guaranteed to be present
		}
	}

	var fakeEcsEbsStorageSetsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"names.#":   "0",
			"storages.#": "0",
		}
	}

	var ecsEbsStorageSetsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEcsEbsStorageSetsMapFunc,
		fakeMapFunc:  fakeEcsEbsStorageSetsMapFunc,
	}
	ecsEbsStorageSetsCheckInfo.dataSourceTestCheck(t, rand, nameConf, zoneIdConf, storageSetIdConf, allConf)
}

func dataSourceEcsEbsStorageSetsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

resource "alibabacloudstack_ecs_ebs_storage_set" "default" {
  storage_set_name      = var.name
  maxpartition_number   = "2"
  zone_id               = data.alibabacloudstack_zones.default.zones[0].id
}

`, name, DataZoneCommonTestCase)
}
