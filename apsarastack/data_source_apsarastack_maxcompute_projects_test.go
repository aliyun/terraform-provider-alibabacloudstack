package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackAscmMaxcomputeProjectDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: datasourceApsarastackMaxcomputeProjects,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_maxcompute_projects.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_maxcompute_projects.default", "projects.id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_maxcompute_projects.default", "projects.name"),
				),
			},
		},
	})
}

const datasourceApsarastackMaxcomputeProjects = `

data "apsarastack_maxcompute_projects" "default"{
	name = "tf_testAccApsaraStack5610"
}
`
