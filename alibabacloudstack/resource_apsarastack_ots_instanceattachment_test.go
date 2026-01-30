package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackOtsInstanceAttachmentBasic(t *testing.T) {
	var v ots.VpcInfo

	resourceId := "alibabacloudstack_ots_instance_attachment.default"
	ra := resourceAttrInit(resourceId, otsInstanceAttachmentBasicMap)

	serviceFunc := func() interface{} {
		return &OtsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tftestattch%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOtsInstanceAttachmentConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${alibabacloudstack_ots_instance.default.name}",
					"vpc_name":      "${var.name}",
					"vswitch_id":    "${alibabacloudstack_vpc_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
						"vpc_name":      name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vswitch_id"},
			},
		},
	})
}

func resourceOtsInstanceAttachmentConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}

	%s
	
	data "alibabacloudstack_ots_clusters" "default" {}

	resource "alibabacloudstack_ots_instance" "default" {
	  name = var.name
	  description   = var.name
	  specification  = "HYBRID"
	}
	`, name, VSwitchCommonTestCase)
}

var otsInstanceAttachmentBasicMap = map[string]string{
	"instance_name": CHECKSET,
	"vpc_name":      CHECKSET,
	"vswitch_id":    CHECKSET,
	"vpc_id":        CHECKSET,
}
