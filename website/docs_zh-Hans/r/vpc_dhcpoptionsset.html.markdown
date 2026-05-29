---
subcategory: "Virtual Private Cloud (VPC)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_dhcpoptionsset"
sidebar_current: "docs-Alibabacloudstack-vpc-dhcpoptionsset"
description: |-
  Provides a vpc Dhcpoptionsset resource.
---

# alibabacloudstack\_vpc\_dhcpoptionsset

Provides a vpc Dhcpoptionsset resource.

## 示例用法
```
variable "name" {
  default = "tf-testAccVpcDhcpoptionssetBasic"
}

resource "alibabacloudstack_vpc" "default0" {

  name = "${var.name}0"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc" "default1" {

  name = "${var.name}1"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc" "default2" {

  name = "${var.name}2"
  cidr_block = "172.16.0.0/16"
}



resource "alibabacloudstack_vpc_dhcp_options_set" "default" {
  domain_name = "aliyun.com"
  domain_name_servers = "10.82.0.12,10.82.0.13"
  associate_vpcs = [
                     "${alibabacloudstack_vpc.default0.id}",
                     "${alibabacloudstack_vpc.default1.id}",
                     "${alibabacloudstack_vpc.default2.id}"
                   ]
  dhcp_options_set_name = "${var.name}"
  dhcp_options_set_description = "${var.name}"
}
```

## 参数参考

支持以下参数：
  * `associate_vpcs` - (选填) - 要关联到此 DHCP 选项集的 VPC ID 列表。
  * `domain_name` - (选填) - 域名后缀，例如 example.com。将 DHCP 选项集与 VPC 关联后，VPC 内的 ECS 实例会自动使用该域名后缀。
  * `domain_name_servers` - (选填) - DNS 服务器 IP 地址。最多可指定 4 个 DNS 服务器 IP，IP 地址之间需用半角逗号（,）分隔。未指定时，VPC 内所有 ECS 实例将使用阿里云 DNS 服务器（100.100.2.136 和 100.100.2.138）。

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `dhcp_options_set_id` - DHCP 选项集的 ID
  * `status` - DHCP 选项集的状态。取值：`Available`、`Pending`、`Deleting`

## 导入

DHCP 选项集可以使用 DhcpOptionsSetId 导入，例如：

```
$ terraform import alibabacloudstack_vpc_dhcp_options_set.example dos-12345678
```
