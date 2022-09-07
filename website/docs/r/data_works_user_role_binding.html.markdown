---
subcategory: "Data Works"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_data_works_user_role_binding"
sidebar_current: "docs-alibabacloudstack-resource-data-works-user-role-binding"
description: |- Provides a AlibabacloudStack Data Works UserRoleBinding resource.
---

# alibabacloudstack\_data\_works\_connection

Provides a Data Works UserRoleBinding resource.

For information about Data Works Connection and how to use it,
see [What is UserRoleBinding](https://help.aliyun.com/apsara/enterprise/v_3_14_0_20210519/dide/enterprise-ascm-developer-guide/AddProjectMemberToRole-1-2.html?spm=a2c4g.14484438.10001.559).

## Example Usage

Basic Usage

```terraform
resource "alibabacloudstack_data_works_user_role_binding" "default" {
  project_id = "10060"
  user_id = "5225501456060119238"
  role_code = "role_project_guest"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.
* `user_id` - (Required) Alibaba Cloud Account ID.
* `role_code` - (Required) Code of DataWorks workspace role.
