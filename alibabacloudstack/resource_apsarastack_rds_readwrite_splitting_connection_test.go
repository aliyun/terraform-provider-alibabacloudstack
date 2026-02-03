package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var DBReadWriteMap = map[string]string{
	"port":              "3306",
	"distribution_type": "Standard",
	"max_delay_time":    "30",
	"instance_id":       CHECKSET,
	"connection_string": CHECKSET,
}

func TestAccAlibabacloudStackDBReadWriteSplittingConnection_update(t *testing.T) {
	var connection map[string]interface{}

	resourceId := "alibabacloudstack_db_read_write_splitting_connection.default"
	ra := resourceAttrInit(resourceId, DBReadWriteMap)

	rc := resourceCheckInitWithDescribeMethod(resourceId, &connection, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDBReadWriteSplittingConnection")
	rand := getAccTestRandInt(10000, 999999)

	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	prefix := fmt.Sprintf("tfproxy%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, prefix, resourceDBReadWriteSplittingConfigDependence)
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
					"instance_id":       "${alibabacloudstack_db_readonly_instance.default.master_db_instance_id}",
					"connection_id":     "${alibabacloudstack_db_connection.default.id}",
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
					"weight":            `${local.dynamic_weight}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_delay_time":    "300",
						"distribution_type": "Custom",
					}),
				),
			},
		},
	})
}

func resourceDBReadWriteSplittingConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	

	%s
	
	%s

	resource "alibabacloudstack_db_readonly_instance" "default" {
		master_db_instance_id = "${alibabacloudstack_db_instance.default.id}"
		zone_id = "${alibabacloudstack_db_instance.default.zone_id}"
		engine_version = "${alibabacloudstack_db_instance.default.engine_version}"
		instance_type = "${alibabacloudstack_db_instance.default.instance_type}"
		instance_storage = "${alibabacloudstack_db_instance.default.instance_storage}"
		instance_name = "${var.name}_ro"
		vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
		db_instance_storage_type = "${alibabacloudstack_db_instance.default.storage_type}"
	}
	
	locals {
	  dynamic_weight = {
		(alibabacloudstack_db_instance.default.id): "0"
		(alibabacloudstack_db_readonly_instance.default.id): "500"
	  }
	}
	
	resource "alibabacloudstack_db_connection" "default" {
		instance_id       = "${alibabacloudstack_db_instance.default.id}"
		connection_prefix = "${var.name}"
	}
	`, name, VSwitchCommonTestCase, RdsMysqlCommonTestCase())
}
