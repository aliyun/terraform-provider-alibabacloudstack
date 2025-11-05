---
subcategory: "API 网关"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_api_gateway_v2_domains"
sidebar_current: "docs-alibabacloudstack-datasource-api-gateway-v2-domains"
description: |-
  提供 API 网关 V2 域名列表给用户。
---

# alibabacloudstack\_api\_gateway\_v2\_domains

该数据源提供当前阿里云用户拥有的 API 网关 V2 域名列表。


## 示例用法

```hcl
variable "name" {
  default = "tf-example"
}

data "alibabacloudstack_api_gateway_v2_domains" "example" {
  instance_id = "example-instance-id"
  domain      = "example.com"
  protocol    = "HTTPS"
}

output "first_api_gateway_v2_domain_id" {
  value = data.alibabacloudstack_api_gateway_v2_domains.example.domains.0.id
}
```

## 参数说明

以下参数是可支持的：

* `instance_id` - (必选) API 网关实例的 ID。
* `domain` - (可选) API 网关的自定义域名。
* `domain_regex` - (可选) 用于按域名过滤结果的正则表达式字符串。
* `ids` - (可选) 域名 ID 列表。
* `protocol` - (可选) 域名使用的协议。有效值：`HTTP`、`HTTPS`。

## 属性说明

以下属性会被导出：

* `ids` - 域名 ID 列表。
* `domains` - API 网关 V2 域名列表。每个元素包含以下属性：
  * `id` - 域名的 ID。格式为 `<instance_id>:<domain_id>`。
  * `instance_id` - API 网关实例的 ID。
  * `domain` - API 网关的自定义域名。
  * `domain_id` - 域名的 ID。
  * `protocol` - 域名使用的协议。
  * `certificate_id` - 证书的 ID。
  * `ca_certificate_id` - 客户端证书认证的 CA 证书 ID。
  * `client_auth` - 是否启用客户端证书认证。有效值：`0` (禁用)、`1` (启用)。
  * `subject_dn` - 客户端证书的 subject DN。
  * `issuer_dn` - 客户端证书的 issuer DN。