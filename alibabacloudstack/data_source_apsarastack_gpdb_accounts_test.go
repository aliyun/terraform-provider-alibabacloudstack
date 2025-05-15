package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackGpdbAccountsDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000, 9999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackGpdbAccountsDataSourceName(rand, map[string]string{
			"ids": `["${alibabacloudstack_gpdb_account.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackGpdbAccountsDataSourceName(rand, map[string]string{
			"ids": `["${alibabacloudstack_gpdb_account.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackGpdbAccountsDataSourceName(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_gpdb_account.default.account_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackGpdbAccountsDataSourceName(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_gpdb_account.default.account_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackGpdbAccountsDataSourceName(rand, map[string]string{
			"ids":    `["${alibabacloudstack_gpdb_account.default.id}"]`,
			"status": `"Active"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackGpdbAccountsDataSourceName(rand, map[string]string{
			"ids":    `["${alibabacloudstack_gpdb_account.default.id}"]`,
			"status": `"Creating"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackGpdbAccountsDataSourceName(rand, map[string]string{
			"ids":        `["${alibabacloudstack_gpdb_account.default.id}"]`,
			"name_regex": `"${alibabacloudstack_gpdb_account.default.account_name}"`,
			"status":     `"Active"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackGpdbAccountsDataSourceName(rand, map[string]string{
			"ids":        `["${alibabacloudstack_gpdb_account.default.id}_fake"]`,
			"name_regex": `"${alibabacloudstack_gpdb_account.default.account_name}_fake"`,
			"status":     `"Creating"`,
		}),
	}
	var existAlibabacloudStackGpdbAccountsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"names.#":                        "1",
			"accounts.#":                     "1",
			"accounts.0.id":                  CHECKSET,
			"accounts.0.account_name":        fmt.Sprintf("tftestacc%d", rand),
			"accounts.0.account_description": fmt.Sprintf("tftestacc%d", rand),
			"accounts.0.db_instance_id":      CHECKSET,
			"accounts.0.status":              "Active",
		}
	}
	var fakeAlibabacloudStackGpdbAccountsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alibabacloudstackGpdbAccountsCheckInfo = dataSourceAttr{
		resourceId:   "data.alibabacloudstack_gpdb_accounts.default",
		existMapFunc: existAlibabacloudStackGpdbAccountsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlibabacloudStackGpdbAccountsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alibabacloudstackGpdbAccountsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlibabacloudStackGpdbAccountsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
  default = "tftestacc%d"
}
data "alibabacloudstack_gpdb_zones" "default" {}

data "alibabacloudstack_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alibabacloudstack_vswitches" "default" {
  vpc_id  = data.alibabacloudstack_vpcs.default.ids.0
  zone_id = data.alibabacloudstack_gpdb_zones.default.zones.2.id
}

resource "alibabacloudstack_vswitch" "default" {
  count        = length(data.alibabacloudstack_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alibabacloudstack_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alibabacloudstack_vpcs.default.vpcs[0].cidr_block, 8, 8)
  availability_zone      = data.alibabacloudstack_gpdb_zones.default.zones.3.id

}

resource "alibabacloudstack_gpdb_elastic_instance" "default" {
  engine                   = "gpdb"
  engine_version           = "6.0"
  seg_storage_type         = "cloud_essd"
  seg_node_num             = 4
  storage_size             = 50
  instance_spec            = "2C16G"
  db_instance_description  = var.name
  instance_network_type    = "VPC"
  payment_type             = "PayAsYouGo"
  vswitch_id               = length(data.alibabacloudstack_vswitches.default.ids) > 0 ? data.alibabacloudstack_vswitches.default.ids[0] : concat(alibabacloudstack_vswitch.default.*.id, [""])[0]
}

resource "alibabacloudstack_gpdb_account" "default" {
  account_name        = var.name
  db_instance_id      = alibabacloudstack_gpdb_elastic_instance.default.id
  account_password    = "%s"
  account_description = var.name
}

data "alibabacloudstack_gpdb_accounts" "default" {	
	db_instance_id = alibabacloudstack_gpdb_elastic_instance.default.id
	%s
}
`, rand, GeneratePassword(), strings.Join(pairs, " \n "))
	return config
}
