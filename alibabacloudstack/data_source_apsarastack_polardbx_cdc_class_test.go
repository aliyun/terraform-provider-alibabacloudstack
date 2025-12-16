package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackPolardbxCdcClassesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_polardbx_cdc_classes.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccPolardbxInstanceTypesDataSource_%d", rand),
		dataSourcePolardbxCdcClassesConfigDependence)

	baseConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${local.polardbx_instance.id}",
			"sorted_by":   "Memory",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "xxxx",

			"sorted_by": "Memory",
		}),
	}

	testAccConfig = dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccPolardbxCdcClassesDataSource_%d", rand),
		dataSourcePolardbxCdcClassesPresetDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${local.polardbx_instance.id}",
			"ids":         []string{"${data.alibabacloudstack_polardbx_cdc_classes.preset.cdc_classes.0.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${local.polardbx_instance.id}",
			"ids":         []string{"xxxxx"},
		}),
	}
	cpuConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${local.polardbx_instance.id}",
			"cpu":         "4",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${local.polardbx_instance.id}",
			"cpu":         "800",
		}),
	}
	memoryConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${local.polardbx_instance.id}",
			"memory":      "8192",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${local.polardbx_instance.id}",
			"memory":      "1",
		}),
	}

	var existPolardbxCdcClassesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                CHECKSET,
			"ids.0":                CHECKSET,
			"cdc_classes.#":        CHECKSET,
			"cdc_classes.0.id":     CHECKSET,
			"cdc_classes.0.cpu":    CHECKSET,
			"cdc_classes.0.memory": CHECKSET,
		}
	}

	var fakePolardbxCdcClassesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":         "0",
			"cdc_classes.#": "0",
		}
	}

	var PolardbxCdcClassesCheckInfo = dataSourceAttr{
		resourceId:        resourceId,
		existMapFunc:      existPolardbxCdcClassesMapFunc,
		fakeMapFunc:       fakePolardbxCdcClassesMapFunc,
		ExternalProviders: testAccExternalProviders,
	}

	PolardbxCdcClassesCheckInfo.dataSourceTestCheck(t, rand, baseConf, idsConf, cpuConf, memoryConf)
}

func dataSourcePolardbxCdcClassesConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}

	%s

	%s

	%s

	 `, name, RandomPasswordTestCase(12, 2), VSwitchCommonTestCase, PolardbxReadOrCreateCommonTestCase())
}

func dataSourcePolardbxCdcClassesPresetDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}

	%s

	%s

	%s
	
	data "alibabacloudstack_polardbx_cdc_classes" "preset" {
		instance_id = local.polardbx_instance.id
	}

	 `, name, RandomPasswordTestCase(12, 2), VSwitchCommonTestCase, PolardbxReadOrCreateCommonTestCase())
}
