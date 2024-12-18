package alibabacloudstack

import (
	"fmt"
	"testing"

	

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var DBReadWriteMap = map[string]string{
	"port":              "3306",
	"distribution_type": "Standard",
	"weight":            NOSET,
	"max_delay_time":    "30",
	"instance_id":       CHECKSET,
	"connection_string": CHECKSET,
}

func TestAccAlibabacloudStackDBReadWriteSplittingConnection_update(t *testing.T) {
	var connection = &rds.DBInstanceNetInfo{}
	var primary = &rds.DBInstanceAttribute{}
	var readonly = &rds.DBInstanceAttribute{}

	resourceId := "alibabacloudstack_db_read_write_splitting_connection.default"
	ra := resourceAttrInit(resourceId, DBReadWriteMap)

	rc_connection := resourceCheckInitWithDescribeMethod(resourceId, &connection, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDBReadWriteSplittingConnection")
	rc_primary := resourceCheckInitWithDescribeMethod("alibabacloudstack_db_instance.default", &primary, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDBInstance")
	rc_readonly := resourceCheckInitWithDescribeMethod("alibabacloudstack_db_readonly_instance.default", &readonly, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDBReadonlyInstance")
	rand := getAccTestRandInt(10000, 999999)

	rac := resourceAttrCheckInit(rc_connection, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	prefix := fmt.Sprintf("t-con-%d", rand)
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
						"${alibabacloudstack_db_instance.default.id}", "0",
						"${alibabacloudstack_db_readonly_instance.default.id}", "500"
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
					"instance_id":       "${alibabacloudstack_db_readonly_instance.default.master_db_instance_id}",
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

func resourceDBReadWriteSplittingConfigDependence(prefix string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
	}
	variable "multi_az" {
		default = "false"
	}
	variable "name" {
		default = "tf-testAccDBInstance_vpc"
	}

	variable "prefix" {
		default = "%s"
	}

	resource "alibabacloudstack_db_instance" "default" {
		engine = "MySQL"
		engine_version = "5.6"
		instance_type = "rds.mysql.s2.large"
		instance_storage = "30"
		instance_charge_type = "Postpaid"
		instance_name = "${var.name}"
		vswitch_id = "${alibabacloudstack_vswitch.default.id}"
		storage_type = "local_ssd"
		security_ips = ["10.168.1.12", "100.69.7.112"]
	}

	resource "alibabacloudstack_db_readonly_instance" "default" {
		master_db_instance_id = "${alibabacloudstack_db_instance.default.id}"
		zone_id = "${alibabacloudstack_db_instance.default.zone_id}"
		engine_version = "${alibabacloudstack_db_instance.default.engine_version}"
		instance_type = "${alibabacloudstack_db_instance.default.instance_type}"
		instance_storage = "${alibabacloudstack_db_instance.default.instance_storage}"
		instance_name = "${var.name}_ro"
		vswitch_id = "${alibabacloudstack_vswitch.default.id}"
		db_instance_storage_type = "${alibabacloudstack_db_instance.default.storage_type}"
	}
`, RdsCommonTestCase, prefix)
}
