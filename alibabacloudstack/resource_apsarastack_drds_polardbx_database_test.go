package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func TestAccAlibabacloudStackDrdsPolardbxDatabase_basic0(t *testing.T) {
	var v *PolardbxDatabase

	resourceId := "alibabacloudstack_drds_polardbx_database.default"
	ra := resourceAttrInit(resourceId, map[string]string{})

	serviceFunc := func() interface{} {
		return &PolardbXService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DoPolardbxDescribeDbListRequest")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf_acc_polardbx_db_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDrdsPolardbxDatabaseDependence)

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
					"instance_id":       "${alibabacloudstack_drds_polardbx_instance.default.id}",
					"database_name":     "${var.name}",
					"encode":            "utf8mb4",
					"description":       "${var.name}",
					"account_name":      "${alibabacloudstack_drds_polardbx_account.default.account_name}",
					"account_privilege": "ReadWrite",
					"mode":              "auto",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"database_name":     name,
						"encode":            "utf8mb4",
						"description":       name,
						"account_name":      CHECKSET,
						"account_privilege": "ReadWrite",
						"mode":              "auto",
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
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// mode 无法回读
				ImportStateVerifyIgnore: []string{"mode"},
			},
		},
	})
}

func resourceDrdsPolardbxDatabaseDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "password" {
  default = "%s"
}

resource "alibabacloudstack_drds_polardbx_instance" "default" {
 description = "testtf1111"
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
	engine_version = "5.7"
	storage = "50"
	network_type = "vpc"
	vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
	vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
	cn_node_class = "polarx.x4.medium.2e"
	cn_node_count = "2"
	dn_node_class = "mysql.n4.medium.25"
	dn_node_count = "2"
}

resource "alibabacloudstack_drds_polardbx_account" "default" {
    instance_id  = alibabacloudstack_drds_polardbx_instance.default.id
	account_name = "testtf"
	account_type = "Normal"
	password     = "${var.password}"
	description  = "Normal user"
}

 `, name, getAccTestPassword(12))
}
