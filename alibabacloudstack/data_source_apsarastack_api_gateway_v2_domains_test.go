package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackAPIGateWayV2DomainsDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_api_gateway_v2_domains.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf.apigwv2.domain%d", rand),
		dataSourceAPIGateWayV2DomainsDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":  "${alibabacloudstack_api_gateway_v2_domain.default.instance_id}",
			"domain_regex": "${alibabacloudstack_api_gateway_v2_domain.default.domain}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":  "${alibabacloudstack_api_gateway_v2_domain.default.instance_id}",
			"domain_regex": "${alibabacloudstack_api_gateway_v2_domain.default.domain}-fakeTestAcccc",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_api_gateway_v2_domain.default.instance_id}",
			"ids":         []string{"${alibabacloudstack_api_gateway_v2_domain.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_api_gateway_v2_domain.default.instance_id}",
			"ids":         []string{"${alibabacloudstack_api_gateway_v2_domain.default.id}-fakeTestAcccc"},
		}),
	}

	protocolConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_api_gateway_v2_domain.default.instance_id}",
			"protocol":    "HTTPS",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_api_gateway_v2_domain.default.instance_id}",
			"protocol":    "HTTP",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":  "${alibabacloudstack_api_gateway_v2_domain.default.instance_id}",
			"domain_regex": "${alibabacloudstack_api_gateway_v2_domain.default.domain}",
			"protocol":     "HTTPS",
			"ids":          []string{"${alibabacloudstack_api_gateway_v2_domain.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":  "${alibabacloudstack_api_gateway_v2_domain.default.instance_id}",
			"domain_regex": "${alibabacloudstack_api_gateway_v2_domain.default.domain}-fakeTestAcccc",
			"protocol":     "HTTP",
			"ids":          []string{"${alibabacloudstack_api_gateway_v2_domain.default.id}-fakeTestAcccc"},
		}),
	}

	var existAPIGateWayV2DomainsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "1",
			"ids.0":                    CHECKSET,
			"domains.#":                "1",
			"domains.0.domain":         fmt.Sprintf("tf.apigwv2.domain%d.com", rand),
			"domains.0.certificate_id": CHECKSET,
			"domains.0.protocol":       "HTTPS",
			"domains.0.domain_id":      CHECKSET,
		}
	}

	var fakeAPIGateWayV2DomainsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"domains.#": "0",
		}
	}

	var APIGateWayV2DomainsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAPIGateWayV2DomainsMapFunc,
		fakeMapFunc:  fakeAPIGateWayV2DomainsMapFunc,
	}

	APIGateWayV2DomainsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, protocolConf, allConf)
}

func dataSourceAPIGateWayV2DomainsDependence(name string) string {
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

variable "private_key" {
  default = %s
}

resource "alibabacloudstack_api_gateway_v2_certificate" "default" {
  	certificate_name = "${var.name}"
  	instance_id = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
  	cert_type = "0"
  	certificates = "${var.certificates}"
  	private_key = "${var.private_key}"
}

resource "alibabacloudstack_api_gateway_v2_domain" "default" {
  	domain = "${var.name}.com"
	instance_id = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
	protocol = "HTTPS"
	certificate_id = "${alibabacloudstack_api_gateway_v2_certificate.default.certificate_id}"
	client_auth = "0"
}

`, name, ServerCertificateTestCase(), RsaPrivateKeyTestCase())
}
