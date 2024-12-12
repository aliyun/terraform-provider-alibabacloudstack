package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackEcsDeploymentSetsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000,20000)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackEcsDeploymentSetsDataSourceName(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_deployment_set.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackEcsDeploymentSetsDataSourceName(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_deployment_set.default.id}_fake"]`,
		}),
	}
	deploymentSetNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackEcsDeploymentSetsDataSourceName(rand, map[string]string{
			"ids":                 `["${alibabacloudstack_ecs_deployment_set.default.id}"]`,
			"deployment_set_name": `"${alibabacloudstack_ecs_deployment_set.default.deployment_set_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackEcsDeploymentSetsDataSourceName(rand, map[string]string{
			"ids":                 `["${alibabacloudstack_ecs_deployment_set.default.id}"]`,
			"deployment_set_name": `"${alibabacloudstack_ecs_deployment_set.default.deployment_set_name}_fake"`,
		}),
	}
	strategyConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackEcsDeploymentSetsDataSourceName(rand, map[string]string{
			"ids":      `["${alibabacloudstack_ecs_deployment_set.default.id}"]`,
			"strategy": `"Availability"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackEcsDeploymentSetsDataSourceName(rand, map[string]string{
			"ids":      `["${alibabacloudstack_ecs_deployment_set.default.id}_fake"]`,
			"strategy": `"Availability"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackEcsDeploymentSetsDataSourceName(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_ecs_deployment_set.default.deployment_set_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackEcsDeploymentSetsDataSourceName(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_ecs_deployment_set.default.deployment_set_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackEcsDeploymentSetsDataSourceName(rand, map[string]string{
			"deployment_set_name": `"${alibabacloudstack_ecs_deployment_set.default.deployment_set_name}"`,
			"ids":                 `["${alibabacloudstack_ecs_deployment_set.default.id}"]`,
			"name_regex":          `"${alibabacloudstack_ecs_deployment_set.default.deployment_set_name}"`,
			"strategy":            `"Availability"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackEcsDeploymentSetsDataSourceName(rand, map[string]string{
			"deployment_set_name": `"${alibabacloudstack_ecs_deployment_set.default.deployment_set_name}_fake"`,
			"ids":                 `["${alibabacloudstack_ecs_deployment_set.default.id}_fake"]`,
			"name_regex":          `"${alibabacloudstack_ecs_deployment_set.default.deployment_set_name}_fake"`,
			"strategy":            `"Availability"`,
		}),
	}
	var existDataAlibabacloudStackAlbAclsSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      "1",
			"names.#":                    "1",
			"sets.#":                     "1",
			"sets.0.deployment_set_name": fmt.Sprintf("tf-testAccDeploymentSet-%d", rand),
			"sets.0.strategy":            "Availability",
			"sets.0.domain":              "Default",
			"sets.0.granularity":         "Host",
		}
	}
	var fakeAlibabacloudStackEcsDeploymentSetsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alibabacloudstackEcsDeploymentSetsCheckInfo = dataSourceAttr{
		resourceId:   "data.alibabacloudstack_ecs_deployment_sets.default",
		existMapFunc: existDataAlibabacloudStackAlbAclsSourceNameMapFunc,
		fakeMapFunc:  fakeAlibabacloudStackEcsDeploymentSetsDataSourceNameMapFunc,
	}
	alibabacloudstackEcsDeploymentSetsCheckInfo.dataSourceTestCheck(t, rand, idsConf, deploymentSetNameConf, strategyConf, nameRegexConf, allConf)
}
func testAccCheckAlibabacloudStackEcsDeploymentSetsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccDeploymentSet-%d"
}

resource "alibabacloudstack_ecs_deployment_set" "default" {
  strategy            = "Availability"
  domain              = "default"
  granularity         = "host"
  deployment_set_name = var.name
  description         = var.name
}

data "alibabacloudstack_ecs_deployment_sets" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
