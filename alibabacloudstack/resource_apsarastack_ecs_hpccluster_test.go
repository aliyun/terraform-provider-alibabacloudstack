package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEcsHpccluster0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_hpccluster.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsHpcclusterCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribehpcclustersRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secshpc_cluster%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsHpcclusterBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "Test For Terraform",

					"hpc_cluster_name": "rdktest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "Test For Terraform",

						"hpc_cluster_name": "rdktest",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"hpc_cluster_name": "rdkTestUpdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"hpc_cluster_name": "rdkTestUpdate",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "Test For Terraform Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "Test For Terraform Update",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccEcsHpcclusterCheckmap = map[string]string{

	"description": CHECKSET,

	"hpc_cluster_id": CHECKSET,

	"hpc_cluster_name": CHECKSET,
}

func AlibabacloudTestAccEcsHpcclusterBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
