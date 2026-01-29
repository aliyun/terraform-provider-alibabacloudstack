package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackCmsAlarmBasic(t *testing.T) {
	var v cms.AlarmInDescribeMetricRuleList
	resourceId := "alibabacloudstack_cms_alarm.default"
	ra := resourceAttrInit(resourceId, testAccCheckAlarm)
	serviceFunc := func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf_testacc_cmsalarm%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccCheckCmsAlarm)
	ResourceTest(t, resource.TestCase{
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
				Config: testAccConfig(map[string]interface{}{
					"name":    "${var.name}",
					"project": "acs_slb_dashboard",
					"metric":  "ActiveConnection",
					"dimensions": map[string]interface{}{
						"instanceId": "${alibabacloudstack_slb.basic.id}",
					},
					"escalations_critical": []map[string]interface{}{{
						"statistics":          "Average",
						"comparison_operator": "<=",
						"threshold":           30,
						"times":               2,
					}},
					"escalations_warn": []map[string]interface{}{{
						"statistics":          "Average",
						"comparison_operator": "<=",
						"threshold":           40,
						"times":               2,
					}},
					"escalations_info": []map[string]interface{}{{
						"statistics":          "Average",
						"comparison_operator": "<=",
						"threshold":           50,
						"times":               2,
					}},
					"enabled":            true,
					"contact_groups":     []string{"test-group"},
					"effective_interval": "0:00-2:00",
					"webhook":            "http://www.alibabacloud.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":    name,
						"project": "acs_slb_dashboard",
						"metric":  "ActiveConnection",
						"enabled": "true",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":    "${var.name}_update",
					"project": "acs_slb_dashboard",
					"metric":  "ActiveConnection",
					"dimensions": map[string]interface{}{
						"instanceId": "${alibabacloudstack_slb.basic.id}",
					},
					"escalations_critical": []map[string]interface{}{{
						"statistics":          "Average",
						"comparison_operator": "<=",
						"threshold":           35,
						"times":               3,
					}},
					"escalations_warn": []map[string]interface{}{{
						"statistics":          "Average",
						"comparison_operator": "<=",
						"threshold":           45,
						"times":               3,
					}},
					"escalations_info": []map[string]interface{}{{
						"statistics":          "Average",
						"comparison_operator": "<=",
						"threshold":           55,
						"times":               3,
					}},
					"enabled":            true,
					"contact_groups":     []string{"test-group"},
					"effective_interval": "0:00-2:00",
					"webhook":            "http://www.alibabacloud.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":    name + "_update",
						"project": "acs_slb_dashboard",
						"metric":  "ActiveConnection",
						"enabled": "true",
					}),
				),
			},
		},
	})

}

func testAccCheckCmsAlarm_Destroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	cmsService := CmsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "alibabacloudstack_cms_alarm" || rs.Type != "alibabacloudstack_cms_alarm" {
			continue
		}
		cms, err := cmsService.DescribeCmsAlarm(rs.Primary.ID)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
		if cms.RuleName != "" {
			return errmsgs.WrapError(errmsgs.Error("resource  still exist"))
		}
	}

	return nil
}

func testAccCheckCmsAlarm(name string) string {
	return fmt.Sprintf(`
variable "name" {
 default = "%s"
}

resource "alibabacloudstack_slb" "basic" {
 name          = "${var.name}"
}
`, name)
}

var testAccCheckAlarm = map[string]string{
	"name":                   CHECKSET,
	"project":                CHECKSET,
	"metric":                 CHECKSET,
	"dimensions.%":           "1",
	"enabled":                CHECKSET,
	"contact_groups.0":       "test-group",
	"effective_interval":     CHECKSET,
	"escalations_critical.#": "1",
	"webhook":                "http://www.alibabacloud.com",
}
