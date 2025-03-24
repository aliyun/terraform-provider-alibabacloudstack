---
subcategory: "CloudFirewall"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudfirewall_controlpolicyorder"
sidebar_current: "docs-Alibabacloudstack-cloudfirewall-controlpolicyorder"
description: |- 
  编排云防火墙控制策略顺序
---

# alibabacloudstack_cloudfirewall_controlpolicyorder
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_cloud_firewall_control_policy_order`

使用Provider配置的凭证在指定的资源集下编排云防火墙控制策略顺序。

## 示例用法

以下是一个完整的示例，展示如何创建一个云防火墙控制策略顺序资源：

```hcl
variable "name" {
  default = "terraform-example"
}

# 创建一个云防火墙访问控制策略
resource "alibabacloudstack_cloud_firewall_control_policy" "default" {
  direction        = "in"
  application_name = "ANY"
  description      = var.name
  acl_action       = "accept"
  source           = "127.0.0.1/32"
  source_type      = "net"
  destination      = "127.0.0.2/32"
  destination_type = "net"
  proto            = "ANY"
}

# 设置访问控制策略的优先级顺序
resource "alibabacloudstack_cloudfirewall_controlpolicyorder" "default" {
  acl_uuid  = alibabacloudstack_cloud_firewall_control_policy.default.acl_uuid
  direction = alibabacloudstack_cloud_firewall_control_policy.default.direction
  order     = 1
}
```

## 参数参考

支持以下参数：

* `acl_uuid` - (必填, 变更时重建) 安全访问控制策略的唯一标识ID。这是在创建策略时分配给该策略的唯一标识符。
* `direction` - (必填, 变更时重建) 安全访问控制策略适用的流量方向。有效值为：
  * `in`：表示入站流量。
  * `out`：表示出站流量。
* `order` - (必填) 安全访问控制策略生效的优先级。优先级数字从1开始顺序递增，优先级数字越小，优先级越高。值为 `-1` 表示优先级最低。**注意：** 从版本1.227.1起，此字段必须设置。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - Terraform中的资源唯一标识符。它格式化为 `<acl_uuid>:<direction>`。
```