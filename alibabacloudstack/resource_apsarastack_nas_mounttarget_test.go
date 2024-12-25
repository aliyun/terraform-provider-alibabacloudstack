package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackNasMounttarget0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_nas_mounttarget.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccNasMounttargetCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoNasDescribemounttargetsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snasmount_target%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccNasMounttargetBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"vswitch_id":        "${alibabacloudstack_vswitch.default.id}",
					"access_group_name": "${alibabacloudstack_nas_access_group.default.access_group_name}",
					"file_system_id":    "${alibabacloudstack_nas_file_system.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"vswitch_id": CHECKSET,

						"access_group_name": CHECKSET,

						"file_system_id": CHECKSET,
					}),
				),
			},

			// {
			// 	Config: testAccConfig(map[string]interface{}{

			// 		"status": "Inactive",
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{

			// 			"status": "Inactive",
			// 		}),
			// 	),
			// },
		},
	})
}

var AlibabacloudTestAccNasMounttargetCheckmap = map[string]string{

	"status": CHECKSET,

	"access_group_name": CHECKSET,

	"vswitch_id": CHECKSET,

	"file_system_id": CHECKSET,
}

func AlibabacloudTestAccNasMounttargetBasicdependence(name string) string {
	rand := getAccTestRandInt(10000, 99999)
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alibabacloudstack_vpc" "default" {
	cidr_block = "172.16.0.0/16"
	name = "${var.name}"
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
variable "storage_type" {
default = "Capacity"
}
data "alibabacloudstack_nas_protocols" "default" {
		type = "${var.storage_type}"
}
resource "alibabacloudstack_nas_file_system" "default" {
description = "${var.name}"
storage_type = "${var.storage_type}"
protocol_type = "${data.alibabacloudstack_nas_protocols.default.protocols.0}"
}
resource "alibabacloudstack_nas_access_group" "default" {
			access_group_name = "tf-testAccNasConfig-resource-test%d"
			access_group_type = "Vpc"
			description = "tf-testAccNasConfig"
}

`, name, rand)
}
