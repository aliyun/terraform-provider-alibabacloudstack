package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAscmMaxcomputeProjectDataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: datasourceAlibabacloudstackMaxcomputeProjects,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_maxcompute_projects.default"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_maxcompute_projects.default", "projects.id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_maxcompute_projects.default", "projects.name"),
				),
			},
		},
	})
}

const datasourceAlibabacloudstackMaxcomputeProjects = `

data "alibabacloudstack_maxcompute_projects" "default"{
	name = "testttt"
}
`
