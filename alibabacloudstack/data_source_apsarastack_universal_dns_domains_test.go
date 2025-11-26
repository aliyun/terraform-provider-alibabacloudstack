package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackUniversalDnsDomainsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alibabacloudstack_universal_dns_domains.default"

	// Define the dataSourceAttr with exist and fake check functions
	dsa := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                      "1",
				"names.#":                    "1",
				"domains.#":                  "1",
				"domains.0.name":             fmt.Sprintf("tf-testacc%d.example.", rand),
				"domains.0.remark":           "Created by Terraform",
				"domains.0.record_count":     CHECKSET,
				"domains.0.id":               CHECKSET,
				"domains.0.create_timestamp": CHECKSET,
				"domains.0.update_timestamp": CHECKSET,
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":     "0",
				"names.#":   "0",
				"domains.#": "0",
			}
		},
	}

	// Run tests using customized test method
	dsa.dataSourceTestCheck(t, rand,
		dataSourceTestAccConfig{
			existConfig: AlibabacloudTestAccUniversalDnsDomainsConfigDependence(rand, map[string]string{
				"name_regex": `alibabacloudstack_universal_dns_domain.default.name`,
			}),
			fakeConfig: AlibabacloudTestAccUniversalDnsDomainsConfigDependence(rand, map[string]string{
				"name_regex": `"^fake-name$"`,
			}),
		},
		dataSourceTestAccConfig{
			existConfig: AlibabacloudTestAccUniversalDnsDomainsConfigDependence(rand, map[string]string{
				"ids": `[alibabacloudstack_universal_dns_domain.default.id]`,
			}),
			fakeConfig: AlibabacloudTestAccUniversalDnsDomainsConfigDependence(rand, map[string]string{
				"ids": `["nonexistent-id"]`,
			}),
		},
		dataSourceTestAccConfig{
			existConfig: AlibabacloudTestAccUniversalDnsDomainsConfigDependence(rand, map[string]string{
				"ids":        `[alibabacloudstack_universal_dns_domain.default.id]`,
				"name_regex": `alibabacloudstack_universal_dns_domain.default.name`,
			}),
			fakeConfig: AlibabacloudTestAccUniversalDnsDomainsConfigDependence(rand, map[string]string{
				"ids":        `["nonexistent-id"]`,
				"name_regex": `"^fake-name$"`,
			}),
		},
	)
}

// Dependence template generation function
func AlibabacloudTestAccUniversalDnsDomainsConfigDependence(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testacc%d"
}

resource "alibabacloudstack_universal_dns_domain" "default" {
	name   = "${var.name}.example."
	remark = "Created by Terraform"
}

data "alibabacloudstack_universal_dns_domains" "default" {
%s
}
`, rand, strings.Join(pairs, "\n	"))
	return config
}
