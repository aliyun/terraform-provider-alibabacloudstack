package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func TestAccAlibabacloudStackPolardbxInstanceReadonlyInstance(t *testing.T) {
	var v *PolardbxDescribedbinstanceattributeResponse

	resourceId := "alibabacloudstack_polardbx_readonly_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{"private_connection_string":CHECKSET})

	serviceFunc := func() interface{} {
		return &PolardbXService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DoPolardbxDescribedbinstanceattributeRequest")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf_acc_drds_polardb_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbxReadOnlyInstanceDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"primary_db_instance_id": "${local.polardbx_instance.id}",
					"description":            "${var.name}",
					"zone_id":                "${data.alibabacloudstack_zones.default.zones.0.id}",
					"engine_version":         "${local.polardbx_instance.engine_version}",
					"storage":                "${local.polardbx_instance.storage}",
					"vswitch_id":             "${alibabacloudstack_vpc_vswitch.default.id}",
					"cn_node_class":          "${data.alibabacloudstack_polardbx_instance_types.cn.instance_types.0.id}",
					"cn_node_count":          "2",
					"dn_node_class":          "${data.alibabacloudstack_polardbx_instance_types.dn.instance_types.0.id}",
					"dn_node_count":          "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":         name,
						"cn_node_count":       "2",
						"dn_node_count":       "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"primary_db_instance_id", "topology_type"},
			},
		},
	})
}

func resourcePolardbxReadOnlyInstanceDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

%s

 `, name, VSwitchCommonTestCase, PolardbxReadOrCreateCommonTestCase())
}
