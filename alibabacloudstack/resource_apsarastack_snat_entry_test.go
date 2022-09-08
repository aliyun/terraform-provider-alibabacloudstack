package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccCheckSnatEntryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_snat_entry" {
			continue
		}

		_, err := vpcService.DescribeSnatEntry(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		return WrapError(Error("Snat entry still exist"))
	}

	return nil
}

func TestAccAlibabacloudStackSnatEntryBasic(t *testing.T) {
	var v vpc.SnatTableEntry

	resourceId := "alibabacloudstack_snat_entry.default"
	ra := resourceAttrInit(resourceId, testAccCheckSnatEntryBasicMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alibabacloudstack_snat_entry.default",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSnatEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSnatEntryConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSnatEntryConfig_snatIp(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func TestAccAlibabacloudStackSnatEntryMulti(t *testing.T) {
	var v vpc.SnatTableEntry

	resourceId := "alibabacloudstack_snat_entry.default.9"
	ra := resourceAttrInit(resourceId, testAccCheckSnatEntryBasicMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alibabacloudstack_snat_entry.default.9",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSnatEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSnatEntryConfigMulti(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func testAccSnatEntryConfigBasic(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccSnatEntryConfig%d"
}
data "alibabacloudstack_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alibabacloudstack_nat_gateway" "default" {
	depends_on = [alibabacloudstack_vpc.default]
	vpc_id = "${alibabacloudstack_vswitch.default.vpc_id}"
	specification = "Small"
	name = "${var.name}"
}

resource "alibabacloudstack_eip" "default" {
	name = "${var.name}"
}

resource "alibabacloudstack_eip_association" "default" {
	depends_on = [alibabacloudstack_eip.default, alibabacloudstack_nat_gateway.default]
	allocation_id = "${alibabacloudstack_eip.default.id}"
	instance_id = "${alibabacloudstack_nat_gateway.default.id}"
}

resource "alibabacloudstack_snat_entry" "default"{
	depends_on = [alibabacloudstack_eip_association.default, alibabacloudstack_nat_gateway.default]
	snat_table_id = "${alibabacloudstack_nat_gateway.default.snat_table_ids}"
	source_vswitch_id = "${alibabacloudstack_vswitch.default.id}"
	snat_ip = "${alibabacloudstack_eip.default.ip_address}"
  
}

resource "alibabacloudstack_snat_entry" "ecs"{
	depends_on = [alibabacloudstack_eip_association.default, alibabacloudstack_nat_gateway.default]
	snat_table_id = "${alibabacloudstack_nat_gateway.default.snat_table_ids}"
	source_cidr = "172.16.0.10/32"
	snat_ip = "${alibabacloudstack_eip.default.ip_address}"
}
`, rand)
}

func testAccSnatEntryConfig_snatIp(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccSnatEntryConfig%d"
}
data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alibabacloudstack_nat_gateway" "default" {
	depends_on = [alibabacloudstack_vpc.default]
	vpc_id = "${alibabacloudstack_vswitch.default.vpc_id}"
	specification = "Small"
	name = "${var.name}"
}

resource "alibabacloudstack_eip" "default" {
	name = "${var.name}"
}

resource "alibabacloudstack_eip_association" "default" {
	depends_on = [alibabacloudstack_eip.default, alibabacloudstack_nat_gateway.default]
	allocation_id = "${alibabacloudstack_eip.default.id}"
	instance_id = "${alibabacloudstack_nat_gateway.default.id}"
}

resource "alibabacloudstack_snat_entry" "default"{
	depends_on = [alibabacloudstack_eip_association.default, alibabacloudstack_nat_gateway.default]
	snat_table_id = "${alibabacloudstack_nat_gateway.default.snat_table_ids}"
	source_vswitch_id = "${alibabacloudstack_vswitch.default.id}"
	snat_ip = "${alibabacloudstack_eip.default.ip_address}"

}

resource "alibabacloudstack_snat_entry" "ecs"{
	depends_on = [alibabacloudstack_eip_association.default, alibabacloudstack_nat_gateway.default]
	snat_table_id = "${alibabacloudstack_nat_gateway.default.snat_table_ids}"
	source_cidr = "172.16.0.10/32"
	snat_ip = "${alibabacloudstack_eip.default.ip_address}"
}

`, rand)
}

func testAccSnatEntryConfigMulti(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccSnatEntryMulti%d"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "10.1.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
    count = 10
    vpc_id            = "${alibabacloudstack_vpc.default.id}"
    cidr_block        = "10.1.${count.index + 1}.0/24"
    availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
    name = "${var.name}"
}

resource "alibabacloudstack_nat_gateway" "default" {
	depends_on = [alibabacloudstack_vpc.default]
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	specification = "Small"
	name = "${var.name}"
}

resource "alibabacloudstack_eip" "default" {
	name = "${var.name}"
}

resource "alibabacloudstack_eip_association" "default" {
	depends_on = [alibabacloudstack_eip.default, alibabacloudstack_nat_gateway.default]
	allocation_id = "${alibabacloudstack_eip.default.id}"
	instance_id = "${alibabacloudstack_nat_gateway.default.id}"
}

resource "alibabacloudstack_snat_entry" "default"{
	depends_on = [alibabacloudstack_eip_association.default, alibabacloudstack_nat_gateway.default]
	count = "10"
	snat_table_id = "${alibabacloudstack_nat_gateway.default.snat_table_ids}"
	source_vswitch_id = "${element(alibabacloudstack_vswitch.default.*.id, count.index)}"
	snat_ip = "${alibabacloudstack_eip.default.ip_address}"

}

resource "alibabacloudstack_snat_entry" "ecs"{
	depends_on = [alibabacloudstack_eip_association.default, alibabacloudstack_nat_gateway.default]
	snat_table_id = "${alibabacloudstack_nat_gateway.default.snat_table_ids}"
	source_cidr = "10.1.2.10/32"
	snat_ip = "${alibabacloudstack_eip.default.ip_address}"
}
`, rand)
}

var testAccCheckSnatEntryBasicMap = map[string]string{
	"snat_table_id":     CHECKSET,
	"source_vswitch_id": CHECKSET,
	"snat_ip":           CHECKSET,
	"snat_entry_id":     CHECKSET,
}
