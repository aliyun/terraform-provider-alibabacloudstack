package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackPolarDBReadWriteSplittingConnection_update(t *testing.T) {
	var connection map[string]interface {}

	resourceId := "alibabacloudstack_polardb_readwrite_splitting_connection.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"port":              "3306",
		"distribution_type": "Standard",
		"max_delay_time":    "30",
		"instance_id":       CHECKSET,
		"connection_string": CHECKSET,
	})

	rc := resourceCheckInitWithDescribeMethod(resourceId, &connection, func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDBReadWriteSplittingConnection")
	rand := getAccTestRandInt(10000, 999999)

	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	prefix := fmt.Sprintf("t-con-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, prefix, resourcePolarDBReadWriteSplittingConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":       "${alibabacloudstack_polardb_readonly_instance.default.master_db_instance_id}",
					"connection_id":     "${alibabacloudstack_polardb_proxy.default.db_proxy_endpoint_name}",
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
						"weight.%":          "2",
						"distribution_type": "Custom",
					}),
				),
			},
		},
	})
}

func resourcePolarDBReadWriteSplittingConfigDependence(name string) string {
	os.Unsetenv("ALIBABACLOUDSTACK_TEST_POLARDB_INSTNCE_ID")
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

%s
	resource "alibabacloudstack_polardb_readonly_instance" "default" {
		master_db_instance_id = "${alibabacloudstack_polardb_dbinstance.default.0.id}"
		zone_id = "${alibabacloudstack_polardb_dbinstance.default.0.zone_id}"
		engine_version = "${alibabacloudstack_polardb_dbinstance.default.0.engine_version}"
		instance_type = "${alibabacloudstack_polardb_dbinstance.default.0.instance_type}"
		instance_storage = "${alibabacloudstack_polardb_dbinstance.default.0.instance_storage}"
		instance_name = "${var.name}"
		db_instance_storage_type = "${alibabacloudstack_polardb_dbinstance.default.0.storage_type}"
	}
	
	locals {
	  dynamic_weight = {
		(alibabacloudstack_polardb_dbinstance.default.0.id): "0"
		(alibabacloudstack_polardb_readonly_instance.default.id): "500"
	  }
	}

	resource "alibabacloudstack_polardb_proxy" "default" {
		db_instance_id        = "${alibabacloudstack_polardb_dbinstance.default.0.id}"
		db_proxy_instance_num = "1"
		lifecycle {
			ignore_changes = [
			read_only_instance_max_delay_time
			]
		}
	}
`, name, PolarDBCommonTestCase("MySQL", false))
}
