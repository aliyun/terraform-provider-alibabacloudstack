package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_drds_account", &resource.Sweeper{
		Name: "alibabacloudstack_drds_account",
		F:    testSweepDRDSInstances,
	})
}

func TestAccAlibabacloudStackDrdsAccount_basic0(t *testing.T) {
	var v *DrdsDescribeinstanceAccount

	resourceId := "alibabacloudstack_drds_account.default"
	ra := resourceAttrInit(resourceId, drdsAccountbasicMap)

	serviceFunc := func() interface{} {
		return &DrdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf_acc_drds_db_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDrdsAccountDependence)

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
					"instance_id":       "${local.drds_instance_id}",
					"drds_account_name": "${var.name}",
					"password":          "${random_password.password.0.result}",
					"description":       name,
					"db_privileges": []map[string]string{
						{"db_name": "${alibabacloudstack_drds_database.default.0.drds_database_name}", "privilege": "R"},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"drds_account_name": name,
						"description":       name,
						"db_privileges.#":   "1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceId,
						"db_privileges.*",
						map[string]string{
							"db_name":   fmt.Sprintf("%s_%d", name, 0),
							"privilege": "R",
						},
					),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// password无法回读
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "${random_password.password.1.result}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_privileges": []map[string]string{
						{"db_name": "${alibabacloudstack_drds_database.default.0.drds_database_name}", "privilege": "RW"},
						{"db_name": "${alibabacloudstack_drds_database.default.1.drds_database_name}", "privilege": "R"},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_privileges.#": "2",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceId,
						"db_privileges.*",
						map[string]string{
							"db_name":   fmt.Sprintf("%s_%d", name, 0),
							"privilege": "RW",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceId,
						"db_privileges.*",
						map[string]string{
							"db_name":   fmt.Sprintf("%s_%d", name, 1),
							"privilege": "R",
						},
					),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_privileges": []map[string]string{
						{"db_name": "${alibabacloudstack_drds_database.default.1.drds_database_name}", "privilege": "R"},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_privileges.#": "1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceId,
						"db_privileges.*",
						map[string]string{
							"db_name":   fmt.Sprintf("%s_%d", name, 1),
							"privilege": "R",
						},
					),
				),
			},
		},
	})
}

func resourceDrdsAccountDependence(name string) string {
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
		zone_id             = data.alibabacloudstack_zones.default.zones.0.id
		db_instance_storage = "20"
		storage_type        = "local_ssd"
		category            = "HighAvailability"
		db_instance_class   = "rds.mysql.s1.small"
		drds_instance_id    = local.drds_instance_id
	}
	
	resource "alibabacloudstack_drds_database" "default" {
		count              = 2
		instance_id        = "${local.drds_instance_id}"
		drds_database_name = "${var.name}_${count.index}"
		password           = "${random_password.password.0.result}"
		rds_instance_ids   = [alibabacloudstack_drds_rds_instance.default.rds_instance_id,]
	}
	
`, name, os.Getenv("ALIBABACLOUDSTACK_TEST_EXISTED_DRDS_ID"), VSwitchCommonTestCase)
}

var drdsAccountbasicMap = map[string]string{
	"instance_id":       CHECKSET,
	"drds_account_name": CHECKSET,
	"host":              "%",
	"account_type":      "1",
}
