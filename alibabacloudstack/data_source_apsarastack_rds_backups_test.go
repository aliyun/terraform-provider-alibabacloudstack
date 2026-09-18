package alibabacloudstack

import (
	"fmt"
	"testing"
	"time"
)

func TestAccAlibabacloudStackRdsBackupsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_rds_backups.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-RdsBackups%d", rand),
		dataSourceRdsBackupsConfigDependence)

	createTime := time.Now().UTC()
	tomorrowTime := createTime.AddDate(0, 0, 1)
	twoDayAgoTime := createTime.AddDate(0, 0, 2)

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
				"backup_ids":  []string{"${alibabacloudstack_rds_backup.default.backup_id}"},
				"instance_id": "${alibabacloudstack_rds_backup.default.instance_id}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"backup_ids":  []string{"${alibabacloudstack_rds_backup.default.backup_id}_fake"},
				"instance_id": "${alibabacloudstack_rds_backup.default.instance_id}",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"instance_id": "${alibabacloudstack_rds_backup.default.instance_id}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"instance_id": "${alibabacloudstack_rds_backup.default.instance_id}_fake",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"instance_id": "${alibabacloudstack_rds_backup.default.instance_id}",
				"start_time":  createTime.Format("2006-01-02T15:04Z"),
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"instance_id": "${alibabacloudstack_rds_backup.default.instance_id}",
				"start_time":  tomorrowTime.Format("2006-01-02T15:04Z"),
				"end_time":    twoDayAgoTime.Format("2006-01-02T15:04Z"),
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"instance_id": "${alibabacloudstack_rds_backup.default.instance_id}",
				"end_time":    tomorrowTime.Format("2006-01-02T15:04Z"),
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"instance_id": "${alibabacloudstack_rds_backup.default.instance_id}",
				"end_time":    createTime.Format("2006-01-02T15:04Z"),
			}),
		},
	)
}

func dataSourceRdsBackupsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s

%s

resource "alibabacloudstack_rds_backup" "default" {
	backup_method = "Physical"
	instance_id   = alibabacloudstack_db_instance.default.id
}
`, name, VSwitchCommonTestCase, RdsMysqlCommonTestCase())
}
