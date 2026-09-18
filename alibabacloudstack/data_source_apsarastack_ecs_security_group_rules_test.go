package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackSecurityGroupRulesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf_acc_sg_rule_%d", rand)
	resourceId := "data.alibabacloudstack_security_group_rules.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceSecurityGroupRulesConfigDependence)

	GroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"group_id": "${alibabacloudstack_security_group_rule.default.security_group_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"group_id": "${alibabacloudstack_security_group_rule.default.security_group_id}_fake",
		}),
	}
	DirectionConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"direction": "ingress",
			"group_id":  "${alibabacloudstack_security_group_rule.default.security_group_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"direction": "egress",
			"group_id":  "${alibabacloudstack_security_group_rule.default.security_group_id}",
		}),
	}
	NicTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"nic_type": "intranet",
			"group_id": "${alibabacloudstack_security_group_rule.default.security_group_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"nic_type": "internet",
			"group_id": "${alibabacloudstack_security_group_rule.default.security_group_id}",
		}),
	}
	IpProtocolConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ip_protocol": "tcp",
			"group_id":    "${alibabacloudstack_security_group_rule.default.security_group_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ip_protocol": "udp",
			"group_id":    "${alibabacloudstack_security_group_rule.default.security_group_id}",
		}),
	}
	PolicyConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"policy":   "accept",
			"group_id": "${alibabacloudstack_security_group_rule.default.security_group_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"policy":   "drop",
			"group_id": "${alibabacloudstack_security_group_rule.default.security_group_id}",
		}),
	}

	var exisMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"group_name":          name + "_sg",
			"rules.#":             "1",
			"rules.0.direction":   "ingress",
			"rules.0.ip_protocol": "tcp",
			"rules.0.nic_type":    "intranet",
			"rules.0.port_range":  "22/22",
			"rules.0.policy":      "accept",
		}
	}
	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rules.#": "0",
		}
	}

	var CheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: exisMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}
	CheckInfo.dataSourceTestCheck(t, rand, GroupIdConf, DirectionConf, NicTypeConf, IpProtocolConf, PolicyConf)
}

func dataSourceSecurityGroupRulesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s
`, name, SecurityGroupCommonTestCase)
}
