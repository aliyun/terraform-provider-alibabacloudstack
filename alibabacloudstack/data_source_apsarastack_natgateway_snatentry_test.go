package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackAlibabacloudstackNatgatewaySnatEntriesDataSource(t *testing.T) {

	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackNatgatewaySnatEntriesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_natgateway_snat_entries.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackNatgatewaySnatEntriesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_natgateway_snat_entries.default.id}_fake"]`,
		}),
	}

	snat_entry_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackNatgatewaySnatEntriesSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_natgateway_snat_entries.default.id}"]`,
			"snat_entry_id": `"${alibabacloudstack_natgateway_snat_entries.default.SnatEntryId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackNatgatewaySnatEntriesSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_natgateway_snat_entries.default.id}_fake"]`,
			"snat_entry_id": `"${alibabacloudstack_natgateway_snat_entries.default.SnatEntryId}_fake"`,
		}),
	}

	snat_entry_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackNatgatewaySnatEntriesSourceConfig(rand, map[string]string{
			"ids":             `["${alibabacloudstack_natgateway_snat_entries.default.id}"]`,
			"snat_entry_name": `"${alibabacloudstack_natgateway_snat_entries.default.SnatEntryName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackNatgatewaySnatEntriesSourceConfig(rand, map[string]string{
			"ids":             `["${alibabacloudstack_natgateway_snat_entries.default.id}_fake"]`,
			"snat_entry_name": `"${alibabacloudstack_natgateway_snat_entries.default.SnatEntryName}_fake"`,
		}),
	}

	snat_ipConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackNatgatewaySnatEntriesSourceConfig(rand, map[string]string{
			"ids":     `["${alibabacloudstack_natgateway_snat_entries.default.id}"]`,
			"snat_ip": `"${alibabacloudstack_natgateway_snat_entries.default.SnatIp}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackNatgatewaySnatEntriesSourceConfig(rand, map[string]string{
			"ids":     `["${alibabacloudstack_natgateway_snat_entries.default.id}_fake"]`,
			"snat_ip": `"${alibabacloudstack_natgateway_snat_entries.default.SnatIp}_fake"`,
		}),
	}

	snat_table_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackNatgatewaySnatEntriesSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_natgateway_snat_entries.default.id}"]`,
			"snat_table_id": `"${alibabacloudstack_natgateway_snat_entries.default.SnatTableId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackNatgatewaySnatEntriesSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_natgateway_snat_entries.default.id}_fake"]`,
			"snat_table_id": `"${alibabacloudstack_natgateway_snat_entries.default.SnatTableId}_fake"`,
		}),
	}

	source_cidrConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackNatgatewaySnatEntriesSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_natgateway_snat_entries.default.id}"]`,
			"source_cidr": `"${alibabacloudstack_natgateway_snat_entries.default.SourceCidr}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackNatgatewaySnatEntriesSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_natgateway_snat_entries.default.id}_fake"]`,
			"source_cidr": `"${alibabacloudstack_natgateway_snat_entries.default.SourceCidr}_fake"`,
		}),
	}

	source_vswitch_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackNatgatewaySnatEntriesSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_natgateway_snat_entries.default.id}"]`,
			"source_vswitch_id": `"${alibabacloudstack_natgateway_snat_entries.default.SourceVSwitchId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackNatgatewaySnatEntriesSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_natgateway_snat_entries.default.id}_fake"]`,
			"source_vswitch_id": `"${alibabacloudstack_natgateway_snat_entries.default.SourceVSwitchId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackNatgatewaySnatEntriesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_natgateway_snat_entries.default.id}"]`,

			"snat_entry_id":     `"${alibabacloudstack_natgateway_snat_entries.default.SnatEntryId}"`,
			"snat_entry_name":   `"${alibabacloudstack_natgateway_snat_entries.default.SnatEntryName}"`,
			"snat_ip":           `"${alibabacloudstack_natgateway_snat_entries.default.SnatIp}"`,
			"snat_table_id":     `"${alibabacloudstack_natgateway_snat_entries.default.SnatTableId}"`,
			"source_cidr":       `"${alibabacloudstack_natgateway_snat_entries.default.SourceCidr}"`,
			"source_vswitch_id": `"${alibabacloudstack_natgateway_snat_entries.default.SourceVSwitchId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackNatgatewaySnatEntriesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_natgateway_snat_entries.default.id}_fake"]`,

			"snat_entry_id":     `"${alibabacloudstack_natgateway_snat_entries.default.SnatEntryId}_fake"`,
			"snat_entry_name":   `"${alibabacloudstack_natgateway_snat_entries.default.SnatEntryName}_fake"`,
			"snat_ip":           `"${alibabacloudstack_natgateway_snat_entries.default.SnatIp}_fake"`,
			"snat_table_id":     `"${alibabacloudstack_natgateway_snat_entries.default.SnatTableId}_fake"`,
			"source_cidr":       `"${alibabacloudstack_natgateway_snat_entries.default.SourceCidr}_fake"`,
			"source_vswitch_id": `"${alibabacloudstack_natgateway_snat_entries.default.SourceVSwitchId}_fake"`}),
	}

	AlibabacloudstackNatgatewaySnatEntriesCheckInfo.dataSourceTestCheck(t, rand, idsConf, snat_entry_idConf, snat_entry_nameConf, snat_ipConf, snat_table_idConf, source_cidrConf, source_vswitch_idConf, allConf)
}

var existAlibabacloudstackNatgatewaySnatEntriesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"entries.#":    "1",
		"entries.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackNatgatewaySnatEntriesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"entries.#": "0",
	}
}

var AlibabacloudstackNatgatewaySnatEntriesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_natgateway_snat_entries.default",
	existMapFunc: existAlibabacloudstackNatgatewaySnatEntriesMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackNatgatewaySnatEntriesMapFunc,
}

func testAccCheckAlibabacloudstackNatgatewaySnatEntriesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackNatgatewaySnatEntries%d"
}






data "alibabacloudstack_natgateway_snat_entries" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
