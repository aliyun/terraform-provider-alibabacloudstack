package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackAlibabacloudstackEcsNetworkInterfacesDataSource(t *testing.T) {

	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_network_interfaces.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_network_interfaces.default.id}_fake"]`,
		}),
	}

	instance_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_ecs_network_interfaces.default.id}"]`,
			"instance_id": `"${alibabacloudstack_ecs_network_interfaces.default.InstanceId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_ecs_network_interfaces.default.id}_fake"]`,
			"instance_id": `"${alibabacloudstack_ecs_network_interfaces.default.InstanceId}_fake"`,
		}),
	}

	network_interface_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig(rand, map[string]string{
			"ids":                  `["${alibabacloudstack_ecs_network_interfaces.default.id}"]`,
			"network_interface_id": `"${alibabacloudstack_ecs_network_interfaces.default.NetworkInterfaceId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig(rand, map[string]string{
			"ids":                  `["${alibabacloudstack_ecs_network_interfaces.default.id}_fake"]`,
			"network_interface_id": `"${alibabacloudstack_ecs_network_interfaces.default.NetworkInterfaceId}_fake"`,
		}),
	}

	network_interface_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig(rand, map[string]string{
			"ids":                    `["${alibabacloudstack_ecs_network_interfaces.default.id}"]`,
			"network_interface_name": `"${alibabacloudstack_ecs_network_interfaces.default.NetworkInterfaceName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig(rand, map[string]string{
			"ids":                    `["${alibabacloudstack_ecs_network_interfaces.default.id}_fake"]`,
			"network_interface_name": `"${alibabacloudstack_ecs_network_interfaces.default.NetworkInterfaceName}_fake"`,
		}),
	}

	primary_ip_addressConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig(rand, map[string]string{
			"ids":                `["${alibabacloudstack_ecs_network_interfaces.default.id}"]`,
			"primary_ip_address": `"${alibabacloudstack_ecs_network_interfaces.default.PrimaryIpAddress}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig(rand, map[string]string{
			"ids":                `["${alibabacloudstack_ecs_network_interfaces.default.id}_fake"]`,
			"primary_ip_address": `"${alibabacloudstack_ecs_network_interfaces.default.PrimaryIpAddress}_fake"`,
		}),
	}

	resource_group_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_ecs_network_interfaces.default.id}"]`,
			"resource_group_id": `"${alibabacloudstack_ecs_network_interfaces.default.ResourceGroupId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_ecs_network_interfaces.default.id}_fake"]`,
			"resource_group_id": `"${alibabacloudstack_ecs_network_interfaces.default.ResourceGroupId}_fake"`,
		}),
	}

	service_managedConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_network_interfaces.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_network_interfaces.default.id}_fake"]`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_ecs_network_interfaces.default.id}"]`,
			"status": `"${alibabacloudstack_ecs_network_interfaces.default.Status}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_ecs_network_interfaces.default.id}_fake"]`,
			"status": `"${alibabacloudstack_ecs_network_interfaces.default.Status}_fake"`,
		}),
	}

	typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig(rand, map[string]string{
			"ids":  `["${alibabacloudstack_ecs_network_interfaces.default.id}"]`,
			"type": `"${alibabacloudstack_ecs_network_interfaces.default.Type}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig(rand, map[string]string{
			"ids":  `["${alibabacloudstack_ecs_network_interfaces.default.id}_fake"]`,
			"type": `"${alibabacloudstack_ecs_network_interfaces.default.Type}_fake"`,
		}),
	}

	vswitch_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_ecs_network_interfaces.default.id}"]`,
			"vswitch_id": `"${alibabacloudstack_ecs_network_interfaces.default.VSwitchId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_ecs_network_interfaces.default.id}_fake"]`,
			"vswitch_id": `"${alibabacloudstack_ecs_network_interfaces.default.VSwitchId}_fake"`,
		}),
	}

	vpc_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_ecs_network_interfaces.default.id}"]`,
			"vpc_id": `"${alibabacloudstack_ecs_network_interfaces.default.VpcId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_ecs_network_interfaces.default.id}_fake"]`,
			"vpc_id": `"${alibabacloudstack_ecs_network_interfaces.default.VpcId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_network_interfaces.default.id}"]`,

			"instance_id":            `"${alibabacloudstack_ecs_network_interfaces.default.InstanceId}"`,
			"network_interface_id":   `"${alibabacloudstack_ecs_network_interfaces.default.NetworkInterfaceId}"`,
			"network_interface_name": `"${alibabacloudstack_ecs_network_interfaces.default.NetworkInterfaceName}"`,
			"primary_ip_address":     `"${alibabacloudstack_ecs_network_interfaces.default.PrimaryIpAddress}"`,
			"resource_group_id":      `"${alibabacloudstack_ecs_network_interfaces.default.ResourceGroupId}"`,
			"status":                 `"${alibabacloudstack_ecs_network_interfaces.default.Status}"`,
			"type":                   `"${alibabacloudstack_ecs_network_interfaces.default.Type}"`,
			"vswitch_id":             `"${alibabacloudstack_ecs_network_interfaces.default.VSwitchId}"`,
			"vpc_id":                 `"${alibabacloudstack_ecs_network_interfaces.default.VpcId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_network_interfaces.default.id}_fake"]`,

			"instance_id":            `"${alibabacloudstack_ecs_network_interfaces.default.InstanceId}_fake"`,
			"network_interface_id":   `"${alibabacloudstack_ecs_network_interfaces.default.NetworkInterfaceId}_fake"`,
			"network_interface_name": `"${alibabacloudstack_ecs_network_interfaces.default.NetworkInterfaceName}_fake"`,
			"primary_ip_address":     `"${alibabacloudstack_ecs_network_interfaces.default.PrimaryIpAddress}_fake"`,
			"resource_group_id":      `"${alibabacloudstack_ecs_network_interfaces.default.ResourceGroupId}_fake"`,
			"status":                 `"${alibabacloudstack_ecs_network_interfaces.default.Status}_fake"`,
			"type":                   `"${alibabacloudstack_ecs_network_interfaces.default.Type}_fake"`,
			"vswitch_id":             `"${alibabacloudstack_ecs_network_interfaces.default.VSwitchId}_fake"`,
			"vpc_id":                 `"${alibabacloudstack_ecs_network_interfaces.default.VpcId}_fake"`}),
	}

	AlibabacloudstackEcsNetworkInterfacesCheckInfo.dataSourceTestCheck(t, rand, idsConf, instance_idConf, network_interface_idConf, network_interface_nameConf, primary_ip_addressConf, resource_group_idConf, service_managedConf, statusConf, typeConf, vswitch_idConf, vpc_idConf, allConf)
}

var existAlibabacloudstackEcsNetworkInterfacesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"interfaces.#":    "1",
		"interfaces.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackEcsNetworkInterfacesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"interfaces.#": "0",
	}
}

var AlibabacloudstackEcsNetworkInterfacesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_ecs_network_interfaces.default",
	existMapFunc: existAlibabacloudstackEcsNetworkInterfacesMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackEcsNetworkInterfacesMapFunc,
}

func testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackEcsNetworkInterfaces%d"
}






data "alibabacloudstack_ecs_network_interfaces" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
