package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackAscmMaxcomputeProjectDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
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
	name = "tf_testAccAlibabacloudStack5610"
}
`
