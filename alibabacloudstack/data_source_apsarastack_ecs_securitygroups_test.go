package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAlibabacloudstackEcsSecurityGroupsDataSource(t *testing.T) {
	// 根据test_meta自动生成的tasecase

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsSecurityGroupsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_security_groups.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsSecurityGroupsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_security_groups.default.id}_fake"]`,
		}),
	}

	resource_group_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsSecurityGroupsDataSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_ecs_security_groups.default.id}"]`,
			"resource_group_id": `"${alibabacloudstack_ecs_security_groups.default.ResourceGroupId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsSecurityGroupsDataSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_ecs_security_groups.default.id}_fake"]`,
			"resource_group_id": `"${alibabacloudstack_ecs_security_groups.default.ResourceGroupId}_fake"`,
		}),
	}

	security_group_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsSecurityGroupsDataSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_ecs_security_groups.default.id}"]`,
			"security_group_id": `"${alibabacloudstack_ecs_security_groups.default.SecurityGroupId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsSecurityGroupsDataSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_ecs_security_groups.default.id}_fake"]`,
			"security_group_id": `"${alibabacloudstack_ecs_security_groups.default.SecurityGroupId}_fake"`,
		}),
	}

	security_group_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsSecurityGroupsDataSourceConfig(rand, map[string]string{
			"ids":                 `["${alibabacloudstack_ecs_security_groups.default.id}"]`,
			"security_group_name": `"${alibabacloudstack_ecs_security_groups.default.SecurityGroupName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsSecurityGroupsDataSourceConfig(rand, map[string]string{
			"ids":                 `["${alibabacloudstack_ecs_security_groups.default.id}_fake"]`,
			"security_group_name": `"${alibabacloudstack_ecs_security_groups.default.SecurityGroupName}_fake"`,
		}),
	}

	security_group_typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsSecurityGroupsDataSourceConfig(rand, map[string]string{
			"ids":                 `["${alibabacloudstack_ecs_security_groups.default.id}"]`,
			"security_group_type": `"${alibabacloudstack_ecs_security_groups.default.SecurityGroupType}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsSecurityGroupsDataSourceConfig(rand, map[string]string{
			"ids":                 `["${alibabacloudstack_ecs_security_groups.default.id}_fake"]`,
			"security_group_type": `"${alibabacloudstack_ecs_security_groups.default.SecurityGroupType}_fake"`,
		}),
	}

	vpc_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsSecurityGroupsDataSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_ecs_security_groups.default.id}"]`,
			"vpc_id": `"${alibabacloudstack_ecs_security_groups.default.VpcId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsSecurityGroupsDataSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_ecs_security_groups.default.id}_fake"]`,
			"vpc_id": `"${alibabacloudstack_ecs_security_groups.default.VpcId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsSecurityGroupsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_security_groups.default.id}"]`,

			"resource_group_id":   `"${alibabacloudstack_ecs_security_groups.default.ResourceGroupId}"`,
			"security_group_id":   `"${alibabacloudstack_ecs_security_groups.default.SecurityGroupId}"`,
			"security_group_name": `"${alibabacloudstack_ecs_security_groups.default.SecurityGroupName}"`,
			"security_group_type": `"${alibabacloudstack_ecs_security_groups.default.SecurityGroupType}"`,
			"vpc_id":              `"${alibabacloudstack_ecs_security_groups.default.VpcId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackEcsSecurityGroupsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_security_groups.default.id}_fake"]`,

			"resource_group_id":   `"${alibabacloudstack_ecs_security_groups.default.ResourceGroupId}_fake"`,
			"security_group_id":   `"${alibabacloudstack_ecs_security_groups.default.SecurityGroupId}_fake"`,
			"security_group_name": `"${alibabacloudstack_ecs_security_groups.default.SecurityGroupName}_fake"`,
			"security_group_type": `"${alibabacloudstack_ecs_security_groups.default.SecurityGroupType}_fake"`,
			"vpc_id":              `"${alibabacloudstack_ecs_security_groups.default.VpcId}_fake"`}),
	}

	AlibabacloudstackEcsSecurityGroupsDataCheckInfo.dataSourceTestCheck(t, rand, idsConf, resource_group_idConf, security_group_idConf, security_group_nameConf, security_group_typeConf, vpc_idConf, allConf)
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
	resourceId:   "data.alibabacloudstack_ecs_security_groups.default",
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






data "alibabacloudstack_ecs_security_groups" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}


func TestAccAlibabacloudStackSecurityGroupsDataSourceBasic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackSecurityGroupsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_security_groups.default"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_groups.default", "groups.#", "1"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_security_groups.default", "ids.#"),
				),
			},
		},
	})
}

const testAccCheckAlibabacloudStackSecurityGroupsDataSourceConfig = `

variable "name" {
  default = "tf-securityGroupdatasource"
}
data "alibabacloudstack_zones" "default" {
	available_resource_creation = "VSwitch"
}
resource "alibabacloudstack_vpc" "vpc" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vswitch" "vswitch" {
  vpc_id            = alibabacloudstack_vpc.vpc.id
  cidr_block        = "172.16.0.0/24"
  availability_zone =  data.alibabacloudstack_zones.default.zones.0.id
  name              = "test45"
}
resource "alibabacloudstack_security_group" "group" {
  name        = var.name
  description = "foo"
  vpc_id      = alibabacloudstack_vpc.vpc.id
}
data "alibabacloudstack_security_groups" "default" {
  ids = [alibabacloudstack_security_group.group.id]
}
`