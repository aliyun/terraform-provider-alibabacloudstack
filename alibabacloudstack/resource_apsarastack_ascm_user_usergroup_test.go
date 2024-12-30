package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackAscm_UserGroup_User_Basic(t *testing.T) {
	var v *User
	resourceId := "alibabacloudstack_ascm_usergroup_user.default"
	ra := resourceAttrInit(resourceId, testAccCheckUserGroupUserBinding)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rand := getAccTestRandInt(100, 999)
	name := fmt.Sprintf("tf-ascmusergroup%v", rand)
	org_id := os.Getenv("ALIBABACLOUDSTACK_DEPARTMENT")
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		CheckDestroy: testAccCheckAscmUserGroupUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAscmUserGroupUserRoleBinding, org_id, name, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func testAccCheckAscmUserGroupUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "alibabacloudstack_ascm_usergroup_user" || rs.Type != "alibabacloudstack_ascm_usergroup_user" {
			continue
		}
		ascm, err := ascmService.DescribeAscmUsergroupUser(rs.Primary.ID)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
		if ascm.Message != "" {
			return errmsgs.WrapError(errmsgs.Error("user still exist"))
		}
	}

	return nil
}

const testAccCheckAscmUserGroupUserRoleBinding = `

variable org_id {
 default = "%s"
}

resource "alibabacloudstack_ascm_user_group" "default" {
 group_name =      "%s"
 organization_id = "${var.org_id}"
}

variable name {
 default = "tftest%d"
}

resource "alibabacloudstack_ascm_user" "default" {
 cellphone_number = "13900000000"
 email = "test@gmail.com"
 display_name = "${var.name}"
 organization_id = "${var.org_id}"
 mobile_nation_code = "86"
 login_name = "${var.name}"
 login_policy_id = 1
}


resource "alibabacloudstack_ascm_usergroup_user" "default" {
  login_names = [alibabacloudstack_ascm_user.default.login_name, ]
  user_group_id = alibabacloudstack_ascm_user_group.default.user_group_id
}
`

var testAccCheckUserGroupUserBinding = map[string]string{
	"user_group_id": CHECKSET,
	//"login_name": CHECKSET,
}
