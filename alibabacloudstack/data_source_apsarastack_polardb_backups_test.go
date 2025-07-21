package alibabacloudstack

import (
	"fmt"
	"testing"
	"time"
)

func TestAccAlibabacloudStackPolardbBackupsDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_polardb_backups.default"
	name := fmt.Sprintf("tf-testAccPolardbBackups%v", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourcePolardbBackupsDependence)
	createTime := time.Now().UTC()
	tomorrowTime := createTime.AddDate(0, 0, 1)
	twoDayAgoTime := createTime.AddDate(0, 0, 2)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${alibabacloudstack_polardb_dbinstance.default.id}",
			"ids":            []string{"${alibabacloudstack_polardb_backup.default.backup_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${alibabacloudstack_polardb_dbinstance.default.id}",
			"ids":            []string{"${alibabacloudstack_polardb_backup.default.backup_id}_fake"},
		}),
	}
	startTimeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${alibabacloudstack_polardb_dbinstance.default.id}",
			"start_time":     createTime.Format("2006-01-02T15:04Z"),
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${alibabacloudstack_polardb_dbinstance.default.id}",
			"start_time":     tomorrowTime.Format("2006-01-02T15:04Z"),
			"end_time":       twoDayAgoTime.Format("2006-01-02T15:04Z"),
		}),
	}
	endTimeConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${alibabacloudstack_polardb_dbinstance.default.id}",
			"end_time":       tomorrowTime.Format("2006-01-02T15:04Z"),
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${alibabacloudstack_polardb_dbinstance.default.id}",
			"end_time":       createTime.Format("2006-01-02T15:04Z"),
		}),
	}
	var existDBbackupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"backups.#":                   "1",
			"backups.0.id":                CHECKSET,
			"backups.0.backup_method":     CHECKSET,
			"backups.0.backup_id":         CHECKSET,
			"backups.0.backup_mode":       CHECKSET,
			"backups.0.backup_status":     CHECKSET,
			"backups.0.backup_size":       CHECKSET,
			"backups.0.slave_status":      CHECKSET,
			"backups.0.host_instance_id":  CHECKSET,
			"backups.0.backup_db_names":   CHECKSET,
			"backups.0.store_status":      CHECKSET,
			"backups.0.backup_end_time":   CHECKSET,
			"backups.0.backup_start_time": CHECKSET,
			"backups.0.meta_status":       CHECKSET,
			"backups.0.backup_scale":      CHECKSET,
			"backups.0.backup_location":   CHECKSET,
		}
	}

	var fakeDBbackupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"backups.#": "0",
		}
	}

	var DBbackupsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDBbackupsMapFunc,
		fakeMapFunc:  fakeDBbackupsMapFunc,
	}

	DBbackupsCheckInfo.dataSourceTestCheck(t, rand, idsConf, startTimeConf, endTimeConfig)
}

func dataSourcePolardbBackupsDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_polardb_dbinstance" "default" {
  instance_storage = "5"
  instance_name = "${var.name}"
  storage_type = "local_ssd"
  engine = "MySQL"
  engine_version = "5.7"
  instance_type = "rds.mysql.t1.small"
}
  
resource "alibabacloudstack_polardb_backup" "default" {
  db_instance_id = "${alibabacloudstack_polardb_dbinstance.default.id}"
  backup_method = "Physical"
}

`, name)
}
