package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCloudmonitorserviceAlarmcontactgroup0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_cloudmonitorservice_alarmcontactgroup.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccCloudmonitorserviceAlarmcontactgroupCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoCmsDescribecontactgrouplistRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloud_monitor_servicealarm_contact_group%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccCloudmonitorserviceAlarmcontactgroupBasicdependence)
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

					"alarm_contact_group_name": "AlarmContactGroupNameTest",

					"describe": "Describe",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"alarm_contact_group_name": "AlarmContactGroupNameTest",

						"describe": "Describe",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"describe": "Describe33",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"describe": "Describe33",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"alarm_contact_group_name": "AlarmContactGroupNameTest",

					"describe": "Describe",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"alarm_contact_group_name": "AlarmContactGroupNameTest",

						"describe": "Describe",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccCloudmonitorserviceAlarmcontactgroupCheckmap = map[string]string{

	"describe": CHECKSET,

	"contact_names": CHECKSET,

	"alarm_contact_group_name": CHECKSET,
}

func AlibabacloudTestAccCloudmonitorserviceAlarmcontactgroupBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
