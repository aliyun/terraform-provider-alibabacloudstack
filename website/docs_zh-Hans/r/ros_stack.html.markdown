---
subcategory: "资源编排 ROS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ros_stack"
sidebar_current: "docs-Alibabacloudstack-resource-ros-stack"
description: |- 
  编排资源编排（ROS）资源栈
---

# alibabacloudstack_ros_stack

使用Provider配置的凭证在指定的资源集编排资源编排（ROS）资源栈。

关于ROS Stack及其使用方法的更多信息，请参阅 [什么是Stack](https://www.alibabacloud.com/help/en/doc-detail/132086.htm)。



## 示例用法

以下是一个完整的示例，展示如何创建一个ROS Stack资源：

```terraform
resource "alibabacloudstack_ros_stack" "example" {
  stack_name = "tf-testaccstack"

  template_body = <<EOF
    {
      "ROSTemplateFormatVersion": "2015-09-01",
      "Parameters": {
        "VpcName": {
          "Type": "String"
        },
        "InstanceType": {
          "Type": "String"
        }
      }
    }
    EOF

  stack_policy_body = <<EOF
    {
      "Statement": [
        {
          "Action": "Update:Delete",
          "Resource": "*",
          "Effect": "Allow",
          "Principal": "*"
        }
      ]
    }
    EOF

  tags = {
    Created = "TF"
    For     = "ROS"
  }

  parameters {
    parameter_key   = "VpcName"
    parameter_value = "MyVpc"
  }

  parameters {
    parameter_key   = "InstanceType"
    parameter_value = "ecs.t5-lc2m1.nano"
  }

  timeout_in_minutes = 90

  create_option = "KeepStackOnCreationComplete"
  deletion_protection = "Enabled"
  disable_rollback = false

  notification_urls = ["http://example.com/callback1", "http://example.com/callback2"]

  ram_role_name = "myRAMRole"

  replacement_option = "Disabled"

  retain_all_resources = true

  use_previous_parameters = true
}
```

## 参数说明

支持以下参数：

* `create_option` - (可选，变更时重建，已废弃) **已废弃：** `create_option` 字段在 ApsaraStack 中不支持，将在 3.21.0 版本中移除。指定在堆栈创建后是否删除堆栈。默认值：`KeepStackOnCreationComplete`。有效值：
  * `KeepStackOnCreationComplete`: 在堆栈创建后保留堆栈及其所有资源。
  * `AbandonStackOnCreationComplete`: 在堆栈创建完成后删除堆栈但保留其所有资源。这确保不会达到允许创建的最大堆栈数。如果堆栈创建失败，将保留堆栈。
  * `AbandonStackOnCreationRollbackComplete`: 在堆栈创建失败后的回滚完成时删除堆栈。这确保不会达到允许创建的最大堆栈数。如果堆栈已创建或回滚失败，则保留堆栈。

* `deletion_protection` - (可选，变更时重建) 指定是否为堆栈启用删除保护。有效值：
  * `Enabled`: 为堆栈启用了删除保护。
  * `Disabled`: 删除保护对于堆栈是禁用的。您可以通过 Resource Orchestration Service (ROS) 控制台或调用 DeleteStack 操作来删除堆栈。

  > 嵌套堆栈的删除保护与根堆栈相同。

* `disable_rollback` - (可选) 指定在堆栈创建失败时是否禁用堆栈回滚。默认值：`false`。

* `notification_urls` - (可选，已废弃) **已废弃：** `notification_urls` 字段在 ApsaraStack 中不支持，将在 3.21.0 版本中移除。接收堆栈事件通知的回调 URL 列表。仅支持 HTTP POST。最多可以指定 5 个 URL。

* `parameters` - (可选) 参数列表。每个参数支持以下内容：
  * `parameter_key` - (可选) 参数的键。
  * `parameter_value` - (必填) 参数的值。

* `ram_role_name` - (可选，已废弃) **已废弃：** `ram_role_name` 字段在 ApsaraStack 中不支持，将在 3.21.0 版本中移除。RAM 角色的名称。ROS 假设指定的 RAM 角色来创建堆栈并使用角色的凭据调用 API 操作。RAM 角色名称最大长度为 64 个字符。

* `replacement_option` - (可选) 是否使用替换更新。当资源属性不支持修改更新时，您可以使用替换更新来更改资源属性。替换更新将删除资源并重新创建资源。新资源的物理 ID 将会改变。有效值：
  * `Enabled`: 允许替换更新。
  * `Disabled` (默认): 不允许替换更新。

  > 修改更新的优先级高于替换更新。

* `retain_all_resources` - (可选) 指定在删除期间是否保留堆栈中的所有资源。

* `retain_resources` - (可选) 指定在删除期间是否保留堆栈中的特定资源。

* `stack_name` - (必填，变更时重建) 堆栈的名称。名称长度最多为 255 个字符，可以包含数字、字母、连字符 (-) 和下划线 (_)。必须以数字或字母开头。

* `stack_policy_body` - (可选，已废弃) **已废弃：** `stack_policy_body` 字段在 ApsaraStack 中不支持，将在 3.21.0 版本中移除。包含堆栈策略主体的结构。堆栈策略主体必须为 1 到 16,384 字节长。

* `stack_policy_during_update_body` - (可选) 临时覆盖资源堆栈策略主体的结构。长度为 1~16,384 字节。如果您希望在更新期间更新受保护的资源，请在更新期间指定一个临时覆盖资源堆栈策略。如果没有指定资源堆栈策略，则将使用当前与资源堆栈关联的策略。

* `stack_policy_during_update_url` - (可选) 更新资源堆栈策略的文件位置。URL 必须指向位于 Web 服务器 (HTTP 或 HTTPS) 或 Alibaba Cloud OSS 上的存储空间(例如，`oss:// ros/stack-policy/demo`, `oss:// ros/stack-policy/demo? RegionId = cn-hangzhou`)。策略文件的最大值为 16,384 字节。

* `stack_policy_url` - (可选) 包含资源堆栈策略的文件位置。URL 必须指向位于 Web 服务器 (HTTP 或 HTTPS) 或 Alibaba Cloud OSS 上的存储空间(例如，`oss:// ros/stack-policy/demo`, `oss:// ros/stack-policy/demo? RegionId = cn-hangzhou`)。策略文件的最大长度为 16,384 字节。

* `template_body` - (可选) 模板主体的结构。长度为 1~524,288 字节。如果长度较长，我们建议您通过 HTTP POST + Body Param 在请求正文中传递参数，以避免由于 URL 过长而导致请求失败。

* `template_url` - (可选，已废弃) **已废弃：** `template_url` 字段在 ApsaraStack 中不支持，将在 3.21.0 版本中移除。包含模板主体的文件位置。URL 必须指向位于 Web 服务器 (HTTP 或 HTTPS) 或 Alibaba Cloud OSS 上的存储空间(例如，`oss:// ros/template/demo`, `oss:// ros/template/demo? RegionId = cn-hangzhou`)。模板的最大长度为 524,288 字节。

* `template_version` - (可选，已废弃) **已废弃：** `template_version` 字段在 ApsaraStack 中不支持，将在 3.21.0 版本中移除。模板的版本。

* `timeout_in_minutes` - (可选) 为堆栈创建请求指定的超时期限。默认值：`60`。

* `use_previous_parameters` - (可选) 未传递的参数是否使用上次传递的值。有效值：
  * `true`: 未传递的参数使用上次传递的值。
  * `false`: 未传递的参数不使用上次传递的值。

* `tags` - (可选) 要分配给资源的标签映射。

## 更新行为

以下参数支持在不重建资源的情况下更新：

* `disable_rollback` - 支持更新。指定是否禁用回滚。
* `parameters` - 支持更新。堆栈参数列表。
* `timeout_in_minutes` - 支持更新。堆栈操作的超时期限。
* `replacement_option` - (仅更新) 当资源属性不支持修改时，是否使用替换更新。
* `stack_policy_during_update_body` - (仅更新) 更新期间的临时堆栈策略。
* `stack_policy_during_update_url` - (仅更新) 更新期间临时堆栈策略的 URL。
* `stack_policy_url` - (仅更新) 堆栈策略文件的 URL。
* `template_body` - (仅更新) 用于更新的模板主体结构。
* `template_url` - (仅更新，已废弃) 用于更新的模板文件 URL。
* `template_version` - (仅更新，已废弃) 用于更新的模板版本。
* `use_previous_parameters` - (仅更新) 对于未指定的参数是否使用之前的值。

以下参数仅在删除时生效：

* `retain_all_resources` - (仅删除) 删除堆栈时是否保留所有资源。
* `retain_resources` - (仅删除) 删除堆栈时要保留的特定资源。

## 属性说明

除了上述参数列出的内容外，还导出以下属性：

* `id` - Terraform 中 Stack 的资源 ID。值为 `stack_id`。
* `status` - 堆栈的状态。

## 导入

ROS 资源栈可以使用 StackId 导入，例如：

```
$ terraform import alibabacloudstack_ros_stack.example 5f5e9b2a-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```