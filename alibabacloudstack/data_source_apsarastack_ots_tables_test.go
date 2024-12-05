package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackOtsTablesDataSource_basic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alibabacloudstack_ots_tables.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("testAcc%d", rand),
		dataSourceOtsTablesConfigDependence)

	instanceNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alibabacloudstack_ots_table.default.instance_name}",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alibabacloudstack_ots_table.default.instance_name}",
			"ids":           []string{"${alibabacloudstack_ots_table.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alibabacloudstack_ots_table.default.instance_name}",
			"ids":           []string{"${alibabacloudstack_ots_table.default.id}-fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alibabacloudstack_ots_table.default.instance_name}",
			"name_regex":    "${alibabacloudstack_ots_table.default.table_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alibabacloudstack_ots_table.default.instance_name}",
			"name_regex":    "${alibabacloudstack_ots_table.default.table_name}-fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alibabacloudstack_ots_table.default.instance_name}",
			"ids":           []string{"${alibabacloudstack_ots_table.default.id}"},
			"name_regex":    "${alibabacloudstack_ots_table.default.table_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alibabacloudstack_ots_table.default.instance_name}",
			"ids":           []string{"${alibabacloudstack_ots_table.default.id}"},
			"name_regex":    "${alibabacloudstack_ots_table.default.table_name}-fake",
		}),
	}

	var existOtsTablesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                "1",
			"names.0":                CHECKSET,
			"tables.#":               "1",
			"tables.0.table_name":    CHECKSET,
			"tables.0.instance_name": CHECKSET,
			"tables.0.primary_key.#": "2",
			"tables.0.time_to_live":  "-1",
			"tables.0.max_version":   "1",
		}
	}

	var fakeOtsTablesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":  "0",
			"tables.#": "0",
		}
	}

	var otsTablesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOtsTablesMapFunc,
		fakeMapFunc:  fakeOtsTablesMapFunc,
	}

	otsTablesCheckInfo.dataSourceTestCheck(t, rand, instanceNameConf, idsConf, nameRegexConf, allConf)
}

func dataSourceOtsTablesConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}
	resource "alibabacloudstack_ots_instance" "default" {
	  name = "tf-${var.name}"
	  description = "${var.name}"
	  accessed_by = "Any"
	  instance_type = "Capacity"
	  tags = {
	    Created = "TF"
	    For = "acceptance test"
	  }
	}

	resource "alibabacloudstack_ots_table" "default" {
	  instance_name = "${alibabacloudstack_ots_instance.default.name}"
	  table_name = "${var.name}"
	  primary_key {
          name = "pk1"
	      type = "Integer"
	  }
	  primary_key {
          name = "pk2"
          type = "String"
      }
	  time_to_live = -1
	  max_version = 1
	}
	`, name)
}
