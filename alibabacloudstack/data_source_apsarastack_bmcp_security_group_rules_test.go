package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackBmcpSecurityGroupRulesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_bcmp_security_group_rules.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_acc_bmcp_sg_rule_%d", rand),
		dataSourceBmcpSecurityGroupRulesConfigDependence)

	SecurityGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"type":            "ingress",
			"security_group_id": "${alibabacloudstack_bcmp_security_group_rule.default.security_group_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"type":            "ingress",
			"security_group_id": "${alibabacloudstack_bcmp_security_group_rule.default.security_group_id}_fake",
		}),
	}
	IpProtocolConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"type":            "ingress",
			"ip_protocol":     "tcp",
			"security_group_id": "${alibabacloudstack_bcmp_security_group_rule.default.security_group_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"type":            "ingress",
			"ip_protocol":     "udp",
			"security_group_id": "${alibabacloudstack_bcmp_security_group_rule.default.security_group_id}",
		}),
	}
	PolicyConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"type":            "ingress",
			"policy":          "accept",
			"security_group_id": "${alibabacloudstack_bcmp_security_group_rule.default.security_group_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"type":            "ingress",
			"policy":          "drop",
			"security_group_id": "${alibabacloudstack_bcmp_security_group_rule.default.security_group_id}",
		}),
	}

	var exisMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rules.#":             "1",
			"rules.0.type":        "ingress",
			"rules.0.ip_protocol": "tcp",
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
	CheckInfo.dataSourceTestCheck(t, rand, SecurityGroupIdConf, IpProtocolConf, PolicyConf)
}

func dataSourceBmcpSecurityGroupRulesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s
`, name, BmcpSecurityGroupCommonTestCase)
}

var BmcpSecurityGroupCommonTestCase = `
resource "alibabacloudstack_vpc" "default" {
	name = "tf_acc_bmcp_sg_${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_bcmp_security_group" "default" {
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  name = "tf_acc_bmcp_sg1_${var.name}"
}

resource "alibabacloudstack_bcmp_security_group_rule" "default" {
  type              = "ingress"
  ip_protocol       = "tcp"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = "${alibabacloudstack_bcmp_security_group.default.id}"
  cidr_ip           = "0.0.0.0/0"
  description       = "test_${var.name}"
}
`
