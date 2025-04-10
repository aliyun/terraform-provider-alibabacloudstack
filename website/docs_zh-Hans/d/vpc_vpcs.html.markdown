---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_vpcs"
sidebar_current: "docs-Alibabacloudstack-datasource-vpc-vpcs"
description: |- 
  查询专有网络（VPC）实例
---

# alibabacloudstack_vpc_vpcs
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_vpcs`

根据指定过滤条件列出当前凭证权限可以访问的专有网络（VPC）实例列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAlibabacloudstackVpcVpcs20210"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details             = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name        = "${var.name}_vsw"
  vpc_id      = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block  = "172.16.0.0/24"
  zone_id    = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

data "alibabacloudstack_vpc_vpcs" "default" {
  ids         = ["${alibabacloudstack_vpc_vpc.default.id}"]
  cidr_block  = "172.16.0.0/16"
  status      = "Available"
  name_regex  = "^${var.name}_vpc"
  is_default  = false
  output_file = "vpc_list.txt"

  dhcp_options_set_id = ""
  dry_run             = false
  resource_group_id   = ""
  vpc_name            = ""
  vpc_owner_id        = ""
  enable_details      = true
}

output "first_vpc_id" {
  value = data.alibabacloudstack_vpc_vpcs.default.vpcs[0].id
}
```

## 参数说明

以下参数是支持的：

* `cidr_block` - (选填, 变更时重建) VPC的网段。建议您使用以下RFC标准私有网段或其子集作为VPC的主IPv4网段：
  * `192.168.0.0/16`
  * `172.16.0.0/12`
  * `10.0.0.0/8`
  网段掩码的有效范围为8到28位。
  您也可以使用除以下网段及其子网外的自定义地址段作为VPC的主IPv4网段：
  * `100.64.0.0/10`
  * `224.0.0.0/4`
  * `127.0.0.0/8`
  * `169.254.0.0/16`

* `status` - (选填, 变更时重建) VPC的状态。有效值包括：
  * `Pending`：VPC正在配置中。
  * `Available`：VPC可用。

* `name_regex` - (选填, 变更时重建) 用于按名称筛选VPC的正则表达式字符串。

* `is_default` - (选填, 变更时重建) 是否创建指定地域下的默认VPC。有效值：
  * `true`
  * `false`(默认)

* `vswitch_id` - (选填, 变更时重建) 通过指定的交换机过滤结果。

* `ids` - (选填, 变更时重建) VPC ID列表。

* `dhcp_options_set_id` - (选填, 变更时重建) DHCP选项集的ID。

* `dry_run` - (选填, 变更时重建) 是否执行模拟运行。有效值：
  * `true`：仅执行模拟运行。系统会检查所需的参数、请求语法和限制。如果请求未能通过模拟运行，则返回错误消息。如果请求通过了模拟运行，则返回 `DryRunOperation` 错误代码。
  * `false`(默认)：执行模拟运行并发送请求。如果请求通过了模拟运行，则返回 2xx HTTP 状态码并执行操作。

* `resource_group_id` - (选填, 变更时重建) 需要移入云资源实例的资源组ID。

* `vpc_name` - (选填, 变更时重建) 要修改的VPC名称。名称长度为1～128个字符，不能以`http://`或`https://`开头。

* `vpc_owner_id` - (选填, 变更时重建) VPC的所有者ID。

* `enable_details` - (选填) 默认为 `true`。将其设置为 `true` 以输出 `route_table_id`。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - VPC ID列表。
* `names` - VPC名称列表。
* `vpcs` - VPC列表。每个元素包含以下属性：
  * `id` - VPC的ID。
  * `region_id` - VPC所在区域的ID。
  * `resource_group_id` - VPC所属资源组的ID。
  * `status` - VPC的状态。
  * `vpc_name` - VPC的名称。
  * `vswitch_ids` - 指定VPC中的交换机ID列表。
  * `cidr_block` - VPC的CIDR块。
  * `vrouter_id` - 路由器的ID。
  * `route_table_id` - 路由器的路由表ID。
  * `description` - VPC的描述。
  * `is_default` - 是否为该区域中的默认VPC。
  * `creation_time` - VPC创建的时间。
  * `tags` - 分配给VPC的标签映射。
  * `ipv6_cidr_block` - VPC的IPv6 CIDR块。
  * `router_id` - 路由器的ID。
  * `secondary_cidr_blocks` - VPC的次级IPv4 CIDR块列表。
  * `user_cidrs` - 用户CIDR列表。
  * `vpc_id` - VPC的ID。
  * `available_ip_address_count` - VPC中可用IP地址的数量。