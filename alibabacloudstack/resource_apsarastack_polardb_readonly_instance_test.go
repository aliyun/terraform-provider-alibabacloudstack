package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackPolarDBReadonlyInstance_update(t *testing.T) {
	var instance *PolardbDescribedbinstanceattributeResponse
	resourceId := "alibabacloudstack_polardb_readonly_instance.default"
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testAccDBInstance%d", rand)
	var PolarDBReadonlyMap = map[string]string{}
	ra := resourceAttrInit(resourceId, PolarDBReadonlyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoPolardbDescribedbinstanceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBReadonlyInstanceConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"master_db_instance_id":    "${alibabacloudstack_polardb_dbinstance.default.id}",
					"zone_id":                  "${alibabacloudstack_polardb_dbinstance.default.zone_id}",
					"engine_version":           "${alibabacloudstack_polardb_dbinstance.default.engine_version}",
					"instance_type":            "${alibabacloudstack_polardb_dbinstance.default.instance_type}",
					"instance_storage":         "${alibabacloudstack_polardb_dbinstance.default.instance_storage}",
					"instance_name":            "${var.name}",
					"db_instance_storage_type": "${alibabacloudstack_polardb_dbinstance.default.storage_type}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":         name,
						"instance_type":         CHECKSET,
						"master_db_instance_id": CHECKSET,
						"zone_id":               CHECKSET,
						"connection_string":     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ssl": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ssl": "true",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${var.name}_ro",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_ro",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_storage": "${local.new_instance_storage}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_storage": CHECKSET,
					}),
				),
			},
//			{
//				Config: testAccConfig(map[string]interface{}{
//					"instance_type": "${data.alibabacloudstack_polardb_instance_types.update.instance_types.1.id}",
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(map[string]string{
//						"instance_storage": CHECKSET,
//					}),
//				),
//			},
		},
	})

}

func resourcePolarDBReadonlyInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
%s

	locals {
		new_instance_storage = data.alibabacloudstack_polardb_instance_types.default.instance_types.0.storage_min + 10
		new_cpu = data.alibabacloudstack_polardb_instance_types.default.instance_types.0.cpu * 2
	}

	data "alibabacloudstack_polardb_instance_types" "update" {
	  engine               = data.alibabacloudstack_polardb_instance_types.default.instance_types.0.engine
	  engine_version       = data.alibabacloudstack_polardb_instance_types.default.instance_types.0.engine_version
	  sorted_by            = "Memory"
	  series               = data.alibabacloudstack_polardb_instance_types.default.instance_types.0.series
	  cpu                  = local.new_cpu
	}

`, name, PolarDBMysqlCommonTestCase(false))
}
