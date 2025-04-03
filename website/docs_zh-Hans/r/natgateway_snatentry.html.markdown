---
subcategory: "NATGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_natgateway_snatentry"
sidebar_current: "docs-Alibabacloudstack-natgateway-snatentry"
description: |- 
  编排专有网络的NAT网关SNAT表规则
---

# alibabacloudstack_natgateway_snatentry
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_snat_entry`

使用Provider配置的凭证在指定的资源集编排专有网络的NAT网关SNAT表规则。

## 示例用法

```hcl
variable "name" {
    default = "tf-testaccnat_gatewaysnat_entry66949"
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
    vpc_id        = "${alibabacloudstack_vpc.default.id}"
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

resource "alibabacloudstack_natgateway_snatentry" "default" {
    snat_table_id     = "${alibabacloudstack_nat_gateway.default.snat_table_ids}"
    source_vswitch_id = "${alibabacloudstack_vswitch.default.id}"
    snat_ip           = "${alibabacloudstack_eip.default.ip_address}"
    source_cidr       = "${alibabacloudstack_vswitch.default.cidr_block}"
}
```

## 参数参考

支持以下参数：

* `snat_table_id` - (必填，变更时重建) SNAT条目所属的SNAT表ID。此值可以从`alibabacloudstack_nat_gateway`资源的`snat_table_ids`属性中获取。
* `source_vswitch_id` - (选填，变更时重建) 与SNAT条目关联的交换机ID。此参数与`source_cidr`参数互斥。
* `snat_ip` - (必填) 用于SNAT条目的公网IP地址。此IP必须属于与NAT网关关联的弹性IP (EIP)。
* `source_cidr` - (选填，变更时重建) ECS实例的私有网络段。此参数与`source_vswitch_id`参数互斥。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - SNAT条目的ID。格式为`<snat_table_id>:<snat_entry_id>`。
* `snat_entry_id` - 服务器上的SNAT条目的唯一标识符。