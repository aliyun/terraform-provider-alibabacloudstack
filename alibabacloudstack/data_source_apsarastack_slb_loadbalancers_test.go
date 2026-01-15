package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackSlbLoadbalancersDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_slbs.default"
	name := fmt.Sprintf("tf-testacc-slbs%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceSlbsConfigDependence)

	// Test with IDs filter
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":                      []string{"${alibabacloudstack_slb_loadbalancer.default.id}"},
			"master_availability_zone": "${alibabacloudstack_vpc_vswitch.default.zone_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":                     []string{"${alibabacloudstack_slb_loadbalancer.default.id}_fake"},
			"slave_availability_zone": "${alibabacloudstack_vpc_vswitch.default.zone_id}_fake",
		}),
	}

	// Test with name regex filter
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_slb_loadbalancer.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_*",
		}),
	}
	otherconfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"network_type": "${alibabacloudstack_slb_loadbalancer.default.network_type}",
			"vpc_id":       "${alibabacloudstack_vpc_vswitch.default.vpc_id}",
			"vswitch_id":   "${alibabacloudstack_vpc_vswitch.default.id}",
			"address":      "${alibabacloudstack_slb_loadbalancer.default.address}",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Test",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"network_type": "${alibabacloudstack_slb_loadbalancer.default.network_type}_fake",
			"vpc_id":       "${alibabacloudstack_vpc_vswitch.default.vpc_id}_fake",
			"vswitch_id":   "${alibabacloudstack_vpc_vswitch.default.id}_fake",
			"address":      "${alibabacloudstack_slb_loadbalancer.default.address}_fake",
			"tags": map[string]string{
				"Created": "TF_fake",
				"For":     "Test_fake",
			},
		}),
	}

	// Test with both filters
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_slb_loadbalancer.default.id}"},
			"name_regex": "${alibabacloudstack_slb_loadbalancer.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_slb_loadbalancer.default.id}_fake"},
			"name_regex": "${alibabacloudstack_slb_loadbalancer.default.name}_fake",
		}),
	}

	var existSlbsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"slbs.#":                          "1",
			"slbs.0.id":                       CHECKSET,
			"slbs.0.name":                     name,
			"slbs.0.address":                  CHECKSET,
			"slbs.0.address_type":             "intranet",
			"slbs.0.internet_charge_type":     CHECKSET,
			"slbs.0.network_type":             "vpc",
			"slbs.0.vpc_id":                   CHECKSET,
			"slbs.0.vswitch_id":               CHECKSET,
			"slbs.0.master_availability_zone": CHECKSET,
			"slbs.0.load_balancer_spec":       "slb.s1.small",
			"slbs.0.create_time":              CHECKSET,
		}
	}

	var fakeSlbsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":  "0",
			"slbs.#": "0",
		}
	}

	var slbsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existSlbsMapFunc,
		fakeMapFunc:  fakeSlbsMapFunc,
	}
	slbsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, otherconfig, allConf)
}

func dataSourceSlbsConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	
	%s
	
	resource "alibabacloudstack_slb_loadbalancer" "default" {
		name = "${var.name}"
		network_type = "vpc"
		specification = "slb.s1.small"
		address = "172.16.1.3"
		vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
		address_type = "intranet"
		tags = {
			Created = "TF"
			For = "Test"
		}
	}
	`, name, VSwitchCommonTestCase)
}
