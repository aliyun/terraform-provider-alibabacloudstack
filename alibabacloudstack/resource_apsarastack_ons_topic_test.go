package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackOnsTopic_basic(t *testing.T) {
	var v *Topic
	resourceId := "alibabacloudstack_ons_topic.default"
	ra := resourceAttrInit(resourceId, onsTopicBasicMap)
	serviceFunc := func() interface{} {
		return &OnsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-topicbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccOnsTopicConfigBasic)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alibabacloudstack_ons_instance.default.id}",
					"topic":        "${var.name}",
					"remark":       "Ons_topic",
					"message_type": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topic":  name,
						"remark": "Ons_topic",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

func testAccOnsTopicConfigBasic(name string) string {
	return fmt.Sprintf(`

variable "name" {
 default = "%s"
}

%s
`, name, OnsCommonTestCase)
}

var onsTopicBasicMap = map[string]string{
	"instance_id":  CHECKSET,
	"topic":        CHECKSET,
	"message_type": CHECKSET,
	"remark":       CHECKSET,
}
