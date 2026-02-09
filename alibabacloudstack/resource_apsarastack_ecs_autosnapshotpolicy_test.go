package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEcsAutosnapshotpolicy0(t *testing.T) {
	var v *ecs.AutoSnapshotPolicy
	resourceId := "alibabacloudstack_ecs_autosnapshotpolicy.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsAutosnapshotpolicyCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribeautosnapshotpolicyexRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(1000, 9999)
	name := fmt.Sprintf("tf_auto_snapshot_policy%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsAutosnapshotpolicyBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"auto_snapshot_policy_name": "${var.name}",
					"repeat_weekdays":           []string{"1", "2", "3"},
					"retention_days":            -1,
					"time_points":               []string{"1", "22", "23"},
					"disk_ids":                  []string{"${alibabacloudstack_ecs_disk.default.0.id}", "${alibabacloudstack_ecs_disk.default.1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"auto_snapshot_policy_name": name,
						"repeat_weekdays.#":         "3",
						"retention_days":            "-1",
						"time_points.#":             "3",
						"disk_ids.#":                "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_ids": []string{"${alibabacloudstack_ecs_disk.default.1.id}"},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_ids.#":   "1",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_ids": []string{},
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_ids.#":   "0",
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_ids":                  []string{"${alibabacloudstack_ecs_disk.default.0.id}"},
					"auto_snapshot_policy_name": "${var.name}-update",
					"repeat_weekdays":           []string{"1", "2", "3", "4", "5"},
					"retention_days":            2,
					"time_points":               []string{"22", "23"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_ids.#":                "1",
						"auto_snapshot_policy_name": name + "-update",
						"repeat_weekdays.#":         "5",
						"retention_days":            "2",
						"time_points.#":             "2",
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

var AlibabacloudTestAccEcsAutosnapshotpolicyCheckmap = map[string]string{

	"auto_snapshot_policy_name": CHECKSET,

	"time_points.#": CHECKSET,

	"retention_days": CHECKSET,

	"repeat_weekdays.#": CHECKSET,
}

func AlibabacloudTestAccEcsAutosnapshotpolicyBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

%s

resource "alibabacloudstack_ecs_disk" "default" {
  count = 2
  disk_name = "${var.name}${count.index}"
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  size = 20
}

`, name, DataZoneCommonTestCase)
}
