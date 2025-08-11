package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackPolardbxInstancesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_polardbx_instances.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccPolardbxInstancesDataSource-%d", rand),
		dataSourcePolardbxInstancesDependence)

	descriptionRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_polardbx_instance.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_polardbx_instance.default.description}-fakeTestAcccc",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_polardbx_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_polardbx_instance.default.id}-fakeTestAcccc"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_polardbx_instance.default.description}",
			"ids":        []string{"${alibabacloudstack_polardbx_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_polardbx_instance.default.description}-fakeTestAcccc",
			"ids":        []string{"${alibabacloudstack_polardbx_instance.default.id}-fakeTestAcccc"},
		}),
	}

	var existPolardbxInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"ids.0":                               CHECKSET,
			"polardbx_instances.#":                "1",
			"polardbx_instances.0.description":    CHECKSET,
			"polardbx_instances.0.storage":        CHECKSET,
			"polardbx_instances.0.cpu_type":       CHECKSET,
			"polardbx_instances.0.cn_node_class":  CHECKSET,
			"polardbx_instances.0.cn_node_count":  CHECKSET,
			"polardbx_instances.0.dn_node_class":  CHECKSET,
			"polardbx_instances.0.dn_node_count":  CHECKSET,
			"polardbx_instances.0.engine_version": CHECKSET,
			"polardbx_instances.0.resource_type":  CHECKSET,
		}
	}

	var fakePolardbxInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                "0",
			"polardbx_instances.#": "0",
		}
	}

	var PolardbxInstancesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existPolardbxInstancesMapFunc,
		fakeMapFunc:  fakePolardbxInstancesMapFunc,
	}

	PolardbxInstancesCheckInfo.dataSourceTestCheck(t, rand, descriptionRegexConf, idsConf, allConf)
}

func dataSourcePolardbxInstancesDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}

%s

%s
 `, name,VSwitchCommonTestCase, PolardbxCommonTestCase)
}
