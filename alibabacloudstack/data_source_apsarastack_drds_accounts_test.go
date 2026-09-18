package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackDrdsAccountsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_drds_accounts.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_acc_drds_db_%d", rand),
		dataSourceDrdsAccountsConfigDependence)

	exisMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"accounts.#":                   CHECKSET,
			"accounts.0.drds_account_name": CHECKSET,
			"accounts.0.instance_id":       CHECKSET,
		}
	}
	fakeMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"accounts.#": "0",
			"ids.#":      "0",
			"names.#":    "0",
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
				"instance_id": "${alibabacloudstack_drds_database.default.instance_id}",
				"names":       []string{"${var.name}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"instance_id": "${alibabacloudstack_drds_database.default.instance_id}",
				"names":       []string{"${var.name}_fake"},
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"instance_id":  "${alibabacloudstack_drds_database.default.instance_id}",
				"names":        []string{"${var.name}_db"},
				"account_type": TfRawString("0"),
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"instance_id":  "${alibabacloudstack_drds_database.default.instance_id}",
				"names":        []string{"${var.name}_db"},
				"account_type": TfRawString("1"),
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"instance_id":  "${alibabacloudstack_drds_database.default.instance_id}",
				"names":        []string{"${var.name}"},
				"account_type": TfRawString("1"),
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"instance_id":  "${alibabacloudstack_drds_database.default.instance_id}",
				"names":        []string{"${var.name}"},
				"account_type": TfRawString("0"),
			}),
		},
	)
}

func dataSourceDrdsAccountsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

variable "instance_series" {
	default = "drds.sn2.4c16g"
}

resource "random_password" "password" {
	length           = 12
	special          = true
	override_special = "_"
	min_lower        = 1
	min_upper        = 1
	min_numeric      = 1
}

%s

resource "alibabacloudstack_drds_instance" "default" {
	description          = "${var.name}"
	zone_id              = "${alibabacloudstack_vpc_vswitch.default.availability_zone}"
	instance_series      = "${var.instance_series}"
	instance_charge_type = "PostPaid"
	vswitch_id           = "${alibabacloudstack_vpc_vswitch.default.id}"
	specification        = "drds.sn2.4c16g.8C32G"
}

resource "alibabacloudstack_drds_rds_instance" "default" {
	zone_id             = data.alibabacloudstack_zones.default.zones.0.id
	db_instance_storage = "20"
	storage_type        = "local_ssd"
	category            = "HighAvailability"
	db_instance_class   = "rds.mysql.s1.small"
	drds_instance_id    = alibabacloudstack_drds_instance.default.id
}

resource "alibabacloudstack_drds_database" "default" {
	instance_id        = "${alibabacloudstack_drds_instance.default.id}"
	drds_database_name = "${var.name}_db"
	password           = random_password.password.result
	rds_instance_ids   = [alibabacloudstack_drds_rds_instance.default.rds_instance_id,]
}

resource "alibabacloudstack_drds_account" "default" {
	instance_id       = alibabacloudstack_drds_instance.default.id
	drds_account_name = var.name
	password          = random_password.password.result
	description       = var.name
	db_privileges {
		db_name   = alibabacloudstack_drds_database.default.drds_database_name
		privilege = "R"
	}
}
`, name, VSwitchCommonTestCase)
}
