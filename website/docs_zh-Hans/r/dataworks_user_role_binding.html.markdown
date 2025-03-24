---
subcategory: "DataWorks"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_data_works_user_role_binding"
sidebar_current: "docs-alibabacloudstack-resource-data-works-user-role-binding"
description: |-
  编排绑定Data Works 用户和角色的关系
---

# alibabacloudstack_data_works_user_role_binding

使用Provider配置的凭证在指定的资源集下编排绑定Data Works用户和角色的关系

关于 Data Works Connection 和如何使用它，
请参阅 [什么是UserRoleBinding](https://help.aliyun.com/apsara/enterprise/v_3_14_0_20210519/dide/enterprise-ascm-developer-guide/AddProjectMemberToRole-1-2.html?spm=a2c4g.14484438.10001.559)。

## 示例用法

### 基础用法

```terraform
resource "alibabacloudstack_data_works_user_role_binding" "default" {
  project_id = "10060"
  user_id = "5225501456060119238"
  role_code = "role_project_guest"
}
```

## 参数参考

支持以下参数：

* `project_id` - (必填) 项目的ID。
* `user_id` - (必填) 阿里云账号ID。
* `role_code` - (必填) DataWorks工作区角色代码。 

## 属性参考

导出以下属性：


* `project_id` - 项目的ID。
* `user_id` - 阿里云账号ID。 

---