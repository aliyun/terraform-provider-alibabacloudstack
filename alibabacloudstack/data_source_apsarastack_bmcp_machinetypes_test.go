package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackBmcpMachineTypesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	testAcc := dataSourceAttr{
		resourceId: "data.alibabacloudstack_bmcp_machinetypes.default",
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
		existConfig: bmcpMachineTypesDependenceNew(rand, map[string]string{
			"name_regex": `"PG"`,
		}),
		fakeConfig: bmcpMachineTypesDependenceNew(rand, map[string]string{
			"name_regex": `"fake-name"`,
		}),
	}

	// Test arch_regex filter
	archRegexConf := dataSourceTestAccConfig{
		existConfig: bmcpMachineTypesDependenceNew(rand, map[string]string{
			"arch_regex": `"x86"`,
		}),
		fakeConfig: bmcpMachineTypesDependenceNew(rand, map[string]string{
			"arch_regex": `"fake-arch"`,
		}),
	}

	// Test min_standard_instance_count filter using dynamic data
	minStandardInstanceCountConf := dataSourceTestAccConfig{
		existConfig: bmcpMachineTypesDependenceWithAllData(rand, map[string]string{
			"min_standard_instance_count": `0`,
		}),
		fakeConfig: bmcpMachineTypesDependenceWithAllData(rand, map[string]string{
			"min_standard_instance_count": `99999`,
		}),
	}

	// Test max_standard_instance_count filter using dynamic data
	maxStandardInstanceCountConf := dataSourceTestAccConfig{
		existConfig: bmcpMachineTypesDependenceWithAllData(rand, map[string]string{
			"max_standard_instance_count": `99999`,
		}),
		fakeConfig: bmcpMachineTypesDependenceWithAllData(rand, map[string]string{
			"max_standard_instance_count": `-1`,
		}),
	}

	// Test combined standard_instance_count range filter using dynamic data
	standardInstanceCountRangeConf := dataSourceTestAccConfig{
		existConfig: bmcpMachineTypesDependenceWithAllData(rand, map[string]string{
			"min_standard_instance_count": `0`,
			"max_standard_instance_count": `99999`,
		}),
		fakeConfig: bmcpMachineTypesDependenceWithAllData(rand, map[string]string{
			"min_standard_instance_count": `99999`,
			"max_standard_instance_count": `999999`,
		}),
	}

	// Test combined filters
	allConf := dataSourceTestAccConfig{
		existConfig: bmcpMachineTypesDependenceWithAllData(rand, map[string]string{
			"name_regex":                  `"PG"`,
			"arch_regex":                  `"x86"`,
			"min_standard_instance_count": `0`,
		}),
		fakeConfig: bmcpMachineTypesDependenceWithAllData(rand, map[string]string{
			"name_regex":                  `"fake-name"`,
			"arch_regex":                  `"fake-arch"`,
			"min_standard_instance_count": `99999`,
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, archRegexConf, 
		minStandardInstanceCountConf, maxStandardInstanceCountConf, 
		standardInstanceCountRangeConf, allConf)
}

// Dependence template generation method for basic filters
func bmcpMachineTypesDependenceNew(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
data "alibabacloudstack_bmcp_machinetypes" "default" {
  %s
}
`, strings.Join(pairs, "\n  "))
}

// Dependence template with pre-query all data for dynamic filtering
func bmcpMachineTypesDependenceWithAllData(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
# First query all machine types to get real data
data "alibabacloudstack_bmcp_machinetypes" "all" {
}

# Use dynamic values from real data for filtering
data "alibabacloudstack_bmcp_machinetypes" "default" {
  %s
}
`, strings.Join(pairs, "\n  "))
}
