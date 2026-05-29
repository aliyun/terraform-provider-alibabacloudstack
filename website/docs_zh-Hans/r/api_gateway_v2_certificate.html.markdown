---
subcategory: "API 网关（API Gateway）V2 版"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_api_gateway_v2_certificate"
sidebar_current: "docs-Alibabacloudstack-resource-api-gateway-v2-certificate"
description: |-
  提供 AlibabacloudStack API 网关 V2 证书资源。
---

# alibabacloudstack\_api\_gateway\_v2\_certificate

提供 API 网关 V2 证书资源。


## 示例用法

基础用法

```hcl
resource "alibabacloudstack_api_gateway_v2_certificate" "example" {
  cert_type        = "0"
  instance_id      = "instance-id"
  certificates     = "-----BEGIN CERTIFICATE-----\n**********\n-----END CERTIFICATE-----"
  private_key      = "-----BEGIN PRIVATE KEY-----\n**********\n-----END PRIVATE KEY-----"
  certificate_name = "example-cert"
}
```

## 参数说明

* `cert_type` - (必选, 强制新建) 证书类型。可选值: `0` (服务器证书), `1` (CA证书)。
* `instance_id` - (必选, 强制新建) API网关实例的ID。
* `certificates` - (必选) 证书内容。当 `cert_type` 为 `0` 时，此参数表示服务器证书；当 `cert_type` 为 `1` 时，此参数表示CA证书。
* `private_key` - (可选) 私钥。当 `cert_type` 为 `0` 时，此参数为必填项。
* `certificate_name` - (必选) 证书名称。

## 属性说明

* `id` - 资源ID。值为证书ID。
* `certificate_id` - 证书ID。
* `expire_time` - 证书的过期时间。
* `create_time` - 证书的创建时间。
* `update_time` - 证书的最后更新时间。

## 导入方式

API网关V2证书可以使用证书ID进行导入，例如：

```shell
$ terraform import alibabacloudstack_api_gateway_v2_certificate.example <instance_id:certificate_id>
```
