package apsarastack

import (
	"testing"

	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackSecurityGroupRulesDataSourceWithDirection(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackSecurityGroupRulesDataSourceConfigDirection,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_security_group_rules.ingress"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.ingress", "group_name", "tf-testAccCheckApsaraStackSecurityGroupRulesDataSourceConfig_1"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.ingress", "rules.#", "1"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.ingress", "rules.0.direction", "ingress"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.ingress", "rules.0.ip_protocol", "tcp"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.ingress", "rules.0.port_range", "5000/5001"),
					resource.TestCheckNoResourceAttr("data.apsarastack_security_group_rules.ingress", "rules.0.cidr_ip"),
					resource.TestCheckResourceAttrSet("data.apsarastack_security_group_rules.ingress", "rules.0.source_cidr_ip"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.ingress", "rules.0.source_group_id", ""),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.ingress", "rules.0.source_group_owner_account", ""),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.ingress", "rules.0.dest_cidr_ip", ""),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.ingress", "rules.0.dest_group_id", ""),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.ingress", "rules.0.dest_group_owner_account", ""),
					resource.TestCheckResourceAttrSet("data.apsarastack_security_group_rules.ingress", "rules.0.priority"),
					resource.TestCheckResourceAttrSet("data.apsarastack_security_group_rules.ingress", "rules.0.nic_type"),
				),
			},
		},
	})
}

func TestAccApsaraStackSecurityGroupRulesDataSourceWithGroupId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackSecurityGroupRulesDataSourceConfigGroup_id,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_security_group_rules.egress"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "group_name", "tf-testAccCheckApsaraStackSecurityGroupRulesDataSourceConfig0"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.#", "1"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.direction", "egress"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.ip_protocol", "udp"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.nic_type", "intranet"),
					resource.TestCheckNoResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.cidr_ip"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.port_range", "6000/6001"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.source_cidr_ip", ""),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.source_group_id", ""),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.source_group_owner_account", ""),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.dest_cidr_ip", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.dest_group_id", ""),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.dest_group_owner_account", ""),
					resource.TestCheckResourceAttrSet("data.apsarastack_security_group_rules.egress", "rules.0.priority"),
				),
			},
		},
	})
}

func TestAccApsaraStackSecurityGroupRulesDataSourceWithNic_Type(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.EcsClassicSupportedRegions)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackSecurityGroupRulesDataSourceConfigNic_Type,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_security_group_rules.egress"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.#", "1"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "group_name", "tf-testAccCheckApsaraStackSecurityGroupRulesDataSourceConfig1"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.direction", "egress"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.ip_protocol", "udp"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.nic_type", "intranet"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.port_range", "6000/6001"),
					resource.TestCheckNoResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.cidr_ip"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.source_cidr_ip", ""),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.source_group_id", ""),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.source_group_owner_account", ""),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.dest_cidr_ip", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.dest_group_id", ""),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.dest_group_owner_account", ""),
					resource.TestCheckResourceAttrSet("data.apsarastack_security_group_rules.egress", "rules.0.priority"),
				),
			},
		},
	})
}

func TestAccApsaraStackSecurityGroupRulesDataSourceWithPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackSecurityGroupRulesDataSourceConfigPolicy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_security_group_rules.egress"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "group_name", "tf-testAccCheckApsaraStackSecurityGroupRulesDataSourceConfig3"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.#", "1"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.direction", "egress"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.ip_protocol", "udp"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.nic_type", "intranet"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.port_range", "6000/6001"),
					resource.TestCheckNoResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.cidr_ip"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.policy", "drop"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.source_cidr_ip", ""),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.source_group_id", ""),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.source_group_owner_account", ""),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.dest_cidr_ip", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.dest_group_id", ""),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.dest_group_owner_account", ""),
					resource.TestCheckResourceAttrSet("data.apsarastack_security_group_rules.egress", "rules.0.priority"),
				),
			},
		},
	})
}

