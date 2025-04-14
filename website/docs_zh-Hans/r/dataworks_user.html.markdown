---
subcategory: "DataWorks"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_data_works_user"
sidebar_current: "docs-alibabacloudstack-resource-data-works-user"
description: |-
  编排Data Works用户
---

# alibabacloudstack_data_works_user

使用Provider配置的凭证在指定的资源集下编排Data Works用户。

有关 Data Works 用户及其使用方法的信息，
请参阅 [什么是用户](https://help.aliyun.com/apsara/enterprise/v_3_14_0_20210519/dide/enterprise-ascm-developer-guide/CreateProjectMember-1-2.html?spm=a2c4g.14484438.10001.561)。

## 示例用法

### 基础用法

```terraform
resource "alibabacloudstack_data_works_user" "default" {
  user_id = "5225501456060119238"
  project_id = "10060"
}
```

## 参数说明

支持以下参数：

* `project_id` - (必填) 项目的 ID。
* `user_id` - (必填) 要添加的用户 ID。
* `role_code` - (可选) 如果不为空，用户将被添加到此角色。 

## 属性说明

目前该资源没有列出任何属性。如果有需要添加的计算属性，可以在此处包含它们，并附上适当的描述和标记为 ``。