---
subcategory: "ApiGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_apigateway_apigroups"
sidebar_current: "docs-Alibabacloudstack-datasource-apigateway-apigroups"
description: |- 
  查询接口网关分组
---

# alibabacloudstack_apigateway_apigroups
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_api_gateway_groups`

根据指定过滤条件列出当前凭证权限可以访问的接口网关分组列表。

## 示例用法

```hcl
variable "name" {
  default = "tf_testAccGroup_6224656"
}

variable "description" {
  default = "tf_testAcc api gateway description"
}

# 创建一个 API 网关分组
resource "alibabacloudstack_api_gateway_group" "default" {
  name        = var.name
  description = var.description
}

# 使用数据源查询 API 网关分组
data "alibabacloudstack_api_gateway_groups" "default" {
  name_regex = alibabacloudstack_api_gateway_group.default.name
  ids        = [alibabacloudstack_api_gateway_group.default.id]
  output_file = "apigroups_output.txt"
}

output "first_group_id" {
  value = data.alibabacloudstack_api_gateway_groups.default.groups.0.id
}
```

## 参数说明

以下参数是支持的：

* `name_regex` - (可选) 用于通过名称过滤 API 网关分组的正则表达式字符串。例如，可以使用 `"example-group"` 来匹配所有以 `"example-group"` 开头的分组名称。
* `ids` - (可选) 用于过滤结果的 API 网关分组 ID 列表。例如，可以使用 `["group1", "group2"]` 来限制返回的分组范围。

## 属性说明

除了上述参数外，还导出以下属性：

* `names` - 匹配过滤条件的所有 API 网关分组的名称列表。
* `groups` - 匹配过滤条件的所有 API 网关分组的详细信息列表。每个元素包含以下属性：
  * `id` - API 网关分组的唯一标识符。
  * `region_id` - API 网关分组所在的区域 ID。
  * `name` - API 网关分组的名称。
  * `sub_domain` - 系统为 API 分组分配的二级域名，用于 API 调用测试。
  * `description` - API 网关分组的描述，不超过 180 个字符，不传递表示不修改。
  * `created_time` - API 网关分组的创建时间(格林尼治标准时间)。
  * `modified_time` - API 网关分组的最后修改时间(格林尼治标准时间)。
  * `traffic_limit` - API 网关分组的最大 QPS 限制，默认值为 500，但可以通过提交申请来增加。
  * `billing_status` - API 网关分组的计费状态。可能的值包括：
    - `NORMAL`: 表示 API 网关分组正常。
    - `LOCKED`: 表示因未支付账单而锁定。
  * `illegal_status` - API 网关分组的违规状态。可能的值包括：
    - `NORMAL`: 表示 API 网关分组正常。
    - `LOCKED`: 表示因违规而锁定。