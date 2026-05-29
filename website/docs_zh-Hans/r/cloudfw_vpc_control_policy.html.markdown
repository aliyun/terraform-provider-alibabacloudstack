---
subcategory: "Cloud Firewall"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudfw_vpc_control_policy"
sidebar_current: "docs-Alibabacloudstack-cloudfw-vpc-control-policy"
description: |-
  管理云防火墙VPC控制策略。
---

# alibabacloudstack_cloudfw_vpc_control_policy

管理云防火墙VPC控制策略，用于定义VPC网络流量的访问控制规则。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tf-testacc_vpc_control_policy92894"
}


resource "alibabacloudstack_cloudfw_address_book" "default" {
  group_type   = "ip"
  group_name   = var.name
  address_list = ["100.100.100.100/30"]
  description  = "test address book"
}



resource "alibabacloudstack_cloudfw_vpc_control_policy" "default" {
  dest_port        = "33/33"
  acl_action       = "log"
  release          = true
  destination      = "0.0.0.0/16"
  proto            = "UDP"
  application_id   = "0"
  application_name = "ANY"
  source_type      = "net"
  destination_type = "net"
  source           = "0.0.0.0/16"
  description      = var.name
  dest_port_type   = "port"
}
```

## 参数说明

支持以下参数：

* `acl_action` - (必填) 访问控制策略的动作。取值：`accept`（允许）、`drop`（拒绝）、`log`（记录日志）。
* `application_id` - (必填) 应用ID。例如，`0`表示ANY（所有应用）。
* `application_name` - (必填) 应用名称。例如，`ANY`表示所有应用。
* `description` - (必填) 策略描述信息，用于标识策略用途。
* `destination` - (必填) 目的地址。支持CIDR格式（如`0.0.0.0/0`）或地址簿名称。
* `destination_type` - (必填) 目的地址类型。取值：`net`（网段）、`group`（地址簿）。
* `proto` - (必填) 协议类型。取值：`TCP`、`UDP`、`ANY`。
* `source` - (必填) 源地址。支持CIDR格式（如`192.168.1.0/24`）或地址簿名称。
* `source_type` - (必填) 源地址类型。取值：`net`（网段）、`group`（地址簿）。
* `dest_port` - (可选) 目的端口范围。格式：`起始端口/结束端口`（如`80/80`）。
* `dest_port_type` - (可选) 目的端口类型。取值：`port`（端口）、`group`（端口簿）。
* `new_order` - (可选) 新的策略顺序。默认值：`-1`（表示添加到策略列表末尾）。
* `release` - (可选) 是否发布策略。取值：`true`（发布）、`false`（不发布）。
* `vpc_firewall_id` - (可选) VPC防火墙实例ID。当未指定时，使用默认防火墙实例。
* `direction` - (可选，可回读) 策略方向。取值：`inout`（双向流量）（默认）、`in`（入站流量）、`out`（出站流量）。

## 属性说明

以下属性导出为资源属性：

* `id` - 策略ID，格式为 `<acl_uuid>:<direction>`。
* `acl_uuid` - 策略唯一标识符（AclUuid）。
* `dest_port_group` - 目的端口组ID（由API返回）。
* `dest_port_group_ports` - 目的端口组端口列表（当`dest_port_type`为`group`时返回）。
* `destination_group_cidrs` - 目的地址组CIDR列表（当`destination_type`为`group`时返回）。
* `hit_times` - 策略命中次数（自创建以来匹配的流量次数）。
* `order` - 策略当前顺序值（用于策略优先级排序）。
* `source_group_cidrs` - 源地址组CIDR列表（当`source_type`为`group`时返回）。

## Import

云防火墙VPC控制策略可以使用ID（格式：`<acl_uuid>:<direction>`）导入，例如：

```
$ terraform import alibabacloudstack_cloudfw_vpc_control_policy.example acl-12345678:inout
```