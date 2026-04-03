package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackBmcpImagesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	testAcc := dataSourceAttr{
		resourceId: "data.alibabacloudstack_bmcp_images.default",
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"images.#":        CHECKSET,
				"images.0.id":     CHECKSET,
				"images.0.name":   CHECKSET,
				"images.0.type":   CHECKSET,
				"images.0.status": CHECKSET,
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"images.#": "0",
			}
		},
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: BmcpImagesDependenceNew(rand, map[string]string{
			"name_regex": `"ubuntu"`,
		}),
		fakeConfig: BmcpImagesDependenceNew(rand, map[string]string{
			"name_regex": `"fake-name"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: BmcpImagesDependenceNew(rand, map[string]string{
			"ids": `["i-d"]`,
		}),
		fakeConfig: BmcpImagesDependenceNew(rand, map[string]string{
			"ids": `["fake-id"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: BmcpImagesDependenceNew(rand, map[string]string{
			"name_regex": `"ubuntu"`,
			"ids":        `["i-d"]`,
		}),
		fakeConfig: BmcpImagesDependenceNew(rand, map[string]string{
			"name_regex": `"fake-name"`,
			"ids":        `["fake-id"]`,
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

// Dependence template generation method
func BmcpImagesDependenceNew(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
data "alibabacloudstack_bmcp_images" "default" {
  type = "public"
  %s
}
`, strings.Join(pairs, "\n  "))
}
