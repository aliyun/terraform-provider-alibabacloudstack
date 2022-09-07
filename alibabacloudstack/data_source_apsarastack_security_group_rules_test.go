package alibabacloudstack

import (
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackSecurityGroupRulesDataSourceWithDirection(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackSecurityGroupRulesDataSourceConfigDirection,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_security_group_rules.ingress"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.ingress", "group_name", "tf-testAccCheckAlibabacloudStackSecurityGroupRulesDataSourceConfig_1"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.ingress", "rules.#", "1"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.ingress", "rules.0.direction", "ingress"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.ingress", "rules.0.ip_protocol", "tcp"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.ingress", "rules.0.port_range", "5000/5001"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_security_group_rules.ingress", "rules.0.cidr_ip"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_security_group_rules.ingress", "rules.0.source_cidr_ip"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.ingress", "rules.0.source_group_id", ""),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.ingress", "rules.0.source_group_owner_account", ""),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.ingress", "rules.0.dest_cidr_ip", ""),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.ingress", "rules.0.dest_group_id", ""),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.ingress", "rules.0.dest_group_owner_account", ""),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_security_group_rules.ingress", "rules.0.priority"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_security_group_rules.ingress", "rules.0.nic_type"),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackSecurityGroupRulesDataSourceWithGroupId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackSecurityGroupRulesDataSourceConfigGroup_id,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_security_group_rules.egress"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "group_name", "tf-testAccCheckAlibabacloudStackSecurityGroupRulesDataSourceConfig0"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.#", "1"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.direction", "egress"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.ip_protocol", "udp"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.nic_type", "intranet"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.cidr_ip"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.port_range", "6000/6001"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.source_cidr_ip", ""),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.source_group_id", ""),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.source_group_owner_account", ""),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.dest_cidr_ip", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.dest_group_id", ""),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.dest_group_owner_account", ""),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_security_group_rules.egress", "rules.0.priority"),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackSecurityGroupRulesDataSourceWithNic_Type(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.EcsClassicSupportedRegions)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackSecurityGroupRulesDataSourceConfigNic_Type,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_security_group_rules.egress"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.#", "1"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "group_name", "tf-testAccCheckAlibabacloudStackSecurityGroupRulesDataSourceConfig1"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.direction", "egress"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.ip_protocol", "udp"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.nic_type", "intranet"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.port_range", "6000/6001"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.cidr_ip"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.source_cidr_ip", ""),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.source_group_id", ""),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.source_group_owner_account", ""),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.dest_cidr_ip", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.dest_group_id", ""),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.dest_group_owner_account", ""),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_security_group_rules.egress", "rules.0.priority"),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackSecurityGroupRulesDataSourceWithPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackSecurityGroupRulesDataSourceConfigPolicy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_security_group_rules.egress"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "group_name", "tf-testAccCheckAlibabacloudStackSecurityGroupRulesDataSourceConfig3"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.#", "1"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.direction", "egress"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.ip_protocol", "udp"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.nic_type", "intranet"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.port_range", "6000/6001"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.cidr_ip"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.policy", "drop"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.source_cidr_ip", ""),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.source_group_id", ""),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.source_group_owner_account", ""),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.dest_cidr_ip", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.dest_group_id", ""),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.dest_group_owner_account", ""),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_security_group_rules.egress", "rules.0.priority"),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackSecurityGroupRulesDataSourceWithIp_Protocol(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackSecurityGroupRulesDataSourceConfigIp_Protocol,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_security_group_rules.egress"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "group_name", "tf-testAccCheckAlibabacloudStackSecurityGroupRulesDataSourceConfig2"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.#", "1"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.direction", "egress"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.ip_protocol", "udp"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.nic_type", "intranet"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.port_range", "6000/6001"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.cidr_ip"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.source_cidr_ip", ""),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.source_group_id", ""),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.source_group_owner_account", ""),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.dest_cidr_ip", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.dest_group_id", ""),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.egress", "rules.0.dest_group_owner_account", ""),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_security_group_rules.egress", "rules.0.priority"),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackSecurityGroupRulesDataSourceEmpty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackSecurityGroupRulesDataSourceConfigEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_security_group_rules.empty"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.empty", "group_name", "tf-testAccCheckAlibabacloudStackSecurityGroupRulesDataSourceConfigEgress"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_security_group_rules.empty", "rules.#", "0"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_security_group_rules.empty", "rules.0.direction"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_security_group_rules.empty", "rules.0.ip_protocol"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_security_group_rules.empty", "rules.0.nic_type"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_security_group_rules.empty", "rules.0.policy"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_security_group_rules.empty", "rules.0.port_range"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_security_group_rules.empty", "rules.0.priority"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_security_group_rules.empty", "rules.0.source_cidr_ip"),
				),
			},
		},
	})
}

