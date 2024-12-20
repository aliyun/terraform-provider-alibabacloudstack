package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackRouteEntriesDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(1000, 9999)

	instanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackRouteEntriesDataSourceConfig(rand, map[string]string{
			"instance_id":    `"${alibabacloudstack_route_entry.default.nexthop_id}"`,
			"route_table_id": `"${alibabacloudstack_route_entry.default.route_table_id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackRouteEntriesDataSourceConfig(rand, map[string]string{
			"instance_id":    `"${alibabacloudstack_route_entry.default.nexthop_id}"`,
			"route_table_id": `"${alibabacloudstack_route_entry.default.route_table_id}_fake"`,
		}),
	}

	typeConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackRouteEntriesDataSourceConfig(rand, map[string]string{
			"instance_id":    `"${alibabacloudstack_route_entry.default.nexthop_id}"`,
			"route_table_id": `"${alibabacloudstack_route_entry.default.route_table_id}"`,
			"type":           `"Custom"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackRouteEntriesDataSourceConfig(rand, map[string]string{
			"instance_id":    `"${alibabacloudstack_route_entry.default.nexthop_id}"`,
			"route_table_id": `"${alibabacloudstack_route_entry.default.route_table_id}"`,
			"type":           `"System"`,
		}),
	}

	cidrBlockConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackRouteEntriesDataSourceConfig(rand, map[string]string{
			"route_table_id": `"${alibabacloudstack_route_entry.default.route_table_id}"`,
			"cidr_block":     `"${alibabacloudstack_route_entry.default.destination_cidrblock}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackRouteEntriesDataSourceConfig(rand, map[string]string{
			"route_table_id": `"${alibabacloudstack_route_entry.default.route_table_id}"`,
			"cidr_block":     `"${alibabacloudstack_route_entry.default.destination_cidrblock}_fake"`,
		}),
	}

	allConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackRouteEntriesDataSourceConfig(rand, map[string]string{
			"instance_id":    `"${alibabacloudstack_instance.default.id}"`,
			"route_table_id": `"${alibabacloudstack_route_entry.default.route_table_id}"`,
			"type":           `"Custom"`,
			"cidr_block":     `"${alibabacloudstack_route_entry.default.destination_cidrblock}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackRouteEntriesDataSourceConfig(rand, map[string]string{
			"instance_id":    `"${alibabacloudstack_instance.default.id}"`,
			"route_table_id": `"${alibabacloudstack_route_entry.default.route_table_id}"`,
			"type":           `"Custom"`,
			"cidr_block":     `"${alibabacloudstack_route_entry.default.destination_cidrblock}_fake"`,
		}),
	}

	routeEntriesCheckInfo.dataSourceTestCheck(t, rand, instanceIdConf, typeConfig, cidrBlockConfig, allConfig)
}

func testAccCheckAlibabacloudStackRouteEntriesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
%s

%s

%s

variable "name" {
   default = "tf-testAcc-for-route-entries-datasource%d"
}
resource "alibabacloudstack_vpc" "default" {
   name = "${var.name}"
   cidr_block = "10.1.0.0/21"
}
resource "alibabacloudstack_vswitch" "default" {
   vpc_id = "${alibabacloudstack_vpc.default.id}"
   cidr_block = "10.1.1.0/24"
   availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
   name = "${var.name}"
}
resource "alibabacloudstack_security_group" "default" {
   name = "${var.name}"
   description = "${var.name}"
   vpc_id = "${alibabacloudstack_vpc.default.id}"
}
resource "alibabacloudstack_security_group_rule" "default" {
   type = "ingress"
   ip_protocol = "tcp"
   nic_type = "intranet"
   policy = "accept"
   port_range = "22/22"
   priority = 1
   security_group_id = "${alibabacloudstack_security_group.default.id}"
   cidr_ip = "0.0.0.0/0"
}
resource "alibabacloudstack_instance" "default" {
   # cn-beijing
   security_groups = ["${alibabacloudstack_security_group.default.id}"]
   vswitch_id = "${alibabacloudstack_vswitch.default.id}"
   # series III
   instance_type = "${local.default_instance_type_id}"
   internet_max_bandwidth_out = 5
   system_disk_category = "cloud_pperf"
   image_id = "${data.alibabacloudstack_images.default.images.0.id}"
   instance_name = "${var.name}"
}
resource "alibabacloudstack_route_entry" "default" {
   route_table_id = "${alibabacloudstack_vpc.default.route_table_id}"
   destination_cidrblock = "172.11.1.1/32"
   nexthop_type = "Instance"
   nexthop_id = "${alibabacloudstack_instance.default.id}"
}
data "alibabacloudstack_route_entries" "default" {
  %s
}`, DataAlibabacloudstackVswitchZones, DataAlibabacloudstackInstanceTypes, DataAlibabacloudstackImages, rand, strings.Join(pairs, "\n  "))
	return config
}

var existRouteEntriesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"entries.#":                "1",
		"entries.0.route_table_id": CHECKSET,
		"entries.0.cidr_block":     CHECKSET,
		"entries.0.instance_id":    CHECKSET,
		"entries.0.status":         CHECKSET,
		"entries.0.type":           "Custom",
		"entries.0.next_hop_type":  "Instance",
	}
}

var fakeRouteEntriesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"entries.#": "0",
	}
}

var routeEntriesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_route_entries.default",
	existMapFunc: existRouteEntriesMapFunc,
	fakeMapFunc:  fakeRouteEntriesMapFunc,
}
