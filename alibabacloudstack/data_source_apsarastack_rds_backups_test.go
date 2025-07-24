package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestAccAlibabacloudStackRdsBackupsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)
	createTime := time.Now().UTC()
	tomorrowTime := createTime.AddDate(0, 0, 1)
	twoDayAgoTime := createTime.AddDate(0, 0, 2)

	instanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRdsBackupsDataSourceConfig(rand, map[string]string{
			"instance_id": `"${alibabacloudstack_rds_backup.default.instance_id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRdsBackupsDataSourceConfig(rand, map[string]string{
			"instance_id": `"${alibabacloudstack_rds_backup.default.instance_id}_fake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRdsBackupsDataSourceConfig(rand, map[string]string{
			"backup_ids":  `["${alibabacloudstack_rds_backup.default.backup_id}"]`,
			"instance_id": `"${alibabacloudstack_rds_backup.default.instance_id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRdsBackupsDataSourceConfig(rand, map[string]string{
			"backup_ids":  `["${alibabacloudstack_rds_backup.default.backup_id}_fake"]`,
			"instance_id": `"${alibabacloudstack_rds_backup.default.instance_id}"`,
		}),
	}
	startTimeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRdsBackupsDataSourceConfig(rand, map[string]string{
			"instance_id": `"${alibabacloudstack_rds_backup.default.instance_id}"`,
			"start_time":  fmt.Sprintf(`"%s"`, createTime.Format("2006-01-02T15:04Z")),
		}),
		fakeConfig: testAccCheckAlibabacloudstackRdsBackupsDataSourceConfig(rand, map[string]string{
			"instance_id": `"${alibabacloudstack_rds_backup.default.instance_id}"`,
			"start_time":  fmt.Sprintf(`"%s"`, tomorrowTime.Format("2006-01-02T15:04Z")),
			"end_time":    fmt.Sprintf(`"%s"`, twoDayAgoTime.Format("2006-01-02T15:04Z")),
		}),
	}
	endTimeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRdsBackupsDataSourceConfig(rand, map[string]string{
			"instance_id": `"${alibabacloudstack_rds_backup.default.instance_id}"`,
			"end_time":    fmt.Sprintf(`"%s"`, tomorrowTime.Format("2006-01-02T15:04Z")),
		}),
		fakeConfig: testAccCheckAlibabacloudstackRdsBackupsDataSourceConfig(rand, map[string]string{
			"instance_id": `"${alibabacloudstack_rds_backup.default.instance_id}"`,
			"end_time":    fmt.Sprintf(`"%s"`, createTime.Format("2006-01-02T15:04Z")),
		}),
	}

	AlibabacloudstackRdsBackupsDataCheckInfo.dataSourceTestCheck(t, rand, idsConf, instanceIdConf, startTimeConf, endTimeConf)
}

var existAlibabacloudstackRdsBackupsDataMapFunc = func(rand int) map[string]string {
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

var fakeAlibabacloudstackRdsBackupsDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"backups.#": "0",
	}
}

var AlibabacloudstackRdsBackupsDataCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_rds_backups.default",
	existMapFunc: existAlibabacloudstackRdsBackupsDataMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackRdsBackupsDataMapFunc,
}

func testAccCheckAlibabacloudstackRdsBackupsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testacc-RdsBackups%d"
}

%s
%s
resource "alibabacloudstack_rds_backup" "default" {
	backup_method = "Physical"
	instance_id = alibabacloudstack_db_instance.default.id
}

data "alibabacloudstack_rds_backups" "default" {
	%s
}

`, rand, VSwitchCommonTestCase, RdsMysqlCommonTestCase(), strings.Join(pairs, "\n   "))
	return config
}
