package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackVpcRoutetableattachment0(t *testing.T) {
	var v vpc.RouterTableListType

	resourceId := "alibabacloudstack_vpc_routetableattachment.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpcRoutetableattachmentCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeRouteTableAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcroute_table_attachment%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpcRoutetableattachmentBasicdependence)
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

					"route_table_id": "${alibabacloudstack_route_table.default.id}",

					"vswitch_id": "${alibabacloudstack_vpc_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"route_table_id": CHECKSET,

						"vswitch_id": CHECKSET,
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccVpcRoutetableattachmentCheckmap = map[string]string{

	// "status": CHECKSET,

	"route_table_id": CHECKSET,

	"vswitch_id": CHECKSET,
}

func AlibabacloudTestAccVpcRoutetableattachmentBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
  
%s
  
resource "alibabacloudstack_route_table" "default" {
vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
name = "${var.name}"
description = "${var.name}_description"
}

`, name, VSwitchCommonTestCase)
}