func TestAccApsaraStackSecurityGroupRulesDataSourceWithIp_Protocol(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackSecurityGroupRulesDataSourceConfigIp_Protocol,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_security_group_rules.egress"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "group_name", "tf-testAccCheckApsaraStackSecurityGroupRulesDataSourceConfig2"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.#", "1"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.direction", "egress"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.ip_protocol", "udp"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.nic_type", "intranet"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.port_range", "6000/6001"),
					resource.TestCheckNoResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.cidr_ip"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.source_cidr_ip", ""),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.source_group_id", ""),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.source_group_owner_account", ""),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.dest_cidr_ip", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.dest_group_id", ""),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.egress", "rules.0.dest_group_owner_account", ""),
					resource.TestCheckResourceAttrSet("data.apsarastack_security_group_rules.egress", "rules.0.priority"),
				),
			},
		},
	})
}

func TestAccApsaraStackSecurityGroupRulesDataSourceEmpty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackSecurityGroupRulesDataSourceConfigEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_security_group_rules.empty"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.empty", "group_name", "tf-testAccCheckApsaraStackSecurityGroupRulesDataSourceConfigEgress"),
					resource.TestCheckResourceAttr("data.apsarastack_security_group_rules.empty", "rules.#", "0"),
					resource.TestCheckNoResourceAttr("data.apsarastack_security_group_rules.empty", "rules.0.direction"),
					resource.TestCheckNoResourceAttr("data.apsarastack_security_group_rules.empty", "rules.0.ip_protocol"),
					resource.TestCheckNoResourceAttr("data.apsarastack_security_group_rules.empty", "rules.0.nic_type"),
					resource.TestCheckNoResourceAttr("data.apsarastack_security_group_rules.empty", "rules.0.policy"),
					resource.TestCheckNoResourceAttr("data.apsarastack_security_group_rules.empty", "rules.0.port_range"),
					resource.TestCheckNoResourceAttr("data.apsarastack_security_group_rules.empty", "rules.0.priority"),
					resource.TestCheckNoResourceAttr("data.apsarastack_security_group_rules.empty", "rules.0.source_cidr_ip"),
				),
			},
		},
	})
}

const testAccCheckApsaraStackSecurityGroupRulesDataSourceConfigDirection = `
variable "name" {
	default = "tf-testAccCheckApsaraStackSecurityGroupRulesDataSourceConfig_1"
}
resource "apsarastack_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name = "${var.name}"
}

resource "apsarastack_security_group" "group" {
  name = "${var.name}"
  vpc_id      = "${apsarastack_vpc.foo.id}"
}

resource "apsarastack_security_group_rule" "rule_ingress" {
  type              = "ingress"
  ip_protocol       = "tcp"
  port_range        = "5000/5001"
  security_group_id = "${apsarastack_security_group.group.id}"
  cidr_ip           = "0.0.0.0/0"
}

resource "apsarastack_security_group_rule" "rule_egress" {
  type              = "egress"
  ip_protocol       = "udp"
  port_range        = "6000/6001"
  security_group_id = "${apsarastack_security_group.group.id}"
  cidr_ip           = "0.0.0.0/0"
}

data "apsarastack_security_group_rules" "ingress" {
  direction   = "ingress"
  group_id    = "${apsarastack_security_group_rule.rule_ingress.security_group_id}"
}
`

const testAccCheckApsaraStackSecurityGroupRulesDataSourceConfigGroup_id = `
variable "name" {
	default = "tf-testAccCheckApsaraStackSecurityGroupRulesDataSourceConfig0"
}
resource "apsarastack_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name = "${var.name}"
}

resource "apsarastack_security_group" "group" {
  name = "${var.name}"
  vpc_id      = "${apsarastack_vpc.foo.id}"
}

resource "apsarastack_security_group" "bar" {
  name = "tf-testAccCheckApsaraStackSecurityGroupRules"
  vpc_id      = "${apsarastack_vpc.foo.id}"
}

resource "apsarastack_security_group_rule" "rule_ingress" {
  type              = "ingress"
  ip_protocol       = "tcp"
  port_range        = "5000/5001"
  security_group_id = "${apsarastack_security_group.bar.id}"
  cidr_ip           = "0.0.0.0/0"
}

resource "apsarastack_security_group_rule" "rule_egress" {
  type              = "egress"
  ip_protocol       = "udp"
  port_range        = "6000/6001"
  security_group_id = "${apsarastack_security_group.group.id}"
  cidr_ip           = "0.0.0.0/0"
}

data "apsarastack_security_group_rules" "egress" {
  group_id    = "${apsarastack_security_group_rule.rule_egress.security_group_id}"
}
`

