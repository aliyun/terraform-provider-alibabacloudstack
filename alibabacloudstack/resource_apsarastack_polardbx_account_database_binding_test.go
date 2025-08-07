package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func TestAccAlibabacloudStackPolardbxAccountDatabaseBinding_basic0(t *testing.T) {
	var v *PolardbxAccount

	resourceId := "alibabacloudstack_polardbx_account_database_binding.default.0"
	ra := resourceAttrInit(resourceId, PolardbxAccountbasicMap)

	serviceFunc := func() interface{} {
		return &PolardbXService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolardbXAccountDBPrivilege")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("accdbbind%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbxAccountDatabaseBindingDependence)

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
					"db_name":      "${alibabacloudstack_polardbx_database.default[4].database_name}",
					"privilege":    "ReadWrite",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_name": name,
						"db_name":      CHECKSET,
						"privilege":    "ReadWrite",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"privilege": "ReadOnly",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"privilege": "ReadOnly",
					}),
				),
			},
		},
	})
}

func resourcePolardbxAccountDatabaseBindingDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "password" {
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

resource "alibabacloudstack_polardbx_database" "default" {
	count = 5
    instance_id  = "pxc-unrxhglptl45ih"
	database_name = "%s${count.index + 1}"
	encode = "utf8mb4"
	mode     = "auto"
}

resource "alibabacloudstack_polardbx_account" "default" {
    instance_id  =  "pxc-unrxhglptl45ih"
	account_name = "${var.name}"
	account_type = "Normal"
	password     = "${var.password}"
	description  = "Normal user"
}

resource "alibabacloudstack_polardbx_account_database_binding" "binding" {
    count = 4
    instance_id = "pxc-unrxhglptl45ih"
	account_name = "${var.name}"
	db_name = "${alibabacloudstack_polardbx_database.default[count.index].database_name}"
	privilege = "ReadWrite"
}


 `, name, getAccTestPassword(12), name)
}
