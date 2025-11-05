---
subcategory: "API 网关"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_api_gateway_v2_domain"
sidebar_current: "docs-alibabacloudstack-resource-api-gateway-v2-domain"
description: |-
  提供 AlibabacloudStack API 网关 V2 域名资源。
---

# alibabacloudstack\_api\_gateway\_v2\_domain

提供 API 网关 V2 域名资源。


## 示例用法

基础用法

```hcl
variable "name" {
  default = "tf-example"
}

data "alibabacloudstack_api_gateway_v2_instance_types" "default" {
  sorted_by = "CPU"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name  = var.name
  node_number    = "1"
  instance_class = data.alibabacloudstack_api_gateway_v2_instance_types.default.instance_types.0.id
  broker_engine_type = "SCG"
  deploy_mode    = "custom"
}

resource "alibabacloudstack_api_gateway_v2_certificate" "default" {
  certificate_name = var.name
  instance_id      = alibabacloudstack_api_gateway_v2_instance.default.id
  cert_type        = "0"
  certificates     = "-----BEGIN CERTIFICATE-----\n******\n-----END CERTIFICATE-----"
  private_key      = "-----BEGIN RSA PRIVATE KEY-----\n******\n-----END RSA PRIVATE KEY-----"
}

resource "alibabacloudstack_api_gateway_v2_domain" "default" {
  domain         = "${var.name}.com"
  instance_id    = alibabacloudstack_api_gateway_v2_instance.default.id
  protocol       = "HTTPS"
  certificate_id = alibabacloudstack_api_gateway_v2_certificate.default.certificate_id
  client_auth    = "0"
}
```

## 参数说明

以下参数是可支持的：

* `instance_id` - (必选, ForceNew) API 网关实例的 ID。
* `domain` - (必选, ForceNew) API 网关的自定义域名。
* `protocol` - (必选) 域名使用的协议。有效值：`HTTP`、`HTTPS`。
* `certificate_id` - (可选) 证书的 ID。当协议为 `HTTPS` 时必选。
* `ca_certificate_id` - (可选) 客户端证书认证的 CA 证书 ID, `client_auth`为 `1` 时必选。。
* `client_auth` - (可选) 是否启用客户端证书认证。有效值：`0` (禁用)、`1` (启用)。默认值：`0`。
* `subject_dn` - (可选) 客户端证书的 subject DN。
* `issuer_dn` - (可选) 客户端证书的 issuer DN。

## 属性说明

以下属性会被导出：

* `id` - 域名的 ID。格式为 `<instance_id>:<domain_id>`。
* `domain_id` - 域名的 ID。

## 导入方式

API 网关 V2 域名可以使用 ID 导入，例如：

```shell
$ terraform import alibabacloudstack_api_gateway_v2_domain.example <instance_id>:<domain_id>
```