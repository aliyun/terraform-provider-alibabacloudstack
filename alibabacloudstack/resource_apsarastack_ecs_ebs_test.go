package alibabacloudstack

import (
	"fmt"
	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEcsEbsStorageSets_basic(t *testing.T) {
	var v *datahub.EcsDescribeEcsEbsStorageSetsResult
	resourceId := "alibabacloudstack_ecs_ebs_storage_set.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackEcsEbsMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeEcsEbsStorageSet")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlibabacloudStackEcsCommand%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackEcsEbsBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: providerCommon + testAccConfig(map[string]interface{}{
					"storage_set_name":    name,
					"maxpartition_number": "2",
					"zone_id":             "${data.alibabacloudstack_zones.default.zones.0.id}",
					//"name":            name,
					//"type":            "RunShellScript",
					//"working_dir":     "/root",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						//"command_content": "bHMK",
						//"description":     "For Terraform Test",
						//"name":            name,
						//"type":            "RunShellScript",
						"storage_set_name": name,
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

var AlibabacloudStackEcsEbsMap = map[string]string{
	//"enable_parameter": "false",
}

//func AlibabacloudStackEcsEbsBasicDependence(name string) string {
//	return ""
//}

func AlibabacloudStackEcsEbsBasicDependence(name string) string {
	return fmt.Sprintf(`
provider "alibabacloudstack" {
	assume_role {}
}
variable "name" {
	default = "%s"
}
data "alibabacloudstack_zones" "default" {}

//data "alibabacloudstack_vpcs" "default" {
//	name_regex = "default-NODELETING"
//}
//resource "alibabacloudstack_vpc" "default" {
//name       = var.name
//cidr_block = "172.16.0.0/16"
//}
//resource "alibabacloudstack_vswitch" "default" {
//  vpc_id            = "${alibabacloudstack_vpc.default.id}"
//  cidr_block        = "172.16.0.0/24"
//  availability_zone = data.alibabacloudstack_hbase_zones.default.ids.0
//  name              = "${var.name}"
//}
//
//resource "alibabacloudstack_security_group" "default" {
//	count = 2
//	vpc_id = alibabacloudstack_vpc.default.id
//	name = var.name
//}
`, name)
}
