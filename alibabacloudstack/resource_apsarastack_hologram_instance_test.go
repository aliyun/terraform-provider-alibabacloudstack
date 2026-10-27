package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackHologramInstance_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_hologram_instance.default"

	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HologramService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeHologramInstance")
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf_hologram_instance%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceHologramInstanceDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t, resourceId)
		},

		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"compute_type":  "Standard",
					"zone_id":       "${data.alibabacloudstack_hologram_clusters.default.zone_id}",
					"cpu":           "intel",
					"node":          "2",
					"vpc_id":        "${data.alibabacloudstack_vpcs.default.vpcs.0.id}",
					"vswitch_id":    "${data.alibabacloudstack_vswitches.default.vswitches.0.id}",
					"instance_name": "${var.name}",
					"cluster":       "${data.alibabacloudstack_hologram_clusters.default.clusters.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compute_type":  "Standard",
						"cpu":           "intel",
						"node":          "2",
						"instance_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node":          "4",
					"instance_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node":          "4",
						"instance_name": fmt.Sprintf("%s_update", name),
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

func TestAccAlibabacloudStackHologramInstance_readOnly(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_hologram_instance.default"

	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HologramService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeHologramInstance")
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf_hologram_instance%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceHologramReadOnlyInstanceDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t, resourceId)
		},

		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"compute_type":       "Follower",
					"zone_id":            "${alibabacloudstack_hologram_instance.standard.0.zone_id}",
					"cpu":                "intel",
					"node":               "2",
					"vpc_id":             "${data.alibabacloudstack_vpcs.default.vpcs.0.id}",
					"vswitch_id":         "${data.alibabacloudstack_vswitches.default.vswitches.0.id}",
					"instance_name":      "${var.name}",
					"cluster":            "${data.alibabacloudstack_hologram_clusters.default.clusters.0.id}",
					"leader_instance_id": "${alibabacloudstack_hologram_instance.standard.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compute_type":  "Follower",
						"cpu":           "intel",
						"node":          "2",
						"instance_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"leader_instance_id": "${alibabacloudstack_hologram_instance.standard.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"leader_instance_id": CHECKSET,
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

func resourceHologramInstanceDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

data "alibabacloudstack_vpcs" "default" {}

data "alibabacloudstack_vswitches" "default" {
  vpc_id = data.alibabacloudstack_vpcs.default.vpcs.0.id
}

data "alibabacloudstack_hologram_clusters" "default" {
  zone_id = "cn-wulan-env2g1-amtest2001-a"
}

`, name)
}

func resourceHologramReadOnlyInstanceDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

data "alibabacloudstack_vpcs" "default" {}

data "alibabacloudstack_vswitches" "default" {
  vpc_id = data.alibabacloudstack_vpcs.default.vpcs.0.id
}

data "alibabacloudstack_hologram_clusters" "default" {
  zone_id = "cn-wulan-env2g1-amtest2001-a"
}

resource "alibabacloudstack_hologram_instance" "standard" {
	count = 2
  	compute_type = "Standard"
	zone_id = "${data.alibabacloudstack_hologram_clusters.default.zone_id}"
	cpu = "intel"
	node = "2"
	vpc_id = "${data.alibabacloudstack_vpcs.default.vpcs.0.id}"
	vswitch_id = "${data.alibabacloudstack_vswitches.default.vswitches.0.id}"
	instance_name = "${var.name}${count.index}"
	cluster = "${data.alibabacloudstack_hologram_clusters.default.clusters.0.id}"
}


`, name)
}
