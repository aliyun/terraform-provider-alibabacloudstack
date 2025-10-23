package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackHologramInstance_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_hologram_instance.default"

	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HologramService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeHologramInstance")
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf_hologram_instance%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceHologramInstanceDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},

		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"compute_type":   "Standard",
					"zone_id":        "${data.alibabacloudstack_zones.default.zones.0.id}",
					"cpu":            "intel",
					"is_new_feature": "true",
					"node":           "2",
					"vpc_id":         "${alibabacloudstack_vpc_vpc.default.id}",
					"vswitch_id":     "${alibabacloudstack_vpc_vswitch.default.id}",
					"instance_name":  "testtf",
					"cluster":        "HologresSSDCluster-A-20250927-012d",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compute_type":   "Standard",
						"cpu":            "intel",
						"is_new_feature": "true",
						"node":           "2",
						"instance_name":  "testtf",
						"cluster":        "HologresSSDCluster-A-20250927-012d",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node":          "4",
					"instance_name": "testtf_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node":          "4",
						"instance_name": "testtf_update",
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

func resourceHologramInstanceDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

%s

`, name, VSwitchCommonTestCase)
}
