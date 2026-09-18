package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudstackRouteEntriesDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(1000, 9999)
	resourceId := "data.alibabacloudstack_route_entries.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccRouteEntryConfigNameDatasource%d", rand),
		dataSourceRouteEntriesConfigDependence)

	instanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":    "${alibabacloudstack_route_entry.default.nexthop_id}",
			"route_table_id": "${alibabacloudstack_route_entry.default.route_table_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":    "${alibabacloudstack_route_entry.default.nexthop_id}",
			"route_table_id": "${alibabacloudstack_route_entry.default.route_table_id}_fake",
		}),
	}

	typeConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":    "${alibabacloudstack_route_entry.default.nexthop_id}",
			"route_table_id": "${alibabacloudstack_route_entry.default.route_table_id}",
			"type":           "Custom",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":    "${alibabacloudstack_route_entry.default.nexthop_id}",
			"route_table_id": "${alibabacloudstack_route_entry.default.route_table_id}",
			"type":           "System",
		}),
	}

	cidrBlockConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"route_table_id": "${alibabacloudstack_route_entry.default.route_table_id}",
			"cidr_block":     "${alibabacloudstack_route_entry.default.destination_cidrblock}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"route_table_id": "${alibabacloudstack_route_entry.default.route_table_id}",
			"cidr_block":     "${alibabacloudstack_route_entry.default.destination_cidrblock}_fake",
		}),
	}

	allConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":    "${alibabacloudstack_ecs_instance.default.id}",
			"route_table_id": "${alibabacloudstack_route_entry.default.route_table_id}",
			"type":           "Custom",
			"cidr_block":     "${alibabacloudstack_route_entry.default.destination_cidrblock}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":    "${alibabacloudstack_ecs_instance.default.id}",
			"route_table_id": "${alibabacloudstack_route_entry.default.route_table_id}",
			"type":           "Custom",
			"cidr_block":     "${alibabacloudstack_route_entry.default.destination_cidrblock}_fake",
		}),
	}

	routeEntriesCheckInfo.dataSourceTestCheck(t, rand, instanceIdConf, typeConfig, cidrBlockConfig, allConfig)
}

func dataSourceRouteEntriesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s

resource "alibabacloudstack_route_entry" "default" {
   route_table_id = "${alibabacloudstack_vpc_vpc.default.route_table_id}"
   destination_cidrblock = "172.11.1.1/32"
   nexthop_type = "Instance"
   nexthop_id = "${alibabacloudstack_ecs_instance.default.id}"
}
`, name, ECSInstanceCommonTestCase)
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
