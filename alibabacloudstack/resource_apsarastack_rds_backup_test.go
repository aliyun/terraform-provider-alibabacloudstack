package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackRdsBackup0(t *testing.T) {
	var v *RdsDescribebackupsResponse

	resourceId := "alibabacloudstack_rds_backup.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRdsBackupCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoDescribebackupsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tfrds_backup%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRdsBackupBasicdependence)
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

					"instance_id": "${alibabacloudstack_db_instance.default.id}",
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
				ImportStateVerifyIgnore: []string{"instance_id"},
			},
		},
	})
}

func TestAccAlibabacloudStackRdsBackup1(t *testing.T) {
	var v *RdsDescribebackupsResponse

	resourceId := "alibabacloudstack_rds_backup.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRdsBackupCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoDescribebackupsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tfrds_backup%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRdsBackupBasicdependence)
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

					"instance_id": "${alibabacloudstack_db_instance.default.id}",
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
				ImportStateVerifyIgnore: []string{"instance_id"},
			},
		},
	})
}

var AlibabacloudTestAccRdsBackupCheckmap = map[string]string{

	"backup_method": CHECKSET,
	// "account_type":   CHECKSET,
	// "character_type": CHECKSET,
	// "instance_id":    CHECKSET,
}

func AlibabacloudTestAccRdsBackupBasicdependence(name string) string {
	return fmt.Sprintf(
		`

	variable "name" {
		default = "%s"
	}

%s

resource "alibabacloudstack_db_instance" "default" {
	engine               = "MySQL"
	engine_version       = "5.6"
	instance_type        = "rds.mysql.s2.large"
	instance_storage     = "20"
	instance_name        = "${var.name}"
	vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
	storage_type         = "local_ssd"
  }

`, name, VSwitchCommonTestCase)
}
