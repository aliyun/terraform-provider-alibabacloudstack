---
subcategory: "OOS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_oos_templates"
sidebar_current: "docs-Alibabacloudstack-datasource-oos-templates"
description: |- 
  查询运维编排（OOS）模板
---

# alibabacloudstack_oos_templates

根据指定过滤条件列出当前凭证权限可以访问的运维编排（OOS）模板列表。

## 示例用法

```hcl
# 创建一个OOS模板资源
resource "alibabacloudstack_oos_template" "default" {
  content = <<EOF
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
  template_name = "tf-testAccOosTemplate-9312712"
  version_name = "test"
  tags = {
    "Created" = "TF"
    "For" = "template Test"
  }
}

# 声明数据源以查询OOS模板
data "alibabacloudstack_oos_templates" "example" {
  name_regex = "${alibabacloudstack_oos_template.default.template_name}"
  tags = {
    "Created" = "TF"
    "For" = "template Test"
  }
  share_type = "Private"
  has_trigger = false
}

# 输出第一个模板的名称
output "first_template_name" {
  value = data.alibabacloudstack_oos_templates.example.templates.0.template_name
}
```

## 参数说明

以下参数是支持的：

* `name_regex` - (可选) 用于通过 `template_name` 过滤结果的正则表达式字符串。
* `category` - (可选) 模板的类别。
* `created_by` - (可选) 模板的创建者。
* `created_date` - (可选) 小于或等于指定时间的模板创建时间。格式：`YYYY-MM-DDThh:mm:ssZ`。
* `created_date_after` - (可选) 大于或等于指定时间的模板创建时间。格式：`YYYY-MM-DDThh:mm:ssZ`。
* `has_trigger` - (可选) 模板是否已成功触发。
* `share_type` - (可选) 模板的共享类型。有效值：`Private`，`Public`。
* `sort_field` - (可选) 用于排序的字段。有效值：`TotalExecutionCount`（总执行次数），`Popularity`（流行度），`TemplateName`（模板名称），`CreatedDate`（创建日期）。默认值：`TotalExecutionCount`。
* `sort_order` - (可选) 排序顺序。有效值：`Ascending`（升序），`Descending`（降序）。默认值：`Descending`。
* `template_format` - (可选) 模板的格式。有效值：`JSON`，`YAML`。
* `template_type` - (可选) OOS 模板的类型。
* `ids` - (可选) OOS 模板 ID 列表(`template_name`)。
* `tags` - (可选) 分配给资源的标签映射。

## 属性说明

除了上述所有参数外，还导出以下属性：

* `ids` - OOS 模板 ID 列表。列表中的每个元素与 `template_name` 相同。
* `names` - (v1.114.0+可用) OOS 模板名称列表。
* `templates` - OOS 模板列表。每个元素包含以下属性：
  * `id` - OOS 模板的 ID。与 `template_name` 相同。
  * `template_name` - OOS 模板的名称。
  * `description` - OOS 模板的描述。
  * `template_id` - OOS 模板资源的 ID。
  * `template_version` - OOS 模板的版本。
  * `updated_by` - 最后更新模板的用户。
  * `updated_date` - 模板最后更新的时间。
  * `category` - 模板的类别。
  * `created_by` - 模板的创建者。
  * `created_date` - 模板的创建时间。
  * `has_trigger` - 模板是否已成功触发。
  * `share_type` - 模板的共享类型。有效值：`Private`（私有），`Public`（公共）。
  * `tags` - 分配给资源的标签映射。
  * `template_format` - 模板的格式。有效值：`JSON`，`YAML`。
  * `template_type` - OOS 模板的类型。