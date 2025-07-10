package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackMongodbAccount0(t *testing.T) {
	var v *DdsDescribeaccountsResponse

	resourceId := "alibabacloudstack_mongodb_account.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccMongodbAccountCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoDdsDescribeaccountsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	passward := getAccTestPassword(12)
	modify_passward := getAccTestPassword(12)

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tfaccount%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccMongodbAccountBasicdependence)
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

					"account_name": name,

					"account_password": passward,

					"instance_id": "${alibabacloudstack_mongodb_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"account_name": name,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"account_password": modify_passward,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"account_name": name,
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
			},
		},
	})
}

var AlibabacloudTestAccMongodbAccountCheckmap = map[string]string{

	"status":         CHECKSET,
	"account_type":   CHECKSET,
	"character_type": CHECKSET,
	"instance_id":    CHECKSET,
}

func AlibabacloudTestAccMongodbAccountBasicdependence(name string) string {
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
