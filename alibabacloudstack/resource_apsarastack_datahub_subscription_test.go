package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
)

func TestAccAlibabacloudStackDatahubSubscription0(t *testing.T) {
	var v *datahub.GetSubscriptionResult

	resourceId := "alibabacloudstack_datahub_subscription.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccDatahubSubscriptionCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DatahubService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoDatahubGetsubscriptionoffsetRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf_testacc_datahub_sub%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccDatahubSubscriptionBasicdependence)
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

					"comment": name,

					"project_name": "${alibabacloudstack_datahub_project.default.name}",

					"topic_name": "${alibabacloudstack_datahub_topic.default.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"comment": name,

						"project_name": name,

						"topic_name": name,
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccDatahubSubscriptionCheckmap = map[string]string{

	"comment": CHECKSET,

	"project_name": CHECKSET,

	"topic_name": CHECKSET,

	"sub_id": CHECKSET,
}

func AlibabacloudTestAccDatahubSubscriptionBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alibabacloudstack_datahub_project" "default" {
    comment = "test"
    name = var.name
}

resource "alibabacloudstack_datahub_topic" "default" {
  name = var.name
  comment = "test"
  record_type = "BLOB"
  project_name = "${alibabacloudstack_datahub_project.default.name}"
}

`, name)
}
