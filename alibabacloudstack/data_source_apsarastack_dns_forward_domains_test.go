package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackDnsForwardDomainsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alibabacloudstack_dns_forward_domains.default"

	dsa := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                          "1",
				"forward_domains.#":              "1",
				"forward_domains.0.name":         fmt.Sprintf("tfacc%d.test.", rand),
				"forward_domains.0.remark":       "Created by Terraform",
				"forward_domains.0.forwarders.#": "1",
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":             "0",
				"forward_domains.#": "0",
			}
		},
	}

	// Run tests using customized test method
	dsa.dataSourceTestCheck(t, rand,
		dataSourceTestAccConfig{
			existConfig: AlibabacloudTestAccDnsForwardDomainsConfigDependence(rand, map[string]string{
				"name_regex": `alibabacloudstack_dns_forward_domain.default.name`,
			}),
			fakeConfig: AlibabacloudTestAccDnsForwardDomainsConfigDependence(rand, map[string]string{
				"name_regex": `"^fake-name$"`,
			}),
		},
		dataSourceTestAccConfig{
			existConfig: AlibabacloudTestAccDnsForwardDomainsConfigDependence(rand, map[string]string{
				"ids": `[alibabacloudstack_dns_forward_domain.default.id]`,
			}),
			fakeConfig: AlibabacloudTestAccDnsForwardDomainsConfigDependence(rand, map[string]string{
				"ids": `["nonexistent-id"]`,
			}),
		},
		dataSourceTestAccConfig{
			existConfig: AlibabacloudTestAccDnsForwardDomainsConfigDependence(rand, map[string]string{
				"ids":        `[alibabacloudstack_dns_forward_domain.default.id]`,
				"name_regex": `alibabacloudstack_dns_forward_domain.default.name`,
			}),
			fakeConfig: AlibabacloudTestAccDnsForwardDomainsConfigDependence(rand, map[string]string{
				"ids":        `["nonexistent-id"]`,
				"name_regex": `"^fake-name$"`,
			}),
		},
	)
}

// Dependence template generation function
func AlibabacloudTestAccDnsForwardDomainsConfigDependence(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tfacc%d.test."
}

resource "alibabacloudstack_dns_forward_domain" "default" {
	name   = "${var.name}"
	remark = "Created by Terraform"
	forward_mode = "FORWARD_FIRST"
	forwarders = ["192.168.101.1"]
}

data "alibabacloudstack_dns_forward_domains" "default" {
%s
}
`, rand, strings.Join(pairs, "\n	"))
	return config
}
