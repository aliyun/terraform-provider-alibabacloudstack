package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackLindormInstance_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_lindorm_instance.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &LindormService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeLindormInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-lindorm-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, LindormInstanceCommonTestCase)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":             "${data.alibabacloudstack_zones.default.zones.0.id}",
					"instance_alias":      "${var.name}",
					"cpu_brand":           "Intel",
					"disk_category":       "HDD",
					"local_disk_size":     200,
					"engine_type":         "lindorm",
					"instance_type":       "${data.alibabacloudstack_lindorm_instance_types.sortbycpu.instance_types[0].name}",
					"vpc_id":              "${alibabacloudstack_vpc_vpc.default.id}",
					"vswitch_id":          "${alibabacloudstack_vpc_vswitch.default.id}",
					"lindorm_num":         2,
					"deletion_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_alias":      name,
						"cpu_brand":           "Intel",
						"engine_type":         "lindorm",
						"lindorm_num":         "2",
						"instance_status":     "ACTIVATION",
						"deletion_protection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_alias":      "${var.name}_update",
					"deletion_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_alias":      fmt.Sprintf("%s_update", name),
						"deletion_protection": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lindorm_num":   "3",
					"instance_type": "${data.alibabacloudstack_lindorm_instance_types.sortbycpu.instance_types[1].name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lindorm_num":   "3",
						"instance_type": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"disk_category", "local_disk_num", "local_disk_size", "arch_version"},
			},
		},
	})
}

func TestAccAlibabacloudStackLindormInstance_multiZone(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_lindorm_instance.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &LindormService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeLindormInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-lindorm-mz-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, LindormInstanceCommonTestCaseMultiZone)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
			testAccPreCheckWithMultiAZ(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":            "${local.zone_id_1}",
					"instance_alias":     "${var.name}",
					"cpu_brand":          "Intel",
					"disk_category":      "LV",
					"local_disk_size":    450,
					"engine_type":        "lindorm",
					"instance_type":      "${data.alibabacloudstack_lindorm_instance_types.sortbycpu.instance_types[0].name}",
					"core_spec":          "${data.alibabacloudstack_lindorm_instance_types.store.instance_types[0].name}",
					"vpc_id":             "${alibabacloudstack_vpc_vpc.default.id}",
					"vswitch_id":         "${alibabacloudstack_vpc_vswitch.default.id}",
					"standby_zone_id":    "${local.zone_id_2}",
					"standby_vswitch_id": "${alibabacloudstack_vpc_vswitch.standby.id}",
					"lindorm_num":        3,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_alias":  name,
						"cpu_brand":       "Intel",
						"engine_type":     "lindorm",
						"topology_type":   "multizone",
						"core_spec":       CHECKSET,
						"lindorm_num":     "3",
						"instance_status": "ACTIVATION",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"disk_category", "local_disk_num", "local_disk_size", "arch_version", "standby_vswitch_id", "core_spec"},
			},
		},
	})
}

func LindormInstanceCommonTestCase(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alibabacloudstack_lindorm_instance_types" "sortbycpu" {
	sorted_by = "CPU"
	engine_type = "lindorm"
}

%s

`, name, VSwitchCommonTestCase)
}

func LindormInstanceCommonTestCaseMultiZone(name string) string {
	zone1 := os.Getenv("ALIBABACLOUDSTACK_LINDORM_MULTIZONE1")
	zone2 := os.Getenv("ALIBABACLOUDSTACK_LINDORM_MULTIZONE2")

	zoneExpr1 := "data.alibabacloudstack_zones.default.zones[0].id"
	zoneExpr2 := "data.alibabacloudstack_zones.default.zones[1].id"
	if zone1 != "" {
		zoneExpr1 = fmt.Sprintf("%q", zone1)
	}
	if zone2 != "" {
		zoneExpr2 = fmt.Sprintf("%q", zone2)
	}

	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alibabacloudstack_lindorm_instance_types" "sortbycpu" {
	sorted_by = "CPU"
	engine_type = "lindorm"
}

data "alibabacloudstack_lindorm_instance_types" "store" {
	engine_type = "store"
}

%s

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "10.0.0.0/16"
  lifecycle {
      ignore_changes = [
		secondary_cidr_blocks,
		tags,
      ]
  }
}

locals {
  zone_id_1 = %s
  zone_id_2 = %s
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  vswitch_name = "${var.name}_vsw"
  vpc_id       = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block   = "10.0.1.0/24"
  zone_id      = "${local.zone_id_1}"
  lifecycle {
      ignore_changes = [
        tags
      ]
  }
}

resource "alibabacloudstack_vpc_vswitch" "standby" {
  vswitch_name = "${var.name}_vsw_standby"
  vpc_id       = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block   = "10.0.2.0/24"
  zone_id      = "${local.zone_id_2}"
  lifecycle {
      ignore_changes = [
        tags
      ]
  }
}

`, name, DataZoneCommonTestCase, zoneExpr1, zoneExpr2)
}
