package apsarastack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccApsaraStackGpdbAccountsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackGpdbAccountsDataSourceName(rand, map[string]string{
			"ids": `["${apsarastack_gpdb_account.default.id}"]`,
		}),
		fakeConfig: testAccCheckApsaraStackGpdbAccountsDataSourceName(rand, map[string]string{
			"ids": `["${apsarastack_gpdb_account.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackGpdbAccountsDataSourceName(rand, map[string]string{
			"name_regex": `"${apsarastack_gpdb_account.default.account_name}"`,
		}),
		fakeConfig: testAccCheckApsaraStackGpdbAccountsDataSourceName(rand, map[string]string{
			"name_regex": `"${apsarastack_gpdb_account.default.account_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackGpdbAccountsDataSourceName(rand, map[string]string{
			"ids":    `["${apsarastack_gpdb_account.default.id}"]`,
			"status": `"Active"`,
		}),
		fakeConfig: testAccCheckApsaraStackGpdbAccountsDataSourceName(rand, map[string]string{
			"ids":    `["${apsarastack_gpdb_account.default.id}"]`,
			"status": `"Creating"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackGpdbAccountsDataSourceName(rand, map[string]string{
			"ids":        `["${apsarastack_gpdb_account.default.id}"]`,
			"name_regex": `"${apsarastack_gpdb_account.default.account_name}"`,
			"status":     `"Active"`,
		}),
		fakeConfig: testAccCheckApsaraStackGpdbAccountsDataSourceName(rand, map[string]string{
			"ids":        `["${apsarastack_gpdb_account.default.id}_fake"]`,
			"name_regex": `"${apsarastack_gpdb_account.default.account_name}_fake"`,
			"status":     `"Creating"`,
		}),
	}
	var existApsaraStackGpdbAccountsDataSourceNameMapFunc = func(rand int) map[string]string {
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
	var fakeApsaraStackGpdbAccountsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var apsarastackGpdbAccountsCheckInfo = dataSourceAttr{
		resourceId:   "data.apsarastack_gpdb_accounts.default",
		existMapFunc: existApsaraStackGpdbAccountsDataSourceNameMapFunc,
		fakeMapFunc:  fakeApsaraStackGpdbAccountsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	apsarastackGpdbAccountsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckApsaraStackGpdbAccountsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
  default = "tftestacc%d"
}
data "apsarastack_gpdb_zones" "default" {}

data "apsarastack_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "apsarastack_vswitches" "default" {
  vpc_id  = data.apsarastack_vpcs.default.ids.0
  zone_id = data.apsarastack_gpdb_zones.default.zones.2.id
}

resource "apsarastack_vswitch" "default" {
  count        = length(data.apsarastack_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.apsarastack_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.apsarastack_vpcs.default.vpcs[0].cidr_block, 8, 8)
  availability_zone      = data.apsarastack_gpdb_zones.default.zones.3.id

}

resource "apsarastack_gpdb_elastic_instance" "default" {
  engine                   = "gpdb"
  engine_version           = "6.0"
  seg_storage_type         = "cloud_essd"
  seg_node_num             = 4
  storage_size             = 50
  instance_spec            = "2C16G"
  db_instance_description  = var.name
  instance_network_type    = "VPC"
  payment_type             = "PayAsYouGo"
  vswitch_id               = length(data.apsarastack_vswitches.default.ids) > 0 ? data.apsarastack_vswitches.default.ids[0] : concat(apsarastack_vswitch.default.*.id, [""])[0]
}

resource "apsarastack_gpdb_account" "default" {
  account_name        = var.name
  db_instance_id      = apsarastack_gpdb_elastic_instance.default.id
  account_password    = "inputYourCodeHere"
  account_description = var.name
}

data "apsarastack_gpdb_accounts" "default" {	
	db_instance_id = apsarastack_gpdb_elastic_instance.default.id
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
