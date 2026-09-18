package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackBmsKeypairsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_bms_keypairs.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-bms-keypair%d", rand),
		dataSourceBmsKeypairsConfigDependence)

	testAcc := dataSourceAttr{
		resourceId: resourceId,
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
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_bms_keypair.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_bms_keypair.default.name}fake-name",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_bms_keypair.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_bms_keypair.default.id}fake-id"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_bms_keypair.default.name}",
			"ids":        []string{"${alibabacloudstack_bms_keypair.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_bms_keypair.default.name}fake-name",
			"ids":        []string{"${alibabacloudstack_bms_keypair.default.id}fake-id"},
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceBmsKeypairsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_bms_keypair" "default" {
  name = "${var.name}"
}
`, name)
}
