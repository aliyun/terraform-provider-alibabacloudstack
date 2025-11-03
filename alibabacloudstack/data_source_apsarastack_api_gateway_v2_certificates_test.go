package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackAPIGateWayV2CertificateDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_api_gateway_v2_certificates.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccAPIGateWayV2Certificate-%d", rand),
		dataSourceAPIGateWayV2CertificateDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_api_gateway_v2_certificate.default.instance_id}",
			"name_regex":  "${alibabacloudstack_api_gateway_v2_certificate.default.certificate_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_api_gateway_v2_certificate.default.instance_id}",
			"name_regex":  "${alibabacloudstack_api_gateway_v2_certificate.default.certificate_name}-fakeTestAcccc",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_api_gateway_v2_certificate.default.instance_id}",
			"ids":         []string{"${alibabacloudstack_api_gateway_v2_certificate.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_api_gateway_v2_certificate.default.instance_id}",
			"ids":         []string{"${alibabacloudstack_api_gateway_v2_certificate.default.id}-fakeTestAcccc"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_api_gateway_v2_certificate.default.instance_id}",
			"name_regex":  "${alibabacloudstack_api_gateway_v2_certificate.default.certificate_name}",
			"ids":         []string{"${alibabacloudstack_api_gateway_v2_certificate.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_api_gateway_v2_certificate.default.instance_id}",
			"name_regex":  "${alibabacloudstack_api_gateway_v2_certificate.default.certificate_name}-fakeTestAcccc",
			"ids":         []string{"${alibabacloudstack_api_gateway_v2_certificate.default.id}-fakeTestAcccc"},
		}),
	}

	var existAPIGateWayV2CertificateMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"ids.0":                           CHECKSET,
			"certificates.#":                  "1",
			"certificates.0.certificate_name": fmt.Sprintf("tf-testAccAPIGateWayV2Certificate-%d", rand),
			"certificates.0.certificate_id":   CHECKSET,
			"certificates.0.cert_type":        CHECKSET,
			"certificates.0.expire_time":      CHECKSET,
			"certificates.0.create_time":      CHECKSET,
		}
	}

	var fakeAPIGateWayV2CertificateMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"certificates.#": "0",
		}
	}

	var APIGateWayV2CertificateCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAPIGateWayV2CertificateMapFunc,
		fakeMapFunc:  fakeAPIGateWayV2CertificateMapFunc,
	}

	APIGateWayV2CertificateCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceAPIGateWayV2CertificateDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alibabacloudstack_api_gateway_v2_instance_types" "default" {
	sorted_by = "CPU"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  	instance_name = "${var.name}"
	node_number = "1"
	instance_class = "${data.alibabacloudstack_api_gateway_v2_instance_types.default.instance_types[0].id}"
	broker_engine_type = "SCG"
	deploy_mode = "custom"
}

variable "certificates" {
  default = %s
}

resource "alibabacloudstack_api_gateway_v2_certificate" "default" {
  	certificate_name = "${var.name}"
	instance_id = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
	cert_type = "1"
	certificates = "${var.certificates}"
}

`, name, ServerCertificateTestCase())
}
