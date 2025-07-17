package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_drds_database", &resource.Sweeper{
		Name: "alibabacloudstack_drds_database",
		F:    testSweepDRDSInstances,
	})
}

func TestAccAlibabacloudStackDrdsDatabase_basic0(t *testing.T) {
	var v *DrdsDescribedrdsdbResponse

	resourceId := "alibabacloudstack_drds_database.default"
	ra := resourceAttrInit(resourceId, drdsDatabasebasicMap)

	serviceFunc := func() interface{} {
		return &DrdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf_acc_drds_db_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDrdsDatabaseDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

		},
		// module name
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":      "${local.drds_instance_id}",
					"db_name":          name,
					"password":         "${random_password.password.0.result}",
					"rds_instance_ids": []string{"${alibabacloudstack_drds_rds_instance.default.0.rds_instance_id}"},
					"ip_white_list": map[string]string{
						"test1": "127.0.0.1,192.168.1.1",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rds_instance_ids.#":  "1",
						"ip_white_list.test1": "192.168.1.1, 127.0.0.1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rds_instance_ids": []string{"${alibabacloudstack_drds_rds_instance.default.0.rds_instance_id}",
						"${alibabacloudstack_drds_rds_instance.default.1.rds_instance_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rds_instance_ids.#": "2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: false,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "${random_password.password.1.result}",
					"ip_white_list": map[string]string{
						"test1": "127.0.0.2,192.168.1.1",
						"test2": "192.168.2.1",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_white_list.test1": "192.168.1.1, 127.0.0.2",
						"ip_white_list.test2": "192.168.2.1",
					}),
				),
			},
		},
	})
}

func resourceDrdsDatabaseDependence(name string) string {
	return fmt.Sprintf(`

	variable "name" {
		default = "%s"
	}
	
	variable "existed_drds_instance" {
		default = "%s"
	}
	
	locals {
		create_drds_instance_count = var.existed_drds_instance == "" ? 1: 0
	}
	
	variable "instance_series" {
		default = "drds.sn2.4c16g"
	}
	
	resource "random_password" "password" {
		count            = 2
		length           = 12
		special          = true
		override_special = "_"
		min_lower        = 1
		min_upper        = 1
		min_numeric      = 1
	}
	
	%s
	
	resource "alibabacloudstack_drds_instance" "default" {
		count                = local.create_drds_instance_count
		description          = "${var.name}"
		zone_id              = "${alibabacloudstack_vpc_vswitch.default.availability_zone}"
		instance_series      = "${var.instance_series}"
		instance_charge_type = "PostPaid"
		vswitch_id           = "${alibabacloudstack_vpc_vswitch.default.id}"
		specification        = "drds.sn2.4c16g.8C32G"
	}
	
	locals {
		drds_instance_id = var.existed_drds_instance == "" ? alibabacloudstack_drds_instance.default.0.id: var.existed_drds_instance
	}
	
	resource "alibabacloudstack_drds_rds_instance" "default" {
		count               = 2
		zone_id             = data.alibabacloudstack_zones.default.zones.0.id
		db_instance_storage = "20"
		storage_type        = "local_ssd"
		category            = "HighAvailability"
		db_instance_class   = "rds.mysql.s1.small"
		drds_instance_id    = local.drds_instance_id
	}
	
`, name, os.Getenv("ALIBABACLOUDSTACK_TEST_EXISTED_DRDS_ID"), VSwitchCommonTestCase)
}

var drdsDatabasebasicMap = map[string]string{
	"instance_id":  CHECKSET,
	"db_name":      CHECKSET,
	"split_mode":   "HORIZONTAL",
	"encode":       "utf8",
	"create_time":  CHECKSET,
	"status":       CHECKSET,
	"storage_type": CHECKSET,
}
