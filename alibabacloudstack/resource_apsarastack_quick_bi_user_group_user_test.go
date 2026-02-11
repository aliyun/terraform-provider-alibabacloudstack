package alibabacloudstack

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlicloudQuickBIUserGroupUser_basic0(t *testing.T) {
	//t.Skip()
	var v map[string]interface{}
	resourceId := "alibabacloudstack_quick_bi_user_group_user.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"user_group_id": CHECKSET,
		"account_id":    CHECKSET,
	})
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &QuickbiPublicService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-quickbiusergroup%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudQuickBIUserGroupUserBasicDependence0)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"user_group_id": "${alibabacloudstack_quick_bi_user_group.default.id}",
					"account_id":    "${alibabacloudstack_quick_bi_user.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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

func AlicloudQuickBIUserGroupUserBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_quick_bi_user" "default" {
  nick_name       = var.name
  account_name    = var.name
  admin_user      = "false"
  auth_admin_user = "false"
  user_type       = "Developer"
}

resource "alibabacloudstack_quick_bi_user_group" "default" {
	user_group_name        = var.name
	user_group_description = var.name
}
`, name)
}
