package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCloudmonitorserviceAlarmcontact0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_cloudmonitorservice_alarmcontact.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccCloudmonitorserviceAlarmcontactCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoCmsDescribecontactlistRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloud_monitor_servicealarm_contact%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccCloudmonitorserviceAlarmcontactBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"describe": "报警联系人信息",

					"alarm_contact_name": "Alice122",

					"channels_ali_im": "leo",

					"channels_ding_web_hook": "https://oapi.dingtalk.com/robot/send?access_token=7d49515e8ebf21106a80a9cc4bb3d2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"describe": "报警联系人信息",

						"alarm_contact_name": "Alice122",

						"channels_ali_im": "leo",

						"channels_ding_web_hook": "https://oapi.dingtalk.com/robot/send?access_token=7d49515e8ebf21106a80a9cc4bb3d2",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"describe": "报警联系人信息111",

					"channels_ali_im": "Lon",

					"channels_ding_web_hook": "https://oapi.dingtalk.com/robot/send",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"describe": "报警联系人信息111",

						"channels_ali_im": "Lon",

						"channels_ding_web_hook": "https://oapi.dingtalk.com/robot/send",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccCloudmonitorserviceAlarmcontactCheckmap = map[string]string{

	"describe": CHECKSET,

	"contact_groups": CHECKSET,

	"channels_sms": CHECKSET,

	"create_time": CHECKSET,

	"channels_ding_web_hook": CHECKSET,

	"alarm_contact_name": CHECKSET,

	"update_time": CHECKSET,

	"channels_mail": CHECKSET,

	"channels_ali_im": CHECKSET,
}

func AlibabacloudTestAccCloudmonitorserviceAlarmcontactBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
