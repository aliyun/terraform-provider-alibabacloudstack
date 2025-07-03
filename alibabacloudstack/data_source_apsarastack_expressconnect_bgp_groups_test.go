package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackExpressconnectBgpGroupsDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000, 2000)
	resourceId := "data.alibabacloudstack_expressconnect_bgp_groups.default"
	testAccPreCheckWithEnvVariable(t, "ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID")

	name := fmt.Sprintf("tf-testAcc-ExpressconnectBgpGroupsDataSource-%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceExpressconnectBgpGroupsDependence(rand))

	descriptionRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"router_id":         "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
			"description_regex": "${alibabacloudstack_expressconnect_bgp_group.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"router_id":         "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
			"description_regex": "${alibabacloudstack_expressconnect_bgp_group.default.description}-fakeTestAcccc",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"router_id": "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
			"ids":       []string{"${alibabacloudstack_expressconnect_bgp_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"router_id": "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
			"ids":       []string{"${alibabacloudstack_expressconnect_bgp_group.default.id}-fakeTestAcccc"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"router_id":  "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
			"name_regex": "${alibabacloudstack_expressconnect_bgp_group.default.description}",
			"ids":        []string{"${alibabacloudstack_expressconnect_bgp_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"router_id":  "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
			"name_regex": "${alibabacloudstack_expressconnect_bgp_group.default.description}-fakeTestAcccc",
			"ids":        []string{"${alibabacloudstack_expressconnect_bgp_group.default.id}-fakeTestAcccc"},
		}),
	}

	var existExpressconnectBgpGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"ids.0":                       CHECKSET,
			"bgp_groups.#":                "1",
			"bgp_groups.0.description":    name,
			"bgp_groups.0.bgp_group_name": name,
			"bgp_groups.0.hold":           CHECKSET,
			"bgp_groups.0.ip_version":     CHECKSET,
			"bgp_groups.0.is_fake":        CHECKSET,
			"bgp_groups.0.keepalive":      CHECKSET,
			"bgp_groups.0.local_asn":      CHECKSET,
			"bgp_groups.0.peer_asn":       CHECKSET,
			"bgp_groups.0.route_limit":    CHECKSET,
			"bgp_groups.0.router_id":      CHECKSET,
		}
	}

	var fakeExpressconnectBgpGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":        "0",
			"bgp_groups.#": "0",
		}
	}

	var ExpressconnectBgpGroupsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existExpressconnectBgpGroupsMapFunc,
		fakeMapFunc:  fakeExpressconnectBgpGroupsMapFunc,
	}

	ExpressconnectBgpGroupsCheckInfo.dataSourceTestCheck(t, rand, descriptionRegexConf, idsConf, allConf)
}

func dataSourceExpressconnectBgpGroupsDependence(vlanId int) func(name string) string {
	return func(name string) string {
		return fmt.Sprintf(` 
	variable "name" {
	  default = "%s"
	}

	resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
		physical_connection_id =     "%s"
		vlan_id =                    %d
		local_gateway_ip =           "10.0.0.1"
		peer_gateway_ip =            "10.0.0.2"
		peering_subnet_mask =        "255.255.255.252"
		virtual_border_router_name = "${var.name}"
		description =                "TestAccAlibabacloudStackExpressconnectBgpgroup_basic0"
	}
	
	resource "alibabacloudstack_expressconnect_bgp_group" "default" {
		bgp_group_name = "${var.name}"
		description =    "${var.name}"
		local_asn =      "65534"
		peer_asn =       "10"
		router_id =      "${alibabacloudstack_express_connect_virtual_border_router.default.id}"
	}
	`, name, getAccTestOsEnv("ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID"), vlanId)
	}
}
