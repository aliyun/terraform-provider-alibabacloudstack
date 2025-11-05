package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAPIGatewayV2Certificate_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_api_gateway_v2_certificate.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeApiGatewayV2Certificate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, APIGatewayV2CertificateDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{
					"certificate_name": "${var.name}",
					"instance_id":      "${alibabacloudstack_api_gateway_v2_instance.default.id}",
					"cert_type":        "0",
					"certificates":     "${var.certificates}",
					"private_key":      "${var.private_key}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"certificate_name": name,
						"cert_type":        "0",
						"expire_time":      CHECKSET,
						"create_time":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"certificates": "${var.certificates2}",
					"private_key":  "${var.private_key2}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificates", "private_key"},
			},
		},
	})
}

func APIGatewayV2CertificateDependence(name string) string {
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

variable "certificates2" {
  default = %s
}

variable "private_key2" {
  default = %s
}

`, name, ServerCertificateTestCase(), RsaPrivateKeyTestCase(), UpdateCertificateTestCase(), UpdateRsaPrivateKeyTestCase())
}
