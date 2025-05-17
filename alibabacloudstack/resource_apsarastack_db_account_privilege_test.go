package alibabacloudstack

import (
	"fmt"
	"testing"

	

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDBAccountPrivilege_mysql(t *testing.T) {

	var v *rds.DBInstanceAccount
	rand := getAccTestRandInt(10000,20000)
	name := fmt.Sprintf("tf-testacc%sdnsrecordbasic%v.abc", defaultRegionToTest, rand)
	resourceId := "alibabacloudstack_db_account_privilege.default"
	var basicMap = map[string]string{
		"instance_id":  CHECKSET,
		"account_name": "tftestprivilege",
		"privilege":    "ReadOnly",
		"db_names.#":   "2",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBAccountPrivilege")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBAccountPrivilegeConfigDependenceForMySql)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		// CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alibabacloudstack_db_instance.default.id}",
					"account_name": "${alibabacloudstack_db_account.default.name}",
					"privilege":    "ReadOnly",
					"db_names":     "${alibabacloudstack_db_database.default.*.name}",
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
					"instance_id":  "${alibabacloudstack_db_instance.default.id}",
					"account_name": "${alibabacloudstack_db_account.default.name}",
					"privilege":    "ReadOnly",
					"db_names":     []string{"${alibabacloudstack_db_database.default.0.name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_names.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alibabacloudstack_db_instance.default.id}",
					"account_name": "${alibabacloudstack_db_account.default.name}",
					"privilege":    "ReadOnly",
					"db_names":     "${alibabacloudstack_db_database.default.*.name}",
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

	variable "creation" {
		default = "Rds"
	}

	variable "name" {
		default = "%s"
	}


resource "alibabacloudstack_db_instance" "default" {
		engine               = "MySQL"
        engine_version       = "5.6"
        instance_type        = "rds.mysql.s2.large"
	    instance_storage     = "30"
		vswitch_id = "${alibabacloudstack_vswitch.default.id}"
	    instance_name = "${var.name}"
	    storage_type         = "local_ssd"

	}



	resource "alibabacloudstack_db_database" "default" {
	  count = 2
	  instance_id = "${alibabacloudstack_db_instance.default.id}"
	  name = "tfaccountpri_${count.index}"
	  description = "from terraform"
	  character_set        = "utf8"
	}

	resource "alibabacloudstack_db_account" "default" {
	  instance_id = "${alibabacloudstack_db_instance.default.id}"
	  name = "tftestprivilege"
	  password = "%s"
	  description = "from terraform"
	}
`, RdsCommonTestCase, name, GeneratePassword())
}

//func TestAccAlibabacloudStackDBAccountPrivilege_PostgreSql(t *testing.T) {
//
//	var v *rds.DBInstanceAccount
//	rand := getAccTestRandInt(10000,20000)
//	name := fmt.Sprintf("tf-testacc%sdnsrecordbasic%v.abc", defaultRegionToTest, rand)
//	resourceId := "alibabacloudstack_db_account_privilege.default"
//	var basicMap = map[string]string{
//		"instance_id":  CHECKSET,
//		"account_name": "tftestprivilege",
//		"privilege":    "ReadOnly",
//		"db_names.#":   "1",
//	}
//	ra := resourceAttrInit(resourceId, basicMap)
//	serviceFunc := func() interface{} {
//		return &RdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
//	}
//	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBAccountPrivilege")
//	rac := resourceAttrCheckInit(rc, ra)
//
//	testAccCheck := rac.resourceAttrMapUpdateSet()
//	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBAccountPrivilegeConfigDependenceForPostgreSql)
//
//	ResourceTest(t, resource.TestCase{
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
//					"instance_id":  "${alibabacloudstack_db_instance.default.id}",
//					"account_name": "${alibabacloudstack_db_account.default.name}",
//					"privilege":    "ReadOnly",
//					"db_names":     []string{"${alibabacloudstack_db_database.default.0.name}"},
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
//					"instance_id":  "${alibabacloudstack_db_instance.default.id}",
//					"account_name": "${alibabacloudstack_db_account.default.name}",
//					"privilege":    "ReadOnly",
//					"db_names":     "${alibabacloudstack_db_database.default.*.name}",
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
//	resource "alibabacloudstack_db_instance" "default" {
//		engine = "PostgreSQL"
//		engine_version = "10.0"
//		instance_type = "pg.n2.large.1"
//		instance_storage = "30"
//		vswitch_id = "${alibabacloudstack_vswitch.default.id}"
//		instance_name = "${var.name}"
//	}
//
//	resource "alibabacloudstack_db_database" "default" {
//	  count = 2
//	  instance_id = "${alibabacloudstack_db_instance.default.id}"
//	  name = "tfaccountpri_${count.index}"
//	  description = "from terraform"
//      character_set = "UTF8"
//	}
//
//	resource "alibabacloudstack_db_account" "default" {
//	  instance_id = "${alibabacloudstack_db_instance.default.id}"
//	  name = "tftestprivilege"
//	  password = "inputYourCodeHere"
//	  description = "from terraform"
//	}
//`, RdsCommonTestCase, name)
//}
