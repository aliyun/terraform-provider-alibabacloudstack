package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackUniversalDnsDomainsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	namePrefix := fmt.Sprintf("tf-testacc%d", rand)
	domainName := namePrefix + ".example."
	resourceId := "data.alibabacloudstack_universal_dns_domains.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		namePrefix,
		dataSourceUniversalDnsDomainsConfigDependence)

	checkInfo := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                      "1",
				"names.#":                    "1",
				"domains.#":                  "1",
				"domains.0.name":             domainName,
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

	checkInfo.dataSourceTestCheck(t, rand,
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"${alibabacloudstack_universal_dns_domain.default.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"nonexistent-id"},
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids":        []string{"${alibabacloudstack_universal_dns_domain.default.id}"},
				"name_regex": namePrefix,
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids":        []string{"${alibabacloudstack_universal_dns_domain.default.id}"},
				"name_regex": "^fake-name$",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids":        []string{"${alibabacloudstack_universal_dns_domain.default.id}"},
				"name_regex": namePrefix,
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids":        []string{"nonexistent-id"},
				"name_regex": "^fake-name$",
			}),
		},
	)
}

func dataSourceUniversalDnsDomainsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alibabacloudstack_universal_dns_domain" "default" {
	name   = "${var.name}.example."
	remark = "Created by Terraform"
}
`, name)
}
