package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackEcsSecurityGroupsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testAccEcsSg%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(AlibabacloudstackEcsSecurityGroupsDataCheckInfo.resourceId, name, testAccCheckAlibabacloudstackEcsSecurityGroupsDataSourceConfig)
	
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ecs_securitygroup.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ecs_securitygroup.default.id}_fake"},
		}),
	}
	
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${alibabacloudstack_ecs_securitygroup.default.name}$",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_name",
		}),
	}

	vpc_idConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ecs_securitygroup.default.id}"},
			"vpc_id": "${alibabacloudstack_ecs_securitygroup.default.vpc_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ecs_securitygroup.default.id}_fake"},
			"vpc_id": "${alibabacloudstack_ecs_securitygroup.default.vpc_id}_fake",
		}),
	}

	AlibabacloudstackEcsSecurityGroupsDataCheckInfo.dataSourceTestCheck(t, rand, idsConf,nameRegexConf, vpc_idConf)
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

func testAccCheckAlibabacloudstackEcsSecurityGroupsDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s

`, name, SecurityGroupCommonTestCase)
}
