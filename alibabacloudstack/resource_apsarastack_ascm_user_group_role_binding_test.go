package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackAscmUserGroupRoleBinding(t *testing.T) {
	var v *UserGroup
	resourceId := "alibabacloudstack_ascm_user_group_role_binding.default"
	ra := resourceAttrInit(resourceId, testAccCheckUserGroupRoleBinding)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-ascmgrouprole%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccCheckAscm_UserGroupRoleBinding)
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
		CheckDestroy: testAccCheckAscm_UserGroupRoleBinding_Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"role_ids" : []int{5},
					"user_group_id" : "${alibabacloudstack_ascm_user_group.default.user_group_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_ids.#": "1",
						"role_ids.0": "5",
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
					"role_ids" : []int{6},
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

func testAccCheckAscm_UserGroupRoleBinding_Destroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "alibabacloudstack_ascm_user_group_role_binding" || rs.Type != "alibabacloudstack_ascm_user_group_role_binding" {
			continue
		}
		ascm, err := ascmService.DescribeAscmUserGroup(rs.Primary.ID)
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

func testAccCheckAscm_UserGroupRoleBinding(name string) string{
	return fmt.Sprintf( `
variable name {
 default = "%s"
}

data alibabacloudstack_account current {
	
}

resource "alibabacloudstack_ascm_user_group" "default" {
 group_name =      var.name
 organization_id = data.alibabacloudstack_account.current.organization_id
}

`, name)
}

var testAccCheckUserGroupRoleBinding = map[string]string{
	"user_group_id": CHECKSET,
}
