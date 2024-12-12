package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

var existMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instances.#":              CHECKSET,
		"instances.0.id":           CHECKSET,
		"instances.0.name":         CHECKSET,
		"instances.0.region_id":    CHECKSET,
		"instances.0.zone_id":      CHECKSET,
		"instances.0.status":       CHECKSET,
		"instances.0.tags.%":       "2",
		"instances.0.tags.Created": "TF",
		"instances.0.tags.For":     "acceptance test",
		"ids.#":                    "1",
		"ids.0":                    CHECKSET,
		"names.#":                  "1",
	}
}

var fakeMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instances.#": "0",
		"ids.#":       "0",
		"names.#":     "0",
	}
}

var checkInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_hbase_instances.default",
	existMapFunc: existMapFunc,
	fakeMapFunc:  fakeMapFunc,
}

func TestAccAlibabacloudStackHBaseInstancesDataSourceNewInstance(t *testing.T) {
	rand := getAccTestRandInt(10000,20000)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackHBaseDataSourceConfigNewInstance(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_hbase_instance.default.name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackHBaseDataSourceConfigNewInstance(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_hbase_instance.default.name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackHBaseDataSourceConfigNewInstance(rand, map[string]string{
			"ids": `["${alibabacloudstack_hbase_instance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackHBaseDataSourceConfigNewInstance(rand, map[string]string{
			"ids": `["${alibabacloudstack_hbase_instance.default.id}_fake"]`,
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackHBaseDataSourceConfigNewInstance(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_hbase_instance.default.name}"`,
			"tags":       `{Created = "TF"}`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackHBaseDataSourceConfigNewInstance(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_hbase_instance.default.name}"`,
			"tags":       `{Created = "TF1"}`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackHBaseDataSourceConfigNewInstance(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_hbase_instance.default.name}"`,
			"ids":        `["${alibabacloudstack_hbase_instance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackHBaseDataSourceConfigNewInstance(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_hbase_instance.default.name}"`,
			"ids":        `["${alibabacloudstack_hbase_instance.default.id}_fake"]`,
		}),
	}

	checkInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, tagsConf, allConf)
}

// new a instance config
func testAccCheckAlibabacloudStackHBaseDataSourceConfigNewInstance(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccHBaseInstance_datasource_%d"
}

data "alibabacloudstack_hbase_zones" "default" {}
data "alibabacloudstack_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alibabacloudstack_vswitches" "default" {
  vpc_id = data.alibabacloudstack_vpcs.default.ids.0
  zone_id = data.alibabacloudstack_hbase_zones.default.ids.0
}
resource "alibabacloudstack_vswitch" "vswitch" {
  count             = length(data.alibabacloudstack_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = alibabacloudstack_vpc.vpc.id
  cidr_block        = cidrsubnet(data.alibabacloudstack_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id = data.alibabacloudstack_hbase_zones.default.ids.0
  vswitch_name              = var.name
}

locals {
  vswitch_id = length(data.alibabacloudstack_vswitches.default.ids) > 0 ? data.alibabacloudstack_vswitches.default.ids[0] : concat(alibabacloudstack_vswitch.vswitch.*.id, [""])[0]
}

resource "alibabacloudstack_hbase_instance" "default" {
  name = var.name
  engine_version = "2.0"
  master_instance_type = "hbase.sn1.large"
  core_instance_type = "hbase.sn1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  pay_type = "PostPaid"
  duration = 1
  auto_renew = false
  vswitch_id = local.vswitch_id
  cold_storage_size = 0
  deletion_protection = false
  immediate_delete_flag = true
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}

data "alibabacloudstack_hbase_instances" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
