---
subcategory: "VPNGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_customergateways"
sidebar_current: "docs-Alibabacloudstack-datasource-vpngateway-customergateways"
description: |- 
  查询VPN网关客户网关
---

# alibabacloudstack_vpngateway_customergateways
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_vpn_customer_gateways`

根据指定过滤条件列出当前凭证权限可以访问的VPN网关客户网关列表。

## 示例用法

```hcl
resource "alibabacloudstack_vpn_customer_gateway" "default" {
  name        = "tf-testAccVpnCgwNameDataResource16214"
  ip_address  = "40.104.22.228"
  description = "tf-testAccVpnCgwNameDataResource16214"
}

data "alibabacloudstack_vpngateway_customergateways" "example" {
  ids        = [alibabacloudstack_vpn_customer_gateway.default.id]
  name_regex = "tf-testAccVpnCgwNameDataResource.*"
  output_file = "./customergateways_output.txt"
}

output "customergateway_ids" {
  value = data.alibabacloudstack_vpngateway_customergateways.example.ids
}
```

## 参数说明

以下参数是支持的：

* `ids` - (可选) 客户网关 ID 列表。如果指定，数据源将返回匹配的客户网关。
* `name_regex` - (可选，变更时重建) 用于按名称过滤客户网关的正则表达式字符串。这允许您在客户网关名称中匹配特定模式。

## 属性说明

除了上述参数外，还导出以下属性：

* `names` - 匹配的客户网关名称列表。
* `gateways` - 客户网关对象列表。每个对象包含以下属性：
  * `id` - 客户网关的 ID。
  * `name` - 客户网关的名称。
  * `ip_address` - 客户网关的 IP 地址。
  * `description` - 客户网关的描述信息。
  * `create_time` - 客户网关的创建时间，格式为 ISO8601 字符串(例如，`2023-09-01T12:00:00Z`)。