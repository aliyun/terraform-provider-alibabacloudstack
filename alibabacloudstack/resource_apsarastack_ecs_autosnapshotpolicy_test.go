package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEcsAutosnapshotpolicy0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_autosnapshotpolicy.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsAutosnapshotpolicyCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribeautosnapshotpolicyexRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsauto_snapshot_policy%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsAutosnapshotpolicyBasicdependence)
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

					"auto_snapshot_policy_name": "RDKTest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"auto_snapshot_policy_name": "RDKTest",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
		},
	})
}

var AlibabacloudTestAccEcsAutosnapshotpolicyCheckmap = map[string]string{

	"status": CHECKSET,

	"time_points": CHECKSET,

	"volume_nums": CHECKSET,

	"resource_group_id": CHECKSET,

	"create_time": CHECKSET,

	"auto_snapshot_policy_id": CHECKSET,

	"retention_days": CHECKSET,

	"repeat_weekdays": CHECKSET,

	"disk_nums": CHECKSET,

	"copied_snapshots_retention_days": CHECKSET,

	"target_copy_regions": CHECKSET,

	"enable_cross_region_copy": CHECKSET,

	"region_id": CHECKSET,

	"auto_snapshot_policy_name": CHECKSET,

	"tags": CHECKSET,
}

func AlibabacloudTestAccEcsAutosnapshotpolicyBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
