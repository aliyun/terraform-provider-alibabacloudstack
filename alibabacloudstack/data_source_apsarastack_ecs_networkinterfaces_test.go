package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackEcsNetworkInterfacesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf_testecsnic_%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(AlibabacloudstackEcsNetworkInterfacesCheckInfo.resourceId, name, testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ecs_networkinterface.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ecs_networkinterface.default.id}_fake"},
		}),
	}

	instance_idConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"type":"Primary",
			"instance_id": "${alibabacloudstack_ecs_instance.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"type":"Primary",
			"instance_id": "${alibabacloudstack_ecs_instance.default.id}_fake",
		}),
	}

	network_interface_nameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${alibabacloudstack_ecs_networkinterface.default.network_interface_name}$",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_ecs_networkinterface.default.network_interface_name}_fake",
		}),
	}

	primary_ip_addressConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"private_ip": "${alibabacloudstack_ecs_networkinterface.default.primary_ip_address}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"private_ip": "${alibabacloudstack_ecs_networkinterface.default.primary_ip_address}_fake",
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]interface{}{
				"filter": name,
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]interface{}{
				"filter": name + "_fake",
			},
		}),
	}
	security_groupConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"security_group_id": "${alibabacloudstack_ecs_securitygroup.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"security_group_id": "${alibabacloudstack_ecs_securitygroup.default.id}_fake",
		}),
	}

	vswitch_idConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alibabacloudstack_ecs_networkinterface.default.id}"},
			"vswitch_id": "${alibabacloudstack_ecs_networkinterface.default.vswitch_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alibabacloudstack_ecs_networkinterface.default.id}"},
			"vswitch_id": "${alibabacloudstack_ecs_networkinterface.default.vswitch_id}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alibabacloudstack_ecs_networkinterface.default.id}"},
			"instance_id": "${alibabacloudstack_ecs_instance.default.id}",
			"name_regex":  "${alibabacloudstack_ecs_networkinterface.default.network_interface_name}",
			"private_ip":  "${alibabacloudstack_ecs_networkinterface.default.primary_ip_address}",
			"tags": map[string]interface{}{
				"filter": name,
			},
			"security_group_id": "${alibabacloudstack_ecs_securitygroup.default.id}",
			"vswitch_id":        "${alibabacloudstack_ecs_networkinterface.default.vswitch_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alibabacloudstack_ecs_networkinterface.default.id}_fake"},
			"instance_id": "${alibabacloudstack_ecs_instance.default.id}_fake",
			"name_regex":  "${alibabacloudstack_ecs_networkinterface.default.network_interface_name}_fake",
			"private_ip":  "${alibabacloudstack_ecs_networkinterface.default.primary_ip_address}_fake",
			"tags": map[string]interface{}{
				"filter": name + "_fake",
			},
			"security_group_id": "${alibabacloudstack_ecs_securitygroup.default.id}_fake",
			"vswitch_id":        "${alibabacloudstack_ecs_networkinterface.default.vswitch_id}_fake",
		}),
	}

	AlibabacloudstackEcsNetworkInterfacesCheckInfo.dataSourceTestCheck(t, rand, instance_idConf, idsConf, tagsConf, security_groupConf, network_interface_nameConf, primary_ip_addressConf, vswitch_idConf, allConf)
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
	resourceId:   "data.alibabacloudstack_ecs_networkinterfaces.default",
	existMapFunc: existAlibabacloudstackEcsNetworkInterfacesMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackEcsNetworkInterfacesMapFunc,
}

func testAccCheckAlibabacloudstackEcsNetworkInterfacesSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s

resource "alibabacloudstack_ecs_networkinterface" "default" {
  name            = var.name
  vswitch_id      = alibabacloudstack_vpc_vswitch.default.id
  security_groups = [alibabacloudstack_ecs_securitygroup.default.id]
  tags = {
	filter = var.name
  }
}

resource "alibabacloudstack_network_interface_attachment" "default" {
  instance_id          = alibabacloudstack_ecs_instance.default.id
  network_interface_id = alibabacloudstack_ecs_networkinterface.default.id
}

`, name, ECSInstanceCommonTestCase)
}
