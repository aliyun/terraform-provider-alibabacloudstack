package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackPolarDBReadWriteSplittingConnection_update(t *testing.T) {
	var connection *PolardbDescribedbinstancenetinfoResponse
	var primary *PolardbDescribedbinstanceattributeResponse
	var readonly *PolardbDescribedbinstanceattributeResponse

	resourceId := "alibabacloudstack_polardb_readwrite_splitting_connection.default"
	ra := resourceAttrInit(resourceId, map[string]string{})

	rc_connection := resourceCheckInitWithDescribeMethod(resourceId, &connection, func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoPolardbDescribedbinstancenetinfoRequest")
	rc_primary := resourceCheckInitWithDescribeMethod("alibabacloudstack_polardb_dbinstance.default", &primary, func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "Describedbinstances")
	rc_readonly := resourceCheckInitWithDescribeMethod("alibabacloudstack_polardb_readonly_instance.default", &readonly, func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoPolardbDescribedbinstanceattributeRequest")
	rand := getAccTestRandInt(10000, 999999)

	rac := resourceAttrCheckInit(rc_connection, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	prefix := fmt.Sprintf("t-con-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, prefix, resourcePolarDBReadWriteSplittingConfigDependence)
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
					"instance_id":       "${alibabacloudstack_polardb_readonly_instance.default.master_db_instance_id}",
					"connection_prefix": "${var.prefix}",
					"distribution_type": "Standard",
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
					"max_delay_time":    "300",
					"distribution_type": "Custom",
					"weight": `${map(
						"${alibabacloudstack_polardb_dbinstance.default.id}", "0",
						"${alibabacloudstack_polardb_readonly_instance.default.id}", "500"
					)}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					rc_primary.checkResourceExists(),
					rc_readonly.checkResourceExists(),
					testAccCheck(map[string]string{
						"max_delay_time":    "300",
						"weight.%":          "2",
						"distribution_type": "Custom",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":       "${alibabacloudstack_polardb_readonly_instance.default.master_db_instance_id}",
					"connection_prefix": "${var.prefix}",
					"distribution_type": "Standard",
					"max_delay_time":    "30",
					"weight":            REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port":              "3306",
						"distribution_type": "Standard",
						"weight.%":          REMOVEKEY,
						"max_delay_time":    "30",
						"instance_id":       CHECKSET,
						"connection_string": CHECKSET,
					}),
				),
			},
		},
	})
}

func resourcePolarDBReadWriteSplittingConfigDependence(prefix string) string {
	return fmt.Sprintf(`
	variable "creation" {
		default = "Rds"
	}
	variable "multi_az" {
		default = "false"
	}
	variable "name" {
		default = "tf-testAccDBInstance"
	}

	variable "prefix" {
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

	resource "alibabacloudstack_polardb_readonly_instance" "default" {
		master_db_instance_id = "${alibabacloudstack_polardb_dbinstance.default.id}"
		zone_id = "${alibabacloudstack_polardb_dbinstance.default.zone_id}"
		engine_version = "${alibabacloudstack_polardb_dbinstance.default.engine_version}"
		instance_type = "${alibabacloudstack_polardb_dbinstance.default.instance_type}"
		instance_storage = "${alibabacloudstack_polardb_dbinstance.default.instance_storage}"
		instance_name = "${var.name}"
		db_instance_storage_type = "${alibabacloudstack_polardb_dbinstance.default.storage_type}"
	}
`, prefix)
}
