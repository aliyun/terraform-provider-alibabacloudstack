package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackAscmRamRoleBasic(t *testing.T) {
	var v *AscmRoleData
	resourceId := "alibabacloudstack_ascm_ram_role.default"
	ra := resourceAttrInit(resourceId, testAccCheckAscmRamRole)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tftestrole%d", rand)
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccAscm_RamRole_resource)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAscm_RamRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"role_name":               "${var.name}",
					"description":             "${var.name} desc",
					"organization_visibility": "global",
					"role_range":              "roleRange.userGroup",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name":               name,
						"description":             name + " desc",
						"organization_visibility": "global",
						"role_range":              "roleRange.userGroup",
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
					"role_name":               "${var.name}new",
					"description":             "${var.name} desc update",
					"organization_visibility": "global",
					"role_range":              "roleRange.userGroup",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name":               name + "new",
						"description":             name + " desc update",
						"organization_visibility": "global",
						"role_range":              "roleRange.userGroup",
					}),
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
		object, err := ascmService.DescribeAscmRamRole(rs.Primary.ID)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
		if object != nil {
			return errmsgs.WrapError(errmsgs.Error("ram role still exist"))
		}
	}

	return nil
}

func testAccAscm_RamRole_resource(name string) string {
	return fmt.Sprintf(`
	variable name {
		default = "%s"
	}
	
	resource alibabacloudstack_ascm_ram_role distractor {
		role_name               = "${var.name}-distractor"
		description             = "${var.name} distractor"
		organization_visibility = "global"
		role_range              = "roleRange.userGroup"
	}
	
`, name)
}

var testAccCheckAscmRamRole = map[string]string{
	"role_name":               CHECKSET,
	"organization_visibility": CHECKSET,
}
