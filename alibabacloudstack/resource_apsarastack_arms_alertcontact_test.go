package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackArmsAlertcontact0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_arms_alertcontact.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccArmsAlertcontactCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoArmsSearchalertcontactRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sarmsalert_contact%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccArmsAlertcontactBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"phone_num": "12345678910",

					"alert_contact_name": "test",

					"email": "123@qq.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"phone_num": "12345678910",

						"alert_contact_name": "test",

						"email": "123@qq.com",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"email": "1234@qq.com",

					"alert_contact_name": "rdktest",

					"phone_num": "99999999999",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"email": "1234@qq.com",

						"alert_contact_name": "rdktest",

						"phone_num": "99999999999",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccArmsAlertcontactCheckmap = map[string]string{

	"alert_contact_id": CHECKSET,

	"email": CHECKSET,

	"alert_contact_name": CHECKSET,

	"system_noc": CHECKSET,

	"ding_robot_webhook_url": CHECKSET,

	"phone_num": CHECKSET,
}

func AlibabacloudTestAccArmsAlertcontactBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
