package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_drds_rds", &resource.Sweeper{
		Name: "alibabacloudstack_drds_rds",
		F:    testSweepDrdsInstances,
	})
}

func TestAccAlibabacloudStackDrdsRdsInstance_basic0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_drds_rds_instance.default"
	ra := resourceAttrInit(resourceId, drdsRdsbasicMap)

	serviceFunc := func() interface{} {
		return &DrdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf_acc_drds_rds_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDrdsRdsDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_type":        "local_ssd",
					"category":            "HighAvailability",
					"db_instance_class":   "rds.mysql.s1.small",
					"drds_instance_id":    "${local.drds_instance_id}",
					"zone_id":             "${data.alibabacloudstack_zones.default.zones.0.id}",
					"db_instance_storage": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage": "20",
					}),
				),
			},
			{
				ResourceName: resourceId,
				ImportState:  true,
				// "zone_id", "db_instance_class", "storage_type" do not support read back
				// "force_remove" is a local control attribute
				ImportStateVerifyIgnore: []string{"zone_id", "storage_type", "db_instance_class", "force_remove"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_storage": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage": "30",
					}),
				),
			},
		},
	})
}

func resourceDrdsRdsDependence(name string) string {
	return fmt.Sprintf(`

	variable "name" {
		default = "%s"
	}
	
	variable "existed_drds_instance" {
		default = "%s"
	}
	
	locals {
		create_drds_instance_count = var.existed_drds_instance == "" ? 1: 0
	}
	
	variable "instance_series" {
		default = "drds.sn2.4c16g"
	}
	
	%s
	
	resource "alibabacloudstack_drds_instance" "default" {
		count                = local.create_drds_instance_count
		description          = "${var.name}"
		zone_id              = "${alibabacloudstack_vpc_vswitch.default.availability_zone}"
		instance_series      = "${var.instance_series}"
		instance_charge_type = "PostPaid"
		vswitch_id           = "${alibabacloudstack_vpc_vswitch.default.id}"
		specification        = "drds.sn2.4c16g.8C32G"
	}
	
	locals {
		drds_instance_id = var.existed_drds_instance == "" ? alibabacloudstack_drds_instance.default.0.id: var.existed_drds_instance
	}
	
`, name, os.Getenv("ALIBABACLOUDSTACK_TEST_EXISTED_DRDS_ID"), VSwitchCommonTestCase)
}

var drdsRdsbasicMap = map[string]string{
	"rds_instance_id": CHECKSET,
	"create_time":     CHECKSET,
	"status":          CHECKSET,
}
