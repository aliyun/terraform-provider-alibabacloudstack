package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackNetworkInterfaceAttachmentBasic(t *testing.T) {
	var v ecs.NetworkInterfaceSet
	resourceId := "alibabacloudstack_network_interface_attachment.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"instance_id":          CHECKSET,
		"network_interface_id": CHECKSET,
	})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeNetworkInterfaceAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsnetwork_attachment%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsNetworkInterfaceAttachmentdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNetworkInterfaceAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":          "${alibabacloudstack_ecs_instance.default.id}",
					"network_interface_id": "${alibabacloudstack_network_interface.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
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

func testAccCheckNetworkInterfaceAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_network_interface_Attachment" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No NetworkInterface ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
		ecsService := EcsService{client}
		_, err := ecsService.DescribeNetworkInterfaceAttachment(rs.Primary.ID)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return err
		}
	}

	return nil
}

func AlibabacloudTestAccEcsNetworkInterfaceAttachmentdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

%s

resource "alibabacloudstack_network_interface" "default" {
    name = "${var.name}"
    vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
    security_groups = [ "${alibabacloudstack_ecs_securitygroup.default.id}" ]
}


`, name, ECSInstanceCommonTestCase)
}
