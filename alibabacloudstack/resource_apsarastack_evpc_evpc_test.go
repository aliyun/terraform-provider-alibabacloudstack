package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func TestAccAlibabacloudStackEvpcEvpc_basic(t *testing.T) {

	var v map[string]interface{}

	resourceId := "alibabacloudstack_evpc_evpc.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEvpcEvpcCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &BcmpService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEasyAIListEvpcRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sevpc%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEvpcEvpcBasicdependence)
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
					"evpc_name":   "${var.name}",
					"description": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"evpc_name":   name,
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"evpc_name":   "${var.name}_updated",
					"description": "${var.name}_updated",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"evpc_name":   name + "_updated",
						"description": name + "_updated",
					}),
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

var AlibabacloudTestAccEvpcEvpcCheckmap = map[string]string{
	"status": CHECKSET,
}

func AlibabacloudTestAccEvpcEvpcBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

`, name)
}
