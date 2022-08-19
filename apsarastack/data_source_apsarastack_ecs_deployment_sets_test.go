package apsarastack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccApsaraStackEcsDeploymentSetsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackEcsDeploymentSetsDataSourceName(rand, map[string]string{
			"ids": `["${apsarastack_ecs_deployment_set.default.id}"]`,
		}),
		fakeConfig: testAccCheckApsaraStackEcsDeploymentSetsDataSourceName(rand, map[string]string{
			"ids": `["${apsarastack_ecs_deployment_set.default.id}_fake"]`,
		}),
	}
	deploymentSetNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackEcsDeploymentSetsDataSourceName(rand, map[string]string{
			"ids":                 `["${apsarastack_ecs_deployment_set.default.id}"]`,
			"deployment_set_name": `"${apsarastack_ecs_deployment_set.default.deployment_set_name}"`,
		}),
		fakeConfig: testAccCheckApsaraStackEcsDeploymentSetsDataSourceName(rand, map[string]string{
			"ids":                 `["${apsarastack_ecs_deployment_set.default.id}"]`,
			"deployment_set_name": `"${apsarastack_ecs_deployment_set.default.deployment_set_name}_fake"`,
		}),
	}
	strategyConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackEcsDeploymentSetsDataSourceName(rand, map[string]string{
			"ids":      `["${apsarastack_ecs_deployment_set.default.id}"]`,
			"strategy": `"Availability"`,
		}),
		fakeConfig: testAccCheckApsaraStackEcsDeploymentSetsDataSourceName(rand, map[string]string{
			"ids":      `["${apsarastack_ecs_deployment_set.default.id}_fake"]`,
			"strategy": `"Availability"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackEcsDeploymentSetsDataSourceName(rand, map[string]string{
			"name_regex": `"${apsarastack_ecs_deployment_set.default.deployment_set_name}"`,
		}),
		fakeConfig: testAccCheckApsaraStackEcsDeploymentSetsDataSourceName(rand, map[string]string{
			"name_regex": `"${apsarastack_ecs_deployment_set.default.deployment_set_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackEcsDeploymentSetsDataSourceName(rand, map[string]string{
			"deployment_set_name": `"${apsarastack_ecs_deployment_set.default.deployment_set_name}"`,
			"ids":                 `["${apsarastack_ecs_deployment_set.default.id}"]`,
			"name_regex":          `"${apsarastack_ecs_deployment_set.default.deployment_set_name}"`,
			"strategy":            `"Availability"`,
		}),
		fakeConfig: testAccCheckApsaraStackEcsDeploymentSetsDataSourceName(rand, map[string]string{
			"deployment_set_name": `"${apsarastack_ecs_deployment_set.default.deployment_set_name}_fake"`,
			"ids":                 `["${apsarastack_ecs_deployment_set.default.id}_fake"]`,
			"name_regex":          `"${apsarastack_ecs_deployment_set.default.deployment_set_name}_fake"`,
			"strategy":            `"Availability"`,
		}),
	}
	var existDataApsaraStackAlbAclsSourceNameMapFunc = func(rand int) map[string]string {
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
	var fakeApsaraStackEcsDeploymentSetsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var apsarastackEcsDeploymentSetsCheckInfo = dataSourceAttr{
		resourceId:   "data.apsarastack_ecs_deployment_sets.default",
		existMapFunc: existDataApsaraStackAlbAclsSourceNameMapFunc,
		fakeMapFunc:  fakeApsaraStackEcsDeploymentSetsDataSourceNameMapFunc,
	}
	apsarastackEcsDeploymentSetsCheckInfo.dataSourceTestCheck(t, rand, idsConf, deploymentSetNameConf, strategyConf, nameRegexConf, allConf)
}
func testAccCheckApsaraStackEcsDeploymentSetsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccDeploymentSet-%d"
}

resource "apsarastack_ecs_deployment_set" "default" {
  strategy            = "Availability"
  domain              = "default"
  granularity         = "host"
  deployment_set_name = var.name
  description         = var.name
}

data "apsarastack_ecs_deployment_sets" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
