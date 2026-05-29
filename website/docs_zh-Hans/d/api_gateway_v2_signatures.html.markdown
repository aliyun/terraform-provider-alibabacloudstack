---
subcategory: "API 网关（API Gateway）V2 版"
page_title: "AlibabacloudStack: alibabacloudstack_api_gateway_v2_signatures"
sidebar_current: "docs-Alibabacloudstack-datasource-api-gateway-v2-signatures"
description: |-
  提供API网关V2签名列表。
---

# alibabacloudstack\_api\_gateway\_v2\_signatures

该数据源根据指定的过滤条件提供阿里云账户中的API网关V2签名列表。

## 基础用法

```terraform
data "alibabacloudstack_api_gateway_v2_signatures" "example" {
  gw_instance_id = "api-gateway-v2-instance-id"
}

output "first_signature_id" {
  value = data.alibabacloudstack_api_gateway_v2_signatures.example.signatures.0.id
}
```

## 参数参考

以下参数可用作筛选条件：

* `gw_instance_id` - (必选) API网关V2实例的ID。
* `ids` - (可选) 签名ID列表。
* `name_regex` - (可选) 用于按签名名称过滤结果的正则表达式。
* `names` - 用于按签名名称过滤签名名称列表。

## 属性参考

以下属性会被导出：

* `ids` - 签名ID列表。
* `names` - 签名名称列表。
* `signatures` - 签名列表。每个元素包含以下属性：
  * `id` - 资源ID，格式为 `{gw_instance_id}:{sig_scheme_id}`。
  * `gw_instance_id` - API网关V2实例的ID。
  * `sig_scheme_id` - 签名方案的ID。
  * `sig_scheme_name` - 签名方案的名称。
  * `sig_alg` - 签名算法。
  * `secret_key` - 密钥。
  * `create_time` - 签名创建时间。
  * `update_time` - 签名最后更新时间。
  * `status` - 签名状态。
