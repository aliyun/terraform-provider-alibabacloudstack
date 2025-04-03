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

## 参数参考

以下参数是支持的：

* `forward_table_id` - (必填, 变更时重建) - DNAT条目所属DNAT表的ID。
* `name_regex` - (选填) - 用于通过Forward Entry的名称筛选结果的正则表达式字符串。
* `external_ip` - (选填) - 公网IP地址。公网IP地址用于ECS实例接收来自互联网的请求。
* `internal_ip` - (选填) - 映射到DNAT条目中公网IP地址的私网IP地址。
* `ids` - (选填) - Forward Entry ID列表。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - Forward Entry ID列表。
* `names` - Forward Entry名称列表。
* `entries` - Forward Entries列表。每个元素包含以下属性：
  * `id` - Forward Entry的ID。
  * `external_ip` - 公网IP地址。公网IP地址用于ECS实例接收来自互联网的请求。
  * `internal_ip` - 映射到DNAT条目中公网IP地址的私网IP地址。
  * `external_port` - 公网端口。公网端口用于ECS实例接收来自互联网的请求。
  * `internal_port` - 目标私网端口。取值范围：1-65535。
  * `ip_protocol` - 协议类型。
  * `status` - DNAT条目的状态。