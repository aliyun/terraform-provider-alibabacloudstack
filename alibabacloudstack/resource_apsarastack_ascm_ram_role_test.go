package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAscm_RamRoleBasic(t *testing.T) {
	var v *AscmRoles
	resourceId := "alibabacloudstack_ascm_ram_role.default"
	ra := resourceAttrInit(resourceId, ascmramroleRoleBasicMap)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-ascmramrole%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testascmramroleconfigbasic)
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
					"role_name":               name,
					"description":             "TestRole",
					"organization_visibility": "organizationVisibility.global",
					"role_range":              "roleRange.allOrganizations",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name":               name,
						"description":             "TestRole",
						"organization_visibility": CHECKSET,
						"role_range":              "roleRange.allOrganizations",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
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
					"role_name":               name,
					"description":             "TestRole",
					"organization_visibility": "organizationVisibility.global",
					"role_range":              "roleRange.userGroup",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name":               name,
						"description":             "TestRole",
						"organization_visibility": CHECKSET,
						"role_range":              "roleRange.userGroup",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testascmramroleconfigbasic(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}

var ascmramroleRoleBasicMap = map[string]string{
	"role_name":               CHECKSET,
	"description":             CHECKSET,
	"organization_visibility": CHECKSET,
	"role_range":              CHECKSET,
}
