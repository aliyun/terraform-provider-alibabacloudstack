package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackNasMountTarget_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_nas_mount_target.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackNasMountTarget0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeNasMountTarget")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlibabacloudStackNasMountTarget%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackNasMountTargetBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"access_group_name": "${alibabacloudstack_nas_access_group.example.access_group_name}",
					"file_system_id":    "${alibabacloudstack_nas_file_system.example.id}",
					//"vswitch_id":        "${data.alibabacloudstack_vpcs.example.vpcs.0.vswitch_ids.0}",
					"vswitch_id":        "${alibabacloudstack_vswitch.default.id}",
					"security_group_id": "${alibabacloudstack_security_group.example.id}",
					//"status":            "Active",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_group_name": name,
						"file_system_id":    CHECKSET,
						"vswitch_id":        CHECKSET,
						"security_group_id": CHECKSET,
						//"status":            "Active",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"security_group_id"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Inactive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Inactive",
					}),
				),
			},
		},
	})
}

var AlibabacloudStackNasMountTarget0 = map[string]string{}

func AlibabacloudStackNasMountTargetBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

variable "name1" {
	default = "%schange"
}

data "alibabacloudstack_nas_protocols" "example" {
	type = "Performance"
}

resource "alibabacloudstack_vpc" "default" {
  cidr_block = "172.16.0.0/16"
  name = "${var.name}"
}

data "alibabacloudstack_nas_zones" "default" {
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "alibabacloudstack_security_group" "example" {
	name = var.name
	//vpc_id = "${data.alibabacloudstack_vpcs.example.vpcs.0.id}"
	vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_nas_file_system" "example" {
	protocol_type = "${data.alibabacloudstack_nas_protocols.example.protocols.0}"
	storage_type = "Capacity"
	//zone_id = "${data.alibabacloudstack_nas_zones.default.zones.0.zone_id}"
}

resource "alibabacloudstack_nas_access_group" "example" {
	access_group_name = "${var.name}"
	access_group_type = "Vpc"
}
`, name, name)
}
