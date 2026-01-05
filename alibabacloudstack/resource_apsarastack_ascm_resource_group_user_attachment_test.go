package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAscmResourceGroupUserAttachmentBasic(t *testing.T) {
	var v *ListAscmUsersResponse
	resourceId := "alibabacloudstack_ascm_resource_group_user_attachment.default"
	ra := resourceAttrInit(resourceId, testAccCheckAscmResourceGroupUserAttachment)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccAscmResourceGroupUserAttachment%d", rand)
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccAscmResourceGroupUserAttachmentResource)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"user_id": "${alibabacloudstack_ascm_user.user.user_id}",
					"rg_id":   "${alibabacloudstack_ascm_resource_group.default.rg_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rg_id":   CHECKSET,
						"user_id": CHECKSET,
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

func testAccAscmResourceGroupUserAttachmentResource(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}


resource "alibabacloudstack_ascm_user" "user" {
	cellphone_number = "13900000000"
	email = "test@gmail.com"
	display_name = "C2C-DELTA"
	mobile_nation_code = "91"
	login_name = "${var.name}"
	login_policy_id = 1
}


resource "alibabacloudstack_ascm_resource_group" "default" {
  name = "${var.name}"
}
`, name)
}

var testAccCheckAscmResourceGroupUserAttachment = map[string]string{
	"rg_id":   CHECKSET,
	"user_id": CHECKSET,
}
