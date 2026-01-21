package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackSlbDomainExtension0(t *testing.T) {
	var v *slb.DescribeDomainExtensionAttributeResponse

	resourceId := "alibabacloudstack_slb_domainextension.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbDomainextensionCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoSlbDescribedomainextensionattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslbdomain_extension%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbDomainextensionBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"server_certificate_id": "${alibabacloudstack_slb_listener.new.server_certificate_id}",

					"listener_port":                "${alibabacloudstack_slb_listener.new.frontend_port}",
					"delete_protection_validation": "true",
					"load_balancer_id":             "${alibabacloudstack_slb_listener.new.load_balancer_id}",

					"domain": "test.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"server_certificate_id": CHECKSET,
						"listener_port":         CHECKSET,

						"load_balancer_id": CHECKSET,

						"domain": "test.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"server_certificate_id": "${alibabacloudstack_slb_server_certificate.servercertificate.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"server_certificate_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
		},
	})
}

var AlibabacloudTestAccSlbDomainextensionCheckmap = map[string]string{

	"domain_extension_id": CHECKSET,

	"listener_port": CHECKSET,

	"server_certificate_id": CHECKSET,

	"load_balancer_id": CHECKSET,

	"domain": CHECKSET,
}

func AlibabacloudTestAccSlbDomainextensionBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

%s

resource "alibabacloudstack_slb_server_certificate" "default" {
	name = "${var.name}"
	server_certificate = %s
	private_key = %s
  }

resource "alibabacloudstack_slb_server_certificate" "servercertificate" {
  name               = "slbservercertificate"
  server_certificate = %s
  private_key        = %s
}


resource "alibabacloudstack_slb_loadbalancer" "default" {
	name = "${var.name}"
	//address_type       = "internet"
  	specification        = "slb.s2.small"
  }


resource "alibabacloudstack_slb_server_group" "default" {
  vserver_group_name = "${var.name}"
  load_balancer_id = "${alibabacloudstack_slb_loadbalancer.default.id}"
  servers {
      server_ids = ["${alibabacloudstack_ecs_instance.default.id}"]
      port = 80
      weight = 100
    }
}

resource "alibabacloudstack_slb_listener" "new" {
  load_balancer_id        = "${alibabacloudstack_slb_loadbalancer.default.id}"
  bandwidth               = "10"
  frontend_port           = "443"
  backend_port            = "80"
  sticky_session          = "off"
  health_check            = "off"
  protocol                = "https"
  server_certificate_id   = "${alibabacloudstack_slb_server_certificate.default.id}"
}


`, name, ECSInstanceCommonTestCase, ServerCertificateTestCase(), RsaPrivateKeyTestCase(), ServerCertificateTestCase(), RsaPrivateKeyTestCase())
}
