package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackUniversalDnsLinesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alibabacloudstack_universal_dns_lines.default"

	// Define the dataSourceAttr with exist and fake check functions
	dsa := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                    "1",
				"lines.#":                  "1",
				"lines.0.name":             fmt.Sprintf("tf-testacc%d", rand),
				"lines.0.priority":         CHECKSET,
				"lines.0.id":               CHECKSET,
				"lines.v4_addresses.#":     "1",
				"lines.v6_addresses.#":     "2",
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

	// Run tests using customized test method
	dsa.dataSourceTestCheck(t, rand,
		dataSourceTestAccConfig{
			existConfig: AlibabacloudTestAccUniversalDnsLinesConfigDependence(rand, map[string]string{
				"name_regex": `alibabacloudstack_universal_dns_line.default.name`,
			}),
			fakeConfig: AlibabacloudTestAccUniversalDnsLinesConfigDependence(rand, map[string]string{
				"name_regex": `"^fake-name$"`,
			}),
		},
		dataSourceTestAccConfig{
			existConfig: AlibabacloudTestAccUniversalDnsLinesConfigDependence(rand, map[string]string{
				"ids": `[alibabacloudstack_universal_dns_line.default.id]`,
			}),
			fakeConfig: AlibabacloudTestAccUniversalDnsLinesConfigDependence(rand, map[string]string{
				"ids": `["nonexistent-id"]`,
			}),
		},
		dataSourceTestAccConfig{
			existConfig: AlibabacloudTestAccUniversalDnsLinesConfigDependence(rand, map[string]string{
				"ids":        `[alibabacloudstack_universal_dns_line.default.id]`,
				"name_regex": `alibabacloudstack_universal_dns_line.default.name`,
			}),
			fakeConfig: AlibabacloudTestAccUniversalDnsLinesConfigDependence(rand, map[string]string{
				"ids":        `["nonexistent-id"]`,
				"name_regex": `"^fake-name$"`,
			}),
		},
	)
}

// Dependence template generation function
func AlibabacloudTestAccUniversalDnsLinesConfigDependence(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testacc%d"
}

resource "alibabacloudstack_universal_dns_line" "default" {
	name   = "${var.name}.example."
	v4_addresses = ["192.168.0.1"]
	v6_addresses = ["2020:148:2:28::", "2020:148:3:28::"]
}

data "alibabacloudstack_universal_dns_lines" "default" {
%s
}
`, rand, strings.Join(pairs, "\n	"))
	return config
}
