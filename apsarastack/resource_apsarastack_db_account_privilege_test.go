package apsarastack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackDBAccountPrivilege_mysql(t *testing.T) {

	var v *rds.DBInstanceAccount
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sdnsrecordbasic%v.abc", defaultRegionToTest, rand)
	resourceId := "apsarastack_db_account_privilege.default"
	var basicMap = map[string]string{
		"instance_id":  CHECKSET,
		"account_name": "tftestprivilege",
		"privilege":    "ReadOnly",
		"db_names.#":   "2",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBAccountPrivilege")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBAccountPrivilegeConfigDependenceForMySql)

	resource.Test(t, resource.TestCase{
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
					"instance_id":  "${apsarastack_db_instance.default.id}",
					"account_name": "${apsarastack_db_account.default.name}",
					"privilege":    "ReadOnly",
					"db_names":     "${apsarastack_db_database.default.*.name}",
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
					"instance_id":  "${apsarastack_db_instance.default.id}",
					"account_name": "${apsarastack_db_account.default.name}",
					"privilege":    "ReadOnly",
					"db_names":     []string{"${apsarastack_db_database.default.0.name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_names.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${apsarastack_db_instance.default.id}",
					"account_name": "${apsarastack_db_account.default.name}",
					"privilege":    "ReadOnly",
					"db_names":     "${apsarastack_db_database.default.*.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_names.#": "2",
					}),
				),
			},
		},
	})

}

//

func resourceDBAccountPrivilegeConfigDependenceForMySql(name string) string {
	return fmt.Sprintf(`
%s
provider "apsarastack" {
	assume_role {}
}
	variable "creation" {
		default = "Rds"
	}

	variable "name" {
		default = "%s"
	}


resource "apsarastack_db_instance" "default" {
		engine               = "MySQL"
        engine_version       = "5.6"
        instance_type        = "rds.mysql.s2.large"
	    instance_storage     = "30"
		vswitch_id = "${apsarastack_vswitch.default.id}"
	    instance_name = "${var.name}"
	    storage_type         = "local_ssd"

	}



	resource "apsarastack_db_database" "default" {
	  count = 2
	  instance_id = "${apsarastack_db_instance.default.id}"
	  name = "tfaccountpri_${count.index}"
	  description = "from terraform"
	  character_set        = "utf8"
	}

	resource "apsarastack_db_account" "default" {
	  instance_id = "${apsarastack_db_instance.default.id}"
	  name = "tftestprivilege"
	  password = "inputYourCodeHere"
	  description = "from terraform"
	}
`, RdsCommonTestCase, name)
}

//func TestAccApsaraStackDBAccountPrivilege_PostgreSql(t *testing.T) {
//
//	var v *rds.DBInstanceAccount
//	rand := acctest.RandInt()
//	name := fmt.Sprintf("tf-testacc%sdnsrecordbasic%v.abc", defaultRegionToTest, rand)
//	resourceId := "apsarastack_db_account_privilege.default"
//	var basicMap = map[string]string{
//		"instance_id":  CHECKSET,
//		"account_name": "tftestprivilege",
//		"privilege":    "ReadOnly",
//		"db_names.#":   "1",
//	}
//	ra := resourceAttrInit(resourceId, basicMap)
//	serviceFunc := func() interface{} {
//		return &RdsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
//	}
//	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBAccountPrivilege")
//	rac := resourceAttrCheckInit(rc, ra)
//
//	testAccCheck := rac.resourceAttrMapUpdateSet()
//	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBAccountPrivilegeConfigDependenceForPostgreSql)
//
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			testAccPreCheck(t)
//		},
//
//		// module name
//		IDRefreshName: resourceId,
//
//		Providers:    testAccProviders,
//		CheckDestroy: rac.checkResourceDestroy(),
//		Steps: []resource.TestStep{
//			{
//				Config: testAccConfig(map[string]interface{}{
//					"instance_id":  "${apsarastack_db_instance.default.id}",
//					"account_name": "${apsarastack_db_account.default.name}",
//					"privilege":    "ReadOnly",
//					"db_names":     []string{"${apsarastack_db_database.default.0.name}"},
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(nil),
//				),
//			},
//			{
//				ResourceName:      resourceId,
//				ImportState:       true,
//				ImportStateVerify: true,
//			},
//			{
//				Config: testAccConfig(map[string]interface{}{
//					"instance_id":  "${apsarastack_db_instance.default.id}",
//					"account_name": "${apsarastack_db_account.default.name}",
//					"privilege":    "ReadOnly",
//					"db_names":     "${apsarastack_db_database.default.*.name}",
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(map[string]string{
//						"db_names.#": "2",
//					}),
//				),
//			},
//		},
//	})
//
//}

//func resourceDBAccountPrivilegeConfigDependenceForPostgreSql(name string) string {
//	return fmt.Sprintf(`
//%s
//	variable "creation" {
//		default = "Rds"
//	}
//
//	variable "name" {
//		default = "%s"
//	}
//
//	resource "apsarastack_db_instance" "default" {
//		engine = "PostgreSQL"
//		engine_version = "10.0"
//		instance_type = "pg.n2.large.1"
//		instance_storage = "30"
//		vswitch_id = "${apsarastack_vswitch.default.id}"
//		instance_name = "${var.name}"
//	}
//
//	resource "apsarastack_db_database" "default" {
//	  count = 2
//	  instance_id = "${apsarastack_db_instance.default.id}"
//	  name = "tfaccountpri_${count.index}"
//	  description = "from terraform"
//      character_set = "UTF8"
//	}
//
//	resource "apsarastack_db_account" "default" {
//	  instance_id = "${apsarastack_db_instance.default.id}"
//	  name = "tftestprivilege"
//	  password = "inputYourCodeHere"
//	  description = "from terraform"
//	}
//`, RdsCommonTestCase, name)
//}
