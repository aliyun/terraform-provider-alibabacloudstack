package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCspprivateHsmGroup_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_cspprivate_hsm_group.default"

	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CspprivateService{testYundunProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeCspprivateHsmGroup")
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf_hsm_group%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCspprivateHsmGroupDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreYunCheck(t)
		},
		IDRefreshName:     resourceId,
		ExternalProviders: testAccExternalProviders,
		Providers:         testYunDunProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"hsm_count": "3",
					"hsm_list":  []string{"${alibabacloudstack_cspprivate_hsm_instance.default.id}", "${alibabacloudstack_cspprivate_hsm_instance.default1.id}"},
					"zone_ids":  []string{"${alibabacloudstack_cspprivate_hsm_instance.default.zone_id}"},
					"password":  "${random_password.password.0.result}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":     CHECKSET,
						"hsm_count":  "3",
						"hsm_list.#": "2",
						"zone_ids.#": "1",
						"group_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					// "vpc_id":    "${alibabacloudstack_cspprivate_hsm_instance.default[0].vpc_id}",
					"hsm_list": []string{"${alibabacloudstack_cspprivate_hsm_instance.default.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"hsm_list.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

var CspprivateHsmGroupMap = map[string]string{
	"vpc_id":             CHECKSET,
	"hsm_count":          CHECKSET,
	"zone_ids":           CHECKSET,
	"group_name":         CHECKSET,
	"create_time":        CHECKSET,
	"update_time":        CHECKSET,
	"security_level_tag": CHECKSET,
}

func resourceCspprivateHsmGroupDependence(name string) string {
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

resource "alibabacloudstack_vpc" "vpc" {	
	provider = alibabacloudstack-common
	vpc_name = "${var.name}"
	cidr_block = "192.168.0.0/16" # VPC CIDR block
}

resource "alibabacloudstack_vswitch" "vsw" {
	provider = alibabacloudstack-common
	vpc_id = alibabacloudstack_vpc.vpc.id
	cidr_block = "192.168.0.0/24" # VSwitch CIDR block
	availability_zone = data.alibabacloudstack_zones.default.zones.0.id
}

resource "alibabacloudstack_vpc" "vpc1" {	
	provider = alibabacloudstack-common
	vpc_name = "${var.name}"
	cidr_block = "192.168.0.0/16" # VPC CIDR block
}

resource "alibabacloudstack_vswitch" "vsw1" {
	provider = alibabacloudstack-common
	vpc_id = alibabacloudstack_vpc.vpc1.id
	cidr_block = "192.168.0.0/24" # VSwitch CIDR block
	availability_zone = data.alibabacloudstack_zones.default.zones.0.id
}

resource "alibabacloudstack_cspprivate_hsm_instance" "default" {
	product_code = "${data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.products.0.code}"
	vendor_code = "${data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.code}"
	vsm_type = "gvsm"
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
	alias_name = "${var.name}0"
	vpc_id = "${alibabacloudstack_vpc.vpc.id}"
	vpc_cidr_block = "${alibabacloudstack_vpc.vpc.cidr_block}"
	vswitch_id = "${alibabacloudstack_vswitch.vsw.id}"
	ip = "192.168.0.100"
}

resource "alibabacloudstack_cspprivate_hsm_instance" "default1" {
	product_code = "${data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.products.0.code}"
	vendor_code = "${data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.code}"
	vsm_type = "gvsm"
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
	alias_name = "${var.name}1"
	vpc_id = "${alibabacloudstack_vpc.vpc1.id}"
	vpc_cidr_block = "${alibabacloudstack_vpc.vpc1.cidr_block}"
	vswitch_id = "${alibabacloudstack_vswitch.vsw1.id}"
	ip = "192.168.0.101"
}

%s

`, name, RandomPasswordTestCase(12, 1))
}
