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
						"instance_storage":      "5",
						"engine_version":        "5.7",
						"engine":                "MySQL",
						"port":                  "3306",
						"instance_name":         name,
						"instance_type":         CHECKSET,
						"master_db_instance_id": CHECKSET,
						"zone_id":               CHECKSET,
						"connection_string":     CHECKSET,
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
					"instance_storage": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_storage": "10",
					}),
				),
			},
		},
	})

}

func TestAccAlibabacloudStackPolarDBReadonlyInstance_multi(t *testing.T) {
	var instance *PolardbDescribedbinstanceattributeResponse
	resourceId := "alibabacloudstack_polardb_readonly_instance.default.1"
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testAccDBInstance%d", rand)
	var PolarDBReadonlyMap = map[string]string{
		"instance_storage":      "5",
		"engine_version":        "5.7",
		"engine":                "MySQL",
		"port":                  "3306",
		"instance_name":         name,
		"instance_type":         CHECKSET,
		"parameters":            NOSET,
		"master_db_instance_id": CHECKSET,
		"zone_id":               CHECKSET,
		"connection_string":     CHECKSET,
	}
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
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":                    "1",
					"master_db_instance_id":    "${alibabacloudstack_polardb_dbinstance.default.id}",
					"zone_id":                  "${alibabacloudstack_polardb_dbinstance.default.zone_id}",
					"engine_version":           "${alibabacloudstack_polardb_dbinstance.default.engine_version}",
					"instance_type":            "${alibabacloudstack_polardb_dbinstance.default.instance_type}",
					"instance_storage":         "${alibabacloudstack_polardb_dbinstance.default.instance_storage}",
					"instance_name":            "${var.name}",
					"db_instance_storage_type": "${alibabacloudstack_polardb_dbinstance.default.storage_type}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourcePolarDBReadonlyInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
resource "alibabacloudstack_polardb_dbinstance" "default" {
  instance_storage = "5"
  instance_name = "${var.name}"
  storage_type = "local_ssd"
  engine = "MySQL"
  engine_version = "5.7"
  instance_type = "rds.mysql.t1.small"
}
`, name)
}
