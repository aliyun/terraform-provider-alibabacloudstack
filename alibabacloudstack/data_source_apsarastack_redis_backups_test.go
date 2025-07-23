package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackRedisBackupsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRedisBackupsDataSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_redis_backup.default.id}"]`,
			"instance_id": `"${alibabacloudstack_redis_backup.default.instance_id}"`,
			"start_time":  `"${alibabacloudstack_redis_backup.default.start_time}"`,
			"end_time":    `"${alibabacloudstack_redis_backup.default.end_time}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRedisBackupsDataSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_redis_backup.default.id}_fake"]`,
			"instance_id": `"${alibabacloudstack_redis_backup.default.instance_id}"`,
			"start_time":  `"${alibabacloudstack_redis_backup.default.start_time}"`,
			"end_time":    `"${alibabacloudstack_redis_backup.default.end_time}"`,
		}),
	}

	AlibabacloudstackRedisBackupsDataCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}

var existAlibabacloudstackRedisBackupsDataMapFunc = func(rand int) map[string]string {
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

var fakeAlibabacloudstackRedisBackupsDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"backups.#": "0",
	}
}

var AlibabacloudstackRedisBackupsDataCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_redis_backups.default",
	existMapFunc: existAlibabacloudstackRedisBackupsDataMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackRedisBackupsDataMapFunc,
}

func testAccCheckAlibabacloudstackRedisBackupsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackRedisBackups%d"
}

%s
resource "alibabacloudstack_db_instance" "default" {
	engine               = "MySQL"
	engine_version       = "5.6"
	instance_type        = "rds.mysql.s2.large"
	instance_storage     = "20"
	instance_name        = "${var.name}"
	vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
	storage_type         = "local_ssd"
  }

resource "alibabacloudstack_redis_backup" "default" {
	backup_method = "Physical"
	instance_id = alibabacloudstack_db_instance.default.id
}

data "alibabacloudstack_redis_backups" "default" {
	%s
}

`, rand, VSwitchCommonTestCase, strings.Join(pairs, "\n   "))
	return config
}
