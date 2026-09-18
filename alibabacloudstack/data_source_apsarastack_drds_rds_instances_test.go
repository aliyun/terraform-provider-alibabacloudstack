package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"
)

func TestAccAlibabacloudStackDrdsRdsInstancesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_drds_rds_instances.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_acc_drds_rds_%d", rand),
		dataSourceDrdsRdsInstancesConfigDependence)

	exisMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"rds_instances.#":                     CHECKSET,
			"rds_instances.0.db_instance_storage": CHECKSET,
			"rds_instances.0.rds_instance_id":     CHECKSET,
			"rds_instances.0.create_time":         CHECKSET,
		}
	}
	fakeMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"rds_instances.#": "0",
			"ids.#":           "0",
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
				"drds_instance_id": "${local.drds_instance_id}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"drds_instance_id": "drdsusrztw1cfake",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"drds_instance_id": "${local.drds_instance_id}",
				"ids":              []string{"${alibabacloudstack_drds_rds_instance.default.rds_instance_id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"drds_instance_id": "${local.drds_instance_id}",
				"ids":              []string{"${alibabacloudstack_drds_rds_instance.default.rds_instance_id}_fake"},
			}),
		},
	)
}

func dataSourceDrdsRdsInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

variable "existed_drds_instance" {
	default = "%s"
}

locals {
	create_drds_instance_count = var.existed_drds_instance == "" ? 1 : 0
}

variable "instance_series" {
	default = "drds.sn2.4c16g"
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
	zone_id             = data.alibabacloudstack_zones.default.zones.0.id
	db_instance_storage = "20"
	storage_type        = "local_ssd"
	category            = "HighAvailability"
	db_instance_class   = "rds.mysql.s1.small"
	drds_instance_id    = local.drds_instance_id
}
`, name, os.Getenv("ALIBABACLOUDSTACK_TEST_EXISTED_DRDS_ID"), VSwitchCommonTestCase)
}
