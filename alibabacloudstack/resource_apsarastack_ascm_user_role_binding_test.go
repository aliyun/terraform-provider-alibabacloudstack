package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAscmUserRoleBinding_roleIds(t *testing.T) {
	var v *User
	resourceId := "alibabacloudstack_ascm_user_role_binding.default"
	ra := resourceAttrInit(resourceId, testAccCheckUserRoleBinding)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-ascmuserrole1%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccCheckAscm_UserRoleBinding)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"role_ids":   []int{4, 5, 6, 8, 9, 10, 11, 12, 13, 14},
					"login_name": "${alibabacloudstack_ascm_user.default.login_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_ids.#": "10",
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
					"role_ids":   []int{2},
					"login_name": "${alibabacloudstack_ascm_user.default.login_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_ids.#": "1",
						"role_ids.0": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"role_ids":   []int{2, 4, 5},
					"login_name": "${alibabacloudstack_ascm_user.default.login_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_ids.#": "3",
						"role_ids.0": REMOVEKEY,
					}),
				),
			},
		},
	})

}

func TestAccAlibabacloudStackAscmUserRoleBinding_roleId(t *testing.T) {
	var v *User
	resourceId := "alibabacloudstack_ascm_user_role_binding.default"
	ra := resourceAttrInit(resourceId, testAccCheckUserRoleBinding)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-ascmuserrole2%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccCheckAscm_UserRoleBinding)
	testAccCheck := rac.resourceAttrMapUpdateSet()
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
					"role_id":    5,
					"login_name": "${alibabacloudstack_ascm_user.default.login_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_id": "5",
					}),
					resource.TestCheckTypeSetElemAttr(resourceId, "role_ids.*", "5"),
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

func testAccCheckAscm_UserRoleBinding(name string) string {
	return fmt.Sprintf(`
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

`, name)
}

var testAccCheckUserRoleBinding = map[string]string{
	"login_name": CHECKSET,
}
