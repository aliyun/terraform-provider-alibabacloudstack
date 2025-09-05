package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEcsSnapshotGroup_basic0(t *testing.T) {
	var v *EcsSnapshotGroup
	resourceId := "alibabacloudstack_ecs_snapshot_group.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackEcsSnapshotGroupCheckMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribesnapshotgroupsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testaccebs-snapshot-group%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackEcsSnapshotGroupDependence0)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":                   "${var.name}",
					"instance_id":                   "${alibabacloudstack_ecs_instance.default.id}",
					"instant_access":                "true",
					"instant_access_retention_days": "7",
					"snapshot_group_name":           "${var.name}",
					"disk_ids": []string{
						"${data.alibabacloudstack_ecs_disks.disks.disks.0.id}",
						"${data.alibabacloudstack_ecs_disks.disks.disks.1.id}",
						"${data.alibabacloudstack_ecs_disks.disks.disks.2.id}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snapshot_group_name": name,
						"description":         name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"snapshot_group_name": "${var.name}_test",
					"description":         "${var.name}_test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snapshot_group_name": fmt.Sprintf("%s_test", name),
						"description":         fmt.Sprintf("%s_test", name),
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// exclude_disk_ids has no read back
				ImportStateVerifyIgnore: []string{"exclude_disk_ids"},
			},
		},
	})
}

var AlibabacloudStackEcsSnapshotGroupCheckMap = map[string]string{
	"create_time":       CHECKSET,
	"status":            CHECKSET,
	"snapshot_group_id": CHECKSET,
}

func AlibabacloudStackEcsSnapshotGroupDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

%s

%s

%s

resource "alibabacloudstack_ecs_instance" "default" {
  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type        = "${local.default_instance_type_id}"
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  system_disk_size     = 20
  system_disk_name     = "test_sys_disk"
  security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
  instance_name        = "${var.name}_ecs"
  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
  zone_id    		   = data.alibabacloudstack_zones.default.zones.0.id
  is_outdated          = false
  data_disks {
      name                 = "disk1"
      category             = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
      size                 = 20
      delete_with_instance = true
    }
  data_disks {
      name                 = "disk2"
      category             = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
      size                 = 20
	  delete_with_instance = true
    }
  lifecycle {
    ignore_changes = [
      instance_type
    ]
  }
}

data "alibabacloudstack_ecs_disks" "disks" {
	instance_id = "${alibabacloudstack_ecs_instance.default.id}"
}

`, name, SecurityGroupCommonTestCase, DataAlibabacloudstackImages, DataAlibabacloudstackInstanceTypes)
}
