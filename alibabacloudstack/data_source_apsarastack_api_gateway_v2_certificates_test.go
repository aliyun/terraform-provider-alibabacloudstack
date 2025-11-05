package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackAPIGatewayV2CertificateDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_api_gateway_v2_certificates.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccAPIGatewayV2Certificate-%d", rand),
		dataSourceAPIGatewayV2CertificateDependence)

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

	var existAPIGatewayV2CertificateMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"ids.0":                           CHECKSET,
			"certificates.#":                  "1",
			"certificates.0.certificate_name": fmt.Sprintf("tf-testAccAPIGatewayV2Certificate-%d", rand),
			"certificates.0.certificate_id":   CHECKSET,
			"certificates.0.cert_type":        CHECKSET,
			"certificates.0.expire_time":      CHECKSET,
			"certificates.0.create_time":      CHECKSET,
		}
	}

	var fakeAPIGatewayV2CertificateMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"certificates.#": "0",
		}
	}

	var APIGatewayV2CertificateCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAPIGatewayV2CertificateMapFunc,
		fakeMapFunc:  fakeAPIGatewayV2CertificateMapFunc,
	}

	APIGatewayV2CertificateCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceAPIGatewayV2CertificateDependence(name string) string {
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
