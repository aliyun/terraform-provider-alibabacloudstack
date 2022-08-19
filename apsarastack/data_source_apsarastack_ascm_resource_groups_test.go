package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackAscm_Resource_Groups_DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceApsaraStackAscm_Resource_Group_Organization,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_ascm_resource_groups.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_resource_groups.default", "groups.id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_resource_groups.default", "groups.name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_resource_groups.default", "groups.organization_id"),
				),
			},
		},
	})
}

const dataSourceApsaraStackAscm_Resource_Group_Organization = `
resource "apsarastack_ascm_organization" "default" {
  name = "Tf-testingresource-org"
  parent_id = "1"
}
 resource "apsarastack_ascm_resource_group" "default" {
  organization_id = apsarastack_ascm_organization.default.org_id
  name = "apsarastack-Datasource-resourceGroup"
}
data "apsarastack_ascm_resource_groups" "default" {
  name_regex = apsarastack_ascm_resource_group.default.name
}
`
