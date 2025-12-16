package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackHsmInstancesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	testAcc := dataSourceAttr{
		resourceId: "data.alibabacloudstack_hsm_instances.default",
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                    "1",
				"names.#":                  "1",
				"instances.#":              "1",
				"instances.0.instance_id":  CHECKSET,
				"instances.0.remark":       fmt.Sprintf("test-tf-hsm-instance%d", rand),
				"instances.0.vsm_type":     "gvsm",
				"instances.0.zone_no":      CHECKSET,
				"instances.0.vendor_code":  CHECKSET,
				"instances.0.product_code": CHECKSET,
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
		existConfig: testAcc.dataSourceHsmInstancesConfigDependenceNew(rand, map[string]string{
			"name_regex": `"test-tf-hsm-instance"`,
			"instance_id": `"${alibabacloudstack_hsm_instance.default.id}"`,
		}),
		fakeConfig: testAcc.dataSourceHsmInstancesConfigDependenceNew(rand, map[string]string{
			"name_regex": `"fake-name-regex"`,
			"instance_id": `"${alibabacloudstack_hsm_instance.default.id}"`,
		}),
	}

	instanceIdConf := dataSourceTestAccConfig{
		existConfig: testAcc.dataSourceHsmInstancesConfigDependenceNew(rand, map[string]string{
			"instance_id": `"${alibabacloudstack_hsm_instance.default.id}"`,
		}),
		fakeConfig: testAcc.dataSourceHsmInstancesConfigDependenceNew(rand, map[string]string{
			"instance_id": `"fake-instance-id"`,
		}),
	}

	vsmTypeConf := dataSourceTestAccConfig{
		existConfig: testAcc.dataSourceHsmInstancesConfigDependenceNew(rand, map[string]string{
			"vsm_type": `"gvsm"`,
			"instance_id": `"${alibabacloudstack_hsm_instance.default.id}"`,
		}),
		fakeConfig: testAcc.dataSourceHsmInstancesConfigDependenceNew(rand, map[string]string{
			"vsm_type": `"evsm"`,
			"instance_id": `"${alibabacloudstack_hsm_instance.default.id}"`,
		}),
	}

	zoneNoConf := dataSourceTestAccConfig{
		existConfig: testAcc.dataSourceHsmInstancesConfigDependenceNew(rand, map[string]string{
			"zone_no": `"${data.alibabacloudstack_zones.default.zones.0.id}"`,
			"instance_id": `"${alibabacloudstack_hsm_instance.default.id}"`,
		}),
		fakeConfig: testAcc.dataSourceHsmInstancesConfigDependenceNew(rand, map[string]string{
			"zone_no": `"fake-zone-no"`,
			"instance_id": `"${alibabacloudstack_hsm_instance.default.id}"`,
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, instanceIdConf, vsmTypeConf, zoneNoConf)
}

func (dsa *dataSourceAttr) dataSourceHsmInstancesConfigDependenceNew(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "test-tf-hsm-instance%d"
}

data "alibabacloudstack_zones" "default" {
  enable_details = true
}

resource "alibabacloudstack_hsm_instance" "default" {
	product_code = "jnta.SJJ1528"
	vendor_code = "jnta"
	vsm_type = "gvsm"
	zone_no = "${data.alibabacloudstack_zones.default.zones.0.id}"
	remark = "${var.name}"
}

data "alibabacloudstack_hsm_instances" "default" {
	status = "${alibabacloudstack_hsm_instance.default.status}"
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
