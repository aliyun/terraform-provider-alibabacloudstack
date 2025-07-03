package alibabacloudstack

import (
	"os"
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackExpressconnectBgpGroupsDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_expressconnect_bgp_groups.default"
	router_id := os.Getenv("ALIBABACLOUDSTACK_EXCONNECT_ROUTER_ID")
	if router_id == "" {
		t.Skip("Skipping TestAccAlibabacloudStackExpressconnectBgpGroupsDataSource: The Env:ALIBABACLOUDSTACK_EXCONNECT_ROUTER_ID unset!")
		t.Skipped()
	}

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAcc%sExpressconnectBgpGroupsDataSource-%d", defaultRegionToTest, rand),
		dataSourceExpressconnectBgpGroupsDependence)

	descriptionRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"router_id": router_id,
			"description_regex": "${alibabacloudstack_expressconnect_bgp_group.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"router_id": router_id,
			"description_regex": "${alibabacloudstack_expressconnect_bgp_group.default.description}-fakeTestAcccc",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"router_id": router_id,
			"ids": []string{"${alibabacloudstack_expressconnect_bgp_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"router_id": router_id,
			"ids": []string{"${alibabacloudstack_expressconnect_bgp_group.default.id}-fakeTestAcccc"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"router_id": router_id,
			"name_regex": "${alibabacloudstack_expressconnect_bgp_group.default.description}",
			"ids":        []string{"${alibabacloudstack_expressconnect_bgp_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"router_id": router_id,
			"name_regex": "${alibabacloudstack_expressconnect_bgp_group.default.description}-fakeTestAcccc",
			"ids":        []string{"${alibabacloudstack_expressconnect_bgp_group.default.id}-fakeTestAcccc"},
		}),
	}

	var existExpressconnectBgpGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"ids.0":                            CHECKSET,
			"bgp_groups.#":             "1",
			"bgp_groups.0.description": fmt.Sprintf("tf-testAcc%sExpressconnectBgpGroupsDataSource-%d", defaultRegionToTest, rand),
			"bgp_groups.0.bgp_group_name":   fmt.Sprintf("tf-testAcc%sExpressconnectBgpGroupsDataSource-%d", defaultRegionToTest, rand),
			"bgp_groups.0.hold": CHECKSET,
			"bgp_groups.0.ip_version": CHECKSET,
			"bgp_groups.0.is_fake": CHECKSET,
			"bgp_groups.0.keepalive": CHECKSET,
			"bgp_groups.0.local_asn": CHECKSET,
			"bgp_groups.0.peer_asn": CHECKSET,
			"bgp_groups.0.route_limit": CHECKSET,
			"bgp_groups.0.router_id": CHECKSET,
		}
	}

	var fakeExpressconnectBgpGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                "0",
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

func dataSourceExpressconnectBgpGroupsDependence(name string) string {
	router_id := os.Getenv("ALIBABACLOUDSTACK_EXCONNECT_ROUTER_ID")
	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}

resource "alibabacloudstack_expressconnect_bgp_group" "default" {
	bgp_group_name = "${var.name}"
	description =    "${var.name}"
	local_asn =      "65534"
	peer_asn =       "10"
	router_id =      "%s"
}

 `, name, router_id)
}
