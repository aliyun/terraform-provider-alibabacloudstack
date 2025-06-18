---
subcategory: "VPC"
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
  * `associate_vpcs` - (选填) - DHCP选项集关联的VPC的信息。
  * `dhcp_options_set_description` - (选填) - DHCP选项集的描述。
  * `dhcp_options_set_name` - (选填) - 代表资源名称的资源属性字段
  * `domain_name` - (选填) - 主机名后缀
  * `domain_name_servers` - (选填) - DNS服务器IP。最多传入4个DNS服务器IP，DNS服务器IP之间用半角逗号（,）隔开。

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `associate_vpcs` - DHCP选项集关联的VPC的信息。
  * `dhcp_options_set_id` - 代表资源一级ID的资源属性字段
  * `status` - 代表资源状态的资源属性字段
