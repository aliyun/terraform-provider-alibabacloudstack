---
subcategory: "EIP"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_eip_addresses"
sidebar_current: "docs-Alibabacloudstack-datasource-eip-addresses"
description: |- 
  查询弹性公网地址
---

# alibabacloudstack_eip_addresses
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_eips`

根据指定过滤条件列出当前凭证权限可以访问的弹性公网地址列表。

## 示例用法

```hcl
# 创建两个弹性公网 IP 资源
resource "alibabacloudstack_eip" "default" {
  name     = "tf-testAccCheckAlibabacloudstackEipsDataSourceConfig7836"
  count    = 2
  bandwidth = 5
}

# 使用数据源查询弹性公网 IP 列表
data "alibabacloudstack_eip_addresses" "example" {
  ids        = [alibabacloudstack_eip.default[0].id]
  ip_addresses = ["192.168.0.1"]

  output_file = "eips_output.txt"
}

output "first_eip_id" {
  value = data.alibabacloudstack_eip_addresses.example.eips.0.id
}
```

## 参数说明

以下参数是支持的：

* `ids` - (可选) 弹性公网 IP 的 ID 列表。如果指定，数据源将仅返回与提供的 ID 匹配的弹性公网 IP。
* `ip_addresses` - (可选) 弹性公网 IP 的公共 IP 地址列表。如果指定，数据源将仅返回与提供的 IP 地址匹配的弹性公网 IP。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - 匹配指定过滤条件的弹性公网 IP 的 ID 列表。
* `names` - 与过滤后的弹性公网 IP 对应的名称列表。
* `eips` - 弹性公网 IP 列表。每个元素包含以下属性：
  * `id` - 弹性公网 IP 的唯一标识符。
  * `status` - 弹性公网 IP 的当前状态。可能的状态值包括：`Associating`（正在绑定）、`Unassociating`（正在解绑）、`InUse`（已使用）和 `Available`（可用）。
  * `ip_address` - 弹性公网 IP 的公共 IP 地址。
  * `bandwidth` - 弹性公网 IP 的最大互联网带宽，单位为 Mbps。
  * `instance_id` - 当前绑定到该弹性公网 IP 的实例的 ID。如果没有绑定实例，则此字段为空。
  * `instance_type` - 绑定到该弹性公网 IP 的实例类型。例如，`ECS` 表示绑定到云服务器实例。