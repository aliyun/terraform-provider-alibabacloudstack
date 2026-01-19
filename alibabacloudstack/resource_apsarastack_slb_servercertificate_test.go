package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackSlbServerCertificate0(t *testing.T) {
	var v *slb.ServerCertificate

	resourceId := "alibabacloudstack_slb_servercertificate.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbServercertificateCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoSlbDescribeservercertificateRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslbserver_certificate%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbServercertificateBasicdependence)
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

					"server_certificate_name": "test-cert-name",
					"server_certificate":      TfRawString(ServerCertificateTestCase()),
					"private_key":             TfRawString(RsaPrivateKeyTestCase()),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"server_certificate_name": "test-cert-name",
						"server_certificate":      CHECKSET,
						"private_key":             CHECKSET,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"server_certificate_name": "test-cert-name-2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"server_certificate_name": "test-cert-name-2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// server_certificate and private_key are not echoed after being set
				ImportStateVerifyIgnore: []string{"server_certificate", "private_key"},
			},
		},
	})
}

var AlibabacloudTestAccSlbServercertificateCheckmap = map[string]string{

	// "fingerprint": CHECKSET,

	// "expire_time_stamp": CHECKSET,

	// "ali_cloud_certificate_id": CHECKSET,

	// "ali_cloud_certificate_name": CHECKSET,

	// "is_ali_cloud_certificate": CHECKSET,

	// "server_certificate_id": CHECKSET,

	// "server_certificate_name": CHECKSET,

	// "region_id": CHECKSET,

	// "expire_time": CHECKSET,
}

func AlibabacloudTestAccSlbServercertificateBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
