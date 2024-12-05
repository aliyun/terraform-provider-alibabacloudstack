package alibabacloudstack

import (
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackAscm_RamRoleBasic(t *testing.T) {
	var v *AscmRoles
	resourceId := "alibabacloudstack_ascm_ram_role.default"
	ra := resourceAttrInit(resourceId, testAccCheckAscmRamRole)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAscm_RamRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAscm_RamRole_resource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func testAccCheckAscm_RamRoleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)

	ascmService := AscmService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "alibabacloudstack_ascm_ram_role" || rs.Type != "alibabacloudstack_ascm_ram_role" {
			continue
		}
		ascm, err := ascmService.DescribeAscmRamRole(rs.Primary.ID)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
		if ascm.AsapiErrorCode != "200" {
			return errmsgs.WrapError(errmsgs.Error("ram role still exist"))
		}
	}

	return nil
}

const testAccAscm_RamRole_resource = `
resource "alibabacloudstack_ascm_ram_role" "default" {
  role_name = "Test_Ram_Role"
  description = "TestRole"
  organization_visibility = "global"
  role_range = "roleRange.userGroup"
}
`

var testAccCheckAscmRamRole = map[string]string{
	"role_name":               CHECKSET,
	"organization_visibility": CHECKSET,
}
