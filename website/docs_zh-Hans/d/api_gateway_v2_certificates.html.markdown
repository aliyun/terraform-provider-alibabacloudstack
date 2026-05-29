---
subcategory: "API 网关（API Gateway）V2 版"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_api_gateway_v2_certificates"
sidebar_current: "docs-alibabacloudstack-datasource-api-gateway-v2-certificates"
description: |-
  提供 AlibabacloudStack API 网关 V2 证书列表。
---

# alibabacloudstack\_api\_gateway\_v2\_certificates

该数据源提供当前阿里云用户拥有的 API 网关 V2 证书列表。


## 示例用法

```hcl
data "alibabacloudstack_api_gateway_v2_certificates" "example" {
  instance_id = "example-instance-id"
}

output "first_certificate_id" {
  value = data.alibabacloudstack_api_gateway_v2_certificates.example.certificates.0.id
}
```

```hcl
data "alibabacloudstack_api_gateway_v2_certificates" "filtered" {
  instance_id  = "example-instance-id"
  cert_type    = "0"
  name_regex   = "^example"
}

output "certificate_names" {
  value = [for cert in data.alibabacloudstack_api_gateway_v2_certificates.filtered.certificates : cert.certificate_name]
}
```

## 参数说明

* `instance_id` - (必选, 强制新建) API网关实例的ID。
* `cert_type` - (可选, 强制新建) 证书类型。可选值: `0` (服务器证书), `1` (CA证书)。
* `name_regex` - (可选, 强制新建) 用于按证书名称过滤结果的正则表达式。
* `sni` - (可选, 强制新建) 证书的SNI (Server Name Indication)。
* `ids` - (可选, 强制新建) 证书ID列表。

## 属性说明

* `id` - 数据源的ID。
* `ids` - 证书ID列表。
* `certificates` - 证书列表。每个元素包含以下属性：
  * `id` - 资源ID。值为 `<instance_id>:<certificate_id>`。
  * `certificate_id` - 证书ID。
  * `certificate_name` - 证书名称。
  * `cert_type` - 证书类型。
  * `expire_time` - 证书的过期时间。
  * `create_time` - 证书的创建时间。
  * `update_time` - 证书的最后更新时间。
  * `snis` - 证书的SNI列表。
```