package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAscm_UserGroup_User_Basic(t *testing.T) {
	var v *User
	resourceId := "alibabacloudstack_ascm_usergroup_user.default"
	ra := resourceAttrInit(resourceId, testAccCheckUserGroupUserBinding)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tfuser%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccCheckAscmUserGroupUserconfigbasic)
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
					"login_names":   []string{"${alibabacloudstack_ascm_user.default[0].login_name}"},
					"user_group_id": "${alibabacloudstack_ascm_user_group.default.user_group_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"login_names.#": "1",
						"login_names.0": name + "0",
						"user_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"login_names":   []string{"${alibabacloudstack_ascm_user.default[0].login_name}", "${alibabacloudstack_ascm_user.default[1].login_name}"},
					"user_group_id": "${alibabacloudstack_ascm_user_group.default.user_group_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"login_names.#": "2",
						"user_group_id": CHECKSET,
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

func testAccCheckAscmUserGroupUserconfigbasic(name string) string {
	return fmt.Sprintf(`
variable name{
 default = "%s"
}

resource "alibabacloudstack_ascm_organization" "default" {
  name = "${var.name}"
  parent_id = "1"
} 


resource "alibabacloudstack_ascm_user_group" "default" {
 group_name =      "${var.name}"
 organization_id = "${alibabacloudstack_ascm_organization.default.id}"
}

resource "alibabacloudstack_ascm_user" "default" {
 count = 2
 cellphone_number = "13900000000"
 email = "test@gmail.com"
 display_name = "${var.name}${count.index}"
 organization_id = "${alibabacloudstack_ascm_organization.default.id}"
 mobile_nation_code = "86"
 login_name = "${var.name}${count.index}"
 login_policy_id = 1
}

`, name)
}

var testAccCheckUserGroupUserBinding = map[string]string{
	"user_group_id": CHECKSET,
	//"login_name": CHECKSET,
}
