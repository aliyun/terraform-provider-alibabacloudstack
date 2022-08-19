package apsarastack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccApsaraStackForwardEntriesDataSourceBasic(t *testing.T) {
	rand := acctest.RandInt()
	forwardTableIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackForwardEntriesDataSourceConfigBasic(rand, map[string]string{
			"forward_table_id": `"${apsarastack_forward_entry.default.forward_table_id}"`,
		}),
		fakeConfig: testAccCheckApsaraStackForwardEntriesDataSourceConfigBasic(rand, map[string]string{
			"forward_table_id": `"${apsarastack_forward_entry.default.forward_table_id}_fake"`,
		}),
	}

	externalIpConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackForwardEntriesDataSourceConfigBasic(rand, map[string]string{
			"forward_table_id": `"${apsarastack_forward_entry.default.forward_table_id}"`,
			"external_ip":      `"${apsarastack_forward_entry.default.external_ip}"`,
		}),
		fakeConfig: testAccCheckApsaraStackForwardEntriesDataSourceConfigBasic(rand, map[string]string{
			"forward_table_id": `"${apsarastack_forward_entry.default.forward_table_id}"`,
			"external_ip":      ` "${apsarastack_forward_entry.default.external_ip}_fake" `,
		}),
	}

	internalIpConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackForwardEntriesDataSourceConfigBasic(rand, map[string]string{
			"internal_ip":      `"${apsarastack_forward_entry.default.internal_ip}"`,
			"forward_table_id": `"${apsarastack_forward_entry.default.forward_table_id}"`,
		}),
		fakeConfig: testAccCheckApsaraStackForwardEntriesDataSourceConfigBasic(rand, map[string]string{
			"internal_ip":      `"${apsarastack_forward_entry.default.internal_ip}_fake"`,
			"forward_table_id": `"${apsarastack_forward_entry.default.forward_table_id}"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackForwardEntriesDataSourceConfigBasic(rand, map[string]string{
			"forward_table_id": `"${apsarastack_forward_entry.default.forward_table_id}"`,
			"ids":              `[ "${apsarastack_forward_entry.default.forward_entry_id}" ]`,
		}),
		fakeConfig: testAccCheckApsaraStackForwardEntriesDataSourceConfigBasic(rand, map[string]string{
			"forward_table_id": `"${apsarastack_forward_entry.default.forward_table_id}"`,
			"ids":              `[ "${apsarastack_forward_entry.default.forward_entry_id}_fake" ]`,
		}),
	}

	//nameRegexConf := dataSourceTestAccConfig{
	//	existConfig: testAccCheckApsaraStackForwardEntriesDataSourceConfigBasic(rand, map[string]string{
	//		"forward_table_id": `"${apsarastack_forward_entry.default.forward_table_id}"`,
	//		"name_regex":       `"${apsarastack_forward_entry.default.name}"`,
	//	}),
	//	fakeConfig: testAccCheckApsaraStackForwardEntriesDataSourceConfigBasic(rand, map[string]string{
	//		"forward_table_id": `"${apsarastack_forward_entry.default.forward_table_id}"`,
	//		"name_regex":       `"${apsarastack_forward_entry.default.name}_fake"`,
	//	}),
	//}

	//allConf := dataSourceTestAccConfig{
	//	existConfig: testAccCheckApsaraStackForwardEntriesDataSourceConfigBasic(rand, map[string]string{
	//		"forward_table_id": `"${apsarastack_forward_entry.default.forward_table_id}"`,
	//		"external_ip":      `"${apsarastack_forward_entry.default.external_ip}"`,
	//		"internal_ip":      `"${apsarastack_forward_entry.default.internal_ip}"`,
	//	}),
	//	fakeConfig: testAccCheckApsaraStackForwardEntriesDataSourceConfigBasic(rand, map[string]string{
	//		"forward_table_id": `"${apsarastack_forward_entry.default.forward_table_id}"`,
	//		"external_ip":      `"${apsarastack_forward_entry.default.external_ip}"`,
	//		"internal_ip":      `"${apsarastack_forward_entry.default.internal_ip}"`,
	//	}),
	//}
	forwardEntriesCheckInfo.dataSourceTestCheck(t, rand, forwardTableIdConf, externalIpConf, internalIpConf, idsConf /*,nameRegexConf/*, allConf*/)

}

func testAccCheckApsaraStackForwardEntriesDataSourceConfigBasic(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccForwardEntryConfig%d"
}

data "apsarastack_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "apsarastack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "apsarastack_vswitch" "default" {
	vpc_id = "${apsarastack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "apsarastack_nat_gateway" "default" {
	vpc_id = "${apsarastack_vswitch.default.vpc_id}"
	specification = "Small"
	name = "${var.name}"
}

resource "apsarastack_eip" "default" {
	name = "${var.name}"
}

resource "apsarastack_eip_association" "default" {
	allocation_id = "${apsarastack_eip.default.id}"
	instance_id = "${apsarastack_nat_gateway.default.id}"
}

resource "apsarastack_forward_entry" "default"{
	forward_table_id = "${apsarastack_nat_gateway.default.forward_table_ids}"
	external_ip = "${apsarastack_eip.default.ip_address}"
	external_port = "80"
	ip_protocol = "tcp"
	internal_ip = "172.16.0.3"
	internal_port = "8080"
}

data "apsarastack_forward_entries" "default" {
	%s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existForwardEntriesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                   "1",
		"entries.#":               "1",
		"entries.0.id":            CHECKSET,
		"entries.0.external_ip":   CHECKSET,
		"entries.0.external_port": "80",
		"entries.0.internal_ip":   "172.16.0.3",
		"entries.0.internal_port": "8080",
		"entries.0.ip_protocol":   "tcp",
		"entries.0.status":        "Available",
	}
}

var fakeForwardEntriesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":     "0",
		"entries.#": "0",
	}
}

var forwardEntriesCheckInfo = dataSourceAttr{
	resourceId:   "data.apsarastack_forward_entries.default",
	existMapFunc: existForwardEntriesMapFunc,
	fakeMapFunc:  fakeForwardEntriesMapFunc,
}
