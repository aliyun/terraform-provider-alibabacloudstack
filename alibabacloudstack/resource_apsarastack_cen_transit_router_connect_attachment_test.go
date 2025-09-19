package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCenTransitRouterConnectAttachment0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_cen_transit_router_connect_attachment.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccCenTransitRouterConnectAttachmentCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeCenTransitRouterConnectAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srouter_connect_attachment%d", defaultRegionToTest, rand)
	modify_name := fmt.Sprintf("tf-testacc%srouter_connect_attachment%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccCenTransitRouterConnectAttachmentBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		// CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{
					"cen_id":            "cen-y5hsaev9m0fyefocz5",
					"transit_router_id": "tr-xxxvhy5cbe7oniom3o1nm",
					"transit_router_attachment_name": "${var.name}"
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"transit_router_attachment_name": name,
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cen_id"},
			},
		},
	})
}

var AlibabacloudTestAccCenTransitRouterConnectAttachmentCheckmap = map[string]string{

	"status": CHECKSET,
}

func AlibabacloudTestAccCenTransitRouterConnectAttachmentBasicdependence(name string) string {
	vlan_id := getAccTestRandInt(1000, 2000)
	return fmt.Sprintf(
		`
variable "name" {
	default = "%s"
}
`, name)
}
