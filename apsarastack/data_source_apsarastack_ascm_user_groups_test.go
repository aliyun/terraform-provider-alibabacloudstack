package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackAscm_User_GroupsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceApsaraStackAscm_User_Group_Organization,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_ascm_user_groups.default"),
				),
			},
		},
	})
}

const dataSourceApsaraStackAscm_User_Group_Organization = `
data "apsarastack_ascm_user_groups" "default" {
	name_regex = "cxt"
	
}
output "groups" {
	  value = "${data.apsarastack_ascm_user_groups.default.groups}"
	}
`
