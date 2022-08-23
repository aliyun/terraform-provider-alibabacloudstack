package apsarastack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackDBDatabaseUpdate(t *testing.T) {
	var database *rds.Database
	resourceId := "apsarastack_db_database.default"

	var dbDatabaseBasicMap = map[string]string{
		"instance_id":   CHECKSET,
		"name":          "tftestdatabase",
		"character_set": "utf8",
		"description":   "",
	}

	ra := resourceAttrInit(resourceId, dbDatabaseBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &database, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeDBDatabase")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBdatabase_basic"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBDatabaseConfigDependence)
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
					"instance_id":   "${apsarastack_db_instance.instance.id}",
					"name":          "tftestdatabase",
					"character_set": "utf8",
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
					"description": "from terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"description": "from terraform"}),
				),
			},
		},
	})

}

func resourceDBDatabaseConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
provider "apsarastack" {
	assume_role {}
}
	variable "creation" {
		default = "Rds"
	}

	variable "name" {
		default = "%s"
	}



	resource "apsarastack_db_instance" "instance" {
	     engine               = "MySQL"
        engine_version       = "5.6"
        instance_type        = "rds.mysql.s2.large"
	    instance_storage     = "30"
		vswitch_id = "${apsarastack_vswitch.default.id}"
		instance_name = "${var.name}"
		storage_type         = "local_ssd"
	}`, RdsCommonTestCase, name)
}
