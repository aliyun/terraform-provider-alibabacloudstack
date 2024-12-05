package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackAlibabacloudstackSlbDomainExtensionsDataSource(t *testing.T) {

	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbDomainExtensionsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_domain_extensions.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbDomainExtensionsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_domain_extensions.default.id}_fake"]`,
		}),
	}

	domain_extension_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbDomainExtensionsSourceConfig(rand, map[string]string{
			"ids":                 `["${alibabacloudstack_slb_domain_extensions.default.id}"]`,
			"domain_extension_id": `"${alibabacloudstack_slb_domain_extensions.default.DomainExtensionId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbDomainExtensionsSourceConfig(rand, map[string]string{
			"ids":                 `["${alibabacloudstack_slb_domain_extensions.default.id}_fake"]`,
			"domain_extension_id": `"${alibabacloudstack_slb_domain_extensions.default.DomainExtensionId}_fake"`,
		}),
	}

	listener_portConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbDomainExtensionsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_domain_extensions.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbDomainExtensionsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_domain_extensions.default.id}_fake"]`,
		}),
	}

	load_balancer_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbDomainExtensionsSourceConfig(rand, map[string]string{
			"ids":              `["${alibabacloudstack_slb_domain_extensions.default.id}"]`,
			"load_balancer_id": `"${alibabacloudstack_slb_domain_extensions.default.LoadBalancerId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbDomainExtensionsSourceConfig(rand, map[string]string{
			"ids":              `["${alibabacloudstack_slb_domain_extensions.default.id}_fake"]`,
			"load_balancer_id": `"${alibabacloudstack_slb_domain_extensions.default.LoadBalancerId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbDomainExtensionsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_domain_extensions.default.id}"]`,

			"domain_extension_id": `"${alibabacloudstack_slb_domain_extensions.default.DomainExtensionId}"`,
			"load_balancer_id":    `"${alibabacloudstack_slb_domain_extensions.default.LoadBalancerId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackSlbDomainExtensionsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_domain_extensions.default.id}_fake"]`,

			"domain_extension_id": `"${alibabacloudstack_slb_domain_extensions.default.DomainExtensionId}_fake"`,
			"load_balancer_id":    `"${alibabacloudstack_slb_domain_extensions.default.LoadBalancerId}_fake"`}),
	}

	AlibabacloudstackSlbDomainExtensionsCheckInfo.dataSourceTestCheck(t, rand, idsConf, domain_extension_idConf, listener_portConf, load_balancer_idConf, allConf)
}

var existAlibabacloudstackSlbDomainExtensionsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"extensions.#":    "1",
		"extensions.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackSlbDomainExtensionsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"extensions.#": "0",
	}
}

var AlibabacloudstackSlbDomainExtensionsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_slb_domain_extensions.default",
	existMapFunc: existAlibabacloudstackSlbDomainExtensionsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackSlbDomainExtensionsMapFunc,
}

func testAccCheckAlibabacloudstackSlbDomainExtensionsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackSlbDomainExtensions%d"
}






data "alibabacloudstack_slb_domain_extensions" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
