package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackCRNamespacesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceCRNamespacesConfigDependence,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_cr_namespaces.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cr_namespaces.default", "namespaces.default_visibility"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cr_namespaces.default", "instance_types.auto_create"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cr_namespaces.default", "instance_types.name"),
					resource.TestCheckResourceAttrSet("data.apsarastack_cr_namespaces.default", "ids.#"),
				),
			},
		},
	})
}

const dataSourceCRNamespacesConfigDependence = `
  resource "apsarastack_cr_namespace" "default" {
  name               = "testing-db-nspace"
  auto_create        = false
  default_visibility = "PUBLIC"
}

  data "apsarastack_cr_namespaces" "default" {
  name_regex    = apsarastack_cr_namespace.default.name
}
`
