---
subcategory: "API 网关（API Gateway）V2 版"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_signature"
sidebar_current: "docs-Alibabacloudstack-api-gateway-v2-signature"
description: |-
  api网关 v2 版本 签名管理
---

# alibabacloudstack_api_gateway_v2_signature

管理API网关v2版本的签名方案，用于配置API请求的签名验证机制。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tf-testacc-sign79175"
}

variable "signature_algorithm" {
  default = "HmacSM3"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name      = var.name
  node_number        = "1"
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

resource "alibabacloudstack_api_gateway_v2_signature" "default" {
  sig_scheme_name = var.name
  sig_alg         = var.signature_algorithm
  gw_instance_id  = alibabacloudstack_api_gateway_v2_instance.default.id
}

```

## 参数说明

支持以下参数：

* `gw_instance_id` - (必填, 变更时重建) API网关实例ID。用于指定签名方案所属的网关实例。
* `sig_alg` - (必填, 变更时重建) 签名算法。可选值：`HmacSHA256`、`HmacSHA1`、`HmacSM3`。
* `sig_scheme_name` - (必填) 签名方案名称。用于标识签名方案的唯一名称，长度限制由API网关服务决定。
* `status` - (可选) 签名方案状态。可选值：`0`（停用）、`1`（启用）。默认值为空，表示保持当前状态。
* `cascade_link_ids` - (可选, 变更时重建) 级联链路ID列表，用于CSB跨域认证。设置此属性将创建源签名方案。

## 属性说明

以下属性导出为资源属性：

* `id` - 资源ID，格式为 `{idpre}:{gwInstanceId}:{sigSchemeId}`。普通签名的 `idpre` 为 `sig`，级联源签名的 `idpre` 为 `sourceSig`。
* `create_time` - 签名方案创建时间，格式为`YYYY-MM-DD HH:mm:ss`。
* `secret_key` - 签名方案的密钥，用于生成和验证签名。
* `sig_scheme_id` - 签名方案唯一标识ID。
* `update_time` - 签名方案最后更新时间，格式为`YYYY-MM-DD HH:mm:ss`。

## Import

API 网关 V2 签名可以使用 `{idpre}:{gwInstanceId}:{sigSchemeId}` 格式的资源 ID 进行导入，例如：

```
$ terraform import alibabacloudstack_api_gateway_v2_signature.example sig:gw-12345678:signature123
```

对于级联源签名，`idpre` 为 `sourceSig`：

```
$ terraform import alibabacloudstack_api_gateway_v2_signature.example sourceSig:gw-12345678:signature456
```