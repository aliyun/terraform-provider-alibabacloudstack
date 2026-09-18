package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackDnsForwardDomainsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alibabacloudstack_dns_forward_domains.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tfacc%d.test.", rand),
		dataSourceDnsForwardDomainsConfigDependence)

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
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_dns_forward_domain.default.name}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex": "^fake-name$",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"${alibabacloudstack_dns_forward_domain.default.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"nonexistent-id"},
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids":        []string{"${alibabacloudstack_dns_forward_domain.default.id}"},
				"name_regex": "${alibabacloudstack_dns_forward_domain.default.name}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids":        []string{"nonexistent-id"},
				"name_regex": "^fake-name$",
			}),
		},
	)
}

func dataSourceDnsForwardDomainsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alibabacloudstack_dns_forward_domain" "default" {
	name   = "${var.name}"
	remark = "Created by Terraform"
	forward_mode = "FORWARD_FIRST"
	forwarders = ["192.168.101.1"]
}
`, name)
}
