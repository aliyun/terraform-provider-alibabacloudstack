package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackHologramInstanceBackupPolicy_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_hologram_instance_backup_policy.default"

	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HologramService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeHologramInstanceBackupPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf_hologram_backup_config%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceHologramInstanceBackupPolicyDependence)

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
					"instance_id":        "${alibabacloudstack_hologram_instance.default.id}",
					"week":               []string{"0", "1", "2", "3", "4", "5", "6"},
					"hour":               16,
					"data_keep_quantity": 5,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"week.#":             "7",
						"hour":               "16",
						"data_keep_quantity": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"week":               []string{"0", "1", "2", "3"},
					"hour":               "11",
					"data_keep_quantity": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"week.#":             "4",
						"hour":               "11",
						"data_keep_quantity": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "false",
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

func resourceHologramInstanceBackupPolicyDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

%s

data "alibabacloudstack_hologram_clusters" "default" {
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_hologram_instance" "default" {
  	compute_type = "Standard"
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
	cpu = "intel"
	node = "2"
	vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
	vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
	instance_name = "${var.name}"
	cluster = "${data.alibabacloudstack_hologram_clusters.default.clusters.0.id}"
}

`, name, VSwitchCommonTestCase)
}
