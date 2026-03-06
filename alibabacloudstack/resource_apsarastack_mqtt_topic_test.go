package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackMqttTopic_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_mqtt_topic.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &OnsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeOnsMqttTopic")
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(0, 1000)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tftestacc%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEMqttTopicDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"topic":             "${var.name}",
					"order_type":        "1",
					"remark":            "test",
					"store_instance_id": "${alibabacloudstack_mqtt_instance.default.store_instance_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topic":             name,
						"order_type":        "1",
						"remark":            "test",
						"store_instance_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// read api have not attr: "cluster_name"
				ImportStateVerifyIgnore: []string{"cluster_name"},
			},
		},
	})
}

func resourceEMqttTopicDependence(name string) string {

	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

%s

		`, name, MqttCommonTestCase)
}
