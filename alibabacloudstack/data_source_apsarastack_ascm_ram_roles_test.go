package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAscmRamRoles_DataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceAlibabacloudStackAscm_Roles(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_ascm_roles.default"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_roles.default", "roles.id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_roles.default", "roles.name"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_roles.default", "roles.role_level"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_roles.default", "roles.role_type"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_roles.default", "roles.user_count"),
				),
			},
		},
	})
}

func dataSourceAlibabacloudStackAscm_Roles() string {
	return fmt.Sprintf(`
resource "alibabacloudstack_ascm_ram_role" "default" {
  role_name = "tftestrole%d"
  description = "TestingRam"
  organization_visibility = "global"
  role_range = "roleRange.userGroup"
}

data "alibabacloudstack_ascm_roles" "default" {
  name_regex = alibabacloudstack_ascm_ram_role.default.role_name
}

`, getAccTestRandInt(1000000, 9999999))
}
