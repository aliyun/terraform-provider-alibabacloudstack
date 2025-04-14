package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackPolardbConnectionConfigUpdate(t *testing.T) {
	var v *PolardbDescribedbinstancenetinfoResponse
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-testAccDBconnection%d", rand)

	var basicMap = map[string]string{
		"instance_id":       CHECKSET,
		"connection_string": CHECKSET,
		"port":              "3306",
		"ip_address":        CHECKSET,
	}
	resourceId := "alibabacloudstack_polardb_dbconnection.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbConnectionConfigDependence)
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
					"instance_id":       "${alibabacloudstack_polardb_dbinstance.instance.id}",
					"connection_prefix": "tftest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
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

func resourcePolardbConnectionConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s

	variable "creation" {
		default = "PolarDB"
	}

	variable "name" {
		default = "%s"
	}
	resource "alibabacloudstack_polardb_dbinstance" "instance" {
		engine            = "MySQL"
		engine_version    = "5.7"
		instance_name = "${var.name}"
		db_instance_storage_type= "local_ssd"
		db_instance_storage = 5
		db_instance_class = "rds.mysql.t1.small"
		zone_id= "${data.alibabacloudstack_zones.default.zones.0.id}"
		vswitch_id = "${alibabacloudstack_vswitch.default.id}"
	}
	`, RdsCommonTestCase, name)
}
