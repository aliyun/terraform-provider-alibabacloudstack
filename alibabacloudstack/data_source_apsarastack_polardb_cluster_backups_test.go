package alibabacloudstack

import (
	"fmt"
	"testing"
	"time"
)

func TestAccAlibabacloudStackPolardbClusterBackupsDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_polardb_cluster_backups.default"
	name := fmt.Sprintf("tf-testAccPolardbClusterBackups%v", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourcePolardbClusterBackupsDependence)
	createTime := time.Now().UTC()
	twoDayLaterTime := createTime.AddDate(0, 0, 2)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_cluster_id": "${local.polardb_cluster_id}",
			"ids":           []string{"${alibabacloudstack_polardb_cluster_backup.default.backup_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"db_cluster_id": "${local.polardb_cluster_id}",
			"ids":           []string{"${alibabacloudstack_polardb_cluster_backup.default.backup_id}_fake"},
		}),
	}
	startTimeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_cluster_id": "${local.polardb_cluster_id}",
			"start_time":    createTime.Format(polardbClusterBackupTimeLayout),
			"ids":           []string{"${alibabacloudstack_polardb_cluster_backup.default.backup_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"db_cluster_id": "${local.polardb_cluster_id}",
		// The service rejects a future StartTime, so the negative window is
		// built by inverting the bounds instead of shifting them ahead: a
		// start later than the end deterministically matches no backup set.
			"start_time":    createTime.Format(polardbClusterBackupTimeLayout),
			"end_time":      createTime.AddDate(0, 0, -30).Format(polardbClusterBackupTimeLayout),
		}),
	}
	endTimeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_cluster_id": "${local.polardb_cluster_id}",
			"end_time":      twoDayLaterTime.Format(polardbClusterBackupTimeLayout),
			"ids":           []string{"${alibabacloudstack_polardb_cluster_backup.default.backup_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"db_cluster_id": "${local.polardb_cluster_id}",
			"end_time":      createTime.AddDate(-3, 0, 0).Format(polardbClusterBackupTimeLayout),
		}),
	}
	var existPolardbClusterBackupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"backups.#":                   CHECKSET,
			"backups.0.id":                CHECKSET,
			"backups.0.backup_id":         CHECKSET,
			"backups.0.db_cluster_id":     CHECKSET,
			"backups.0.backup_method":     CHECKSET,
			"backups.0.backup_mode":       CHECKSET,
			"backups.0.backup_type":       CHECKSET,
			"backups.0.backup_status":     CHECKSET,
			"backups.0.backups_level":     CHECKSET,
			"backups.0.store_status":      CHECKSET,
			"backups.0.backup_set_size":   CHECKSET,
			"backups.0.backup_start_time": CHECKSET,
			"backups.0.backup_end_time":   CHECKSET,
		}
	}

	var fakePolardbClusterBackupsMapFunc = func(rand int) map[string]string {
		// The cluster is shared with other test cases and keeps its earlier backup
		// sets, so only the ids-filtered count is deterministic for a fake filter;
		// the unfiltered list size is not.
		return map[string]string{
			"ids.#": "0",
		}
	}

	var polardbClusterBackupsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existPolardbClusterBackupsMapFunc,
		fakeMapFunc:  fakePolardbClusterBackupsMapFunc,
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
			testAccPreCheckApiAvailable(t, "GET", "polardb", "2017-08-01", "DescribeBackups", "")
		},
	}

	polardbClusterBackupsCheckInfo.dataSourceTestCheck(t, rand, idsConf, startTimeConf, endTimeConf)
}

func dataSourcePolardbClusterBackupsDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}

%s

resource "alibabacloudstack_polardb_cluster_backup" "default" {
  db_cluster_id = "${local.polardb_cluster_id}"
}

`, name, VSwitchCommonTestCase+PolardbClusterCommonTestCase("PostgreSQL", "16", true))
}
