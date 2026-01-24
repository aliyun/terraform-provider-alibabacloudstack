package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackDnsDomainsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_dns_domains.default"
	name := fmt.Sprintf("tftestdomain%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceDnsDomainsConfigDependence)

	domainNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alibabacloudstack_dns_domain.default.domain_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "fake-domain.com.",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_dns_domain.default.domain_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"fake-domain-id"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alibabacloudstack_dns_domain.default.domain_name}",
			"ids":         []string{"${alibabacloudstack_dns_domain.default.domain_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "fake-domain.com.",
			"ids":         []string{"fake-domain-id"},
		}),
	}

	var existDnsDomainsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"domains.#":        "1",
			"ids.#":            "1",
			"names.#":          "1",
			"domains.0.domain_id":   CHECKSET,
			"domains.0.domain_name": name+".",
			"domains.0.dns_servers.#": CHECKSET,
		}
	}

	var fakeDnsDomainsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"domains.#": "0",
			"ids.#":     "0",
			"names.#":   "0",
		}
	}

	var dnsDomainsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDnsDomainsMapFunc,
		fakeMapFunc:  fakeDnsDomainsMapFunc,
	}
	dnsDomainsCheckInfo.dataSourceTestCheck(t, rand, domainNameConf, idsConf, allConf)
}

func dataSourceDnsDomainsConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alibabacloudstack_dns_domain" "default" {
  domain_name = "%s."
}
`, name)
}