const testAccCheckApsaraStackSecurityGroupRulesDataSourceConfigNic_Type = `
variable "name" {
	default = "tf-testAccCheckApsaraStackSecurityGroupRulesDataSourceConfig1"
}
resource "apsarastack_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name = "${var.name}"
}

resource "apsarastack_security_group" "group" {
  name = "${var.name}"
  vpc_id      = "${apsarastack_vpc.foo.id}"
}

resource "apsarastack_security_group_rule" "rule_egress" {
  type              = "egress"
  ip_protocol       = "udp"
  port_range        = "6000/6001"
  security_group_id = "${apsarastack_security_group.group.id}"
  cidr_ip           = "0.0.0.0/0"
  nic_type          = "intranet"
}

data "apsarastack_security_group_rules" "egress" {
  nic_type   = "intranet"
  group_id    = "${apsarastack_security_group_rule.rule_egress.security_group_id}"
}
`

const testAccCheckApsaraStackSecurityGroupRulesDataSourceConfigIp_Protocol = `
variable "name" {
	default = "tf-testAccCheckApsaraStackSecurityGroupRulesDataSourceConfig2"
}
resource "apsarastack_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name = "${var.name}"
}

resource "apsarastack_security_group" "group" {
  name = "${var.name}"
  vpc_id      = "${apsarastack_vpc.foo.id}"
}

resource "apsarastack_security_group_rule" "rule_ingress" {
  type              = "ingress"
  ip_protocol       = "tcp"
  port_range        = "5000/5001"
  security_group_id = "${apsarastack_security_group.group.id}"
  cidr_ip           = "0.0.0.0/0"
  nic_type          = "intranet"
}

resource "apsarastack_security_group_rule" "rule_egress" {
  type              = "egress"
  ip_protocol       = "udp"
  port_range        = "6000/6001"
  security_group_id = "${apsarastack_security_group.group.id}"
  cidr_ip           = "0.0.0.0/0"
  nic_type          = "intranet"
}

data "apsarastack_security_group_rules" "egress" {
  ip_protocol   = "udp"
  group_id    = "${apsarastack_security_group_rule.rule_egress.security_group_id}"
}
`

const testAccCheckApsaraStackSecurityGroupRulesDataSourceConfigPolicy = `
variable "name" {
	default = "tf-testAccCheckApsaraStackSecurityGroupRulesDataSourceConfig3"
}
resource "apsarastack_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name = "${var.name}"
}

resource "apsarastack_security_group" "group" {
  name = "${var.name}"
  vpc_id      = "${apsarastack_vpc.foo.id}"
}

resource "apsarastack_security_group_rule" "rule_ingress" {
  type              = "ingress"
  ip_protocol       = "tcp"
  port_range        = "5000/5001"
  security_group_id = "${apsarastack_security_group.group.id}"
  cidr_ip           = "0.0.0.0/0"
  nic_type          = "intranet"
}

resource "apsarastack_security_group_rule" "rule_egress" {
  type              = "egress"
  ip_protocol       = "udp"
  port_range        = "6000/6001"
  policy            = "drop"
  security_group_id = "${apsarastack_security_group.group.id}"
  cidr_ip           = "0.0.0.0/0"
  nic_type          = "intranet"
}

data "apsarastack_security_group_rules" "egress" {
  policy   = "drop"
  group_id   ="${apsarastack_security_group_rule.rule_egress.security_group_id}"
}
`

const testAccCheckApsaraStackSecurityGroupRulesDataSourceConfigEmpty = `
variable "name" {
	default = "tf-testAccCheckApsaraStackSecurityGroupRulesDataSourceConfigEgress"
}
resource "apsarastack_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  name = "${var.name}"
}

resource "apsarastack_security_group" "group" {
  name = "${var.name}"
  vpc_id      = "${apsarastack_vpc.foo.id}"
}

data "apsarastack_security_group_rules" "empty" {
  group_id    = "${apsarastack_security_group.group.id}"
}
`
