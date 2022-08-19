package apsarastack

import (
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccApsaraStackAscm_RamPolicyForRoleBasic(t *testing.T) {
	var v *RamPolicies

	resourceId := "apsarastack_ascm_ram_policy_for_role.default"
	ra := resourceAttrInit(resourceId, testAccCheckAscmRamPolicyForRole)
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

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAscm_RamPolicyForRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAscm_RamPolicyForRole_resource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func testAccCheckAscm_RamPolicyForRoleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)

	ascmService := AscmService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "apsarastack_ascm_ram_policy_for_role" || rs.Type != "apsarastack_ascm_ram_policy_for_role" {
			continue
		}
		ascm, err := ascmService.DescribeAscmRamPolicy(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		if ascm.AsapiErrorCode != "200" {
			return WrapError(Error("ram policy still exist"))
		}
	}

	return nil
}

const testAccAscm_RamPolicyForRole_resource = `
resource "apsarastack_ascm_ram_policy" "default" {
  name = "TestPolicyRole"
  description = "Testing Complete"
  policy_document = "{\"Statement\":[{\"Action\":\"ecs:*\",\"Effect\":\"Allow\",\"Resource\":\"*\"}],\"Version\":\"1\"}"

}

resource "apsarastack_ascm_ram_role" "default" {
  role_name = "TestPolicyRole"
  description = "TestingRole"
  organization_visibility = "global"
}

resource "apsarastack_ascm_ram_policy_for_role" "default" {
  ram_policy_id = apsarastack_ascm_ram_policy.default.ram_id
  role_id = apsarastack_ascm_ram_role.default.role_id
}

`

var testAccCheckAscmRamPolicyForRole = map[string]string{
	"ram_policy_id": CHECKSET,
	"role_id":       CHECKSET,
}