const testAccCheckAlibabacloudStackSecurityGroupRulesDataSourceConfigDirection = `
variable "name" {
	default = "tf-testAccCheckAlibabacloudStackSecurityGroupRulesDataSourceConfig_1"
}
resource "alibabacloudstack_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name = "${var.name}"
}

resource "alibabacloudstack_security_group" "group" {
  name = "${var.name}"
  vpc_id      = "${alibabacloudstack_vpc.foo.id}"
}

resource "alibabacloudstack_security_group_rule" "rule_ingress" {
  type              = "ingress"
  ip_protocol       = "tcp"
  port_range        = "5000/5001"
  security_group_id = "${alibabacloudstack_security_group.group.id}"
  cidr_ip           = "0.0.0.0/0"
}

resource "alibabacloudstack_security_group_rule" "rule_egress" {
  type              = "egress"
  ip_protocol       = "udp"
  port_range        = "6000/6001"
  security_group_id = "${alibabacloudstack_security_group.group.id}"
  cidr_ip           = "0.0.0.0/0"
}

data "alibabacloudstack_security_group_rules" "ingress" {
  direction   = "ingress"
  group_id    = "${alibabacloudstack_security_group_rule.rule_ingress.security_group_id}"
}
`

const testAccCheckAlibabacloudStackSecurityGroupRulesDataSourceConfigGroup_id = `
variable "name" {
	default = "tf-testAccCheckAlibabacloudStackSecurityGroupRulesDataSourceConfig0"
}
resource "alibabacloudstack_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name = "${var.name}"
}

resource "alibabacloudstack_security_group" "group" {
  name = "${var.name}"
  vpc_id      = "${alibabacloudstack_vpc.foo.id}"
}

resource "alibabacloudstack_security_group" "bar" {
  name = "tf-testAccCheckAlibabacloudStackSecurityGroupRules"
  vpc_id      = "${alibabacloudstack_vpc.foo.id}"
}

resource "alibabacloudstack_security_group_rule" "rule_ingress" {
  type              = "ingress"
  ip_protocol       = "tcp"
  port_range        = "5000/5001"
  security_group_id = "${alibabacloudstack_security_group.bar.id}"
  cidr_ip           = "0.0.0.0/0"
}

resource "alibabacloudstack_security_group_rule" "rule_egress" {
  type              = "egress"
  ip_protocol       = "udp"
  port_range        = "6000/6001"
  security_group_id = "${alibabacloudstack_security_group.group.id}"
  cidr_ip           = "0.0.0.0/0"
}

data "alibabacloudstack_security_group_rules" "egress" {
  group_id    = "${alibabacloudstack_security_group_rule.rule_egress.security_group_id}"
}
`

