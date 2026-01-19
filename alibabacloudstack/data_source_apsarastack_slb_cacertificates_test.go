package alibabacloudstack

import (
	"fmt"

	"testing"
)

func TestAccAlibabacloudStackSlbCACertificatesDataSource_basic(t *testing.T) {
	resourceId:=   "data.alibabacloudstack_slb_ca_certificates.default"
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-testAccSlbCADataSource-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlibabacloudStackSlbCaCertificatesDataSourceConfig)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_slb_ca_certificate.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_slb_ca_certificate.default.name}_fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_slb_ca_certificate.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_slb_ca_certificate.default.id}_fake"},
		}),
	}

	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_slb_ca_certificate.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_slb_ca_certificate.default.id}_fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_slb_ca_certificate.default.id}"},
			"name_regex": "${alibabacloudstack_slb_ca_certificate.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_slb_ca_certificate.default.id}_fake"},
			"name_regex": "${alibabacloudstack_slb_ca_certificate.default.name}",
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
			"certificates.0.region_id":         defaultRegionToTest,
		}
	}

	var fakeDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"certificates.#": "0",
			"ids.#":          "0",
			"names.#":        "0",
		}
	}

	var slbCaCertificatesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDnsRecordsMapFunc,
		fakeMapFunc:  fakeDnsRecordsMapFunc,
	}

	slbCaCertificatesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, resourceGroupIdConf, allConf)

}

func testAccCheckAlibabacloudStackSlbCaCertificatesDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}


resource "alibabacloudstack_slb_ca_certificate" "default" {
  name = "${var.name}"
  ca_certificate = %s
}


`, name, ServerCertificateTestCase())
}
