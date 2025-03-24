---
subcategory: "OOS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_oos_template"
sidebar_current: "docs-Alibabacloudstack-oos-template"
description: |- 
  编排运维编排（OOS）模板
---

# alibabacloudstack_oos_template

使用Provider配置的凭证在指定的资源集编排运维编排（OOS）模板。

## 示例用法

```terraform
variable "name" {
    default = "tf-testaccoostemplate93918"
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
  template_name = var.name
  version_name  = "v1.0"
  tags = {
    "Created" = "TF",
    "For"     = "acceptance Test"
  }
}
```

## 参数参考

支持以下参数：

* `content` - (必填) 模板的内容。模板必须为 JSON 或 YAML 格式，最大大小为 64 KB。此字段定义了模板的结构和逻辑，包括参数、任务和其他配置。
* `auto_delete_executions` - (可选) 删除模板时是否删除其相关执行。默认值为 `false`。
* `template_name` - (必填，变更时重建) 模板名称。模板名称最多可以包含 200 个字符，名称可以包含字母、数字、连字符(-)和下划线(_)。不能以 `ALIYUN`、`ACS`、`ALIBABA` 或 `ALICLOUD` 开头。
* `version_name` - (可选) 模板版本名称。这允许您管理同一模板的不同版本。
* `tags` - (可选) 要分配给资源的标签映射。标签有助于组织和分类您的资源。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - 资源的唯一标识符。它与 `template_name` 相同。
* `created_by` - 模板的创建者。指示谁创建了模板。
* `created_date` - 模板创建的时间。这对跟踪模板的生命周期很有帮助。
* `description` - 模板的描述。提供了模板功能的简要概述。
* `has_trigger` - 指示模板是否已成功触发。此属性有助于了解是否已启动模板执行。
* `share_type` - 模板的共享类型。用户创建的模板设置为 `Private`，OOS 提供的公共模板设置为 `Public`。
* `template_format` - 模板格式。系统自动识别模板是 JSON 还是 YAML 格式。
* `template_id` - OOS 模板的唯一标识符。在其他 API 调用中引用模板时很有用。
* `template_type` - OOS 模板的类型。`Automation` 表示实现阿里巴巴云 API 模板，而 `Package` 表示用于安装软件的模板。
* `template_version` - OOS 模板的版本。有助于管理同一模板的不同迭代。
* `updated_by` - 最后更新模板的用户。这对审计目的很有用。
* `updated_date` - 模板最后一次更新的时间。这有助于跟踪随时间对模板所做的更改。