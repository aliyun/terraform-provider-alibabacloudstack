package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackMongodbBackup0(t *testing.T) {
	var v *DdsDescribebackupsResponse

	resourceId := "alibabacloudstack_mongodb_backup.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccMongodbBackupCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoDdsDescribebackupsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tfmongodb_backup%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccMongodbBackupBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		// CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"backup_method": "Physical",

					"db_instance_id": "${alibabacloudstack_mongodb_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"backup_method": "Physical",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"db_instance_id"},
			},
		},
	})
}

func TestAccAlibabacloudStackMongodbBackup1(t *testing.T) {
	var v *DdsDescribebackupsResponse

	resourceId := "alibabacloudstack_mongodb_backup.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccMongodbBackupCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoDdsDescribebackupsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tfmongodb_backup%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccMongodbBackupBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		// CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"backup_method": "Logical",

					"db_instance_id": "${alibabacloudstack_mongodb_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"backup_method": "Logical",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"db_instance_id"},
			},
		},
	})
}

var AlibabacloudTestAccMongodbBackupCheckmap = map[string]string{

	"backup_method": CHECKSET,
	// "account_type":   CHECKSET,
	// "character_type": CHECKSET,
	// "instance_id":    CHECKSET,
}

func AlibabacloudTestAccMongodbBackupBasicdependence(name string) string {
	return fmt.Sprintf(
		`

	variable "name" {
		default = "%s"
	}

%s

	resource "alibabacloudstack_mongodb_instance" "default" {
		vswitch_id          = alibabacloudstack_vpc_vswitch.default.id
		engine_version      = "3.0"
		db_instance_class   = "dds.mongo.mid"
		db_instance_storage = "10"
		name                = "${var.name}"
		storage_engine      = "WiredTiger"
		instance_charge_type = "PostPaid"
		replication_factor = "3"
	  }

`, name, VSwitchCommonTestCase)
}
