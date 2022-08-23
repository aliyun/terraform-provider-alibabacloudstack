package apsarastack

import (
	"fmt"
	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackEcsEbsStorageSets_basic(t *testing.T) {
	var v *datahub.EcsDescribeEcsEbsStorageSetsResult
	resourceId := "apsarastack_ecs_ebs_storage_set.default"
	ra := resourceAttrInit(resourceId, ApsaraStackEcsEbsMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeEcsEbsStorageSet")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sApsaraStackEcsCommand%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApsaraStackEcsEbsBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: providerCommon + testAccConfig(map[string]interface{}{
					"storage_set_name":    name,
					"maxpartition_number": "2",
					"zone_id":             "${data.apsarastack_zones.default.zones.0.id}",
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

var ApsaraStackEcsEbsMap = map[string]string{
	//"enable_parameter": "false",
}

//func ApsaraStackEcsEbsBasicDependence(name string) string {
//	return ""
//}

func ApsaraStackEcsEbsBasicDependence(name string) string {
	return fmt.Sprintf(`
provider "apsarastack" {
	assume_role {}
}
variable "name" {
	default = "%s"
}
data "apsarastack_zones" "default" {}

//data "apsarastack_vpcs" "default" {
//	name_regex = "default-NODELETING"
//}
//resource "apsarastack_vpc" "default" {
//name       = var.name
//cidr_block = "172.16.0.0/16"
//}
//resource "apsarastack_vswitch" "default" {
//  vpc_id            = "${apsarastack_vpc.default.id}"
//  cidr_block        = "172.16.0.0/24"
//  availability_zone = data.apsarastack_hbase_zones.default.ids.0
//  name              = "${var.name}"
//}
//
//resource "apsarastack_security_group" "default" {
//	count = 2
//	vpc_id = apsarastack_vpc.default.id
//	name = var.name
//}
`, name)
}
