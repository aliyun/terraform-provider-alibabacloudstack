package alibabacloudstack

import (
	"fmt"
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
	rand := getAccTestRandInt(10000,20000)
	name := fmt.Sprintf("tf-ascmusergroup%v", rand)
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
				Config: fmt.Sprintf(testAccCheckAscmUserGroupUserRoleBinding, name, rand),
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
			return errmsgs.WrapError(errmsgs.Error("user  still exist"))
		}
	}

	return nil
}

const testAccCheckAscmUserGroupUserRoleBinding = `
resource "alibabacloudstack_ascm_organization" "default" {
 name = "Test_binder"
 parent_id = "1"
}

resource "alibabacloudstack_ascm_user_group" "default" {
 group_name =      "%s"
 organization_id = alibabacloudstack_ascm_organization.default.org_id
}

resource "alibabacloudstack_ascm_user" "default" {
 cellphone_number = "13900000000"
 email = "test@gmail.com"
 display_name = "C2C-DELTA"
 organization_id = alibabacloudstack_ascm_organization.default.org_id
 mobile_nation_code = "91"
 login_name = "User_Role_Test%d"
 login_policy_id = 1
}

resource "alibabacloudstack_ascm_usergroup_user" "default" {
  //login_name = alibabacloudstack_ascm_user.default.login_name
  login_names = ["User_Role_Test6304175127373178963", "User_Role_Test7233024715252325400"]
  //login_names = ["[\"User_Role_Test929636066677054911\"]"]
  user_group_id = alibabacloudstack_ascm_user_group.default.user_group_id
}
`

var testAccCheckUserGroupUserBinding = map[string]string{
	"user_group_id": CHECKSET,
	//"login_name": CHECKSET,
}
