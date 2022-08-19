package apsarastack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"

	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackAdbConnectionConfig(t *testing.T) {
	var v *adb.Address
	rand := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccAdbConnection%s", rand)
	var basicMap = map[string]string{
		// 已有 实例 使用给定 id 测试
		//"db_cluster_id":    "am-3rq9uva152cn34drs",
		"db_cluster_id":     CHECKSET,
		"connection_string": CHECKSET,
		"ip_address":        CHECKSET,
		"port":              CHECKSET,
	}
	resourceId := "apsarastack_adb_connection.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAdbConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAdbConnectionConfigDependence)
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
					//"db_cluster_id":    "am-3rq9uva152cn34drs",
					"db_cluster_id":     "${apsarastack_adb_db_cluster.cluster.id}",
					"connection_prefix": fmt.Sprintf("tf-testacc%s", rand),
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
		},
	})
}

func resourceAdbConnectionConfigDependence(name string) string {
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
/*func resourceAdbConnectionConfigDependence(name string) string {
	return fmt.Sprintf(`

	variable "creation" {
		default = "ADB"
	}

	variable "name" {
		default = "%s"
	}

	`, name)
}*/
