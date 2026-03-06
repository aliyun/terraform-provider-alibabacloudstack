package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackMqttGroup_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_mqtt_group.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &OnsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeOnsMqttGroup")
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(0, 1000)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("GID_tftestacc%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEMqttGroupDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"group_id":    "${var.name}",
					"instance_id": "${alibabacloudstack_mqtt_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_id":     name,
						"instance_id":  CHECKSET,
						"create_time":  CHECKSET,
						"update_time":  CHECKSET,
						"channel_name": CHECKSET,
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

func resourceEMqttGroupDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s
`, name, MqttCommonTestCase)
}
