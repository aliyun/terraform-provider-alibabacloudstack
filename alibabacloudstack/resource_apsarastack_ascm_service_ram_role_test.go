package alibabacloudstack

import (
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAscmServiceRamRoleBasic(t *testing.T) {
	var v *ListRAMServiceRolesResponse
	resourceId := "alibabacloudstack_ascm_service_ram_role.default"
	ra := resourceAttrInit(resourceId, testAccCheckAscmServiceRamRole)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "ListRAMServiceRoles")
	rac := resourceAttrCheckInit(rc, ra)
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
				Config: testAccAscmServiceRamRoleDependence,
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

const testAccAscmServiceRamRoleDependence = `
resource "alibabacloudstack_ascm_organization" "default" {
  name = "Tf-testingresource-org2"
  parent_id = "1"
} 
 resource "alibabacloudstack_ascm_service_ram_role" "default" {
  organization_id = "${alibabacloudstack_ascm_organization.default.id}"
  product_name = "ECS"
}`

var testAccCheckAscmServiceRamRole = map[string]string{
	"ram_roles.0.id":                          CHECKSET,
	"ram_roles.0.arn":                         CHECKSET,
	"ram_roles.0.region":                      CHECKSET,
	"ram_roles.0.role_id":                     CHECKSET,
	"ram_roles.0.role_name":                   CHECKSET,
	"ram_roles.0.role_type":                   CHECKSET,
	"ram_roles.0.description":                 CHECKSET,
	"ram_roles.0.aliyun_user_id":              CHECKSET,
	"ram_roles.0.assume_role_policy_document": CHECKSET,
	"ram_roles.0.organization_name":           CHECKSET,
	"ram_roles.0.policies.#":                  CHECKSET,
}
