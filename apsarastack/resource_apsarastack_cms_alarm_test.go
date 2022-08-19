package apsarastack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"testing"
)

func TestAccApsaraStackCmsAlarmBasic(t *testing.T) {
	log.Printf("2022.7.29 ")
	testAccPreCheckWithAPIIsNotSupport(t)
	var v cms.Alarm
	resourceId := "apsarastack_cms_alarm.default"
	ra := resourceAttrInit(resourceId, testAccCheckAlarm)
	serviceFunc := func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		CheckDestroy: testAccCheckCmsAlarm_Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCmsAlarm,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})

}

func testAccCheckCmsAlarm_Destroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
	cmsService := CmsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "apsarastack_cms_alarm" || rs.Type != "apsarastack_cms_alarm" {
			continue
		}
		cms, err := cmsService.DescribeCmsAlarm(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		if cms.RuleName != "" {
			return WrapError(Error("resource  still exist"))
		}
	}

	return nil
}

const testAccCheckCmsAlarm = `
resource "apsarastack_slb" "basic" {
 name          = "terraform_omega_1"
}
resource "apsarastack_cms_alarm" "default" {
  name    = "TfAccCmsAlarm_omega_1"
  project = "acs_slb_dashboard"
  metric  = "ActiveConnection"
  dimensions = {
    instanceId = apsarastack_slb.basic.id
  }
  escalations_critical {
    statistics = "Average"
    comparison_operator = "<="
    threshold = 35
    times = 2
  }
  enabled =      true
  contact_groups     = ["test-group"]
  effective_interval = "0:00-2:00"
}
`

var testAccCheckAlarm = map[string]string{
	"name":    CHECKSET,
	"project": CHECKSET,
	"metric":  CHECKSET,
}
