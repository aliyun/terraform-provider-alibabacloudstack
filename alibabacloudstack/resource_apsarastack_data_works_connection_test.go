package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDataWorksConnection_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_data_works_connection.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackDataWorksConnectionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataworksService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDataWorksConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf_testconnection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackDataWorksConnectionBasicDependence0)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project_id":      "${alibabacloudstack_data_works_project.default.id}",
					"connection_type": "rds",
					"content": map[string]string{
						"password":     "${random_password.password.0.result}",
						"instanceName": "${alibabacloudstack_db_instance.default.id}",
						"username":     "${alibabacloudstack_db_account.default.name}",
						"database":     "${alibabacloudstack_db_database.default.0.name}",
						"tag":          "rds",
					},
					"env_type":    "1",
					"sub_type":    "mysql",
					"name":        "${var.name}",
					"description": "${var.name} description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_type": "rds",
						"env_type":        "1",
						"sub_type":        "mysql",
						"name":            name,
						"description":     name + " description",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"content.password"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name} description update",
					"content": map[string]string{
						"password":     "${random_password.password.0.result}",
						"instanceName": "${alibabacloudstack_db_instance.default.id}",
						"username":     "${alibabacloudstack_db_account.default.name}",
						"database":     "${alibabacloudstack_db_database.default.1.name}",
						"tag":          "rds",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + " description update",
					}),
				),
			},
		},
	})
}

var AlibabacloudStackDataWorksConnectionMap0 = map[string]string{}

func AlibabacloudStackDataWorksConnectionBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

%s

%s


%s

resource "alibabacloudstack_db_database" "default" {
  count  = 2
  instance_id = "${alibabacloudstack_db_instance.default.id}"
  name = "tfaccount${count.index}"
  description = "from terraform"
  character_set        = "utf8"
}

resource "alibabacloudstack_db_account" "default" {
  instance_id = "${alibabacloudstack_db_instance.default.id}"
  name = "tftest"
  password = random_password.password.0.result
  description = "from terraform"
}

resource "alibabacloudstack_db_account_privilege" "default" {
	instance_id  = "${alibabacloudstack_db_instance.default.id}"
	account_name = "${alibabacloudstack_db_account.default.name}"
	privilege    = "ReadOnly"
	db_names     = alibabacloudstack_db_database.default.*.name
}


resource "alibabacloudstack_data_works_project" "default" {
	name =           "${var.name}"
	description =    "${var.name}_desc"
	task_auth_type = "PROJECT"
}

`, name, RandomPasswordTestCase(12, 1), VSwitchCommonTestCase, RdsMysqlCommonTestCase())
}
