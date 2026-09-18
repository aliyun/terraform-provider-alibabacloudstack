package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackbmsMachineTypesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_bms_machinetypes.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccBmsMachineTypes%d", rand),
		dataSourceBmsMachineTypesConfigDependence)

	testAcc := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"machinetypes.#":             CHECKSET,
				"machinetypes.0.id":          CHECKSET,
				"machinetypes.0.name":        CHECKSET,
				"machinetypes.0.deploy_type": CHECKSET,
				"machinetypes.0.cpu_arch":    CHECKSET,
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"machinetypes.#": "0",
			}
		},
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "PG",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake-name",
		}),
	}

	archRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"arch_regex": "x86",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"arch_regex": "fake-arch",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "PG",
			"arch_regex": "x86",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake-name",
			"arch_regex": "fake-arch",
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, archRegexConf, allConf)
}

func dataSourceBmsMachineTypesConfigDependence(name string) string {
	return fmt.Sprintf(``)
}
