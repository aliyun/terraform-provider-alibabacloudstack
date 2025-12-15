package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccAlibabacloudStackHsmVendorsDataSource(t *testing.T) {
	yundunProvider := Provider()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: func() map[string]*schema.Provider {
			yundunProvider.Schema["access_key"] = &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_YUNDUN_ACCESS_KEY", ""),
				Description: descriptions["access_key"],
			}
			yundunProvider.Schema["secret_key"] = &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_YUNDUN_SECRET_KEY", ""),
				Description: descriptions["secret_key"],
			}
			yundunProvider.Schema["role_arn"] = &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["assume_role_role_arn"],
				DefaultFunc: schema.EnvDefaultFunc("ALIBABACLOUDSTACK_YUNDUN_ASSUME_ROLE_ARN", ""),
			}
			return map[string]*schema.Provider{
				"alibabacloudstack": yundunProvider,
			}
		}(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackHsmVendorsDataSource,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_hsm_vendors.default"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_hsm_vendors.default", "vendors.#"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_hsm_vendors.default", "vendors.0.code"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_hsm_vendors.default", "vendors.0.name"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_hsm_vendors.default", "vendors.0.products.#"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_hsm_vendors.default", "vendors.0.products.0.code"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_hsm_vendors.default", "vendors.0.products.0.name"),
				),
			},
		},
	})
}

const testAccCheckAlibabacloudStackHsmVendorsDataSource = `
data "alibabacloudstack_hsm_vendors" "default" {
}
`
