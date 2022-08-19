package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackAscm_EnviromentServiceByProduct_DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceApsaraStackAscm_EnviromentServiceByProduct,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_ascm_environment_services_by_product.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_environment_services_by_product.default", "result"),
				),
			},
		},
	})
}

const dataSourceApsaraStackAscm_EnviromentServiceByProduct = `

data "apsarastack_ascm_environment_services_by_product" "default" {
  output_file = "environment"
}

`
