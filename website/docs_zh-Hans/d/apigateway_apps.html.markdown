---
subcategory: "ApiGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_apigateway_apps"
sidebar_current: "docs-Alibabacloudstack-datasource-apigateway-apps"
description: |- 
  查询接口网关应用
---

# alibabacloudstack_apigateway_apps
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_api_gateway_apps`

根据指定过滤条件列出当前凭证权限可以访问的接口网关应用列表。

## 示例用法

```hcl
variable "name" {
  default = "tf_testAccApp_5207142"
}

variable "description" {
  default = "tf_testAcc api gateway description"
}

resource "alibabacloudstack_api_gateway_app" "default" {
  name        = var.name
  description = var.description
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}

data "alibabacloudstack_apigateway_apps" "default" {
  name_regex = alibabacloudstack_api_gateway_app.default.name
  ids        = [alibabacloudstack_api_gateway_app.default.id]
  output_file = "apps_list.txt"
}

output "first_app_id" {
  value = data.alibabacloudstack_apigateway_apps.default.apps.0.id
}
```

## 参数参考

以下参数是支持的：

* `name_regex` - (可选) 用于按名称筛选应用的正则表达式字符串。通过该参数可以匹配符合条件的应用名称。
* `ids` - (可选) 一个应用ID列表，用于过滤结果。如果提供了此参数，则仅返回与这些ID匹配的应用。
* `tags` - (可选) 标签映射，每个标签由键值对组成。可以通过这些标签筛选应用。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - 应用ID列表。表示所有匹配条件的应用的唯一标识符。
* `names` - 应用名称列表。表示所有匹配条件的应用的名称。
* `apps` - API网关应用列表。每个元素包含以下属性：
  * `id` - 应用ID，由系统生成并全局唯一。
  * `name` - 应用名称。
  * `description` - 应用描述。
  * `created_time` - 创建时间(UTC时间)。表示应用创建的时间戳。
  * `modified_time` - 最后修改时间(UTC时间)。表示应用最后一次被修改的时间戳。
  * `app_code` - 应用简单认证密码。这是应用的唯一标识符，用于简单的身份验证。
