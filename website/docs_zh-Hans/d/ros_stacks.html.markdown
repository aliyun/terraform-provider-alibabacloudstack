---
subcategory: "ROS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ros_stacks"
sidebar_current: "docs-Alibabacloudstack-datasource-ros-stacks"
description: |- 
  查询资源编排（ROS）资源栈
---

# alibabacloudstack_ros_stacks

根据指定过滤条件列出当前凭证权限可以访问的资源编排（ROS）资源栈列表。

## 示例用法

### 基础用法

以下示例展示了如何使用 `alibabacloudstack_ros_stacks` 数据源来查询 ROS 堆栈列表，并输出第一个堆栈的 ID。

```terraform
data "alibabacloudstack_ros_stacks" "example" {
  ids        = ["example_value"]
  name_regex = "the_resource_name"
}

output "first_ros_stack_id" {
  value = data.alibabacloudstack_ros_stacks.example.stacks.0.id
}
```

### 高级用法，带过滤器

以下示例展示了如何通过多个过滤条件(如堆栈名称、状态、父堆栈 ID 和标签)来筛选 ROS 堆栈列表，并输出所有匹配堆栈的 ID。

```terraform
resource "alibabacloudstack_ros_stack" "default" {
  stack_name       = "tf-testAccRosStacks1723056"
  template_body    = "{\"ROSTemplateFormatVersion\":\"2015-09-01\"}"
  stack_policy_body = "{\"Statement\": [{\"Action\": \"Update:Delete\", \"Resource\": \"*\", \"Effect\": \"Allow\", \"Principal\": \"*\"}]}"
  tags = {
    "Created" = "TF"
    "For"     = "ROS"
  }
}

data "alibabacloudstack_ros_stacks" "default" {
  name_regex      = "${alibabacloudstack_ros_stack.default.stack_name}"
  enable_details  = true
  parent_stack_id = "parent-stack-id"
  show_nested_stack = true
  status          = "CREATE_COMPLETE"
  tags = jsonencode({
    Created = "TF"
    For     = "ROS"
  })
}

output "stack_ids" {
  value = data.alibabacloudstack_ros_stacks.default.ids
}
```

## 参数说明

以下参数是支持的：

* `enable_details` - (可选) 默认值为 `false`。将其设置为 `true` 可以输出更多关于资源属性的详细信息。
* `ids` - (可选，强制更新) 堆栈 ID 列表。
* `name_regex` - (可选，强制更新) 用于通过堆栈名称筛选结果的正则表达式字符串。
* `parent_stack_id` - (可选，强制更新) 父堆栈的 ID。
* `show_nested_stack` - (可选，强制更新) 指定是否在结果中包含嵌套堆栈。有效值：`true`, `false`。
* `stack_name` - (可选，强制更新) 堆栈的名称。名称长度最多为 255 个字符，可以包含数字、字母、连字符 (-) 和下划线 (_)。必须以数字或字母开头。
* `status` - (可选，强制更新) 堆栈的状态。有效值：`CREATE_COMPLETE`, `CREATE_FAILED`, `CREATE_IN_PROGRESS`, `DELETE_COMPLETE`, `DELETE_FAILED`, `DELETE_IN_PROGRESS`, `ROLLBACK_COMPLETE`, `ROLLBACK_FAILED`, `ROLLBACK_IN_PROGRESS`。
* `tags` - (可选) 查询绑定到标签的实例。传入值的格式为 `json` 字符串，包括 `TagKey` 和 `TagValue`。`TagKey` 不可以为空，而 `TagValue` 可以为空。格式示例 `{"key1":"value1"}`。

## 属性说明

除了上述所有参数外，还导出以下属性：

* `names` - 堆栈名称列表。
* `stacks` - ROS 堆栈列表。每个元素包含以下属性：
  * `deletion_protection` - 是否开启资源栈删除保护。有效值：`Enabled`, `Disabled`。
  * `description` - 堆栈的描述。
  * `disable_rollback` - 当创建堆栈失败时，是否禁用回滚。默认值为 `false`。
  * `drift_detection_time` - 堆栈最近一次成功的偏差检测的时间。
  * `id` - 堆栈的 ID。
  * `parent_stack_id` - 父堆栈的 ID。
  * `ram_role_name` - RAM 角色的名称。ROS 假设 RAM 角色来创建堆栈，并使用角色的凭据调用阿里巴巴云服务的 API。RAM 角色名称最大长度为 64 个字符。
  * `root_stack_id` - 最顶层的资源栈的 ID。当资源栈为嵌套资源栈时，会返回该属性。
  * `stack_drift_status` - 堆栈最近一次成功的偏差检测中的堆栈状态。
  * `stack_id` - 堆栈的 ID。
  * `stack_name` - 堆栈的名称。长度不超过 255 个字符，必须以数字或英文字母开头，可包含数字、英文字母、短划线(-)和下划线(_)。
  * `stack_policy_body` - 包含堆栈策略主体的结构。
  * `status` - 资源状态。
  * `status_reason` - 堆栈处于某个状态的原因。
  * `template_description` - 用于创建堆栈的模板的描述。
  * `timeout_in_minutes` - 创建堆栈的超时时间。单位：分钟。
  * `parameters` - 参数列表。
    * `parameter_key` - 参数的键。
    * `parameter_value` - 参数的值。
  * `tags` - 与堆栈关联的标签。