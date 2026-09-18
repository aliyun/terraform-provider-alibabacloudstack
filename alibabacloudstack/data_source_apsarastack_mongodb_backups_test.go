package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackMongodbBackupsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_mongodb_backups.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAlibabacloudstackMongodbBackups%d", rand),
		dataSourceMongodbBackupsConfigDependence)

	existMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"backups.#":               "1",
			"backups.0.backup_id":     CHECKSET,
			"backups.0.backup_mode":   CHECKSET,
			"backups.0.backup_size":   CHECKSET,
			"backups.0.start_time":    CHECKSET,
			"backups.0.end_time":      CHECKSET,
			"backups.0.backup_type":   CHECKSET,
			"backups.0.backup_method": CHECKSET,
			"backups.0.status":        CHECKSET,
		}
	}
	fakeMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"backups.#": "0",
		}
	}
	checkInfo := dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}

	checkInfo.dataSourceTestCheck(t, rand,
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids":            []string{"${alibabacloudstack_mongodb_backup.default.id}"},
				"db_instance_id": "${alibabacloudstack_mongodb_backup.default.db_instance_id}",
				"start_time":     "${alibabacloudstack_mongodb_backup.default.start_time}",
				"end_time":       "${alibabacloudstack_mongodb_backup.default.end_time}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids":            []string{"${alibabacloudstack_mongodb_backup.default.id}_fake"},
				"db_instance_id": "${alibabacloudstack_mongodb_backup.default.db_instance_id}",
				"start_time":     "${alibabacloudstack_mongodb_backup.default.start_time}",
				"end_time":       "${alibabacloudstack_mongodb_backup.default.end_time}",
			}),
		},
	)
}

func dataSourceMongodbBackupsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s

resource "alibabacloudstack_mongodb_instance" "default" {
	vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
	engine_version       = "3.0"
	db_instance_class    = "dds.mongo.mid"
	db_instance_storage  = "10"
	name                 = "${var.name}"
	storage_engine       = "WiredTiger"
	instance_charge_type = "PostPaid"
	replication_factor   = "3"
}

resource "alibabacloudstack_mongodb_backup" "default" {
	backup_method  = "Physical"
	db_instance_id = alibabacloudstack_mongodb_instance.default.id
}
`, name, VSwitchCommonTestCase)
}
