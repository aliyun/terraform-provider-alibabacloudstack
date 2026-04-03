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

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: BmcpMachineTypesDependenceNew(rand, map[string]string{
			"name_regex": `"PG"`,
		}),
		fakeConfig: BmcpMachineTypesDependenceNew(rand, map[string]string{
			"name_regex": `"fake-name"`,
		}),
	}

	archRegexConf := dataSourceTestAccConfig{
		existConfig: BmcpMachineTypesDependenceNew(rand, map[string]string{
			"arch_regex": `"x86"`,
		}),
		fakeConfig: BmcpMachineTypesDependenceNew(rand, map[string]string{
			"arch_regex": `"fake-arch"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: BmcpMachineTypesDependenceNew(rand, map[string]string{
			"name_regex": `"PG"`,
			"arch_regex": `"x86"`,
		}),
		fakeConfig: BmcpMachineTypesDependenceNew(rand, map[string]string{
			"name_regex": `"fake-name"`,
			"arch_regex": `"fake-arch"`,
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, archRegexConf, allConf)
}

// Dependence template generation method
func BmcpMachineTypesDependenceNew(rand int, attrMap map[string]string) string {
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
