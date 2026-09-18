package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackCenTransitMulticastDomainsDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_cen_transit_router_multicast_domains.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccMulticastDomainsDatasource%d", rand),
		dataSourceCenMulticastDomainsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"transit_router_id": "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_id}",
			"name_regex":        "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"transit_router_id": "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_id}",
			"name_regex":        "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_name}_fake",
		}),
	}
	IdsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"transit_router_id": "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_id}",
			"ids":               []string{"${alibabacloudstack_cen_transit_router_multicast_domain.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"transit_router_id": "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_id}",
			"ids":               []string{"${alibabacloudstack_cen_transit_router_multicast_domain.default.id}_fake"},
		}),
	}

	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"transit_router_id":   "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_id}",
			"description_regex":   "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"transit_router_id":   "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_id}",
			"description_regex":   "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_description}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"transit_router_id":   "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_id}",
			"name_regex":          "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_name}",
			"description_regex":   "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_description}",
			"ids":                 []string{"${alibabacloudstack_cen_transit_router_multicast_domain.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"transit_router_id":   "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_id}",
			"name_regex":          "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_name}",
			"ids":                 []string{"${alibabacloudstack_cen_transit_router_multicast_domain.default.id}"},
			"description_regex":   "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_description}_fake",
		}),
	}

	MulticastDomainsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, IdsConf, descriptionConf, allConf)
}

func dataSourceCenMulticastDomainsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
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
`, name)
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
