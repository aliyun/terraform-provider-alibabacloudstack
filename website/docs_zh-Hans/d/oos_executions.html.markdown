---
subcategory: "OOS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_oos_executions"
sidebar_current: "docs-Alibabacloudstack-datasource-oos-executions"
description: |- 
  查询运维编排（OOS）任务执行
---

# alibabacloudstack_oos_executions

根据指定过滤条件列出当前凭证权限可以访问的运维编排（OOS）任务执行列表。

## 示例用法

以下是一个完整的示例，展示了如何使用 `alibabacloudstack_oos_executions` 数据源来获取 OOS 执行列表，并输出第一个执行的 ID：

```hcl
# 创建一个 OOS 模板
resource "alibabacloudstack_oos_template" "default" {
  content = <<EOF
  {
    "FormatVersion": "OOS-2019-06-01",
    "Description": "Describe instances of given status",
    "Parameters": {
      "Status": {
        "Type": "String",
        "Description": "(Required) The status of the Ecs instance."
      }
    },
    "Tasks": [
      {
        "Properties": {
          "Parameters": {
            "Status": "{{ Status }}"
          },
          "API": "DescribeInstances",
          "Service": "Ecs"
        },
        "Name": "foo",
        "Action": "ACS::ExecuteApi"
      }
    ]
  }
  EOF
  template_name = "tf-testAccOosTemplate-5359742"
  version_name = "test"
  tags = {
    "Created" = "TF",
    "For" = "template Test"
  }
}

# 创建一个 OOS 执行
resource "alibabacloudstack_oos_execution" "default" {
  template_name = alibabacloudstack_oos_template.default.template_name
  description  = "From TF Test"
  parameters   = <<EOF
    {"Status":"Running"}
  EOF
}

# 声明数据源以获取 OOS 执行列表
data "alibabacloudstack_oos_executions" "default" {
  status       = "Success"
  ids          = [alibabacloudstack_oos_execution.default.id]
}

# 输出第一个执行的 ID
output "first_execution_id" {
  value = data.alibabacloudstack_oos_executions.default.executions.0.id
}
```

## 参数参考

以下参数是支持的：

* `category` - (可选) 模板的类别。有效值：`AlarmTrigger`、`EventTrigger`、`Other` 和 `TimerTrigger`。
* `end_date` - (可选) 执行结束的时间。
* `end_date_after` - (可选) 结束时间小于或等于指定时间的执行。
* `executed_by` - (可选) 执行模板的用户。
* `ids` - (可选) OOS 执行 ID 列表。
* `include_child_execution` - (可选) 是否包含子执行。
* `mode` - (可选) OOS 执行的模式。有效值：`Automatic`、`Debug`。
* `parent_execution_id` - (可选) 父 OOS 执行的 ID。
* `ram_role` - (可选) 执行当前模板的角色。
* `sort_field` - (可选) 排序字段。
* `sort_order` - (可选) 排序顺序。
* `start_date_after` - (可选) 开始时间大于或等于指定时间的执行。
* `start_date_before` - (可选) 开始时间小于或等于指定时间的执行。
* `status` - (可选) OOS 执行的状态。有效值：`Cancelled`、`Failed`、`Queued`、`Running`、`Started`、`Success`、`Waiting`。
* `template_name` - (可选) 执行模板的名称。

> **注意**：所有标记为“强制新建”的参数表示如果这些参数被修改，则会触发新资源的创建。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - OOS 执行 ID 列表。
* `executions` - OOS 执行列表。每个元素包含以下属性：
  * `id` - OOS 执行的 ID。
  * `parent_execution_id` - 父 OOS 执行的 ID。
  * `category` - 模板的类别。有效值：`AlarmTrigger`、`EventTrigger`、`Other` 和 `TimerTrigger`。
  * `counters` - OOS 执行的计数器。
  * `create_date` - 执行创建的时间。
  * `end_date` - 执行结束的时间。
  * `executed_by` - 执行模板的用户。
  * `execution_id` - OOS 执行的 ID。
  * `is_parent` - 是否包含子任务。
  * `outputs` - OOS 执行的输出。
  * `parameters` - 模板所需的参数。
  * `mode` - OOS 执行的模式。有效值：`Automatic`、`Debug`。
  * `ram_role` - 执行当前模板的角色。
  * `start_date` - 模板启动的时间。
  * `status_message` - 状态消息。
  * `status_reason` - 状态原因。
  * `template_id` - 执行模板的 ID。
  * `template_name` - 执行模板的名称。
  * `template_version` - 执行模板的版本。
  * `update_date` - 模板更新的时间。
  * `status` - OOS 执行的状态。有效值：`Cancelled`、`Failed`、`Queued`、`Running`、`Started`、`Success`、`Waiting`。