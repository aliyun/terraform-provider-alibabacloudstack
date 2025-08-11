package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func TestAccAlibabacloudStackPolardbxAccountDatabaseBinding_basic0(t *testing.T) {
	var v []map[string]string

	resourceId := "alibabacloudstack_polardbx_account_database_binding.default"
	ra := resourceAttrInit(resourceId, map[string]string{})

	serviceFunc := func() interface{} {
		return &PolardbXService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolardbXAccountDBPrivilege")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
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
					"instance_id":  "${local.polardbx_instance.id}",
					"account_name": "${alibabacloudstack_polardbx_account.default.account_name}",
					"db_privileges": []map[string]interface{}{
						{
							"db_name":   "${alibabacloudstack_polardbx_database.default.0.database_name}",
							"privilege": "ReadOnly",
						},
						{
							"db_name":   "${alibabacloudstack_polardbx_database.default.1.database_name}",
							"privilege": "ReadWrite",
						},
						{
							"db_name":   "${alibabacloudstack_polardbx_database.default.2.database_name}",
							"privilege": "DDLOnly",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_privileges.#":           "3",
						"db_privileges.0.db_name":   fmt.Sprintf("%s0", name),
						"db_privileges.0.privilege": "ReadOnly",
						"db_privileges.1.db_name":   fmt.Sprintf("%s1", name),
						"db_privileges.1.privilege": "ReadWrite",
						"db_privileges.2.db_name":   fmt.Sprintf("%s2", name),
						"db_privileges.2.privilege": "DDLOnly",
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
					"db_privileges": []map[string]interface{}{
						{
							"db_name":   "${alibabacloudstack_polardbx_database.default.0.database_name}",
							"privilege": "DDLOnly",
						},
						{
							"db_name":   "${alibabacloudstack_polardbx_database.default.1.database_name}",
							"privilege": "ReadOnly",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_privileges.#":           "2",
						"db_privileges.0.db_name":   fmt.Sprintf("%s0", name),
						"db_privileges.0.privilege": "DDLOnly",
						"db_privileges.1.db_name":   fmt.Sprintf("%s1", name),
						"db_privileges.1.privilege": "ReadOnly",
						"db_privileges.2.db_name":   REMOVEKEY,
						"db_privileges.2.privilege": REMOVEKEY,
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

%s

%s

%s

resource "alibabacloudstack_polardbx_database" "default" {
	count         = 3
    instance_id   = "${local.polardbx_instance.id}"
	database_name = "${var.name}${count.index}"
	encode        = "utf8mb4"
	mode          = "auto"
}

resource "alibabacloudstack_polardbx_account" "default" {
    instance_id  = "${local.polardbx_instance.id}"
	account_name = "${var.name}"
	password     = "${random_password.password.0.result}"
	description  = "Normal user"
}

 `, name, RandomPasswordTestCase(12,1), VSwitchCommonTestCase, PolardbxReadOrCreateCommonTestCase())
}
