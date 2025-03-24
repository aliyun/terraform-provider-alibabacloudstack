---
subcategory: "DataHub"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_datahub_subscription"
sidebar_current: "docs-Alibabacloudstack-datahub-subscription"
description: |- 
  编排datahub订阅
---

# alibabacloudstack_datahub_subscription

使用Provider配置的凭证在指定的资源集下编排datahub订阅。

## 示例用法

```hcl
variable "name" {
    default = "tf_testacc_datahub_sub71248"
}

resource "alibabacloudstack_datahub_project" "default" {
    comment = "test"
    name = var.name
}

resource "alibabacloudstack_datahub_topic" "default" {
  name = var.name
  comment = "test"
  record_type = "BLOB"
  project_name = alibabacloudstack_datahub_project.default.name
}

resource "alibabacloudstack_datahub_subscription" "default" {
  comment = "Subscription created by Terraform"
  project_name = alibabacloudstack_datahub_project.default.name
  topic_name = alibabacloudstack_datahub_topic.default.name
}
```

## 参数参考

支持以下参数：
  * `project_name` - (必填, 变更时重建) 订阅所属的DataHub项目的名称。其长度限制为3-32个字符，仅允许字母、数字和下划线(`_`)，不区分大小写。
  * `topic_name` - (必填, 变更时重建) 订阅所属的DataHub主题的名称。其长度限制为1-128个字符，仅允许字母、数字和下划线(`_`)，不区分大小写。
  * `comment` - (选填, 变更时重建) DataHub订阅的注释。最大长度为255个字符。

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `id` - 作为Terraform资源的DataHub订阅的唯一标识符。它以 `<project_name>:<topic_name>:<sub_id>` 的格式组成。
  * `sub_id` - 订阅的身份，由服务器端生成。
  * `create_time` - DataHub订阅的创建时间。这是一个人类可读的字符串，而不是64位UTC时间戳。
  * `last_modify_time` - DataHub订阅的最后修改时间。最初，它与 `create_time` 相同。像 `create_time` 一样，它也是一个人类可读的字符串，而不是64位UTC时间戳。
