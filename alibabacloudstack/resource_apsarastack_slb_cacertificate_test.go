package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackSlbCACertificate0(t *testing.T) {
	var v *slb.CACertificate

	resourceId := "alibabacloudstack_slb_cacertificate.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbCacertificateCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoSlbDescribecacertificatesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslbca_certificate%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbCacertificateBasicdependence)
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

					"name":           name,
					"ca_certificate": TfRawString(ServerCertificateTestCase()),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"name":           name,
						"ca_certificate": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// ca_certificate does not echo back after being set
				ImportStateVerifyIgnore: []string{"ca_certificate"},
			},
		},
	})
}

var AlibabacloudTestAccSlbCacertificateCheckmap = map[string]string{

	// "fingerprint": CHECKSET,

	// "region_id": CHECKSET,

	// "ca_certificate_name": CHECKSET,

	// "ca_certificate_id": CHECKSET,

	// "tags": CHECKSET,
}

func AlibabacloudTestAccSlbCacertificateBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
	
}
`, name)
}
