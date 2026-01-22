package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackVpcIpv6AddressesDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_vpc_ipv6_addresses.default"
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-vpcipv6address-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceVpcIpv6AddressesDependence)

	associatedInstanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"associated_instance_id": "${alibabacloudstack_ecs_instance.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"associated_instance_id": "${alibabacloudstack_ecs_instance.default.id}_fake",
		}),
	}
	vswitchIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"vswitch_id": "${alibabacloudstack_vpc_vswitch.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"vswitch_id": "${alibabacloudstack_vpc_vswitch.default.id}_fake",
		}),
	}
	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"vpc_id": "${alibabacloudstack_vpc_vpc.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"vpc_id": "${alibabacloudstack_vpc_vpc.default.id}_fake",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"status": "Available",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"status": "Pending",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"vpc_id":                 "${alibabacloudstack_vpc_vpc.default.id}",
			"vswitch_id":             "${alibabacloudstack_vpc_vswitch.default.id}",
			"associated_instance_id": "${alibabacloudstack_ecs_instance.default.id}",
			"status":                 "Available",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"vpc_id":                 "${alibabacloudstack_vpc_vpc.default.id}_fake",
			"vswitch_id":             "${alibabacloudstack_vpc_vswitch.default.id}_fake",
			"associated_instance_id": "${alibabacloudstack_ecs_instance.default.id}_fake",
			"status":                 "Pending",
		}),
	}
	var existVpcIpv6AddressMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                CHECKSET,
			"ids.0":                                CHECKSET,
			"addresses.#":                          CHECKSET,
			"addresses.0.id":                       CHECKSET,
			"addresses.0.status":                   "Available",
			"addresses.0.associated_instance_id":   CHECKSET,
			"addresses.0.associated_instance_type": CHECKSET,
			"addresses.0.ipv6_address":             CHECKSET,
			"addresses.0.ipv6_address_id":          CHECKSET,
			"addresses.0.ipv6_address_name":        "",
//			"addresses.0.ipv6_gateway_id":          CHECKSET,
			"addresses.0.network_type":             CHECKSET,
			"addresses.0.create_time":              CHECKSET,
			"addresses.0.vswitch_id":               CHECKSET,
			"addresses.0.vpc_id":                   CHECKSET,
		}
	}

	var fakeVpcIpv6AddressMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"addresses.#": "0",
		}
	}

	var VpcIpv6AddressCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existVpcIpv6AddressMapFunc,
		fakeMapFunc:  fakeVpcIpv6AddressMapFunc,
		PreCheck: func() {
			testAccPreCheck(t)
		},
	}

	VpcIpv6AddressCheckInfo.dataSourceTestCheck(t, rand, associatedInstanceIdConf, vswitchIdConf, vpcIdConf, statusConf, allConf)
}

func dataSourceVpcIpv6AddressesDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

`, name, ECSInstanceCommonTestCase)
}
