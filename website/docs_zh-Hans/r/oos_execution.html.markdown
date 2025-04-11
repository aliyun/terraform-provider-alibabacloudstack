---
subcategory: "OOS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_oos_execution"
sidebar_current: "docs-Alibabacloudstack-oos-execution"
description: |- 
  编排运维编排（OOS）任务执行
---

# alibabacloudstack_oos_execution

使用Provider配置的凭证在指定的资源集编排运维编排（OOS）任务执行。

## 示例用法

```terraform
variable "name" {
    default = "tf-testaccoosexecution34941"
}

resource "alibabacloudstack_oos_template" "default" {
  content       = <<EOF
  {
    "FormatVersion": "OOS-2019-06-01",
    "Description": "Update Describe instances of given status",
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
  template_name = "test-name"
  version_name  = "test"
  tags = {
    "Created" = "TF"
    "For"     = "acceptance Test"
  }
}

resource "alibabacloudstack_oos_execution" "default" {
  template_name = alibabacloudstack_oos_template.default.template_name
  description   = "From TF Test"
  parameters    = "{\"Status\":\"Running\"}"
  mode          = "Automatic"
  template_version = "test"
}
```

## 参数说明

支持以下参数：

* `template_name` - (必填, 变更时重建) 执行模板的名称。这是要执行的OOS模板的标识符。
* `description` - (选填, 变更时重建) OOS Execution的简要描述。这有助于识别执行的目的或上下文。
* `loop_mode` - (选填, 变更时重建) 指定执行的循环模式。这决定了模板中的任务是并行还是顺序执行。
* `mode` - (选填, 变更时重建) 指定执行模式。有效值包括：
  * `Automatic`: 自动执行所有任务，无需手动干预。
  * `Debug`: 以调试模式执行任务，允许逐步执行。
  
  默认值为 `Automatic`。
* `parameters` - (选填, 变更时重建) JSON格式的字符串，包含OOS模板所需的参数。这些参数用于自定义模板在执行期间的行为。默认值为 `{}`。
* `parent_execution_id` - (选填, 变更时重建) 如果此执行是较大工作流或子任务的一部分，则为父执行的ID。
* `safety_check` - (选填, 变更时重建) 指定执行的安全检查模式。这确保在继续执行之前满足某些条件。
* `template_version` - (选填, 变更时重建) 正在执行的OOS模板的版本。如果不指定，将使用模板的最新版本。
* `template_content` - (选填, 变更时重建) OOS模板的原始内容。当从自定义模板而不是现有模板创建执行时，这非常有用。
* `ram_role` - (选填, 变更时重建) 指定与执行关联的RAM角色，该角色授予任务执行所需的权限。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - OOS Execution的唯一标识符。
* `counters` - 与执行相关的计数器摘要，例如已执行、成功或失败的任务数量。
* `create_date` - 表示执行创建时间的时间戳。
* `end_date` - 表示执行完成时间的时间戳。
* `executed_by` - 发起执行的用户或系统。
* `is_parent` - 指示执行是否包含子任务或子执行。
* `outputs` - JSON格式的字符串，包含执行生成的输出。这些输出可用于进一步处理或报告。
* `ram_role` - 与执行关联的RAM角色，该角色授予任务执行所需的权限。
* `start_date` - 表示执行开始时间的时间戳。
* `status` - 执行的当前状态。可能的值包括 `Pending`, `Running`, `Success`, `Failed` 等。
* `status_message` - 描述执行当前状态的详细消息。
* `template_id` - 用于执行的OOS模板的唯一标识符。
* `template_version` - 用于执行的OOS模板的具体版本号。
* `update_date` - 表示执行最后一次更新时间的时间戳。
* `safety_check` - 指定执行的安全检查模式。这确保在继续执行之前满足某些条件。
