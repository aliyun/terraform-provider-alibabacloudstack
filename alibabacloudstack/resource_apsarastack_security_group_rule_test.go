package alibabacloudstack

import (
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackSecurityGroupRuleBasic(t *testing.T) {
	var v ecs.Permission
	resourceId := "alibabacloudstack_security_group_rule.default"
	ra := resourceAttrInit(resourceId, testAccCheckSecurityGroupRuleBasicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "abc",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSecurityGroupRule_cidrIp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_security_group_id": REMOVEKEY,
						"cidr_ip":                  "0.0.0.0/0",
						"description":              "abcd",
					}),
				),
			},
			{
				Config: testAccSecurityGroupRule_description,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "description",
					}),
				),
			},
			{
				Config: testAccSecurityGroupRule_all,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "abcd",
					}),
				),
			},
		},
	})

}

func TestAccAlibabacloudStackSecurityGroupEgressRule(t *testing.T) {
	var v ecs.Permission
	resourceId := "alibabacloudstack_security_group_rule.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"type":        "egress",
		"policy":      "accept",
		"description": "SHDRP-7513",
		"port_range":  "443/443",
		"priority":    "1",
		"cidr_ip":     "182.254.11.243/32",
	})
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupEgressRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "SHDRP-7513",
					}),
				),
			},
			{
				Config: testAccSecurityGroupEgressRule_description,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "SHDRP-7512",
					}),
				),
			},
		},
	})

}

func TestAccAlibabacloudStackSecurityGroupRuleMulti(t *testing.T) {
	var v ecs.Permission
	resourceId := "alibabacloudstack_security_group_rule.default.2"
	ra := resourceAttrInit(resourceId, testAccCheckSecurityGroupRuleBasicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr_ip":                  "45.20.250.240/32",
						"source_security_group_id": REMOVEKEY,
					}),
				),
			},
		},
	})

}

const testAccSecurityGroupRuleBasic = `
variable "name" {
  default = "tf-testAccSecurityGroupRuleBasic"
}

resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "192.168.0.0/16"
}

resource "alibabacloudstack_security_group" "default" {
  count  = 2
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  name = "${var.name}"
}

resource "alibabacloudstack_security_group_rule" "default" {
  type = "ingress"
  ip_protocol = "tcp"
  nic_type = "intranet"
  policy = "drop"
  port_range = "22/22"
  priority = 100
  security_group_id = "${alibabacloudstack_security_group.default.0.id}"
  source_security_group_id = "${alibabacloudstack_security_group.default.1.id}"
  description = "abc"
}
`
const testAccSecurityGroupRule_cidrIp = `

variable "name" {
  default = "tf-testAccSecurityGroupRuleBasic"
}

resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "192.168.0.0/16"
}

resource "alibabacloudstack_security_group" "default" {
  count = 2
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  name = "${var.name}"
}

resource "alibabacloudstack_security_group_rule" "default" {
  type = "ingress"
  ip_protocol = "tcp"
  nic_type = "intranet"
  policy = "drop"
  port_range = "22/22"
  priority = 100
  security_group_id = "${alibabacloudstack_security_group.default.0.id}"
  cidr_ip = "0.0.0.0/0"
  description = "abcd"
}
`

const testAccSecurityGroupRule_description = `

variable "name" {
  default = "tf-testAccSecurityGroupRuleBasic"
}

resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "192.168.0.0/16"
}

resource "alibabacloudstack_security_group" "default" {
  count = 2
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  name = "${var.name}"
}

resource "alibabacloudstack_security_group_rule" "default" {
  type = "ingress"
  ip_protocol = "tcp"
  nic_type = "intranet"
  policy = "drop"
  port_range = "22/22"
  priority = 100
  security_group_id = "${alibabacloudstack_security_group.default.0.id}"
  cidr_ip = "0.0.0.0/0"
  description = "description"
}
`

const testAccSecurityGroupRule_all = `

variable "name" {
  default = "tf-testAccSecurityGroupRuleBasic"
}

resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "192.168.0.0/16"
}

resource "alibabacloudstack_security_group" "default" {
  count = 2
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  name = "${var.name}"
}

resource "alibabacloudstack_security_group_rule" "default" {
  type = "ingress"
  ip_protocol = "tcp"
  nic_type = "intranet"
  policy = "drop"
  port_range = "22/22"
  priority = 100
  security_group_id = "${alibabacloudstack_security_group.default.0.id}"
  cidr_ip = "0.0.0.0/0"
  description = "abcd"
}
`

const testAccSecurityGroupRuleMulti = `
variable "name" {
  default = "tf-testAccSecurityGroupRuleBasic"
}

variable "cidr_ip_list" {
  type = list(string)
  default = ["50.255.255.255/32", "75.250.250.250/32", "45.20.250.240/32"]
}

resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alibabacloudstack_security_group" "default" {
  name = "${var.name}"
  vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_security_group_rule" "default" {
  count = "${length(compact(var.cidr_ip_list))}"
  security_group_id = "${alibabacloudstack_security_group.default.id}"
  type = "ingress"
  policy = "drop"
  port_range = "22/22"
  ip_protocol = "tcp"
  nic_type = "intranet"
  priority = 100
  cidr_ip = "${element(var.cidr_ip_list, count.index)}"
}
`

var testAccCheckSecurityGroupRuleBasicMap = map[string]string{
	"type":                     "ingress",
	"ip_protocol":              "tcp",
	"nic_type":                 "intranet",
	"policy":                   "drop",
	"port_range":               "22/22",
	"priority":                 "100",
	"security_group_id":        CHECKSET,
	"source_security_group_id": CHECKSET,
	"cidr_ip":                  "",
}

func testAccCheckSecurityGroupRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_security_group_rule" {
			continue
		}
		_, err := ecsService.DescribeSecurityGroupRule(rs.Primary.ID)

		// Verify the error is what we want
		if err != nil && !errmsgs.NotFoundError(err) {
			return errmsgs.WrapError(err)
		}
	}

	return nil
}

const testAccSecurityGroupEgressRule = `
variable "name" {
  default = "tf-testAccSecurityGroupRuleBasic"
}

resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "192.168.0.0/16"
}

resource "alibabacloudstack_security_group" "default" {
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  name = "${var.name}"
}

resource "alibabacloudstack_security_group_rule" "default" {
  type = "egress"
  ip_protocol = "tcp"
  nic_type = "intranet"
  policy = "accept"
  port_range = "443/443"
  priority = "1"
  security_group_id = "${alibabacloudstack_security_group.default.id}"
  cidr_ip = "182.254.11.243/32"
  description = "SHDRP-7513"
}
`

const testAccSecurityGroupEgressRule_description = `
variable "name" {
  default = "tf-testAccSecurityGroupRuleBasic"
}

resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "192.168.0.0/16"
}

resource "alibabacloudstack_security_group" "default" {
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  name = "${var.name}"
}

resource "alibabacloudstack_security_group_rule" "default" {
  type = "egress"
  ip_protocol = "tcp"
  nic_type = "intranet"
  policy = "accept"
  port_range = "443/443"
  priority = "1"
  security_group_id = "${alibabacloudstack_security_group.default.id}"
  cidr_ip = "182.254.11.243/32"
  description = "SHDRP-7512"
}
`
