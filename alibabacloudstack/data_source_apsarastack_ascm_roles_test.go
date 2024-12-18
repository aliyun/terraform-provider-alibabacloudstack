package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackAscm_Roles_DataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceAlibabacloudStackAscm_Roles,
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

const dataSourceAlibabacloudStackAscm_Roles = `
resource "alibabacloudstack_ascm_ram_role" "default" {
  role_name = "TestRamRoles"
  description = "TestingRam"
  organization_visibility = "global"
role_range = "roleRange.allOrganizations"
}

data "alibabacloudstack_ascm_roles" "default" {
  name_regex = alibabacloudstack_ascm_ram_role.default.role_name
}

`
