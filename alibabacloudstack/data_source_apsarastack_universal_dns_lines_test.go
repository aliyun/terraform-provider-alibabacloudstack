package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackUniversalDnsLinesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_universal_dns_lines.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tfacc%d", rand),
		dataSourceUniversalDnsLinesConfigDependence)

	checkInfo := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                    "1",
				"lines.#":                  "1",
				"lines.0.name":             fmt.Sprintf("tfacc%d", rand),
				"lines.0.priority":         CHECKSET,
				"lines.0.id":               CHECKSET,
				"lines.0.v4_addresses.#":   "1",
				"lines.0.v6_addresses.#":   "2",
				"lines.0.create_timestamp": CHECKSET,
				"lines.0.update_timestamp": CHECKSET,
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":   "0",
				"lines.#": "0",
			}
		},
	}

	checkInfo.dataSourceTestCheck(t, rand,
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_universal_dns_line.default.name}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex": "^fake-name$",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"${alibabacloudstack_universal_dns_line.default.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"nonexistent-id"},
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids":        []string{"${alibabacloudstack_universal_dns_line.default.id}"},
				"name_regex": "${alibabacloudstack_universal_dns_line.default.name}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids":        []string{"nonexistent-id"},
				"name_regex": "^fake-name$",
			}),
		},
	)
}

func dataSourceUniversalDnsLinesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alibabacloudstack_universal_dns_line" "default" {
	name   = "${var.name}"
	v4_addresses = ["192.168.0.1"]
	v6_addresses = ["2020:148:2:28::", "2020:148:3:28::"]
}
`, name)
}
