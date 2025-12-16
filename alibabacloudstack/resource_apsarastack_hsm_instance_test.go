package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackHsmInstance_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_hsm_instance.default"

	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HsmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeHsmInstance")
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf_hsm_instance%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceHsmInstanceDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreYunCheck(t)
		},
		IDRefreshName: resourceId,
		Providers: testYunDunProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"product_code": "${data.alibabacloudstack_hsm_vendors.default.vendors.0.products.0.code}",
					"vendor_code":  "${data.alibabacloudstack_hsm_vendors.default.vendors.0.code}",
					"vsm_type":     "gvsm",
					"zone_no":      "${data.alibabacloudstack_zones.default.zones.0.id}",
					"remark":       "test-tf-hsm-instance",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product_code": CHECKSET,
						"vendor_code":  CHECKSET,
						"vsm_type":     "gvsm",
						"zone_no":      CHECKSET,
						"remark":       "test-tf-hsm-instance",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "updated_remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "updated_remark",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":     "${alibabacloudstack_vpc.vpc.id}",
					"vswitch_id": "${alibabacloudstack_vswitch.vsw.id}",
					"ip":         "192.168.0.100",
					"white_list": "192.168.0.123/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":     CHECKSET,
						"vswitch_id": CHECKSET,
						"ip":         "192.168.0.100",
						"white_list": "192.168.0.123/24",
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

func resourceHsmInstanceDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

data "alibabacloudstack_zones" "default" {
	provider = alibabacloudstack-common
	available_resource_creation = "VSwitch"
}

data "alibabacloudstack_hsm_vendors" "default" {
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

func TestAccAlibabacloudStackHsmInstance_IsAppointDevice(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_hsm_instance.default"

	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HsmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeHsmInstance")
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf_hsm_instance%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceHsmInstanceIsAppointDeviceDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreYunCheck(t)
		},
		IDRefreshName: resourceId,
		Providers: testYunDunProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"product_code": "${data.alibabacloudstack_hsm_vendors.default.vendors.0.products.0.code}",
					"vendor_code":  "${data.alibabacloudstack_hsm_vendors.default.vendors.0.code}",
					"vsm_type":     "gvsm",
					"zone_no":      "${data.alibabacloudstack_zones.default.zones.0.id}",
					"remark":       "test-tf-hsm-instance",
					"hsm_id":       "${data.alibabacloudstack_hsms.default.hsms.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product_code": CHECKSET,
						"vendor_code":  CHECKSET,
						"vsm_type":     "gvsm",
						"zone_no":      CHECKSET,
						"remark":       "test-tf-hsm-instance",
						"hsm_id":       CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "updated_remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "updated_remark",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":     "${alibabacloudstack_vpc.vpc.id}",
					"vswitch_id": "${alibabacloudstack_vswitch.vsw.id}",
					"ip":         "192.168.0.100",
					"white_list": "192.168.0.123/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":     CHECKSET,
						"vswitch_id": CHECKSET,
						"ip":         "192.168.0.100",
						"white_list": "192.168.0.123/24",
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

func resourceHsmInstanceIsAppointDeviceDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

data "alibabacloudstack_zones" "default" {
	provider = alibabacloudstack-common
	available_resource_creation = "VSwitch"
}

data "alibabacloudstack_hsm_vendors" "default" {
}

data "alibabacloudstack_hsms" "default" {
  vendor_code = "${data.alibabacloudstack_hsm_vendors.default.vendors[0].code}"
  product_code = "${data.alibabacloudstack_hsm_vendors.default.vendors[0].products[0].code}"
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
