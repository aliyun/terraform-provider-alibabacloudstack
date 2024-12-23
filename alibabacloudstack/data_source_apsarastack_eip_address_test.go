package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackEipAddressesDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEipAddressesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_eip_addresses.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEipAddressesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_eip_addresses.default.id}_fake"]`,
		}),
	}

	address_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEipAddressesSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_eip_addresses.default.id}"]`,
			"address_name": `"${alibabacloudstack_eip_addresses.default.AddressName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEipAddressesSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_eip_addresses.default.id}_fake"]`,
			"address_name": `"${alibabacloudstack_eip_addresses.default.AddressName}_fake"`,
		}),
	}

	allocation_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEipAddressesSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_eip_addresses.default.id}"]`,
			"allocation_id": `"${alibabacloudstack_eip_addresses.default.AllocationId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEipAddressesSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_eip_addresses.default.id}_fake"]`,
			"allocation_id": `"${alibabacloudstack_eip_addresses.default.AllocationId}_fake"`,
		}),
	}

	instance_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEipAddressesSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_eip_addresses.default.id}"]`,
			"instance_id": `"${alibabacloudstack_eip_addresses.default.InstanceId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEipAddressesSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_eip_addresses.default.id}_fake"]`,
			"instance_id": `"${alibabacloudstack_eip_addresses.default.InstanceId}_fake"`,
		}),
	}

	instance_typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEipAddressesSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_eip_addresses.default.id}"]`,
			"instance_type": `"${alibabacloudstack_eip_addresses.default.InstanceType}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEipAddressesSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_eip_addresses.default.id}_fake"]`,
			"instance_type": `"${alibabacloudstack_eip_addresses.default.InstanceType}_fake"`,
		}),
	}

	ip_addressConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEipAddressesSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_eip_addresses.default.id}"]`,
			"ip_address": `"${alibabacloudstack_eip_addresses.default.IpAddress}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEipAddressesSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_eip_addresses.default.id}_fake"]`,
			"ip_address": `"${alibabacloudstack_eip_addresses.default.IpAddress}_fake"`,
		}),
	}

	ispConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEipAddressesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_eip_addresses.default.id}"]`,
			"isp": `"${alibabacloudstack_eip_addresses.default.Isp}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEipAddressesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_eip_addresses.default.id}_fake"]`,
			"isp": `"${alibabacloudstack_eip_addresses.default.Isp}_fake"`,
		}),
	}

	payment_typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEipAddressesSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_eip_addresses.default.id}"]`,
			"payment_type": `"${alibabacloudstack_eip_addresses.default.PaymentType}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEipAddressesSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_eip_addresses.default.id}_fake"]`,
			"payment_type": `"${alibabacloudstack_eip_addresses.default.PaymentType}_fake"`,
		}),
	}

	resource_group_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEipAddressesSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_eip_addresses.default.id}"]`,
			"resource_group_id": `"${alibabacloudstack_eip_addresses.default.ResourceGroupId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEipAddressesSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_eip_addresses.default.id}_fake"]`,
			"resource_group_id": `"${alibabacloudstack_eip_addresses.default.ResourceGroupId}_fake"`,
		}),
	}

	segment_instance_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEipAddressesSourceConfig(rand, map[string]string{
			"ids":                 `["${alibabacloudstack_eip_addresses.default.id}"]`,
			"segment_instance_id": `"${alibabacloudstack_eip_addresses.default.SegmentInstanceId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEipAddressesSourceConfig(rand, map[string]string{
			"ids":                 `["${alibabacloudstack_eip_addresses.default.id}_fake"]`,
			"segment_instance_id": `"${alibabacloudstack_eip_addresses.default.SegmentInstanceId}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEipAddressesSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_eip_addresses.default.id}"]`,
			"status": `"${alibabacloudstack_eip_addresses.default.Status}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEipAddressesSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_eip_addresses.default.id}_fake"]`,
			"status": `"${alibabacloudstack_eip_addresses.default.Status}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEipAddressesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_eip_addresses.default.id}"]`,

			"address_name":        `"${alibabacloudstack_eip_addresses.default.AddressName}"`,
			"allocation_id":       `"${alibabacloudstack_eip_addresses.default.AllocationId}"`,
			"instance_id":         `"${alibabacloudstack_eip_addresses.default.InstanceId}"`,
			"instance_type":       `"${alibabacloudstack_eip_addresses.default.InstanceType}"`,
			"ip_address":          `"${alibabacloudstack_eip_addresses.default.IpAddress}"`,
			"isp":                 `"${alibabacloudstack_eip_addresses.default.Isp}"`,
			"payment_type":        `"${alibabacloudstack_eip_addresses.default.PaymentType}"`,
			"resource_group_id":   `"${alibabacloudstack_eip_addresses.default.ResourceGroupId}"`,
			"segment_instance_id": `"${alibabacloudstack_eip_addresses.default.SegmentInstanceId}"`,
			"status":              `"${alibabacloudstack_eip_addresses.default.Status}"`}),
		fakeConfig: testAccCheckAlibabacloudstackEipAddressesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_eip_addresses.default.id}_fake"]`,

			"address_name":        `"${alibabacloudstack_eip_addresses.default.AddressName}_fake"`,
			"allocation_id":       `"${alibabacloudstack_eip_addresses.default.AllocationId}_fake"`,
			"instance_id":         `"${alibabacloudstack_eip_addresses.default.InstanceId}_fake"`,
			"instance_type":       `"${alibabacloudstack_eip_addresses.default.InstanceType}_fake"`,
			"ip_address":          `"${alibabacloudstack_eip_addresses.default.IpAddress}_fake"`,
			"isp":                 `"${alibabacloudstack_eip_addresses.default.Isp}_fake"`,
			"payment_type":        `"${alibabacloudstack_eip_addresses.default.PaymentType}_fake"`,
			"resource_group_id":   `"${alibabacloudstack_eip_addresses.default.ResourceGroupId}_fake"`,
			"segment_instance_id": `"${alibabacloudstack_eip_addresses.default.SegmentInstanceId}_fake"`,
			"status":              `"${alibabacloudstack_eip_addresses.default.Status}_fake"`}),
	}

	AlibabacloudstackEipAddressesCheckInfo.dataSourceTestCheck(t, rand, idsConf, address_nameConf, allocation_idConf, instance_idConf, instance_typeConf, ip_addressConf, ispConf, payment_typeConf, resource_group_idConf, segment_instance_idConf, statusConf, allConf)
}

var existAlibabacloudstackEipAddressesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"addresses.#":    "1",
		"addresses.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackEipAddressesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"addresses.#": "0",
	}
}

var AlibabacloudstackEipAddressesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_eip_addresses.default",
	existMapFunc: existAlibabacloudstackEipAddressesMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackEipAddressesMapFunc,
}

func testAccCheckAlibabacloudstackEipAddressesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackEipAddresses%d"
}






data "alibabacloudstack_eip_addresses" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
