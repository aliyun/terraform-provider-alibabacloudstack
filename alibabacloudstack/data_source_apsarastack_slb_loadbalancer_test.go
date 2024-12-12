package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackSlbLoadBalancersDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_load_balancers.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_load_balancers.default.id}_fake"]`,
		}),
	}

	addressConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":     `["${alibabacloudstack_slb_load_balancers.default.id}"]`,
			"address": `"${alibabacloudstack_slb_load_balancers.default.Address}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":     `["${alibabacloudstack_slb_load_balancers.default.id}_fake"]`,
			"address": `"${alibabacloudstack_slb_load_balancers.default.Address}_fake"`,
		}),
	}

	address_ip_versionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":                `["${alibabacloudstack_slb_load_balancers.default.id}"]`,
			"address_ip_version": `"${alibabacloudstack_slb_load_balancers.default.AddressIpVersion}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":                `["${alibabacloudstack_slb_load_balancers.default.id}_fake"]`,
			"address_ip_version": `"${alibabacloudstack_slb_load_balancers.default.AddressIpVersion}_fake"`,
		}),
	}

	address_typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_slb_load_balancers.default.id}"]`,
			"address_type": `"${alibabacloudstack_slb_load_balancers.default.AddressType}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_slb_load_balancers.default.id}_fake"]`,
			"address_type": `"${alibabacloudstack_slb_load_balancers.default.AddressType}_fake"`,
		}),
	}

	internet_charge_typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":                  `["${alibabacloudstack_slb_load_balancers.default.id}"]`,
			"internet_charge_type": `"${alibabacloudstack_slb_load_balancers.default.InternetChargeType}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":                  `["${alibabacloudstack_slb_load_balancers.default.id}_fake"]`,
			"internet_charge_type": `"${alibabacloudstack_slb_load_balancers.default.InternetChargeType}_fake"`,
		}),
	}

	load_balancer_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":              `["${alibabacloudstack_slb_load_balancers.default.id}"]`,
			"load_balancer_id": `"${alibabacloudstack_slb_load_balancers.default.LoadBalancerId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":              `["${alibabacloudstack_slb_load_balancers.default.id}_fake"]`,
			"load_balancer_id": `"${alibabacloudstack_slb_load_balancers.default.LoadBalancerId}_fake"`,
		}),
	}

	load_balancer_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":                `["${alibabacloudstack_slb_load_balancers.default.id}"]`,
			"load_balancer_name": `"${alibabacloudstack_slb_load_balancers.default.LoadBalancerName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":                `["${alibabacloudstack_slb_load_balancers.default.id}_fake"]`,
			"load_balancer_name": `"${alibabacloudstack_slb_load_balancers.default.LoadBalancerName}_fake"`,
		}),
	}

	master_zone_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_slb_load_balancers.default.id}"]`,
			"master_zone_id": `"${alibabacloudstack_slb_load_balancers.default.MasterZoneId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_slb_load_balancers.default.id}_fake"]`,
			"master_zone_id": `"${alibabacloudstack_slb_load_balancers.default.MasterZoneId}_fake"`,
		}),
	}

	network_typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_slb_load_balancers.default.id}"]`,
			"network_type": `"${alibabacloudstack_slb_load_balancers.default.NetworkType}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_slb_load_balancers.default.id}_fake"]`,
			"network_type": `"${alibabacloudstack_slb_load_balancers.default.NetworkType}_fake"`,
		}),
	}

	payment_typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_slb_load_balancers.default.id}"]`,
			"payment_type": `"${alibabacloudstack_slb_load_balancers.default.PaymentType}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_slb_load_balancers.default.id}_fake"]`,
			"payment_type": `"${alibabacloudstack_slb_load_balancers.default.PaymentType}_fake"`,
		}),
	}

	resource_group_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_slb_load_balancers.default.id}"]`,
			"resource_group_id": `"${alibabacloudstack_slb_load_balancers.default.ResourceGroupId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_slb_load_balancers.default.id}_fake"]`,
			"resource_group_id": `"${alibabacloudstack_slb_load_balancers.default.ResourceGroupId}_fake"`,
		}),
	}

	slave_zone_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_slb_load_balancers.default.id}"]`,
			"slave_zone_id": `"${alibabacloudstack_slb_load_balancers.default.SlaveZoneId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_slb_load_balancers.default.id}_fake"]`,
			"slave_zone_id": `"${alibabacloudstack_slb_load_balancers.default.SlaveZoneId}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_slb_load_balancers.default.id}"]`,
			"status": `"${alibabacloudstack_slb_load_balancers.default.Status}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_slb_load_balancers.default.id}_fake"]`,
			"status": `"${alibabacloudstack_slb_load_balancers.default.Status}_fake"`,
		}),
	}

	vswitch_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_slb_load_balancers.default.id}"]`,
			"vswitch_id": `"${alibabacloudstack_slb_load_balancers.default.VSwitchId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_slb_load_balancers.default.id}_fake"]`,
			"vswitch_id": `"${alibabacloudstack_slb_load_balancers.default.VSwitchId}_fake"`,
		}),
	}

	vpc_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_slb_load_balancers.default.id}"]`,
			"vpc_id": `"${alibabacloudstack_slb_load_balancers.default.VpcId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_slb_load_balancers.default.id}_fake"]`,
			"vpc_id": `"${alibabacloudstack_slb_load_balancers.default.VpcId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_load_balancers.default.id}"]`,

			"address":              `"${alibabacloudstack_slb_load_balancers.default.Address}"`,
			"address_ip_version":   `"${alibabacloudstack_slb_load_balancers.default.AddressIpVersion}"`,
			"address_type":         `"${alibabacloudstack_slb_load_balancers.default.AddressType}"`,
			"internet_charge_type": `"${alibabacloudstack_slb_load_balancers.default.InternetChargeType}"`,
			"load_balancer_id":     `"${alibabacloudstack_slb_load_balancers.default.LoadBalancerId}"`,
			"load_balancer_name":   `"${alibabacloudstack_slb_load_balancers.default.LoadBalancerName}"`,
			"master_zone_id":       `"${alibabacloudstack_slb_load_balancers.default.MasterZoneId}"`,
			"network_type":         `"${alibabacloudstack_slb_load_balancers.default.NetworkType}"`,
			"payment_type":         `"${alibabacloudstack_slb_load_balancers.default.PaymentType}"`,
			"resource_group_id":    `"${alibabacloudstack_slb_load_balancers.default.ResourceGroupId}"`,
			"slave_zone_id":        `"${alibabacloudstack_slb_load_balancers.default.SlaveZoneId}"`,
			"status":               `"${alibabacloudstack_slb_load_balancers.default.Status}"`,
			"vswitch_id":           `"${alibabacloudstack_slb_load_balancers.default.VSwitchId}"`,
			"vpc_id":               `"${alibabacloudstack_slb_load_balancers.default.VpcId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_load_balancers.default.id}_fake"]`,

			"address":              `"${alibabacloudstack_slb_load_balancers.default.Address}_fake"`,
			"address_ip_version":   `"${alibabacloudstack_slb_load_balancers.default.AddressIpVersion}_fake"`,
			"address_type":         `"${alibabacloudstack_slb_load_balancers.default.AddressType}_fake"`,
			"internet_charge_type": `"${alibabacloudstack_slb_load_balancers.default.InternetChargeType}_fake"`,
			"load_balancer_id":     `"${alibabacloudstack_slb_load_balancers.default.LoadBalancerId}_fake"`,
			"load_balancer_name":   `"${alibabacloudstack_slb_load_balancers.default.LoadBalancerName}_fake"`,
			"master_zone_id":       `"${alibabacloudstack_slb_load_balancers.default.MasterZoneId}_fake"`,
			"network_type":         `"${alibabacloudstack_slb_load_balancers.default.NetworkType}_fake"`,
			"payment_type":         `"${alibabacloudstack_slb_load_balancers.default.PaymentType}_fake"`,
			"resource_group_id":    `"${alibabacloudstack_slb_load_balancers.default.ResourceGroupId}_fake"`,
			"slave_zone_id":        `"${alibabacloudstack_slb_load_balancers.default.SlaveZoneId}_fake"`,
			"status":               `"${alibabacloudstack_slb_load_balancers.default.Status}_fake"`,
			"vswitch_id":           `"${alibabacloudstack_slb_load_balancers.default.VSwitchId}_fake"`,
			"vpc_id":               `"${alibabacloudstack_slb_load_balancers.default.VpcId}_fake"`}),
	}

	AlibabacloudstackSlbLoadBalancersCheckInfo.dataSourceTestCheck(t, rand, idsConf, addressConf, address_ip_versionConf, address_typeConf, internet_charge_typeConf, load_balancer_idConf, load_balancer_nameConf, master_zone_idConf, network_typeConf, payment_typeConf, resource_group_idConf, slave_zone_idConf, statusConf, vswitch_idConf, vpc_idConf, allConf)
}

var existAlibabacloudstackSlbLoadBalancersMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"balancers.#":    "1",
		"balancers.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackSlbLoadBalancersMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"balancers.#": "0",
	}
}

var AlibabacloudstackSlbLoadBalancersCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_slb_load_balancers.default",
	existMapFunc: existAlibabacloudstackSlbLoadBalancersMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackSlbLoadBalancersMapFunc,
}

func testAccCheckAlibabacloudstackSlbLoadBalancersSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackSlbLoadBalancers%d"
}






data "alibabacloudstack_slb_load_balancers" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
