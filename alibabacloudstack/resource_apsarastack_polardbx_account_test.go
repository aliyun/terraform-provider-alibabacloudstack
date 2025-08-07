package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func TestAccAlibabacloudStackPolardbxAccount_basic0(t *testing.T) {
	var v *PolardbxAccount

	resourceId := "alibabacloudstack_polardbx_account.default"
	ra := resourceAttrInit(resourceId, PolardbxAccountbasicMap)

	serviceFunc := func() interface{} {
		return &PolardbXService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DoPolardbxDescribeAccountListRequest")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf_acc_pldbx_account_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbxAccountDependence)

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
					"instance_id":  "pxc-unrxhglptl45ih",
					"account_name": "${var.name}",
					"password":     "${var.password}",
					"description":  "${var.name}",
					"db_privileges": []map[string]string{
						// {
						// 	"db_name":   "${alibabacloudstack_drds_database.default.0.drds_database_name}",
						// 	"privilege": "ReadOnly",
						// },
						{
							"db_name":   "${alibabacloudstack_polardbx_database.default1.database_name}",
							"privilege": "ReadOnly",
						},
						{
							"db_name":   "${alibabacloudstack_polardbx_database.default2.database_name}",
							"privilege": "ReadWrite",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_name":              name,
						"description":               name,
						"db_privileges.#":           "2",
						"db_privileges.0.db_name":   CHECKSET,
						"db_privileges.0.privilege": "ReadOnly",
						"db_privileges.1.db_name":   CHECKSET,
						"db_privileges.1.privilege": "ReadWrite",
					}),
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
					"password": "${var.password2}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_privileges": []map[string]string{
						{
							"db_name":   "${alibabacloudstack_polardbx_database.default1.database_name}",
							"privilege": "ReadWrite",
						},
						{
							"db_name":   "${alibabacloudstack_polardbx_database.default2.database_name}",
							"privilege": "DDLOnly",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_privileges.0.privilege": "ReadWrite",
						"db_privileges.1.privilege": "DDLOnly",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("%s_update", name),
					}),
				),
			},
		},
	})
}

func resourcePolardbxAccountDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "password" {
  default = "%s"
}

variable "password2" {
  default = "%s"
}

// resource "alibabacloudstack_polardbx_instance" "default" {
//  description = "testtf1111"
// 	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
// 	engine_version = "5.7"
// 	storage = "50"
// 	network_type = "vpc"
// 	vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
// 	vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
// 	cn_node_class = "polarx.x4.medium.2e"
// 	cn_node_count = "2"
// 	dn_node_class = "mysql.n4.medium.25"
// 	dn_node_count = "2"
// }

resource "alibabacloudstack_polardbx_database" "default1" {
    instance_id  = "pxc-unrxhglptl45ih"
	database_name = "testtf1"
	encode = "utf8mb4"
	mode     = "auto"
	description  = "testtf1"
	account_name = "admin"
	account_privilege = "ReadWrite"
}

resource "alibabacloudstack_polardbx_database" "default2" {
    instance_id  = "pxc-unrxhglptl45ih"
	database_name = "testtf2"
	encode = "utf8mb4"
	mode     = "auto"
	description  = "testtf2"
	account_name = "admin"
	account_privilege = "ReadWrite"
}

 `, name, getAccTestPassword(12), getAccTestPassword(10))
}

var PolardbxAccountbasicMap = map[string]string{
	"instance_id":  CHECKSET,
	"account_name": CHECKSET,
	"account_type": CHECKSET,
}
