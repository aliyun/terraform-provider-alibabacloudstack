package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStacMulticastDomainsDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacMulticastDomainsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacMulticastDomainsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_name}_fake"`,
		}),
	}
	IdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacMulticastDomainsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cen_transit_router_multicast_domain.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacMulticastDomainsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cen_transit_router_multicast_domain.default.id}_fake"]`,
		}),
	}

	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacMulticastDomainsDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_description}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacMulticastDomainsDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_description}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacMulticastDomainsDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_name}"`,
			"description_regex": `"${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_description}"`,
			"ids":               `["${alibabacloudstack_cen_transit_router_multicast_domain.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacMulticastDomainsDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_name}"`,
			"ids":               `["${alibabacloudstack_cen_transit_router_multicast_domain.default.id}"]`,
			"description_regex": `"${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_description}_fake"`,
		}),
	}

	MulticastDomainsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, IdsConf, descriptionConf, allConf)
}

func testAccCheckAlibabacloudStacMulticastDomainsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccMulticastDomainsDatasource%d"
}

resource "alibabacloudstack_cen_instance" "default" {
	cen_instance_name = "${var.name}"
	description = "${var.name}"
	transit_router_name = "${var.name}"
	transit_router_description = "${var.name}"
}

resource "alibabacloudstack_cen_transit_router_multicast_domain" "default" {
	transit_router_multicast_domain_description = "${var.name}"
	transit_router_multicast_domain_name = "${var.name}"
	transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
}

data "alibabacloudstack_cen_transit_router_multicast_domains" "default" {
	transit_router_route_id="${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_id}"
	%s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existMulticastDomainsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"transit_router_multicast_domains.#": "1",
		"ids.#":                              "1",
		"transit_router_multicast_domains.0.transit_router_multicast_domain_name":        CHECKSET,
		"transit_router_multicast_domains.0.transit_router_multicast_domain_description": CHECKSET,
		"transit_router_multicast_domains.0.transit_router_id":                           CHECKSET,
		"transit_router_multicast_domains.0.status":                                      CHECKSET,
	}
}

var fakeMulticastDomainsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"transit_router_multicast_domains.#": "0",
		"ids.#":                              "0",
	}
}

var MulticastDomainsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_cen_transit_router_multicast_domains.default",
	existMapFunc: existMulticastDomainsMapFunc,
	fakeMapFunc:  fakeMulticastDomainsMapFunc,
}