const testAccCheckAlibabacloudStackSecurityGroupRulesDataSourceConfigNic_Type = `
variable "name" {
	default = "tf-testAccCheckAlibabacloudStackSecurityGroupRulesDataSourceConfig1"
}
resource "alibabacloudstack_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name = "${var.name}"
}

resource "alibabacloudstack_security_group" "group" {
  name = "${var.name}"
  vpc_id      = "${alibabacloudstack_vpc.foo.id}"
}

resource "alibabacloudstack_security_group_rule" "rule_egress" {
  type              = "egress"
  ip_protocol       = "udp"
  port_range        = "6000/6001"
  security_group_id = "${alibabacloudstack_security_group.group.id}"
  cidr_ip           = "0.0.0.0/0"
  nic_type          = "intranet"
}

data "alibabacloudstack_security_group_rules" "egress" {
  nic_type   = "intranet"
  group_id    = "${alibabacloudstack_security_group_rule.rule_egress.security_group_id}"
}
`

const testAccCheckAlibabacloudStackSecurityGroupRulesDataSourceConfigIp_Protocol = `
variable "name" {
	default = "tf-testAccCheckAlibabacloudStackSecurityGroupRulesDataSourceConfig2"
}
resource "alibabacloudstack_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name = "${var.name}"
}

resource "alibabacloudstack_security_group" "group" {
  name = "${var.name}"
  vpc_id      = "${alibabacloudstack_vpc.foo.id}"
}

resource "alibabacloudstack_security_group_rule" "rule_ingress" {
  type              = "ingress"
  ip_protocol       = "tcp"
  port_range        = "5000/5001"
  security_group_id = "${alibabacloudstack_security_group.group.id}"
  cidr_ip           = "0.0.0.0/0"
  nic_type          = "intranet"
}

resource "alibabacloudstack_security_group_rule" "rule_egress" {
  type              = "egress"
  ip_protocol       = "udp"
  port_range        = "6000/6001"
  security_group_id = "${alibabacloudstack_security_group.group.id}"
  cidr_ip           = "0.0.0.0/0"
  nic_type          = "intranet"
}

data "alibabacloudstack_security_group_rules" "egress" {
  ip_protocol   = "udp"
  group_id    = "${alibabacloudstack_security_group_rule.rule_egress.security_group_id}"
}
`

const testAccCheckAlibabacloudStackSecurityGroupRulesDataSourceConfigPolicy = `
variable "name" {
	default = "tf-testAccCheckAlibabacloudStackSecurityGroupRulesDataSourceConfig3"
}
resource "alibabacloudstack_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name = "${var.name}"
}

resource "alibabacloudstack_security_group" "group" {
  name = "${var.name}"
  vpc_id      = "${alibabacloudstack_vpc.foo.id}"
}

resource "alibabacloudstack_security_group_rule" "rule_ingress" {
  type              = "ingress"
  ip_protocol       = "tcp"
  port_range        = "5000/5001"
  security_group_id = "${alibabacloudstack_security_group.group.id}"
  cidr_ip           = "0.0.0.0/0"
  nic_type          = "intranet"
}

resource "alibabacloudstack_security_group_rule" "rule_egress" {
  type              = "egress"
  ip_protocol       = "udp"
  port_range        = "6000/6001"
  policy            = "drop"
  security_group_id = "${alibabacloudstack_security_group.group.id}"
  cidr_ip           = "0.0.0.0/0"
  nic_type          = "intranet"
}

data "alibabacloudstack_security_group_rules" "egress" {
  policy   = "drop"
  group_id   ="${alibabacloudstack_security_group_rule.rule_egress.security_group_id}"
}
`

const testAccCheckAlibabacloudStackSecurityGroupRulesDataSourceConfigEmpty = `
variable "name" {
	default = "tf-testAccCheckAlibabacloudStackSecurityGroupRulesDataSourceConfigEgress"
}
resource "alibabacloudstack_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name = "${var.name}"
}

resource "alibabacloudstack_security_group" "group" {
  name = "${var.name}"
  vpc_id      = "${alibabacloudstack_vpc.foo.id}"
}

data "alibabacloudstack_security_group_rules" "empty" {
  group_id    = "${alibabacloudstack_security_group.group.id}"
}
`
