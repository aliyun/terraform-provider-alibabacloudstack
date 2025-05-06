package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackForwardEntryBasic(t *testing.T) {
	var v vpc.ForwardTableEntry
	resourceId := "alibabacloudstack_forward_entry.default"

	rand := getAccTestRandInt(10000, 20000)
	testAccForwardEntryCheckMap["name"] = fmt.Sprintf("tf-testAccForwardEntryConfig%d", rand)
	ra := resourceAttrInit(resourceId, testAccForwardEntryCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckForwardEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccForwardEntryConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
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

func TestAccAlibabacloudStackForwardEntryMulti(t *testing.T) {
	var v vpc.ForwardTableEntry
	resourceId := "alibabacloudstack_forward_entry.default.4"
	rand := acctest.RandInt()
	testAccForwardEntryCheckMap["name"] = fmt.Sprintf("tf-testAccForwardEntryConfig%d", rand)
	ra := resourceAttrInit(resourceId, testAccForwardEntryCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckForwardEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccForwardEntryConfig_multi(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"external_port": "84",
						"internal_port": "8084",
					}),
				),
			},
		},
	})
}

func testAccCheckForwardEntryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_forward_entry" {
			continue
		}
		if _, err := vpcService.DescribeForwardEntry(rs.Primary.ID); err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}

		return errmsgs.WrapError(fmt.Errorf("Forward entry %s still exist", rs.Primary.ID))
	}
	return nil
}

func testAccForwardEntryConfigBasic(rand int) string {
	config := fmt.Sprintf(`
%s

resource "alibabacloudstack_forward_entry" "default"{
	name = "${var.name}"
	forward_table_id = "${alibabacloudstack_nat_gateway.default.forward_table_ids}"
	external_ip = "${alibabacloudstack_eip.default.0.ip_address}"
	external_port = "80"
	ip_protocol = "tcp"
	internal_ip = "172.16.0.4"
	internal_port = "8080"


}
`, testAccForwardEntryConfigCommon(rand))
	return config
}

func testAccForwardEntryConfig_multi(rand int) string {
	config := fmt.Sprintf(`
%s

resource "alibabacloudstack_forward_entry" "default"{
	count = 5
	name = "${var.name}"
	forward_table_id = "${alibabacloudstack_nat_gateway.default.forward_table_ids}"
	external_ip = "${alibabacloudstack_eip.default.0.ip_address}"
	external_port = "${80 + count.index}"
	ip_protocol = "tcp"
	internal_ip = "172.16.0.3"
	internal_port = "${8080 + count.index}"
}
`, testAccForwardEntryConfigCommon(rand))
	return config
}

func testAccForwardEntryConfigCommon(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccForwardEntryConfig%d"
}

variable "number" {
	default = "2"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
	
	cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	
}

resource "alibabacloudstack_nat_gateway" "default" {
	vpc_id = "${alibabacloudstack_vswitch.default.vpc_id}"
	specification = "Small"
	
}

resource "alibabacloudstack_eip" "default" {
	count = "${var.number}"
	
}

resource "alibabacloudstack_eip_association" "default" {
	count = "${var.number}"
	allocation_id = "${element(alibabacloudstack_eip.default.*.id,count.index)}"
	instance_id = "${alibabacloudstack_nat_gateway.default.id}"
}
`, rand)
}

var testAccForwardEntryCheckMap = map[string]string{
	"forward_table_id": CHECKSET,
	"external_ip":      CHECKSET,
	"external_port":    "80",
	"ip_protocol":      "tcp",
	"internal_port":    "8080",
	"forward_entry_id": CHECKSET,
}
