package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"
)

func TestAccAlibabacloudStackDrdsDatabasesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_drds_databases.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_acc_drds_db_%d", rand),
		dataSourceDrdsDatabasesConfigDependence)

	exisMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"databases.#":                    CHECKSET,
			"databases.0.drds_database_name": CHECKSET,
			"databases.0.create_time":        CHECKSET,
		}
	}
	fakeMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"databases.#":           "0",
			"drds_database_names.#": "0",
		}
	}

	checkInfo := dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: exisMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}
	checkInfo.dataSourceTestCheck(t, rand,
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"instance_id": "${alibabacloudstack_drds_database.default.instance_id}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"instance_id": "drdsusrztw1cfake",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"instance_id":        "${local.drds_instance_id}",
				"drds_database_name": fmt.Sprintf("%s", fmt.Sprintf("tf_acc_drds_db_%d", rand)),
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"instance_id":        "${local.drds_instance_id}",
				"drds_database_name": "tf_acc_drds_db__fake",
			}),
		},
	)
}

func dataSourceDrdsDatabasesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

variable "existed_drds_instance" {
	default = "%s"
}

locals {
	create_drds_instance_count = var.existed_drds_instance == "" ? 1: 0
}

variable "instance_series" {
	default = "drds.sn2.4c16g"
}

resource "random_password" "password" {
	count            = 1
	length           = 12
	special          = true
	override_special = "_"
	min_lower        = 1
	min_upper        = 1
	min_numeric      = 1
}

%s

resource "alibabacloudstack_drds_instance" "default" {
	count                = local.create_drds_instance_count
	description          = var.name
	zone_id              = alibabacloudstack_vpc_vswitch.default.availability_zone
	instance_series      = var.instance_series
	instance_charge_type = "PostPaid"
	vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
	specification        = "drds.sn2.4c16g.8C32G"
}

locals {
	drds_instance_id = var.existed_drds_instance == "" ? alibabacloudstack_drds_instance.default.0.id : var.existed_drds_instance
}

resource "alibabacloudstack_drds_rds_instance" "default" {
	count               = 1
	zone_id             = data.alibabacloudstack_zones.default.zones.0.id
	db_instance_storage = "20"
	storage_type        = "local_ssd"
	category            = "HighAvailability"
	db_instance_class   = "rds.mysql.s1.small"
	drds_instance_id    = local.drds_instance_id
}

resource "alibabacloudstack_drds_database" "default" {
	instance_id        = local.drds_instance_id
	drds_database_name = var.name
	password = random_password.password.0.result
	rds_instance_ids = [
		"${alibabacloudstack_drds_rds_instance.default.0.rds_instance_id}",
	]
}
`, name, os.Getenv("ALIBABACLOUDSTACK_TEST_EXISTED_DRDS_ID"), VSwitchCommonTestCase)
}
