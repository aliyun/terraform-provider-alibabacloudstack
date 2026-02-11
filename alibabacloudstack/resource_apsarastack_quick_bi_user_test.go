package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlicloudQuickBIUser_basic0(t *testing.T) {
	//t.Skip()
	var v map[string]interface{}
	resourceId := "alibabacloudstack_quick_bi_user.default"
	ra := resourceAttrInit(resourceId, AlicloudQuickBIUserMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &QuickbiPublicService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeQuickBiUser")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%squickbiuser%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudQuickBIUserBasicDependence0)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"nick_name":       name,
					"account_name":    name,
					"admin_user":      "false",
					"auth_admin_user": "false",
					"user_type":       "Developer",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nick_name":       name,
						"account_name":    CHECKSET,
						"admin_user":      "false",
						"auth_admin_user": "false",
						"user_type":       "Developer",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"admin_user": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"admin_user": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_type": "Analyst",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_type": "Analyst",
					}),
				),
			},
			// todo fix  product problem
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"auth_admin_user": "true",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"auth_admin_user": "true",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"admin_user":      "false",
					"auth_admin_user": "false",
					"user_type":       "Developer",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"admin_user":      "false",
						"auth_admin_user": "false",
						"user_type":       "Developer",
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

var AlicloudQuickBIUserMap0 = map[string]string{}

func AlicloudQuickBIUserBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
