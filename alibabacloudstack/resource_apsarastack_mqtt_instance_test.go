package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackMqttInstance_basic(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_mqtt_instance.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &OnsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeMqttInstance")
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(0, 1000)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tftestacc%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEMqttInstanceDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"max_sub":           "10000",
					"instance_name":     "${var.name}",
					"cluster_name":      "${data.alibabacloudstack_mqtt_clusters.anyone.clusters.0.id}",
					"max_conn":          "1000",
					"max_up_tps":        "500",
					"max_down_tps":      "1000",
					"remark":            "test",
					"store_instance_id": "${alibabacloudstack_ons_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name,
						"max_conn":          "1000",
						"max_sub":           "10000",
						"max_up_tps":        "500",
						"max_down_tps":      "1000",
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
			{
				Config: testAccConfig(map[string]interface{}{
					"max_sub":       "1000000",
					"instance_name": "${var.name}-updated",
					"max_conn":      "2000",
					"max_up_tps":    "1000",
					"max_down_tps":  "2000",
					"remark":        "updated test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": fmt.Sprintf("%s-updated", name),
						"max_conn":      "2000",
						"max_sub":       "1000000",
						"max_up_tps":    "1000",
						"max_down_tps":  "2000",
						"remark":        "updated test",
					}),
				),
			},
		},
	})
}

func resourceEMqttInstanceDependence(name string) string {

	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

%s

data "alibabacloudstack_mqtt_clusters" "anyone" {
}

`, name, OnsCommonTestCase)
}
