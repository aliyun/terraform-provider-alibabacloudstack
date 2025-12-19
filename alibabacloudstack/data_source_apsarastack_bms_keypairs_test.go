package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackBmsKeypairsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	testAcc := dataSourceAttr{
		resourceId: "data.alibabacloudstack_bms_keypairs.default",
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"keypairs.#":                      CHECKSET,
				"keypairs.0.id":                   CHECKSET,
				"keypairs.0.name":                 CHECKSET,
				"keypairs.0.key_pair_fingerprint": CHECKSET,
				"keypairs.0.create_time":          CHECKSET,
				"keypairs.0.update_time":          CHECKSET,
				"keypairs.0.shared":               CHECKSET,
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"keypairs.#": "0",
			}
		},
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: BmsKeypairDependenceNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_bms_keypair.default.name}"`,
		}),
		fakeConfig: BmsKeypairDependenceNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_bms_keypair.default.name}fake-name"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: BmsKeypairDependenceNew(rand, map[string]string{
			"ids": `["${alibabacloudstack_bms_keypair.default.id}"]`,
		}),
		fakeConfig: BmsKeypairDependenceNew(rand, map[string]string{
			"ids": `["${alibabacloudstack_bms_keypair.default.id}fake-id"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: BmsKeypairDependenceNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_bms_keypair.default.name}"`,
			"ids":        `["${alibabacloudstack_bms_keypair.default.id}"]`,
		}),
		fakeConfig: BmsKeypairDependenceNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_bms_keypair.default.name}fake-name"`,
			"ids":        `["${alibabacloudstack_bms_keypair.default.id}fake-id"]`,
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

// Dependence template generation method
func BmsKeypairDependenceNew(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
variable "name" {
  default = "tf-bms-keypair%d"
}

resource "alibabacloudstack_bms_keypair" "default" {
  name = "${var.name}"
}

data "alibabacloudstack_bms_keypairs" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
}
