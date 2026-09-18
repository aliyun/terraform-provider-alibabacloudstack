package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackECSDeploymentSetsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_ecs_deployment_sets.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccDeploymentSet-%d", rand),
		dataSourceEcsDeploymentSetsConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ecs_deployment_set.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ecs_deployment_set.default.id}_fake"},
		}),
	}
	deploymentSetNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":                 []string{"${alibabacloudstack_ecs_deployment_set.default.id}"},
			"deployment_set_name": "${alibabacloudstack_ecs_deployment_set.default.deployment_set_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":                 []string{"${alibabacloudstack_ecs_deployment_set.default.id}"},
			"deployment_set_name": "${alibabacloudstack_ecs_deployment_set.default.deployment_set_name}_fake",
		}),
	}
	strategyConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":      []string{"${alibabacloudstack_ecs_deployment_set.default.id}"},
			"strategy": "Availability",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":      []string{"${alibabacloudstack_ecs_deployment_set.default.id}_fake"},
			"strategy": "Availability",
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_ecs_deployment_set.default.deployment_set_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_ecs_deployment_set.default.deployment_set_name}_fake",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"deployment_set_name": "${alibabacloudstack_ecs_deployment_set.default.deployment_set_name}",
			"ids":                 []string{"${alibabacloudstack_ecs_deployment_set.default.id}"},
			"name_regex":          "${alibabacloudstack_ecs_deployment_set.default.deployment_set_name}",
			"strategy":            "Availability",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"deployment_set_name": "${alibabacloudstack_ecs_deployment_set.default.deployment_set_name}_fake",
			"ids":                 []string{"${alibabacloudstack_ecs_deployment_set.default.id}_fake"},
			"name_regex":          "${alibabacloudstack_ecs_deployment_set.default.deployment_set_name}_fake",
			"strategy":            "Availability",
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
		resourceId:   resourceId,
		existMapFunc: existDataAlibabacloudStackAlbAclsSourceNameMapFunc,
		fakeMapFunc:  fakeAlibabacloudStackEcsDeploymentSetsDataSourceNameMapFunc,
	}
	alibabacloudstackEcsDeploymentSetsCheckInfo.dataSourceTestCheck(t, rand, idsConf, deploymentSetNameConf, strategyConf, nameRegexConf, allConf)
}

func dataSourceEcsDeploymentSetsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alibabacloudstack_ecs_deployment_set" "default" {
  strategy            = "Availability"
  domain              = "Default"
  granularity         = "Host"
  deployment_set_name = var.name
  description         = var.name
}
`, name)
}
