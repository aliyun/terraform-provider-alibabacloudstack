package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackAlibabacloudstackEcsDeploymentSetsDataSource(t *testing.T) {

	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDeploymentSetsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_deployment_sets.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDeploymentSetsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_deployment_sets.default.id}_fake"]`,
		}),
	}

	deployment_set_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDeploymentSetsSourceConfig(rand, map[string]string{
			"ids":                 `["${alibabacloudstack_ecs_deployment_sets.default.id}"]`,
			"deployment_set_name": `"${alibabacloudstack_ecs_deployment_sets.default.DeploymentSetName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDeploymentSetsSourceConfig(rand, map[string]string{
			"ids":                 `["${alibabacloudstack_ecs_deployment_sets.default.id}_fake"]`,
			"deployment_set_name": `"${alibabacloudstack_ecs_deployment_sets.default.DeploymentSetName}_fake"`,
		}),
	}

	domainConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDeploymentSetsSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_ecs_deployment_sets.default.id}"]`,
			"domain": `"${alibabacloudstack_ecs_deployment_sets.default.Domain}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDeploymentSetsSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_ecs_deployment_sets.default.id}_fake"]`,
			"domain": `"${alibabacloudstack_ecs_deployment_sets.default.Domain}_fake"`,
		}),
	}

	granularityConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDeploymentSetsSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_ecs_deployment_sets.default.id}"]`,
			"granularity": `"${alibabacloudstack_ecs_deployment_sets.default.Granularity}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDeploymentSetsSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_ecs_deployment_sets.default.id}_fake"]`,
			"granularity": `"${alibabacloudstack_ecs_deployment_sets.default.Granularity}_fake"`,
		}),
	}

	strategyConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDeploymentSetsSourceConfig(rand, map[string]string{
			"ids":      `["${alibabacloudstack_ecs_deployment_sets.default.id}"]`,
			"strategy": `"${alibabacloudstack_ecs_deployment_sets.default.Strategy}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDeploymentSetsSourceConfig(rand, map[string]string{
			"ids":      `["${alibabacloudstack_ecs_deployment_sets.default.id}_fake"]`,
			"strategy": `"${alibabacloudstack_ecs_deployment_sets.default.Strategy}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsDeploymentSetsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_deployment_sets.default.id}"]`,

			"deployment_set_name": `"${alibabacloudstack_ecs_deployment_sets.default.DeploymentSetName}"`,
			"domain":              `"${alibabacloudstack_ecs_deployment_sets.default.Domain}"`,
			"granularity":         `"${alibabacloudstack_ecs_deployment_sets.default.Granularity}"`,
			"strategy":            `"${alibabacloudstack_ecs_deployment_sets.default.Strategy}"`}),
		fakeConfig: testAccCheckAlibabacloudstackEcsDeploymentSetsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_deployment_sets.default.id}_fake"]`,

			"deployment_set_name": `"${alibabacloudstack_ecs_deployment_sets.default.DeploymentSetName}_fake"`,
			"domain":              `"${alibabacloudstack_ecs_deployment_sets.default.Domain}_fake"`,
			"granularity":         `"${alibabacloudstack_ecs_deployment_sets.default.Granularity}_fake"`,
			"strategy":            `"${alibabacloudstack_ecs_deployment_sets.default.Strategy}_fake"`}),
	}

	AlibabacloudstackEcsDeploymentSetsCheckInfo.dataSourceTestCheck(t, rand, idsConf, deployment_set_nameConf, domainConf, granularityConf, strategyConf, allConf)
}

var existAlibabacloudstackEcsDeploymentSetsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"sets.#":    "1",
		"sets.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackEcsDeploymentSetsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"sets.#": "0",
	}
}

var AlibabacloudstackEcsDeploymentSetsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_ecs_deployment_sets.default",
	existMapFunc: existAlibabacloudstackEcsDeploymentSetsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackEcsDeploymentSetsMapFunc,
}

func testAccCheckAlibabacloudstackEcsDeploymentSetsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackEcsDeploymentSets%d"
}






data "alibabacloudstack_ecs_deployment_sets" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
