---
subcategory: "专有网络 VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_dhcp_options_sets"
sidebar_current: "docs-Alibabacloudstack-datasource-vpc-dhcp-options-sets"
description: |-
  提供阿里云账号下拥有的vpc dhcpoptionssets列表�?
---

# alibabacloudstack\_vpc\_dhcp_options_sets

此数据源提供根据指定过滤条件列出的阿里云账号下的vpc dhcpoptionssets资源列表�?

## 示例用法
```
variable "name" {
  default = "tf-testAccVpcDhcpOptionsDataSource-2860522"
}

resource "alibabacloudstack_vpc" "default0" {

  name = "${var.name}0"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc" "default1" {

  name = "${var.name}1"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_dhcp_options_set" "default" {
  	dhcp_options_set_name =        "${var.name}"
	dhcp_options_set_description = "${var.name}"
	domain_name =                  "aliyun.com"
	domain_name_servers =         "10.82.0.12,10.82.0.13"
	associate_vpcs = [
		"${alibabacloudstack_vpc.default0.id}",
		"${alibabacloudstack_vpc.default1.id}"
	]
}

 

data "alibabacloudstack_vpc_dhcp_options_sets" "default" {
  description_regex = "${alibabacloudstack_vpc_dhcp_options_set.default.dhcp_options_set_description}"
}
```

## 参数参�?
以下参数是支持的�?
  * `ids` - (选填) - dhcp options sets ID 列表�?
  * `name_regex` - (选填) - 正则表达式字符串，用于通过名称筛选dhcp options sets�?
  * `description_regex` - (选填) - 正则表达式字符串，用于通过描述筛选dhcp options sets�?
  * `domain_name` - (选填) - 主机名后缀
  * `dhcp_options_set_name` - (选填) - 代表资源名称的资源属性字�?

## Attributes Reference
除了上述参数外，还导出以下属性：
  * `dhcp_options_sets` - dhcp 选项集列表�?
    * `id` - dhcp 选项集列表的ID.
    * `associate_vpcs` - DHCP选项集关联的VPC的信息�?
    * `dhcp_options_set_description` - DHCP选项集的描述�?
    * `dhcp_options_set_id` - 代表资源一级ID的资源属性字�?
    * `dhcp_options_set_name` - 代表资源名称的资源属性字�?
    * `domain_name` - 主机名后缀
    * `domain_name_servers` - DNS服务器IP。最多传�?个DNS服务器IP，DNS服务器IP之间用半角逗号�?）隔开�?
    * `status` - 代表资源状态的资源属性字�?
