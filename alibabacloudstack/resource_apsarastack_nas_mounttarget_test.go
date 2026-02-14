package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackNasMountTarget0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_ .default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccNasMounttargetCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoNasDescribemounttargetsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snasmount_target%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccNasMounttargetBasicdependence)
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

					"vswitch_id":        "${alibabacloudstack_vpc_vswitch.default.id}",
					"access_group_name": "${alibabacloudstack_nas_access_group.default.0.access_group_name}",
					"file_system_id":    "${alibabacloudstack_nas_file_system.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":        CHECKSET,
						"access_group_name": CHECKSET,
						"file_system_id":    CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_group_name": "${alibabacloudstack_nas_access_group.default.1.access_group_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Inactive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Inactive",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccNasMounttargetCheckmap = map[string]string{
	"status":            CHECKSET,
	"access_group_name": CHECKSET,
	"vswitch_id":        CHECKSET,
	"file_system_id":    CHECKSET,
}

func AlibabacloudTestAccNasMounttargetBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

%s

%s

resource "alibabacloudstack_nas_access_group" "default" {
	count = 2
	access_group_name = "${var.name}_${count.index}"
	access_group_type = "Vpc"
	description = "tf-testAccNasConfig"
}

`, name, VSwitchCommonTestCase, NasCommonTestCase)
}
