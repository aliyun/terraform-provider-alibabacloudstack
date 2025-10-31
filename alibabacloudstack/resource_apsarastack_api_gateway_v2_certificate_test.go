package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAPIGateWayV2Certificate_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_api_gateway_v2_certificate.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeApiGatewayV2Certificate")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 2000)
	name := fmt.Sprintf("testtf-apigw-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AIGateWayV2InstanceK8sDepDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"certificate_name": "${var.name}",
					"instance_id":      "i-iv5kgenx2ddlygr8mpao",
					"cert_type":        0,
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
				// sls_enabled, prometheus_enabled  not read
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func APIGateWayV2CertificateDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "certificates" {
  default = "%s"
}

variable "private_key" {
  default = "%s"
}

variable "certificates2" {
  default = "%s"
}

variable "private_key2" {
  default = "%s"
}

// data "alibabacloudstack_api_gateway_v2_instance_types" "default" {
// 	sorted_by = "CPU"
// }
`, name, ServerCertificateTestCase(), RsaPrivateKeyTestCase(), UpdateCertificateTestCase(), UpdateRsaPrivateKeyTestCase())
}
