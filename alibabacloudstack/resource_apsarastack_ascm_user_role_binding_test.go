package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackAscmUserRoleBinding(t *testing.T) {
	var v *User
	resourceId := "alibabacloudstack_ascm_user_role_binding.default"
	ra := resourceAttrInit(resourceId, testAccCheckUserRoleBinding)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-ascmuserrole%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccCheckAscm_UserRoleBinding)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		CheckDestroy: testAccCheckAscm_UserRoleBinding_Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"role_ids" : []int{5},
					"login_name" : "${alibabacloudstack_ascm_user.default.login_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_ids.#": "1",
						"role_ids.0": "5",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"role_ids" : []int{6},
					"login_name" : "${alibabacloudstack_ascm_user.default.login_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_ids.#": "1",
						"role_ids.0": "6",
					}),
				),
			},
		},
	})

}

func testAccCheckAscm_UserRoleBinding_Destroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "alibabacloudstack_ascm_user_role_binding" || rs.Type != "alibabacloudstack_ascm_user_role_binding" {
			continue
		}
		ascm, err := ascmService.DescribeAscmUserRoleBinding(rs.Primary.ID)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
		if ascm.Message != "" {
			return errmsgs.WrapError(errmsgs.Error("resource  still exist"))
		}
	}

	return nil
}

func testAccCheckAscm_UserRoleBinding(name string) string{
	return fmt.Sprintf( `
variable name {
	default = "%s"
}
data alibabacloudstack_account current {
	
}

resource "alibabacloudstack_ascm_user" "default" {
 cellphone_number = "13900000000"
 email = "test@gmail.com"
 display_name = "C2C-DELTA"
 organization_id = data.alibabacloudstack_account.current.organization_id
 mobile_nation_code = "91"
 login_name = var.name
 login_policy_id = 1
}

`, name )
}

var testAccCheckUserRoleBinding = map[string]string{
	"login_name": CHECKSET,
}
