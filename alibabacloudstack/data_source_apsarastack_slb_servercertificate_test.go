package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackAlibabacloudstackSlbServerCertificatesDataSource(t *testing.T) {

	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbServerCertificatesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_server_certificates.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbServerCertificatesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_server_certificates.default.id}_fake"]`,
		}),
	}

	server_certificate_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbServerCertificatesSourceConfig(rand, map[string]string{
			"ids":                   `["${alibabacloudstack_slb_server_certificates.default.id}"]`,
			"server_certificate_id": `"${alibabacloudstack_slb_server_certificates.default.ServerCertificateId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbServerCertificatesSourceConfig(rand, map[string]string{
			"ids":                   `["${alibabacloudstack_slb_server_certificates.default.id}_fake"]`,
			"server_certificate_id": `"${alibabacloudstack_slb_server_certificates.default.ServerCertificateId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbServerCertificatesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_server_certificates.default.id}"]`,

			"server_certificate_id": `"${alibabacloudstack_slb_server_certificates.default.ServerCertificateId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackSlbServerCertificatesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_server_certificates.default.id}_fake"]`,

			"server_certificate_id": `"${alibabacloudstack_slb_server_certificates.default.ServerCertificateId}_fake"`}),
	}

	AlibabacloudstackSlbServerCertificatesCheckInfo.dataSourceTestCheck(t, rand, idsConf, server_certificate_idConf, allConf)
}

var existAlibabacloudstackSlbServerCertificatesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"certificates.#":    "1",
		"certificates.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackSlbServerCertificatesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"certificates.#": "0",
	}
}

var AlibabacloudstackSlbServerCertificatesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_slb_server_certificates.default",
	existMapFunc: existAlibabacloudstackSlbServerCertificatesMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackSlbServerCertificatesMapFunc,
}

func testAccCheckAlibabacloudstackSlbServerCertificatesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackSlbServerCertificates%d"
}






data "alibabacloudstack_slb_server_certificates" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
