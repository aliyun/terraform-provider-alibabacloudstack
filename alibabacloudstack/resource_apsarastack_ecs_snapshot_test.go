package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEcsSnapshot0(t *testing.T) {
	var v *ecs.Snapshot

	resourceId := "alibabacloudstack_ecs_snapshot.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsSnapshotCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribesnapshotsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secssnapshot%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsSnapshotBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "rdk_test_description",

					"snapshot_name": "rdk_test_name",

					"disk_id": "${alibabacloudstack_instance.instance.system_disk_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "rdk_test_description",

						"snapshot_name": "rdk_test_name",

						"disk_id": CHECKSET,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "rdk_test_description-update",

					"snapshot_name": "rdk_test_name-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "rdk_test_description-update",

						"snapshot_name": "rdk_test_name-update",
					}),
				),
			},

			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"tags": map[string]string{
			// 			"Created": "TF",
			// 			"For":     "Test",
			// 		},
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"tags.%":       "2",
			// 			"tags.Created": "TF",
			// 			"tags.For":     "Test",
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"tags": map[string]string{
			// 			"Created": "TF-update",
			// 			"For":     "Test-update",
			// 		},
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"tags.%":       "2",
			// 			"tags.Created": "TF-update",
			// 			"tags.For":     "Test-update",
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"tags": REMOVEKEY,
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"tags.%":       "0",
			// 			"tags.Created": REMOVEKEY,
			// 			"tags.For":     REMOVEKEY,
			// 		}),
			// 	),
			// },
		},
	})
}

var AlibabacloudTestAccEcsSnapshotCheckmap = map[string]string{

	// "instant_access": CHECKSET,

	// "description": CHECKSET,

	// "resource_group_id": CHECKSET,

	// "encrypted": CHECKSET,

	// "instant_access_retention_days": CHECKSET,

	// "snapshot_name": CHECKSET,

	// "snapshot_sn": CHECKSET,

	// "tags": CHECKSET,

	// "status": CHECKSET,

	// "progress": CHECKSET,

	// "usage": CHECKSET,

	// "product_code": CHECKSET,

	// "create_time": CHECKSET,

	// "retention_days": CHECKSET,

	// "source_storage_type": CHECKSET,

	// "snapshot_id": CHECKSET,

	// "source_disk_size": CHECKSET,

	// "snapshot_type": CHECKSET,

	// "remain_time": CHECKSET,

	// "source_disk_type": CHECKSET,

	// "disk_id": CHECKSET,
}

func AlibabacloudTestAccEcsSnapshotBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation = "Instance"
  }

data "alibabacloudstack_images" "default" {
name_regex  = "^ubuntu_18.*64"
most_recent = true
owners      = "system"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  eni_amount        = 2
  sorted_by         = "Memory"
}

resource "alibabacloudstack_vpc" "default" {
    name = "${var.name}"
    cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
    vpc_id = "${alibabacloudstack_vpc.default.id}"
    cidr_block = "172.16.0.0/16"
    availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
    name = "${var.name}"
}

resource "alibabacloudstack_security_group" "new" {
	name = "${var.name}"
  	vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_instance" "instance" {
	image_id = "${data.alibabacloudstack_images.default.ids[0]}"
	instance_type = "${data.alibabacloudstack_instance_types.default.ids[0]}"
	instance_name = "${var.name}"
	security_groups = "${alibabacloudstack_security_group.new.*.id}"
	internet_max_bandwidth_out = "10"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	system_disk_category = "cloud_efficiency"
	vswitch_id = "${alibabacloudstack_vswitch.default.id}"
	system_disk_size = 20
	system_disk_name = "system_disk"
}

`, name)
}
