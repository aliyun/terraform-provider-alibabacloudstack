package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackCms_Alarams_DataSource(t *testing.T) {
	testAccPreCheckWithAPIIsNotSupport(t)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceApsaraStackcms_alarms,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_cms_alarms.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cms_alarms.default", "alarms.group_name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cms_alarms.default", "alarms.metric_name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cms_alarms.default", "alarms.no_effective_interval"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cms_alarms.default", "alarms.silence_time"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cms_alarms.default", "alarms.contact_groups"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cms_alarms.default", "alarms.mail_subject"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cms_alarms.default", "alarms.source_type"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cms_alarms.default", "alarms.rule_id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cms_alarms.default", "alarms.period"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cms_alarms.default", "alarms.dimensions"),
				),
			},
		},
	})
}

const dataSourceApsaraStackcms_alarms = `
data "apsarastack_cms_alarms" "default" {

}
`
