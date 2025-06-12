package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEbsdiskreplicapair_basic0(t *testing.T) {
	var v *EbsDiskReplicaPair
	resourceId := "alibabacloudstack_ebs_diskreplicapair.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackEbsDiskreplicapairCheckMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EbsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "Describediskreplicapairs")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testaccebs-diskreplicapair%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackEbsDiskreplicapairDependence0)
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
					"disk_replica_pair_name": "${var.name}",
					"description":            "ebs_diskreplicapair test",
					"source_zone_id":         "${data.alibabacloudstack_zones.default.zones[0].id}",
					"source_region_id":       "${var.region}",
					"source_disk_id":         "${alibabacloudstack_ecs_disk.disk1.id}",
					"destination_zone_id":    "${data.alibabacloudstack_zones.default.zones[1].id}",
					"destination_region_id":  "${var.region}",
					"destination_disk_id":    "${alibabacloudstack_ecs_disk.disk2.id}",
					"rpo":                    "300",
					// "bandwidth": 				1000,  属性不支持修改
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_replica_pair_name": name,
						"description":            "ebs_diskreplicapair test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_replica_pair_name": "${var.name}_test",
					"description":            "ebs_diskreplicapair test2",
					"rpo":                    "600",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_replica_pair_name": fmt.Sprintf("%s_test", name),
						"description":            "ebs_diskreplicapair test2",
						"rpo":                    "600",
					}),
				),
			},

			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlibabacloudStackEbsDiskreplicapairCheckMap = map[string]string{}

func AlibabacloudStackEbsDiskreplicapairDependence0(name string) string {
	region := os.Getenv("ALIBABACLOUDSTACK_REGION")
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

variable "region" {
  default = "%s"
}
%s

resource "alibabacloudstack_ecs_disk" "disk1" {
	availability_zone = "${data.alibabacloudstack_zones.default.zones[0].id}"
	size = "20"
	name = "${var.name}"
	category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
}

resource "alibabacloudstack_ecs_disk" "disk2" {
	availability_zone = "${data.alibabacloudstack_zones.default.zones[1].id}"
	size = "20"
	name = "${var.name}"
	category = "${data.alibabacloudstack_zones.default.zones.1.available_disk_categories.0}"
}


`, DataAlibabacloudstackVswitchZones, region, name)
}
