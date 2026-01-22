package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackForwardEntriesDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccForwardEntryConfig%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(forwardEntriesCheckInfo.resourceId, name, testAccCheckAlibabacloudStackForwardEntriesDataSourceConfigBasic)
	forwardTableIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"forward_table_id": "${alibabacloudstack_forward_entry.default.forward_table_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"forward_table_id": "${alibabacloudstack_forward_entry.default.forward_table_id}_fake",
		}),
	}

	externalIpConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"forward_table_id": "${alibabacloudstack_forward_entry.default.forward_table_id}",
			"external_ip":      "${alibabacloudstack_forward_entry.default.external_ip}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"forward_table_id": "${alibabacloudstack_forward_entry.default.forward_table_id}",
			"external_ip":      "${alibabacloudstack_forward_entry.default.external_ip}_fake",
		}),
	}

	internalIpConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"internal_ip":      "${alibabacloudstack_forward_entry.default.internal_ip}",
			"forward_table_id": "${alibabacloudstack_forward_entry.default.forward_table_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"internal_ip":      "${alibabacloudstack_forward_entry.default.internal_ip}_fake",
			"forward_table_id": "${alibabacloudstack_forward_entry.default.forward_table_id}",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"forward_table_id": "${alibabacloudstack_forward_entry.default.forward_table_id}",
			"ids":              []string{ "${alibabacloudstack_forward_entry.default.forward_entry_id}" },
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"forward_table_id": "${alibabacloudstack_forward_entry.default.forward_table_id}",
			"ids":              []string{ "${alibabacloudstack_forward_entry.default.forward_entry_id}_fake" },
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"forward_table_id": "${alibabacloudstack_forward_entry.default.forward_table_id}",
			"name_regex":       "${alibabacloudstack_forward_entry.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"forward_table_id": "${alibabacloudstack_forward_entry.default.forward_table_id}",
			"name_regex":       "${alibabacloudstack_forward_entry.default.name}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"forward_table_id": "${alibabacloudstack_forward_entry.default.forward_table_id}",
			"external_ip":      "${alibabacloudstack_forward_entry.default.external_ip}",
			"internal_ip":      "${alibabacloudstack_forward_entry.default.internal_ip}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"forward_table_id": "${alibabacloudstack_forward_entry.default.forward_table_id}",
			"external_ip":      "${alibabacloudstack_forward_entry.default.external_ip}_fake",
			"internal_ip":      "${alibabacloudstack_forward_entry.default.internal_ip}_fake",
		}),
	}
	forwardEntriesCheckInfo.dataSourceTestCheck(t, rand, forwardTableIdConf, externalIpConf, internalIpConf, idsConf ,nameRegexConf, allConf)

}

func testAccCheckAlibabacloudStackForwardEntriesDataSourceConfigBasic(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alibabacloudstack_nat_gateway" "default" {
	vpc_id = "${alibabacloudstack_vswitch.default.vpc_id}"
	specification = "Small"
	name = "${var.name}"
}

resource "alibabacloudstack_eip" "default" {
	name = "${var.name}"
}

resource "alibabacloudstack_eip_association" "default" {
	allocation_id = "${alibabacloudstack_eip.default.id}"
	instance_id = "${alibabacloudstack_nat_gateway.default.id}"
}

resource "alibabacloudstack_forward_entry" "default"{
	forward_table_id = "${alibabacloudstack_nat_gateway.default.forward_table_ids}"
	external_ip = "${alibabacloudstack_eip.default.ip_address}"
	external_port = "80"
	ip_protocol = "tcp"
	internal_ip = "172.16.0.3"
	internal_port = "8080"
}

`, name)
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
	resourceId:   "data.alibabacloudstack_forward_entries.default",
	existMapFunc: existForwardEntriesMapFunc,
	fakeMapFunc:  fakeForwardEntriesMapFunc,
}
