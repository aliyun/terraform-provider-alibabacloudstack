package apsarastack

import (
	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccApsaraStackAscm_RamRoleBasic(t *testing.T) {
	var v *AscmRoles
	resourceId := "apsarastack_ascm_ram_role.default"
	ra := resourceAttrInit(resourceId, testAccCheckAscmRamRole)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
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
	client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)

	ascmService := AscmService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "apsarastack_ascm_ram_role" || rs.Type != "apsarastack_ascm_ram_role" {
			continue
		}
		ascm, err := ascmService.DescribeAscmRamRole(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		if ascm.AsapiErrorCode != "200" {
			return WrapError(Error("ram role still exist"))
		}
	}

	return nil
}

const testAccAscm_RamRole_resource = `
resource "apsarastack_ascm_ram_role" "default" {
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
