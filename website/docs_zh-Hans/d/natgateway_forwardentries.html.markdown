---
subcategory: "NATGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_natgateway_forwardentries"
sidebar_current: "docs-Alibabacloudstack-datasource-natgateway-forwardentries"
description: |- 
  查询专有网络的NAT网关DNAT表规则
---

# alibabacloudstack_natgateway_forwardentries
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_forward_entries`

根据指定过滤条件列出当前凭证权限可以访问的专有网络的NAT网关DNAT表规则列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAccForwardEntryConfig17257"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/21"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alibabacloudstack_nat_gateway" "default" {
  vpc_id        = "${alibabacloudstack_vswitch.default.vpc_id}"
  specification = "Small"
  name          = "${var.name}"
}

resource "alibabacloudstack_eip" "default" {
  name = "${var.name}"
}

resource "alibabacloudstack_eip_association" "default" {
  allocation_id = "${alibabacloudstack_eip.default.id}"
  instance_id   = "${alibabacloudstack_nat_gateway.default.id}"
}

resource "alibabacloudstack_forward_entry" "default" {
  forward_table_id = "${alibabacloudstack_nat_gateway.default.forward_table_ids}"
  external_ip      = "${alibabacloudstack_eip.default.ip_address}"
  external_port    = "80"
  ip_protocol      = "tcp"
  internal_ip     = "172.16.0.3"
  internal_port   = "8080"
}

data "alibabacloudstack_natgateway_forwardentries" "default" {
  forward_table_id = "${alibabacloudstack_forward_entry.default.forward_table_id}"
  external_ip      = "${alibabacloudstack_eip.default.ip_address}"
  internal_ip      = "172.16.0.3"
  name_regex       = "example.*"
  output_file      = "forward_entries_output.txt"
}

output "natgateway_forward_entries" {
  value = "${data.alibabacloudstack_natgateway_forwardentries.default.entries}"
}
```

## 参数说明

以下参数是支持的：

* `forward_table_id` - (必填, 变更时重建) DNAT条目所属DNAT表的ID。此参数用于指定查询范围，限定为特定DNAT表中的条目。
* `name_regex` - (选填) 用于通过Forward Entry名称筛选结果的正则表达式字符串。可以通过此参数匹配符合特定命名规则的DNAT条目。
* `external_ip` - (选填) 公网IP地址。公网IP地址用于ECS实例接收来自互联网的请求。可以通过此参数筛选指定公网IP地址相关的DNAT条目。
* `internal_ip` - (选填) 映射到DNAT条目中公网IP地址的私网IP地址。可以通过此参数筛选指定私网IP地址相关的DNAT条目。
* `ids` - (选填) Forward Entry ID列表。可以通过此参数直接指定需要查询的DNAT条目ID列表。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - Forward Entry ID列表。此属性返回所有符合条件的DNAT条目的ID。
* `names` - Forward Entry名称列表。此属性返回所有符合条件的DNAT条目的名称。
* `entries` - Forward Entries列表。每个元素包含以下属性：
  * `id` - Forward Entry的ID。唯一标识一个DNAT条目。
  * `external_ip` - 公网IP地址。公网IP地址用于ECS实例接收来自互联网的请求。
  * `internal_ip` - 映射到DNAT条目中公网IP地址的私网IP地址。此地址为实际目标服务器的私网IP。
  * `external_port` - 公网端口。公网端口用于ECS实例接收来自互联网的请求。取值范围：1-65535。
  * `internal_port` - 目标私网端口。映射到外部端口的实际私网端口。取值范围：1-65535。
  * `ip_protocol` - 协议类型。支持的协议类型包括TCP、UDP等。
  * `status` - DNAT条目的状态。表示当前DNAT条目的运行状态，例如“active”表示条目已启用。
  * `forward_table_id` - DNAT条目所属DNAT表的ID。此属性用于标识DNAT条目所属的DNAT表。