package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCspprivateHsmInstance_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_cspprivate_hsm_instance.default"

	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CspprivateService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeCspprivateHsmInstance")
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf_hsm_instance%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCspprivateHsmInstanceDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreYunCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testYunDunProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"product_code":   "${data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.products.0.code}",
					"vendor_code":    "${data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.code}",
					"vsm_type":       "gvsm",
					"zone_id":        "${data.alibabacloudstack_zones.default.zones.0.id}",
					"alias_name":     "${var.name}",
					"vpc_id":         "${alibabacloudstack_vpc.vpc0.id}",
					"vpc_cidr_block": "${alibabacloudstack_vpc.vpc0.cidr_block}",
					"vswitch_id":     "${alibabacloudstack_vswitch.vsw0.id}",
					"ip":             "192.168.0.100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product_code": CHECKSET,
						"vendor_code":  CHECKSET,
						"vsm_type":     "gvsm",
						"zone_id":      CHECKSET,
						"alias_name":   name,
						"ip":           "192.168.0.100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alias_name":     "${var.name}_updated",
					"vpc_id":         "${alibabacloudstack_vpc.vpc1.id}",
					"vpc_cidr_block": "${alibabacloudstack_vpc.vpc1.cidr_block}",
					"vswitch_id":     "${alibabacloudstack_vswitch.vsw1.id}",
					"ip":             "192.168.0.101",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alias_name": name + "_updated",
						"ip":         "192.168.0.101",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_code", "vpc_cidr_block"},
			},
		},
	})
}

func resourceCspprivateHsmInstanceDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

data "alibabacloudstack_zones" "default" {
	provider = alibabacloudstack-common
	available_resource_creation = "VSwitch"
}

data "alibabacloudstack_cspprivate_hsm_vendors" "default" {
}

resource "alibabacloudstack_vpc" "vpc0" {	
	provider = alibabacloudstack-common
	vpc_name = "${var.name}0"
	cidr_block = "192.168.0.0/16" # VPC CIDR block
}

resource "alibabacloudstack_vswitch" "vsw0" {
	provider = alibabacloudstack-common
	vpc_id = alibabacloudstack_vpc.vpc0.id
	cidr_block = "192.168.0.0/24" # VSwitch CIDR block
	availability_zone = data.alibabacloudstack_zones.default.zones.0.id
}

resource "alibabacloudstack_vpc" "vpc1" {	
	provider = alibabacloudstack-common
	vpc_name = "${var.name}1"
	cidr_block = "192.168.0.0/16"
}


resource "alibabacloudstack_vswitch" "vsw1" {
	provider = alibabacloudstack-common
	vpc_id = alibabacloudstack_vpc.vpc1.id
	cidr_block = "192.168.0.0/24"
	availability_zone = data.alibabacloudstack_zones.default.zones.0.id
}

`, name)
}

func TestAccAlibabacloudStackCspprivateHsmInstance_IsAppointDevice(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_cspprivate_hsm_instance.default"

	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CspprivateService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeCspprivateHsmInstance")
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf_hsm_instance%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCspprivateHsmInstanceIsAppointDeviceDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreYunCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testYunDunProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"product_code": "${data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.products.0.code}",
					"vendor_code":  "${data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.code}",
					"vsm_type":     "gvsm",
					"zone_id":      "${data.alibabacloudstack_zones.default.zones.0.id}",
					"alias_name":   "${var.name}",
					"device_id":    "${data.alibabacloudstack_cspprivate_hsms.default.hsms.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product_code": CHECKSET,
						"vendor_code":  CHECKSET,
						"vsm_type":     "gvsm",
						"zone_id":      CHECKSET,
						"alias_name":   name,
						"device_id":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alias_name": "${var.name}_updated",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alias_name": name + "_updated",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_code", "vpc_cidr_block"},
			},
		},
	})
}

func resourceCspprivateHsmInstanceIsAppointDeviceDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

data "alibabacloudstack_zones" "default" {
	provider = alibabacloudstack-common
	available_resource_creation = "VSwitch"
}

data "alibabacloudstack_cspprivate_hsm_vendors" "default" {
}

data "alibabacloudstack_cspprivate_hsms" "default" {
  vendor_code = "${data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors[0].code}"
  product_code = "${data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors[0].products[0].code}"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_vpc" "vpc" {	
	provider = alibabacloudstack-common
	vpc_name = var.name
	cidr_block = "192.168.0.0/16" # VPC CIDR block
}

resource "alibabacloudstack_vswitch" "vsw" {
	provider = alibabacloudstack-common
	vpc_id = alibabacloudstack_vpc.vpc.id
	cidr_block = "192.168.0.0/24" # VSwitch CIDR block
	availability_zone = data.alibabacloudstack_zones.default.zones.0.id # Availability zone
}

`, name)
}
