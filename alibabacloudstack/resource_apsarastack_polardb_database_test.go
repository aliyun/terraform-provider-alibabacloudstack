package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackPolardbDBDatabaseUpdate(t *testing.T) {
	var database *PolardbDescribedatabasesResponse
	resourceId := "alibabacloudstack_polardb_database.default"
	name := "tf-testaccdbdatabase_basic"

	var dbDatabaseBasicMap = map[string]string{
		"data_base_instance_id": CHECKSET,
		"data_base_name":        name,
		"character_set_name":    "utf8",
		"data_base_description": "",
	}

	ra := resourceAttrInit(resourceId, dbDatabaseBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &database, func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDBDatabase")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbDatabaseConfigDependence)
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
					"data_base_instance_id": "${alibabacloudstack_polardb_dbinstance.instance.id}",
					"data_base_name":        name,
					"character_set_name":    "utf8",
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
			{
				Config: testAccConfig(map[string]interface{}{
					"data_base_description": "from terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"data_base_description": "from terraform"}),
				),
			},
		},
	})

}

func resourcePolardbDatabaseConfigDependence(name string) string {
	return fmt.Sprintf(`


	variable "name" {
		default = "%s"
	}

	%s
	resource "alibabacloudstack_polardb_dbinstance" "instance" {
		engine            = "MySQL"
		engine_version    = "5.7"
		instance_name = "${var.name}"
		db_instance_storage_type= "local_ssd"
		db_instance_storage = 5
		db_instance_class = "rds.mysql.t1.small"
		zone_id= "${data.alibabacloudstack_zones.default.zones.0.id}"
		vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
	}`, name, VSwitchCommonTestCase)
}
