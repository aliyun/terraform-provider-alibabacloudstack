package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackCRReposDataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceCRReposConfigDataSource,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_cr_repos.default"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_cr_repos.default", "repos.namespace"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_cr_repos.default", "repos.name"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_cr_repos.default", "repos.repo_type"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_cr_repos.default", "repos.summary"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_cr_repos.default", "ids.#"),
				),
			},
		},
	})
}

const dataSourceCRReposConfigDataSource = `
variable "name" {
    default = "cr-repos-datasource"
}

resource "alibabacloudstack_cr_namespace" "default" {
    name = "${var.name}"
    auto_create	= false
    default_visibility = "PUBLIC"
}

resource "alibabacloudstack_cr_repo" "default" {
    namespace = "${alibabacloudstack_cr_namespace.default.name}"
    name = "${var.name}"
    summary = "OLD"
    repo_type = "PUBLIC"
    detail  = "OLD"
}

data "alibabacloudstack_cr_repos" "default" {
  name_regex = alibabacloudstack_cr_repo.default.name
}

`
