package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackApigatewayApiGroupsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackApigatewayApiGroupsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_apigateway_api_groups.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackApigatewayApiGroupsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_apigateway_api_groups.default.id}_fake"]`,
		}),
	}

	api_group_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackApigatewayApiGroupsSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_apigateway_api_groups.default.id}"]`,
			"api_group_name": `"${alibabacloudstack_apigateway_api_groups.default.ApiGroupName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackApigatewayApiGroupsSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_apigateway_api_groups.default.id}_fake"]`,
			"api_group_name": `"${alibabacloudstack_apigateway_api_groups.default.ApiGroupName}_fake"`,
		}),
	}

	group_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackApigatewayApiGroupsSourceConfig(rand, map[string]string{
			"ids":      `["${alibabacloudstack_apigateway_api_groups.default.id}"]`,
			"group_id": `"${alibabacloudstack_apigateway_api_groups.default.GroupId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackApigatewayApiGroupsSourceConfig(rand, map[string]string{
			"ids":      `["${alibabacloudstack_apigateway_api_groups.default.id}_fake"]`,
			"group_id": `"${alibabacloudstack_apigateway_api_groups.default.GroupId}_fake"`,
		}),
	}

	instance_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackApigatewayApiGroupsSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_apigateway_api_groups.default.id}"]`,
			"instance_id": `"${alibabacloudstack_apigateway_api_groups.default.InstanceId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackApigatewayApiGroupsSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_apigateway_api_groups.default.id}_fake"]`,
			"instance_id": `"${alibabacloudstack_apigateway_api_groups.default.InstanceId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackApigatewayApiGroupsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_apigateway_api_groups.default.id}"]`,

			"api_group_name": `"${alibabacloudstack_apigateway_api_groups.default.ApiGroupName}"`,
			"group_id":       `"${alibabacloudstack_apigateway_api_groups.default.GroupId}"`,
			"instance_id":    `"${alibabacloudstack_apigateway_api_groups.default.InstanceId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackApigatewayApiGroupsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_apigateway_api_groups.default.id}_fake"]`,

			"api_group_name": `"${alibabacloudstack_apigateway_api_groups.default.ApiGroupName}_fake"`,
			"group_id":       `"${alibabacloudstack_apigateway_api_groups.default.GroupId}_fake"`,
			"instance_id":    `"${alibabacloudstack_apigateway_api_groups.default.InstanceId}_fake"`}),
	}

	AlibabacloudstackApigatewayApiGroupsCheckInfo.dataSourceTestCheck(t, rand, idsConf, api_group_nameConf, group_idConf, instance_idConf, allConf)
}

var existAlibabacloudstackApigatewayApiGroupsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#":    "1",
		"groups.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackApigatewayApiGroupsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#": "0",
	}
}

var AlibabacloudstackApigatewayApiGroupsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_apigateway_api_groups.default",
	existMapFunc: existAlibabacloudstackApigatewayApiGroupsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackApigatewayApiGroupsMapFunc,
}

func testAccCheckAlibabacloudstackApigatewayApiGroupsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackApigatewayApiGroups%d"
}






data "alibabacloudstack_apigateway_api_groups" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
