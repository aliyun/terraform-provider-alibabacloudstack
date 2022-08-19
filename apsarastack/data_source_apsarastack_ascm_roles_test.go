package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackAscm_Roles_DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceApsaraStackAscm_Roles,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_ascm_roles.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_roles.default", "roles.id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_roles.default", "roles.name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_roles.default", "roles.role_level"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_roles.default", "roles.role_type"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_roles.default", "roles.user_count"),
				),
			},
		},
	})
}

const dataSourceApsaraStackAscm_Roles = `
resource "apsarastack_ascm_ram_role" "default" {
  role_name = "TestRamRoles"
  description = "TestingRam"
  organization_visibility = "global"
}

data "apsarastack_ascm_roles" "default" {
  name_regex = apsarastack_ascm_ram_role.default.role_name
}

`
