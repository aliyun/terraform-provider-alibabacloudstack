package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackEcsSecurityGroupsDataSource(t *testing.T) {
	// 根据test_meta自动生成的tasecase

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsSecurityGroupsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_securitygroup.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsSecurityGroupsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_securitygroup.default.id}_fake"]`,
		}),
	}

	vpc_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsSecurityGroupsDataSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_ecs_securitygroup.default.id}"]`,
			"vpc_id": `"${alibabacloudstack_ecs_securitygroup.default.vpc_id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsSecurityGroupsDataSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_ecs_securitygroup.default.id}_fake"]`,
			"vpc_id": `"${alibabacloudstack_ecs_securitygroup.default.vpc_id}_fake"`,
		}),
	}

	AlibabacloudstackEcsSecurityGroupsDataCheckInfo.dataSourceTestCheck(t, rand, idsConf, vpc_idConf)
}

var existAlibabacloudstackEcsSecurityGroupsDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#":    "1",
		"groups.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackEcsSecurityGroupsDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#": "0",
	}
}

var AlibabacloudstackEcsSecurityGroupsDataCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_ecs_securitygroups.default",
	existMapFunc: existAlibabacloudstackEcsSecurityGroupsDataMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackEcsSecurityGroupsDataMapFunc,
}

func testAccCheckAlibabacloudstackEcsSecurityGroupsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackEcsSecurityGroups%d"
}

%s

data "alibabacloudstack_ecs_securitygroups" "default" {
%s
}
`, rand, SecurityGroupCommonTestCase, strings.Join(pairs, "\n   "))
	return config
}