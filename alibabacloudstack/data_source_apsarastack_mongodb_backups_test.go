package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackMongodbBackupsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackMongodbBackupsDataSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_mongodb_backup.default.id}"]`,
			"db_instance_id": `"${alibabacloudstack_mongodb_backup.default.db_instance_id}"`,
			"start_time":     `"${alibabacloudstack_mongodb_backup.default.start_time}"`,
			"end_time":       `"${alibabacloudstack_mongodb_backup.default.end_time}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackMongodbBackupsDataSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_mongodb_backup.default.id}_fake"]`,
			"db_instance_id": `"${alibabacloudstack_mongodb_backup.default.db_instance_id}"`,
			"start_time":     `"${alibabacloudstack_mongodb_backup.default.start_time}"`,
			"end_time":       `"${alibabacloudstack_mongodb_backup.default.end_time}"`,
		}),
	}

	AlibabacloudstackMongodbBackupsDataCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}

var existAlibabacloudstackMongodbBackupsDataMapFunc = func(rand int) map[string]string {
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

var fakeAlibabacloudstackMongodbBackupsDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"backups.#": "0",
	}
}

var AlibabacloudstackMongodbBackupsDataCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_mongodb_backups.default",
	existMapFunc: existAlibabacloudstackMongodbBackupsDataMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackMongodbBackupsDataMapFunc,
}

func testAccCheckAlibabacloudstackMongodbBackupsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackMongodbBackups%d"
}

%s
resource "alibabacloudstack_mongodb_instance" "default" {
	vswitch_id          = alibabacloudstack_vpc_vswitch.default.id
	engine_version      = "3.0"
	db_instance_class   = "dds.mongo.mid"
	db_instance_storage = "10"
	name                = "${var.name}"
	storage_engine      = "WiredTiger"
	instance_charge_type = "PostPaid"
	replication_factor = "3"
  }

resource "alibabacloudstack_mongodb_backup" "default" {
	backup_method = "Physical"
	db_instance_id = alibabacloudstack_mongodb_instance.default.id
}

data "alibabacloudstack_mongodb_backups" "default" {
	%s
}

`, rand, VSwitchCommonTestCase, strings.Join(pairs, "\n   "))
	return config
}
