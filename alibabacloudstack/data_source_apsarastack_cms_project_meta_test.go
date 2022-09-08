package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackCms_Projectmeta_DataSource(t *testing.T) {
	testAccPreCheckWithAPIIsNotSupport(t)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceAlibabacloudStackcms_project,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_cms_project_meta.default"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_cms_project_meta.default", "resources.description"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_cms_project_meta.default", "resources.labels"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_cms_project_meta.default", "resources.namespace"),
				),
			},
		},
	})
}

const dataSourceAlibabacloudStackcms_project = `
data "alibabacloudstack_cms_project_meta" "default" {
}
`
