---
subcategory: "ROS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ros_templates"
sidebar_current: "docs-Alibabacloudstack-datasource-ros-templates"
description: |- 
  查询资源编排（ROS）模板
---

# alibabacloudstack_ros_templates

根据指定过滤条件列出当前凭证权限可以访问的资源编排（ROS）模板列表。

## 示例用法

### 基础用法：

```terraform
data "alibabacloudstack_ros_templates" "example" {
  ids        = ["example_value"]
  name_regex = "the_resource_name"
}

output "first_ros_template_id" {
  value = data.alibabacloudstack_ros_templates.example.templates.0.id
}
```

高级用法带标签：

```terraform
data "alibabacloudstack_ros_templates" "example_with_tags" {
  template_name = "example_template"
  tags = jsonencode({
    Environment = "Production"
    Owner      = "TeamA"
  })
}

output "ros_template_ids_with_tags" {
  value = data.alibabacloudstack_ros_templates.example_with_tags.ids
}
```

变量定义与使用示例：

```terraform
variable "name" {
  default = "tf-testAlibabacloudstackRosTemplates19947"
}

data "alibabacloudstack_ros_templates" "default" {
  ids = ["${alibabacloudstack_ros_templates.default.id}"]
}
```

## 参数参考

以下参数是支持的：

* `share_type` - (可选，变更时重建) ROS 模板的共享类型。有效值：`Private`(私有模板)，`Shared`(共享模板)。
* `ids` - (可选，变更时重建) 模板 ID 列表，用于筛选特定模板。
* `name_regex` - (可选，变更时重建) 用于按模板名称过滤结果的正则表达式字符串。
* `template_name` - (可选，变更时重建) 模板的名称。长度不超过255个字符，必须以数字或英文字母开头，可包含数字、英文字母、短划线(-)和下划线(_)。
* `enable_details` - (可选) 默认为 `false`。将其设置为 `true` 以输出更多关于资源属性的详细信息。
* `tags` - (可选) 查询绑定到标签的资源。传入值的格式为 `json` 字符串，包括 `TagKey` 和 `TagValue`。`TagKey` 不能为 null，而 `TagValue` 可以为空。格式示例：`{"key1":"value1"}`。

## 属性参考

除了上述参数外，还导出以下属性：

* `names` - 模板名称列表。
* `templates` - ROS 模板列表。每个元素包含以下属性：
  * `change_set_id` - 与模板关联的更改集的 ID。
  * `description` - 模板的描述。最大长度为256个字符。
  * `id` - 模板的 ID。
  * `share_type` - 模板的共享类型(`Private` 或 `Shared`)。
  * `stack_group_name` - 与模板关联的堆栈组的名称。名称在区域中必须唯一，并且最多可以包含255个字符。
  * `stack_id` - 与模板关联的堆栈的 ID。
  * `tags` - 与模板关联的标签。
    * `tag_key` - 资源的第 N 个标签的键。
    * `tag_value` - 资源的第 N 个标签的值。
  * `template_body` - 包含模板主体的结构。模板主体的长度必须为1到524,288字节。如果模板主体的长度过长，建议将参数添加到 HTTP POST 请求主体中，以避免因 URL 过长而导致请求失败。必须指定 `TemplateBody` 和 `TemplateURL` 参数之一，但不能同时指定两者。
  * `template_id` - 模板的 ID。
  * `template_name` - 模板的名称。长度不超过255个字符，必须以数字或英文字母开头，可包含数字、英文字母、短划线(-)和下划线(_)。
  * `template_version` - 模板的版本。