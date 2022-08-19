package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackCRReposDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceCRReposConfigDataSource,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_cr_repos.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cr_repos.default", "repos.namespace"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cr_repos.default", "repos.name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cr_repos.default", "repos.repo_type"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cr_repos.default", "repos.summary"),
					resource.TestCheckResourceAttrSet("data.apsarastack_cr_repos.default", "ids.#"),
				),
			},
		},
	})
}

const dataSourceCRReposConfigDataSource = `
variable "name" {
    default = "cr-repos-datasource"
}

resource "apsarastack_cr_namespace" "default" {
    name = "${var.name}"
    auto_create	= false
    default_visibility = "PUBLIC"
}

resource "apsarastack_cr_repo" "default" {
    namespace = "${apsarastack_cr_namespace.default.name}"
    name = "${var.name}"
    summary = "OLD"
    repo_type = "PUBLIC"
    detail  = "OLD"
}

data "apsarastack_cr_repos" "default" {
  name_regex = apsarastack_cr_repo.default.name
}

`
