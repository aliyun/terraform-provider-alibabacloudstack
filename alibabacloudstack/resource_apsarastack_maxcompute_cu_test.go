package alibabacloudstack

import (
	"fmt"

	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackMaxcomputeCu(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_maxcompute_cu.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &MaxcomputeService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeMaxcomputeCu")
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 9999)
	rac := resourceAttrCheckInit(rc, ra)
	name := fmt.Sprintf("tf_testAcck%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceMaxcomputeCuDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			// Currently does not support creating projects with sub-accounts
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cu_name":      "${var.name}",
					"cu_num":       "1",
					"cluster_name": "${data.alibabacloudstack_maxcompute_clusters.default.clusters.0.cluster}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cu_name":      name,
						"cu_num":       "1",
						"cluster_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cu_num": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cu_num": "2",
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

func resourceMaxcomputeCuDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alibabacloudstack_maxcompute_clusters" "default"{
	name_regex = "HYBRIDODPSCLUSTER-.*"
}

`, name)
}
