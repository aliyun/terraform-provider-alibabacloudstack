package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackBmcpSecurityGroupsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)

	name := fmt.Sprintf("tf_testbmcp_sg_%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(AlibabacloudstackBmcpSecurityGroupsDataCheckInfo.resourceId, name, testAccCheckAlibabacloudstackBmcpSecurityGroupsDataSourceConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_bcmp_security_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_bcmp_security_group.default.id}_fake"},
		}),
	}

	nameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"${alibabacloudstack_bcmp_security_group.default.id}"},
			"name": "${alibabacloudstack_bcmp_security_group.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"${alibabacloudstack_bcmp_security_group.default.id}_fake"},
			"name": "${alibabacloudstack_bcmp_security_group.default.name}_fake",
		}),
	}

	name_regexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_bcmp_security_group.default.id}"},
			"name_regex": "${alibabacloudstack_bcmp_security_group.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_bcmp_security_group.default.id}_fake"},
			"name_regex": "${alibabacloudstack_bcmp_security_group.default.name}_fake",
		}),
	}

	vpc_idConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alibabacloudstack_bcmp_security_group.default.id}"},
			"vpc_id": "${alibabacloudstack_bcmp_security_group.default.vpc_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alibabacloudstack_bcmp_security_group.default.id}_fake"},
			"vpc_id": "vpc-fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"${alibabacloudstack_bcmp_security_group.default.id}"},
			"name": "${alibabacloudstack_bcmp_security_group.default.name}"},
		),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"${alibabacloudstack_bcmp_security_group.default.id}_fake"},
			"name": "${alibabacloudstack_bcmp_security_group.default.name}_fake"}),
	}

	AlibabacloudstackBmcpSecurityGroupsDataCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameConf, name_regexConf, vpc_idConf, allConf)
}

var existAlibabacloudstackBmcpSecurityGroupsDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"security_groups.#":       "1",
		"security_groups.0.sg_id": CHECKSET,
	}
}

var fakeAlibabacloudstackBmcpSecurityGroupsDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"security_groups.#": "0",
	}
}

var AlibabacloudstackBmcpSecurityGroupsDataCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_bcmp_security_groups.default",
	existMapFunc: existAlibabacloudstackBmcpSecurityGroupsDataMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackBmcpSecurityGroupsDataMapFunc,
}

func testAccCheckAlibabacloudstackBmcpSecurityGroupsDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alibabacloudstack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_bcmp_security_group" "default" {
	vpc_id = alibabacloudstack_vpc.default.id
	name = var.name
	description = var.name
}

`, name)
}
