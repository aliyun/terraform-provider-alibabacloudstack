package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackSlbCacertificate0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_slb_cacertificate.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbCacertificateCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoSlbDescribecacertificatesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslbca_certificate%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbCacertificateBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"region_id": "cn-hangzhou",

					"ca_certificate_name": "mycacert01",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"region_id": "cn-hangzhou",

						"ca_certificate_name": "mycacert01",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccSlbCacertificateCheckmap = map[string]string{

	"fingerprint": CHECKSET,

	"region_id": CHECKSET,

	"ca_certificate_name": CHECKSET,

	"ca_certificate_id": CHECKSET,

	"tags": CHECKSET,
}

func AlibabacloudTestAccSlbCacertificateBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
