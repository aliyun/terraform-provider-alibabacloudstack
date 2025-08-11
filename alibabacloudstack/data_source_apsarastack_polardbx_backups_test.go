package alibabacloudstack

import (
	"fmt"
	"testing"
	"time"
)

func TestAccAlibabacloudStackPolardbxBackupsDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_polardbx_backups.default"
	name := fmt.Sprintf("tf-testAccPolardbxBackups%v", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourcePolardbxBackupsDependence)
	createTime := time.Now().UTC()
	yesterDay := createTime.AddDate(0, 0, -1)
	tomorrow := createTime.AddDate(0, 0, 1)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${alibabacloudstack_polardbx_backup.default.instance_id}",
			"ids":            []string{"${alibabacloudstack_polardbx_backup.default.backup_set_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${alibabacloudstack_polardbx_backup.default.instance_id}",
			"ids":            []string{"${alibabacloudstack_polardbx_backup.default.backup_set_id}_fake"},
		}),
	}
	startTimeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${alibabacloudstack_polardbx_backup.default.instance_id}",
			"start_time":     yesterDay.Format("2006-01-02T15:04Z"),
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${alibabacloudstack_polardbx_backup.default.instance_id}",
			"start_time":     tomorrow.Format("2006-01-02T15:04Z"),
		}),
	}
	endTimeConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${alibabacloudstack_polardbx_backup.default.instance_id}",
			"end_time":       tomorrow.Format("2006-01-02T15:04Z"),
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${alibabacloudstack_polardbx_backup.default.instance_id}",
			"end_time":       yesterDay.Format("2006-01-02T15:04Z"),
		}),
	}
	var existDBbackupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     CHECKSET,
			"backups.#":                 CHECKSET,
			"backups.0.id":              CHECKSET,
			"backups.0.backup_model":    CHECKSET,
			"backups.0.backup_set_size": CHECKSET,
			"backups.0.backup_type":     CHECKSET,
			"backups.0.backup_set_id":   CHECKSET,
			"backups.0.status":          CHECKSET,
			"backups.0.end_time":        CHECKSET,
			"backups.0.begin_time":      CHECKSET,
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

func dataSourcePolardbxBackupsDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

%s

resource "alibabacloudstack_polardbx_backup" "default" {
  instance_id = "${local.polardbx_instance.id}"
}
`, name, VSwitchCommonTestCase,PolardbxReadOrCreateCommonTestCase())
}
