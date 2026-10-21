package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackPolardbClusterParameterGroup_basic(t *testing.T) {
	var v map[string]interface{}
	rand := getAccTestRandInt(10000, 999999)
	name := fmt.Sprintf("tf_cluster_pg%d", rand)
	resourceId := "alibabacloudstack_polardb_cluster_parameter_group.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	serviceFunc := func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeClusterParameterGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbClusterParameterGroupConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},

		IDRefreshName: resourceId,

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":               "MySQL",
					"engine_version":       "8.0",
					"parameter_group_name": "${var.name}",
					"parameter_group_desc": "tf-testAccPolardbClusterParameterGroupDesc",
					"parameters": map[string]interface{}{
						"wait_timeout": "3600",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "MySQL",
						"engine_version":       "8.0",
						"parameter_group_name": name,
						"parameter_group_desc": "tf-testAccPolardbClusterParameterGroupDesc",
						"parameters.%":         "1",
						"parameters.wait_timeout": "3600",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parameters"},
			},
		},
	})
}

func resourcePolardbClusterParameterGroupConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

`, name)
}
