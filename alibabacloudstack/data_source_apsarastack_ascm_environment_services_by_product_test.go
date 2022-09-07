package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackAscm_EnviromentServiceByProduct_DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceAlibabacloudStackAscm_EnviromentServiceByProduct,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_ascm_environment_services_by_product.default"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_environment_services_by_product.default", "result"),
				),
			},
		},
	})
}

const dataSourceAlibabacloudStackAscm_EnviromentServiceByProduct = `

data "alibabacloudstack_ascm_environment_services_by_product" "default" {
  output_file = "environment"
}

`
