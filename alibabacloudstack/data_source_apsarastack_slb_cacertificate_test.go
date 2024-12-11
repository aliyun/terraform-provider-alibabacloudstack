package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackSlbCaCertificatesDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbCaCertificatesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_ca_certificates.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbCaCertificatesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_ca_certificates.default.id}_fake"]`,
		}),
	}

	ca_certificate_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbCaCertificatesSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_slb_ca_certificates.default.id}"]`,
			"ca_certificate_id": `"${alibabacloudstack_slb_ca_certificates.default.CaCertificateId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbCaCertificatesSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_slb_ca_certificates.default.id}_fake"]`,
			"ca_certificate_id": `"${alibabacloudstack_slb_ca_certificates.default.CaCertificateId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbCaCertificatesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_ca_certificates.default.id}"]`,

			"ca_certificate_id": `"${alibabacloudstack_slb_ca_certificates.default.CaCertificateId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackSlbCaCertificatesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_ca_certificates.default.id}_fake"]`,

			"ca_certificate_id": `"${alibabacloudstack_slb_ca_certificates.default.CaCertificateId}_fake"`}),
	}

	AlibabacloudstackSlbCaCertificatesCheckInfo.dataSourceTestCheck(t, rand, idsConf, ca_certificate_idConf, allConf)
}

var existAlibabacloudstackSlbCaCertificatesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"certificates.#":    "1",
		"certificates.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackSlbCaCertificatesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"certificates.#": "0",
	}
}

var AlibabacloudstackSlbCaCertificatesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_slb_ca_certificates.default",
	existMapFunc: existAlibabacloudstackSlbCaCertificatesMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackSlbCaCertificatesMapFunc,
}

func testAccCheckAlibabacloudstackSlbCaCertificatesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackSlbCaCertificates%d"
}






data "alibabacloudstack_slb_ca_certificates" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
