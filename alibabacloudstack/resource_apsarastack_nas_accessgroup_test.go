package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackNasAccessgroup0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_nas_accessgroup.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccNasAccessgroupCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoNasDescribeaccessgroupsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snasaccess_group%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccNasAccessgroupBasicdependence)
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

					"access_group_type": "Vpc",

					"description": "test",

					"access_group_name": "accssGroupExtremeVpcTest",

					"file_system_type": "extreme",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"access_group_type": "Vpc",

						"description": "test",

						"access_group_name": "accssGroupExtremeVpcTest",

						"file_system_type": "extreme",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test123",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccNasAccessgroupCheckmap = map[string]string{

	"rule_count": CHECKSET,

	"access_group_type": CHECKSET,

	"description": CHECKSET,

	"access_group_name": CHECKSET,

	"create_time": CHECKSET,

	"file_system_type": CHECKSET,

	"mount_target_count": CHECKSET,
}

func AlibabacloudTestAccNasAccessgroupBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
