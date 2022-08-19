package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackAscm_OrganizationDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceApsaraStackAscm_E_Organization,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_ascm_organizations.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_organizations.default", "organizations.id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_organizations.default", "organizations.name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_organizations.default", "organizations.cuser_id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_organizations.default", "organizations.muser_id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_organizations.default", "organizations.alias"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_organizations.default", "organizations.parent_id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_organizations.default", "organizations.internal"),
				),
			},
		},
	})
}

const dataSourceApsaraStackAscm_E_Organization = `

resource "apsarastack_ascm_organization" "org" {
 name = "Tf-testing-DataSourceascm-organization"
 parent_id = "1"
}

data "apsarastack_ascm_organizations" "default" {
  name_regex = apsarastack_ascm_organization.org.name
}
`
