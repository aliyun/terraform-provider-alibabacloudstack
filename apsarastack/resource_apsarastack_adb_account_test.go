package apsarastack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"

	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackAdbAccount_update_forSuper(t *testing.T) {
	var v *adb.DBAccount
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccadbaccount-%d", rand)
	var basicMap = map[string]string{
		// 已有 实例 使用给定 id 测试
		//"db_cluster_id":    "am-3rqb9q5nk034py521",
		"db_cluster_id":    CHECKSET,
		"account_name":     "tftestsuper",
		"account_password": "inputYourCodeHere",
	}
	resourceId := "apsarastack_adb_account.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAdbAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAdbAccountConfigDependence)
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
					// 已有 实例 使用给定 id 测试
					//"db_cluster_id":    "am-3rqb9q5nk034py521",
					"db_cluster_id":    "${apsarastack_adb_db_cluster.cluster.id}",
					"account_name":     "tftestsuper",
					"account_password": "inputYourCodeHere",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
			},
			// 专有云 没有该接口 ModifyAccountDescription
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"account_description": "from terraform super",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"account_description": "from terraform super",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "inputYourCodeHere",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "inputYourCodeHere",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					//"account_description": "tf test super",
					"account_password": "inputYourCodeHere",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						//"account_description": "tf test super",
						"account_password": "inputYourCodeHere",
					}),
				),
			},
		},
	})

}

func resourceAdbAccountConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "ADB"
	}

	variable "name" {
		default = "%s"
	}

	resource "apsarastack_adb_db_cluster" "cluster" {
		db_cluster_category = "Basic"
		db_cluster_class = "C8"
		vswitch_id     = "${apsarastack_vswitch.default.id}"
		description             = "${var.name}"
		db_node_storage = "200"
	    db_cluster_version = "3.0"
	    db_node_count = "2"
		mode					= "reserver"
		cluster_type =        "analyticdb"
		cpu_type =            "intel"
	}`, AdbCommonTestCase, name)
}

// 已有 实例创建测试使用
/*func resourceAdbAccountConfigDependence(name string) string {
	return fmt.Sprintf(`

	variable "creation" {
		default = "ADB"
	}

	variable "name" {
		default = "%s"
	}

	`, name)
}*/
