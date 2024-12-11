package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDatahubSubscription0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_datahub_subscription.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccDatahubSubscriptionCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DatahubService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoDatahubGetsubscriptionoffsetRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdata_hubsubscription%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccDatahubSubscriptionBasicdependence)
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

					"comment": "test",

					"project_name": "datahub_pop_test",

					"application": "pop",

					"topic_name": "test_topic",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"comment": "test",

						"project_name": "datahub_pop_test",

						"application": "pop",

						"topic_name": "test_topic",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccDatahubSubscriptionCheckmap = map[string]string{

	"status": CHECKSET,

	"comment": CHECKSET,

	"project_name": CHECKSET,

	"creator": CHECKSET,

	"topic_name": CHECKSET,

	"total_count": CHECKSET,

	"type": CHECKSET,

	"subscription_id": CHECKSET,

	"subscription_offset": CHECKSET,

	"application": CHECKSET,
}

func AlibabacloudTestAccDatahubSubscriptionBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
