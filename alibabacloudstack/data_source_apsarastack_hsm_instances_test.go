package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackHsmInstancesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_hsm_instances.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("test-tf-hsm-instance%d", rand),
		dataSourceHsmInstancesConfigDependence)

	testAcc := dataSourceAttr{
		resourceId: resourceId,
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

	testAcc.dataSourceTestCheck(t, rand,
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex":  "test-tf-hsm-instance",
				"instance_id": "${alibabacloudstack_hsm_instance.default.id}",
				"status":      "${alibabacloudstack_hsm_instance.default.status}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex":  "fake-name-regex",
				"instance_id": "${alibabacloudstack_hsm_instance.default.id}",
				"status":      "${alibabacloudstack_hsm_instance.default.status}",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"instance_id": "${alibabacloudstack_hsm_instance.default.id}",
				"status":      "${alibabacloudstack_hsm_instance.default.status}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"instance_id": "fake-instance-id",
				"status":      "${alibabacloudstack_hsm_instance.default.status}",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"vsm_type":    "gvsm",
				"instance_id": "${alibabacloudstack_hsm_instance.default.id}",
				"status":      "${alibabacloudstack_hsm_instance.default.status}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"vsm_type":    "evsm",
				"instance_id": "${alibabacloudstack_hsm_instance.default.id}",
				"status":      "${alibabacloudstack_hsm_instance.default.status}",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"zone_no":     "${data.alibabacloudstack_zones.default.zones.0.id}",
				"instance_id": "${alibabacloudstack_hsm_instance.default.id}",
				"status":      "${alibabacloudstack_hsm_instance.default.status}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"zone_no":     "fake-zone-no",
				"instance_id": "${alibabacloudstack_hsm_instance.default.id}",
				"status":      "${alibabacloudstack_hsm_instance.default.status}",
			}),
		},
	)
}

func dataSourceHsmInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alibabacloudstack_zones" "default" {
	enable_details = true
}

resource "alibabacloudstack_hsm_instance" "default" {
	product_code = "jnta.SJJ1528"
	vendor_code  = "jnta"
	vsm_type     = "gvsm"
	zone_no      = "${data.alibabacloudstack_zones.default.zones.0.id}"
	remark       = "${var.name}"
}
`, name)
}
