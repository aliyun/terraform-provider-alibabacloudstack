package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackLindormLtsInstance_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_lindorm_lts_instance.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &LindormService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeLindormInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-lindorm-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, LindormLtsInstanceCommonTestCase)

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
					"instance_type":       "${data.alibabacloudstack_lindorm_instance_types.sortbycpu.instance_types[0].name}",
					"vpc_id":              "${alibabacloudstack_vpc_vpc.default.id}",
					"vswitch_id":          "${alibabacloudstack_vpc_vswitch.default.id}",
					"lts_num":             2,
					"deletion_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_alias":      name,
						"cpu_brand":           "Intel",
						"service_type":        "lts",
						"lts_num":             "2",
						"instance_status":     "ACTIVATION",
						"deletion_protection": "true",
						"topology_type":       "1azone",
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
					"lts_num":       "3",
					"instance_type": "${data.alibabacloudstack_lindorm_instance_types.sortbycpu.instance_types[1].name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lts_num":       "3",
						"instance_type": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"arch_version", "vpc_id", "vswitch_id"},
			},
		},
	})
}

func LindormLtsInstanceCommonTestCase(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

data "alibabacloudstack_lindorm_instance_types" "sortbycpu" {
	sorted_by = "CPU"
	engine_type = "lts"
}

`, name, VSwitchCommonTestCase)
}

func TestAccAlibabacloudStackLindormLtsInstance_multiZone(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_lindorm_lts_instance.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &LindormService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeLindormInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-lindorm-lts-mz-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, LindormLtsInstanceCommonTestCaseMultiZone)

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
					"zone_id":             "${data.alibabacloudstack_zones.default.zones.0.id}",
					"instance_alias":      "${var.name}",
					"cpu_brand":           "Intel",
					"instance_type":       "${data.alibabacloudstack_lindorm_instance_types.sortbycpu.instance_types[0].name}",
					"vpc_id":              "${alibabacloudstack_vpc_vpc.default.id}",
					"vswitch_id":          "${alibabacloudstack_vpc_vswitch.default.id}",
					"standby_zone_id":     "${data.alibabacloudstack_zones.default.zones.1.id}",
					"arbiter_zone_id":     "${data.alibabacloudstack_zones.default.zones.2.id}",
					"standby_vswitch_id":  "${alibabacloudstack_vpc_vswitch.standby.id}",
					"arbiter_vswitch_id":  "${alibabacloudstack_vpc_vswitch.arbiter.id}",
					"lts_num":             2,
					"deletion_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_alias":      name,
						"cpu_brand":           "Intel",
						"service_type":        "lts",
						"topology_type":       "multizone",
						"lts_num":             "2",
						"instance_status":     "ACTIVATION",
						"deletion_protection": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"arch_version", "vpc_id", "vswitch_id", "standby_zone_id", "standby_vswitch_id", "arbiter_zone_id", "arbiter_vswitch_id"},
			},
		},
	})
}

func LindormLtsInstanceCommonTestCaseMultiZone(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

data "alibabacloudstack_lindorm_instance_types" "sortbycpu" {
	sorted_by = "CPU"
	engine_type = "lts"
}

resource "alibabacloudstack_vpc_vswitch" "standby" {
  vswitch_name = "${var.name}_vsw_standby"
  vpc_id       = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block   = "172.16.2.0/24"
  zone_id      = "${data.alibabacloudstack_zones.default.zones.1.id}"
  lifecycle {
      ignore_changes = [
        tags
      ]
  }
}

resource "alibabacloudstack_vpc_vswitch" "arbiter" {
  vswitch_name = "${var.name}_vsw_arbiter"
  vpc_id       = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block   = "172.16.3.0/24"
  zone_id      = "${data.alibabacloudstack_zones.default.zones.2.id}"
  lifecycle {
      ignore_changes = [
        tags
      ]
  }
}

`, name, VSwitchCommonTestCase)
}
