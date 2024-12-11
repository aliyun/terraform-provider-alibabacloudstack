package alibabacloudstack

import (
	"fmt"
	"testing"

	

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDBReadonlyInstance_update(t *testing.T) {
	var instance *rds.DBInstanceAttribute
	resourceId := "alibabacloudstack_db_readonly_instance.default"
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testAccDBInstance_vpc_%d", rand)
	var DBReadonlyMap = map[string]string{
		"instance_storage":      "5",
		"engine_version":        "5.6",
		"engine":                "MySQL",
		"port":                  "3306",
		"instance_name":         name,
		"instance_type":         CHECKSET,
		"parameters":            NOSET,
		"master_db_instance_id": CHECKSET,
		"zone_id":               CHECKSET,
		"vswitch_id":            CHECKSET,
		"connection_string":     CHECKSET,
	}
	ra := resourceAttrInit(resourceId, DBReadonlyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDBReadonlyInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBReadonlyInstanceConfigDependence)
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
					"master_db_instance_id":    "${alibabacloudstack_db_instance.default.id}",
					"zone_id":                  "${alibabacloudstack_db_instance.default.zone_id}",
					"engine_version":           "${alibabacloudstack_db_instance.default.engine_version}",
					"instance_type":            "${alibabacloudstack_db_instance.default.instance_type}",
					"instance_storage":         "${alibabacloudstack_db_instance.default.instance_storage}",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${alibabacloudstack_vswitch.default.id}",
					"db_instance_storage_type": "${alibabacloudstack_db_instance.default.storage_type}",
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
			// upgrade storage
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_storage": "${alibabacloudstack_db_instance.default.instance_storage + data.alibabacloudstack_db_instance_classes.default.instance_classes.0.storage_range.step}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"instance_storage": "10"}),
				),
			},
			// upgrade instanceType
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "${data.alibabacloudstack_db_instance_classes.default.instance_classes.1.instance_class}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"instance_type": CHECKSET}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${var.name}_ro",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_ro",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       REMOVEKEY,
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"master_db_instance_id": "${alibabacloudstack_db_instance.default.id}",
					"zone_id":               "${alibabacloudstack_db_instance.default.zone_id}",
					"engine_version":        "${alibabacloudstack_db_instance.default.engine_version}",
					"instance_type":         "${alibabacloudstack_db_instance.default.instance_type}",
					"instance_storage":      "${alibabacloudstack_db_instance.default.instance_storage + 2*data.alibabacloudstack_db_instance_classes.default.instance_classes.0.storage_range.step}",
					"instance_name":         "${var.name}",
					"vswitch_id":            "${alibabacloudstack_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":    name,
						"instance_storage": "15",
					}),
				),
			},
		},
	})

}

func TestAccAlibabacloudStackDBReadonlyInstance_multi(t *testing.T) {
	var instance *rds.DBInstanceAttribute
	resourceId := "alibabacloudstack_db_readonly_instance.default.1"
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testAccDBInstance_vpc_%d", rand)
	var DBReadonlyMap = map[string]string{
		"instance_storage":      "5",
		"engine_version":        "5.6",
		"engine":                "MySQL",
		"port":                  "3306",
		"instance_name":         name,
		"instance_type":         CHECKSET,
		"parameters":            NOSET,
		"master_db_instance_id": CHECKSET,
		"zone_id":               CHECKSET,
		"vswitch_id":            CHECKSET,
		"connection_string":     CHECKSET,
	}
	ra := resourceAttrInit(resourceId, DBReadonlyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDBReadonlyInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBReadonlyInstanceConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":                    "1",
					"master_db_instance_id":    "${alibabacloudstack_db_instance.default.id}",
					"zone_id":                  "${alibabacloudstack_db_instance.default.zone_id}",
					"engine_version":           "${alibabacloudstack_db_instance.default.engine_version}",
					"instance_type":            "${alibabacloudstack_db_instance.default.instance_type}",
					"instance_storage":         "${alibabacloudstack_db_instance.default.instance_storage}",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${alibabacloudstack_vswitch.default.id}",
					"db_instance_storage_type": "${alibabacloudstack_db_instance.default.storage_type}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceDBReadonlyInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
%s
	variable "creation" {
		default = "Rds"
	}
	variable "multi_az" {
		default = "false"
	}
	variable "name" {
		default = "%s"
	}
resource "alibabacloudstack_db_instance" "default" {
		engine = "MySQL"
		engine_version = "5.6"
		instance_type = "rds.mysql.s2.large"
		instance_storage = "30"
		instance_charge_type = "Postpaid"
		instance_name = "${var.name}"
		storage_type = "local_ssd"
		vswitch_id = "${alibabacloudstack_vswitch.default.id}"
		security_ips = ["10.168.1.12", "100.69.7.112"]
	}
	
`, RdsCommonTestCase, name)
}
