package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackCms_Projectmetalist_DataSource(t *testing.T) {
	testAccPreCheckWithAPIIsNotSupport(t)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceApsaraStackcms_metalist,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_cms_metric_metalist.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cms_metric_metalist.default", "resources.metric_name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cms_metric_metalist.default", "resources.periods"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cms_metric_metalist.default", "resources.description"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cms_metric_metalist.default", "resources.dimensions"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cms_metric_metalist.default", "resources.labels"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cms_metric_metalist.default", "resources.unit"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cms_metric_metalist.default", "resources.statistics"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cms_metric_metalist.default", "resources.namespace"),
				),
			},
		},
	})
}

const dataSourceApsaraStackcms_metalist = `
data "apsarastack_cms_metric_metalist" "default" {
namespace="acs_slb_dashboard"
}
`
