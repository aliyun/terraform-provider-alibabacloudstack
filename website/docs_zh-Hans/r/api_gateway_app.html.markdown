---
subcategory: "API 网关"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_app"
sidebar_current: "docs-Alibabacloudstack-resource-api-gateway-app"
description: |-
  提供阿里云 API 网关应用资源。
---

# alibabacloudstack_api_gateway_app

提供 API 网关应用资源。该资源用于创建和管理 API 网关的应用，用于 API 的身份认证和授权。

有关 API 网关应用及其使用方法的更多信息，请参阅 [应用管理](https://help.aliyun.com/zh/api-gateway/traditional-api-gateway/developer-reference/api-cloudapi-2016-07-14-dir-applications/)。

-> **注意：** 当使用 `alibabacloudstack_api_gateway_app` 创建资源时，Terraform 将自动创建应用。

## 示例用法

基本用法

```hcl
variable "name" {
  default = "tf_testAccApp_example"
}

variable "description" {
  default = "tf_testAcc api gateway 应用描述"
}

resource "alibabacloudstack_api_gateway_app" "default" {
  name        = var.name
  description = var.description
}
```

## 参数说明

支持以下参数：

* `name` - （必填）API 网关应用的名称。必须以字母或汉字开头，可包含字母、数字和下划线，长度为 4~26 个字符。
* `description` - （可选）应用的描述信息，长度不超过 180 个字符。

## 属性说明

导出以下属性：

* `id` - 应用的 ID，由 API 网关分配的唯一标识符。

## 导入

API 网关应用可以使用应用 ID 导入，例如：

```
$ terraform import alibabacloudstack_api_gateway_app.example 12345678
```
