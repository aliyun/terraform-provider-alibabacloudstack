package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackPolardbParameterGroup_basic(t *testing.T) {
	var v map[string]interface{}
	rand := getAccTestRandInt(10000, 999999)
	name := fmt.Sprintf("tf_paramgroup%d", rand)
	resourceId := "alibabacloudstack_polardb_parameter_group.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	serviceFunc := func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeParameterGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbParameterGroupConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},

		// module name
		IDRefreshName: resourceId,

		Providers: testAccProviders,
		// CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":               "mysql",
					"engine_version":       "8.0",
					"parameter_group_name": "${var.name}",
					"parameter_group_desc": "tf-testAccPolardbParameterGroupDesc",
					"parameters": map[string]interface{}{
						"wait_timeout": "3600",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "mysql",
						"engine_version":       "8.0",
						"parameter_group_name": name,
						"parameter_group_desc": "tf-testAccPolardbParameterGroupDesc",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parameters"},
				// parameters field is not returned by DescribeParameterGroups list API
			},
		},
	})
}

var AlibabacloudStackparameterMap = map[string]interface{}{
	"loose_multi_blocks_ddl_count": "11",
}

func resourcePolardbParameterGroupConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%v"
	}
	
	`, name)
}
