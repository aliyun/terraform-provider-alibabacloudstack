package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEbsDiskreplicagroup_basic0(t *testing.T) {
	var v *EbsReplicaGroup
	resourceId := "alibabacloudstack_ebs_diskreplicagroup.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackEbsDiskreplicagroupCheckMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EbsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeEbsDiskreplicagroups")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testaccebs-diskreplicagroup%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackEbsDiskreplicagroupDependence0)
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
					"disk_replica_group_name": "${var.name}",
					"description":             "ebs_diskreplicagroup test",
					"source_zone_id":          "${data.alibabacloudstack_zones.default.zones[0].id}",
					"source_region_id":        "${var.region}",
					"destination_zone_id":     "${data.alibabacloudstack_zones.default.zones[1].id}",
					"destination_region_id":   "${var.region}",
					"site":                    "production",
					"rpo":                     300,
					// "bandwidth": 				1000,  属性不支持修改
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_replica_group_name": name,
						"description":             "ebs_diskreplicagroup test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_replica_group_name": "${var.name}_test",
					"description":             "ebs_diskreplicagroup test2",
					"rpo":                     600,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_replica_group_name": fmt.Sprintf("%s_test", name),
						"description":             "ebs_diskreplicagroup test2",
						"rpo":                     "600",
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

var AlibabacloudStackEbsDiskreplicagroupCheckMap = map[string]string{}

func AlibabacloudStackEbsDiskreplicagroupDependence0(name string) string {
	region := os.Getenv("ALIBABACLOUDSTACK_REGION")
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
variable "region" {
  default = "%s"
}
%s

`, DataAlibabacloudstackVswitchZones, region, name)
}
