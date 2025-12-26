package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackAscm_RamPolicyForRoleBasic(t *testing.T) {
	var v *RamPolicies

	resourceId := "alibabacloudstack_ascm_ram_policy_for_role.default"
	ra := resourceAttrInit(resourceId, testAccCheckAscmRamPolicyForRole)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccAscm_RamPolicyForRole_resource)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAscm_RamPolicyForRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ram_policy_id": "${alibabacloudstack_ascm_ram_policy.default.ram_id}",
					"role_id":       "${alibabacloudstack_ascm_ram_role.default.role_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
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

func testAccCheckAscm_RamPolicyForRoleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)

	ascmService := AscmService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "alibabacloudstack_ascm_ram_policy_for_role" || rs.Type != "alibabacloudstack_ascm_ram_policy_for_role" {
			continue
		}
		ascm, err := ascmService.DescribeAscmRamPolicy(rs.Primary.ID)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
		if ascm.AsapiErrorCode != "200" {
			return errmsgs.WrapError(errmsgs.Error("ram policy still exist"))
		}
	}

	return nil
}

func testAccAscm_RamPolicyForRole_resource(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
	
resource "alibabacloudstack_ascm_ram_policy" "default" {
  name = var.name
  description = "Testing Complete"
  policy_document = "{\"Statement\":[{\"Action\":\"ecs:*\",\"Effect\":\"Allow\",\"Resource\":\"*\"}],\"Version\":\"1\"}"

}

resource "alibabacloudstack_ascm_ram_role" "default" {
  role_name = var.name
  description = "TestingRole"
  organization_visibility = "global"
role_range = "roleRange.allOrganizations"
}

`, name)
}

var testAccCheckAscmRamPolicyForRole = map[string]string{
	"ram_policy_id": CHECKSET,
	"role_id":       CHECKSET,
}
