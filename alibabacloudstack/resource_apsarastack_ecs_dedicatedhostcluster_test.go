package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEcsdedicatedhostcluster0(t *testing.T) {
	var v *EcsDescribededicatedhostclustersResponse

	resourceId := "alibabacloudstack_ecs_dedicated_host_cluster.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsdedicatedhostclusterCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribededicatedhostclustersRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdedicated_hostcluster%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsdedicatedhostclusterBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		// CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"dedicated_host_cluster_name": "test-name",
					"zone_id":                     "${data.alibabacloudstack_zones.default.zones.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"dedicated_host_cluster_name": "test-name",
						"zone_id":                     CHECKSET,
					}),
				),
			},

//			{
//				Config: testAccConfig(map[string]interface{}{
//					"tags": map[string]string{
//						"Created": "TF",
//						"For":     "Test",
//					},
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(map[string]string{
//						"tags.%":       "2",
//						"tags.Created": "TF",
//						"tags.For":     "Test",
//					}),
//				),
//			},
//			{
//				Config: testAccConfig(map[string]interface{}{
//					"tags": map[string]string{
//						"Created": "TF-update",
//						"For":     "Test-update",
//					},
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(map[string]string{
//						"tags.%":       "2",
//						"tags.Created": "TF-update",
//						"tags.For":     "Test-update",
//					}),
//				),
//			},
//			{
//				Config: testAccConfig(map[string]interface{}{
//					"tags": REMOVEKEY,
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(map[string]string{
//						"tags.%":       "0",
//						"tags.Created": REMOVEKEY,
//						"tags.For":     REMOVEKEY,
//					}),
//				),
//			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlibabacloudTestAccEcsdedicatedhostclusterCheckmap = map[string]string{

	"dedicated_host_cluster_id": CHECKSET,

	// "region_id": CHECKSET,

	"zone_id": CHECKSET,
}

func AlibabacloudTestAccEcsdedicatedhostclusterBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

%s

`, name, DataZoneCommonTestCase)
}
