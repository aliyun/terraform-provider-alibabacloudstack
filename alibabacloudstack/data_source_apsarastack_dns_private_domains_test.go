package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackDnsPrivateDomainsDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alibabacloudstack_dns_private_domains.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc%d.", rand),
		dataSourceDnsPrivateDomainsConfigDependence)

	dsa := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                       "1",
				"domains.#":                   "1",
				"domains.0.name":              fmt.Sprintf("tf-testacc%d.", rand),
				"domains.0.remark":            fmt.Sprintf("tf-testacc%d.", rand),
				"domains.0.caller_uid":        CHECKSET,
				"domains.0.record_count":      "0",
				"domains.0.region_and_vpcs.#": "1",
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":     "0",
				"domains.#": "0",
			}
		},
	}

	// Run tests using customized test method
	dsa.dataSourceTestCheck(t, rand,
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_dns_private_domain.default.name}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex": "^fake-name$",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"${alibabacloudstack_dns_private_domain.default.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"nonexistent-id"},
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"vpc_id": "${alibabacloudstack_vpc_vpc.default.id}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"vpc_id": "fake-id",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids":        []string{"${alibabacloudstack_dns_private_domain.default.id}"},
				"name_regex": "${alibabacloudstack_dns_private_domain.default.name}",
				"vpc_id":     "${alibabacloudstack_vpc_vpc.default.id}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids":        []string{"nonexistent-id"},
				"name_regex": "^fake-name$",
				"vpc_id":     "fake-id",
			}),
		},
	)
}

func dataSourceDnsPrivateDomainsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alibabacloudstack_vpc_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	vpc_name   = "${var.name}_vpc"
}

resource "alibabacloudstack_dns_private_domain" "default" {
	name = var.name
	remark = var.name
	vpc_ids = [alibabacloudstack_vpc_vpc.default.id]
}
`, name)
}
