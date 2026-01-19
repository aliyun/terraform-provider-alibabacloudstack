package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackSlbServerCertificatesDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_slb_server_certificates.default"
	name := fmt.Sprintf("tf-testAccSlbSerCertDataSource-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlibabacloudStackSlbServerCertificatesDataSourceConfig)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_slb_server_certificate.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_slb_server_certificate.default.name}_fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_slb_server_certificate.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_slb_server_certificate.default.id}_fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_slb_server_certificate.default.id}"},
			"name_regex": "${alibabacloudstack_slb_server_certificate.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_slb_server_certificate.default.id}_fake"},
			"name_regex": "${alibabacloudstack_slb_server_certificate.default.name}",
		}),
	}

	var existDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"certificates.#":                   "1",
			"ids.#":                            "1",
			"names.#":                          "1",
			"certificates.0.id":                CHECKSET,
			"certificates.0.name":              name,
			"certificates.0.fingerprint":       CHECKSET,
			"certificates.0.created_timestamp": CHECKSET,
		}
	}

	var fakeDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"certificates.#": "0",
			"ids.#":          "0",
			"names.#":        "0",
		}
	}

	var slbServerCertificatesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDnsRecordsMapFunc,
		fakeMapFunc:  fakeDnsRecordsMapFunc,
	}

	slbServerCertificatesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)

}

func testAccCheckAlibabacloudStackSlbServerCertificatesDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
resource "alibabacloudstack_slb_server_certificate" "default" {
  name = "${var.name}"
  server_certificate = %s
  private_key = %s
}
`, name, ServerCertificateTestCase(), RsaPrivateKeyTestCase())
}
