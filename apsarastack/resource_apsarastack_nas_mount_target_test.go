package apsarastack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackNasMountTarget_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "apsarastack_nas_mount_target.default"
	ra := resourceAttrInit(resourceId, ApsaraStackNasMountTarget0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeNasMountTarget")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sApsaraStackNasMountTarget%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApsaraStackNasMountTargetBasicDependence0)
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
					"access_group_name": "${apsarastack_nas_access_group.example.access_group_name}",
					"file_system_id":    "${apsarastack_nas_file_system.example.id}",
					//"vswitch_id":        "${data.apsarastack_vpcs.example.vpcs.0.vswitch_ids.0}",
					"vswitch_id":        "${apsarastack_vswitch.default.id}",
					"security_group_id": "${apsarastack_security_group.example.id}",
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

var ApsaraStackNasMountTarget0 = map[string]string{}

func ApsaraStackNasMountTargetBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

variable "name1" {
	default = "%schange"
}

data "apsarastack_nas_protocols" "example" {
	type = "Performance"
}

resource "apsarastack_vpc" "default" {
  cidr_block = "172.16.0.0/16"
  name = "${var.name}"
}

data "apsarastack_nas_zones" "default" {
}

data "apsarastack_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "apsarastack_vswitch" "default" {
  vpc_id = "${apsarastack_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "apsarastack_security_group" "example" {
	name = var.name
	//vpc_id = "${data.apsarastack_vpcs.example.vpcs.0.id}"
	vpc_id = "${apsarastack_vpc.default.id}"
}

resource "apsarastack_nas_file_system" "example" {
	protocol_type = "${data.apsarastack_nas_protocols.example.protocols.0}"
	storage_type = "Capacity"
	//zone_id = "${data.apsarastack_nas_zones.default.zones.0.zone_id}"
}

resource "apsarastack_nas_access_group" "example" {
	access_group_name = "${var.name}"
	access_group_type = "Vpc"
}
`, name, name)
}
