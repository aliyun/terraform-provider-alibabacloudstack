package apsarastack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"

	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackAdbBackupPolicy(t *testing.T) {
	var v *adb.DescribeBackupPolicyResponse
	resourceId := "apsarastack_adb_backup_policy.default"
	serverFunc := func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeAdbBackupPolicy")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccAdbBackupPolicy"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAdbBackupPolicyConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//该功能没有 删除接口，只是恢复默认
		//CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					// 已有 实例 使用给定 id 测试
					//"db_cluster_id":    "am-3rq9uva152cn34drs",
					"db_cluster_id":           "${apsarastack_adb_db_cluster.default.id}",
					"preferred_backup_period": []string{"Tuesday", "Wednesday"},
					"preferred_backup_time":   "10:00Z-11:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#": "2",
						"preferred_backup_time":     "10:00Z-11:00Z",
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
					"preferred_backup_period": []string{"Wednesday", "Monday", "Saturday"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_time": "15:00Z-16:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_time": "15:00Z-16:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_period": []string{"Tuesday", "Thursday", "Friday", "Sunday"},
					"preferred_backup_time":   "17:00Z-18:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#": "4",
						"preferred_backup_time":     "17:00Z-18:00Z",
					}),
				),
			}},
	})
}

func resourceAdbBackupPolicyConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "ADB"
	}

	variable "name" {
		default = "%s"
	}

	resource "apsarastack_adb_db_cluster" "default" {
	db_cluster_category = "Basic"
	db_cluster_class = "C8"
	db_node_storage = "200"
	db_cluster_version = "3.0"
	db_node_count = "2"
	mode					= "reserver"
	vswitch_id              = "${apsarastack_vswitch.default.id}"
	description             = "${var.name}"
	cluster_type =        "analyticdb"
	cpu_type =            "intel"

	}`, AdbCommonTestCase, name)
}

// 已有 实例创建测试使用
/*func resourceAdbBackupPolicyConfigDependence(name string) string {
	return fmt.Sprintf(`

	variable "creation" {
		default = "ADB"
	}

	variable "name" {
		default = "%s"
	}

	`, name)
}*/
