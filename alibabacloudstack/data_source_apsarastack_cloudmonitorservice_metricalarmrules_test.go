package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCms_Alarams_DataSource(t *testing.T) {
	// testAccPreCheckWithAPIIsNotSupport(t)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceAlibabacloudStackcms_alarms,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_cms_alarms.default"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_cms_alarms.default", "alarms.group_name"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_cms_alarms.default", "alarms.metric_name"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_cms_alarms.default", "alarms.no_effective_interval"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_cms_alarms.default", "alarms.silence_time"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_cms_alarms.default", "alarms.contact_groups"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_cms_alarms.default", "alarms.mail_subject"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_cms_alarms.default", "alarms.source_type"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_cms_alarms.default", "alarms.rule_id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_cms_alarms.default", "alarms.period"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_cms_alarms.default", "alarms.dimensions"),
				),
			},
		},
	})
}

const dataSourceAlibabacloudStackcms_alarms = VSwitchCommonTestCase + `

%s

resource "alibabacloudstack_slb" "default" {
  name          = "terraform_test111"
  address_type  = "internet"
  specification = "slb.s2.small"
  vswitch_id    = alibabacloudstack_vpc_vswitch.default.id
}


data "alibabacloudstack_cms_alarms" "default" {
  name    = "tf-testAccCmsAlarm_basic"
  project = "acs_slb_dashboard"
  metric  = "ActiveConnection"
  dimensions = {
    instanceId = "${alibabacloudstack_slb.default.id}"
  }
  escalations_critical {
    statistics = "Average"
    comparison_operator = "<="
    threshold = 35
    times = 2
  }
  period  =    300
  enabled =      true
  contact_groups     = ["test-group"]
  effective_interval = "0:00-2:00"
  webhook = "http://www.aliyun.com"
}
`
