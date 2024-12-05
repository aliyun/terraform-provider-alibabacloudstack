package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackSlbDomainextension0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_slb_domainextension.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccSlbDomainextensionCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoSlbDescribedomainextensionattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslbdomain_extension%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccSlbDomainextensionBasicdependence)
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

					"server_certificate_id": "1511928242963727_183698fd346_-318606714_-773395734",

					"load_balancer_id": "lb-bp1jijiyb2hdauenc32zi",

					"domain": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"server_certificate_id": "1511928242963727_183698fd346_-318606714_-773395734",

						"load_balancer_id": "lb-bp1jijiyb2hdauenc32zi",

						"domain": "test",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"server_certificate_id": "1511928242963727_183698fd346_-318606714_-773395734",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"server_certificate_id": "1511928242963727_183698fd346_-318606714_-773395734",
					}),
				),
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



`, name)
}
