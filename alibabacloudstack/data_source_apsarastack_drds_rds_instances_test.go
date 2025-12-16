package alibabacloudstack

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackDrdsRdsInstancesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	drdsInstanceIdRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackDrdsRdsInstancesSourceConfig(rand, map[string]string{
			"drds_instance_id": `"${local.drds_instance_id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackDrdsRdsInstancesSourceConfig(rand, map[string]string{
			"drds_instance_id": `"drdsusrztw1cfake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackDrdsRdsInstancesSourceConfig(rand, map[string]string{
			"drds_instance_id": `"${local.drds_instance_id}"`,
			"ids": `["${alibabacloudstack_drds_rds_instance.default.rds_instance_id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackDrdsRdsInstancesSourceConfig(rand, map[string]string{
			"drds_instance_id": `"${local.drds_instance_id}"`,
			"ids": `["${alibabacloudstack_drds_rds_instance.default.rds_instance_id}_fake"]`,
		}),
	}

	var exisMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rds_instances.#":                     CHECKSET,
			"rds_instances.0.db_instance_storage": CHECKSET,
			"rds_instances.0.rds_instance_id":     CHECKSET,
			"rds_instances.0.create_time":         CHECKSET,
		}
	}
	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rds_instances.#": "0",
			"ids.#":           "0",
		}
	}

	var CheckInfo = dataSourceAttr{
		resourceId:   "data.alibabacloudstack_drds_rds_instances.default",
		existMapFunc: exisMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}
	CheckInfo.dataSourceTestCheck(t, rand, drdsInstanceIdRegexConf, idsConf)
}

func testAccCheckAlibabacloudStackDrdsRdsInstancesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
	variable "name" {
	  default = "tf_acc_drds_rds_%d"
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
	
data "alibabacloudstack_drds_rds_instances" "default" {
  %s
}
`, rand, os.Getenv("ALIBABACLOUDSTACK_TEST_EXISTED_DRDS_ID"), VSwitchCommonTestCase, strings.Join(pairs, "\n  "))
	return config
}
