package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackBmcpMachineTypesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_bmcp_machinetypes.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccBmcpMachineTypes%d", rand),
		dataSourceBmcpMachineTypesConfigDependence)

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

	// Test name_regex filter
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "PG",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake-name",
		}),
	}

	// Test arch_regex filter
	archRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"arch_regex": "x86",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"arch_regex": "fake-arch",
		}),
	}

	// Test min_standard_instance_count filter using dynamic data
	minStandardInstanceCountConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"min_standard_instance_count": TfRawString("0"),
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"min_standard_instance_count": TfRawString("99999"),
		}),
	}

	// Test max_standard_instance_count filter using dynamic data
	maxStandardInstanceCountConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"max_standard_instance_count": TfRawString("99999"),
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"max_standard_instance_count": TfRawString("-1"),
		}),
	}

	// Test combined standard_instance_count range filter using dynamic data
	standardInstanceCountRangeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"min_standard_instance_count": TfRawString("0"),
			"max_standard_instance_count": TfRawString("99999"),
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"min_standard_instance_count": TfRawString("99999"),
			"max_standard_instance_count": TfRawString("999999"),
		}),
	}

	// Test combined filters
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":                  "PG",
			"arch_regex":                  "x86",
			"min_standard_instance_count": TfRawString("0"),
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":                  "fake-name",
			"arch_regex":                  "fake-arch",
			"min_standard_instance_count": TfRawString("99999"),
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, archRegexConf,
		minStandardInstanceCountConf, maxStandardInstanceCountConf,
		standardInstanceCountRangeConf, allConf)
}

func dataSourceBmcpMachineTypesConfigDependence(name string) string {
	return fmt.Sprintf(`
# Pre-query all machine types for dynamic filtering reference
data "alibabacloudstack_bmcp_machinetypes" "all" {
}
`)
}
