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

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsauto_snapshot_policy%d", defaultRegionToTest, rand)

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

					"auto_snapshot_policy_name": "RDKTest",
					"repeat_weekdays":           []string{"1", "2", "3"},
					"retention_days":            -1,
					"time_points":               []string{"1", "22", "23"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"auto_snapshot_policy_name": "RDKTest",
						"repeat_weekdays.#":         "3",
						"retention_days":            "-1",
						"time_points.#":             "3",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"auto_snapshot_policy_name": "RDKTest-update",
					"repeat_weekdays":           []string{"1", "2", "3", "4", "5"},
					"retention_days":            2,
					"time_points":               []string{"22", "23"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_snapshot_policy_name": "RDKTest-update",
						"repeat_weekdays.#":         "5",
						"retention_days":            "2",
						"time_points.#":             "2",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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



`, name)
}
