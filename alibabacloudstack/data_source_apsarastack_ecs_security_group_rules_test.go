package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackSecurityGroupRulesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf_acc_sg_rule_%d", rand)
	GroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccAlibabacloudStackSecurityGroupRulesDataSourceConfig(name, map[string]string{
			"group_id": `"${alibabacloudstack_security_group_rule.default.security_group_id}"`,
		}),
		fakeConfig: testAccAlibabacloudStackSecurityGroupRulesDataSourceConfig(name, map[string]string{
			"group_id": `"${alibabacloudstack_security_group_rule.default.security_group_id}_fake"`,
		}),
	}
	DirectionConf := dataSourceTestAccConfig{
		existConfig: testAccAlibabacloudStackSecurityGroupRulesDataSourceConfig(name, map[string]string{
			"direction": `"ingress"`,
			"group_id":  `"${alibabacloudstack_security_group_rule.default.security_group_id}"`,
		}),
		fakeConfig: testAccAlibabacloudStackSecurityGroupRulesDataSourceConfig(name, map[string]string{
			"direction": `"egress"`,
			"group_id":  `"${alibabacloudstack_security_group_rule.default.security_group_id}"`,
		}),
	}
	NicTypeConf := dataSourceTestAccConfig{
		existConfig: testAccAlibabacloudStackSecurityGroupRulesDataSourceConfig(name, map[string]string{
			"nic_type": `"intranet"`,
			"group_id": `"${alibabacloudstack_security_group_rule.default.security_group_id}"`,
		}),
		fakeConfig: testAccAlibabacloudStackSecurityGroupRulesDataSourceConfig(name, map[string]string{
			"nic_type": `"internet"`,
			"group_id": `"${alibabacloudstack_security_group_rule.default.security_group_id}"`,
		}),
	}
	IpProtocolConf := dataSourceTestAccConfig{
		existConfig: testAccAlibabacloudStackSecurityGroupRulesDataSourceConfig(name, map[string]string{
			"ip_protocol": `"tcp"`,
			"group_id":    `"${alibabacloudstack_security_group_rule.default.security_group_id}"`,
		}),
		fakeConfig: testAccAlibabacloudStackSecurityGroupRulesDataSourceConfig(name, map[string]string{
			"ip_protocol": `"udp"`,
			"group_id":    `"${alibabacloudstack_security_group_rule.default.security_group_id}"`,
		}),
	}
	PolicyConf := dataSourceTestAccConfig{
		existConfig: testAccAlibabacloudStackSecurityGroupRulesDataSourceConfig(name, map[string]string{
			"policy":   `"accept"`,
			"group_id": `"${alibabacloudstack_security_group_rule.default.security_group_id}"`,
		}),
		fakeConfig: testAccAlibabacloudStackSecurityGroupRulesDataSourceConfig(name, map[string]string{
			"policy":   `"drop"`,
			"group_id": `"${alibabacloudstack_security_group_rule.default.security_group_id}"`,
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
		resourceId:   "data.alibabacloudstack_security_group_rules.default",
		existMapFunc: exisMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}
	CheckInfo.dataSourceTestCheck(t, rand, GroupIdConf, DirectionConf, NicTypeConf, IpProtocolConf, PolicyConf)
}

func testAccAlibabacloudStackSecurityGroupRulesDataSourceConfig(name string, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s

data "alibabacloudstack_security_group_rules" "default" {
  %s
}
`, name, SecurityGroupCommonTestCase, strings.Join(pairs, "\n  "))
	return config
}
