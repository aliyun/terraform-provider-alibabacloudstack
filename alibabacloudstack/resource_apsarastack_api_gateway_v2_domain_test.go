package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAPIGateWayV2Domain_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_api_gateway_v2_domain.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeApiGatewayV2Domain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, APIGateWayV2DomainDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{
					"domain":            "${var.name}.com",
					"instance_id":       "${alibabacloudstack_api_gateway_v2_instance.default.id}",
					"protocol":          "HTTPS",
					"certificate_id":    "${alibabacloudstack_api_gateway_v2_certificate.default.certificate_id}",
					"ca_certificate_id": "${alibabacloudstack_api_gateway_v2_certificate.ca.certificate_id}",
					"client_auth":       "1",
					"subject_dn":        "test1",
					"issuer_dn":         "test2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain":         fmt.Sprintf("%s.com", name),
						"protocol":       "HTTPS",
						"certificate_id": CHECKSET,
						"domain_id":      CHECKSET,
						"client_auth":    "1",
						"subject_dn":     CHECKSET,
						"issuer_dn":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol":          "HTTPS",
					"certificate_id":    CHECKSET,
					"client_auth":       "0",
					"ca_certificate_id": REMOVEKEY,
					"subject_dn":        REMOVEKEY,
					"issuer_dn":         REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol":          "HTTPS",
						"certificate_id":    CHECKSET,
						"client_auth":       "0",
						"ca_certificate_id": REMOVEKEY,
						"subject_dn":        REMOVEKEY,
						"issuer_dn":         REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func APIGateWayV2DomainDependence(name string) string {
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

resource "alibabacloudstack_api_gateway_v2_certificate" "ca" {
  	certificate_name = "${var.name}ca"
  	instance_id = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
  	cert_type = "1"
  	certificates = "${var.certificates}"
}

`, name, ServerCertificateTestCase(), RsaPrivateKeyTestCase())
}
