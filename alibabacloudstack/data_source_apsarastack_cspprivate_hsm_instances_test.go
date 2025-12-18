package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackCspprivateHsmInstancesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	testAcc := dataSourceAttr{
		resourceId: "data.alibabacloudstack_cspprivate_hsm_instances.default",
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"id":                      CHECKSET,
				"ids.#":                   "1",
				"names.#":                 "1",
				"instances.#":             "1",
				"instances.0.instance_id": CHECKSET,
				"instances.0.alias_name":  fmt.Sprintf("test-tf-cspprivate-hsm%d", rand),
				"instances.0.vsm_type":    "gvsm",
				"instances.0.zone_id":     CHECKSET,
				"instances.0.vendor_code": CHECKSET,
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":       "0",
				"names.#":     "0",
				"instances.#": "0",
			}
		},
		Providers: testYunDunProviders(),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAcc.dataSourceCspprivateHsmInstancesConfigDependenceNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_cspprivate_hsm_instance.default.alias_name}"`,
		}),
		fakeConfig: testAcc.dataSourceCspprivateHsmInstancesConfigDependenceNew(rand, map[string]string{
			"name_regex": `"fake-name-regex"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAcc.dataSourceCspprivateHsmInstancesConfigDependenceNew(rand, map[string]string{
			"ids": `["${alibabacloudstack_cspprivate_hsm_instance.default.id}"]`,
		}),
		fakeConfig: testAcc.dataSourceCspprivateHsmInstancesConfigDependenceNew(rand, map[string]string{
			"ids": `["fake-instance-id"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAcc.dataSourceCspprivateHsmInstancesConfigDependenceNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_cspprivate_hsm_instance.default.alias_name}"`,
			"ids":        `["${alibabacloudstack_cspprivate_hsm_instance.default.id}"]`,
		}),
		fakeConfig: testAcc.dataSourceCspprivateHsmInstancesConfigDependenceNew(rand, map[string]string{
			"name_regex": `"fake-name-regex"`,
			"ids":        `["fake-instance-id"]`,
		}),
	}
	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func (dsa *dataSourceAttr) dataSourceCspprivateHsmInstancesConfigDependenceNew(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "test-tf-cspprivate-hsm%d"
}

data "alibabacloudstack_zones" "default" {
  enable_details = true
}

data "alibabacloudstack_cspprivate_hsm_vendors" "default" {
}

resource "alibabacloudstack_cspprivate_hsm_instance" "default" {
	product_code = "${data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.products.0.code}"
	vendor_code = "${data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.code}"
	vsm_type = "gvsm"
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
	alias_name = "${var.name}"
}

data "alibabacloudstack_cspprivate_hsm_instances" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
