---
subcategory: "Data Works"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_data_works_user"
sidebar_current: "docs-alibabacloudstack-resource-data-works-user"
description: |- 
  Provides a AlibabacloudStack Data Works User resource.
---

# alibabacloudstack\_data\_works\_user

Provides a Data Works User resource.

For information about Data Works User and how to use it,
see [What is User](https://help.aliyun.com/apsara/enterprise/v_3_14_0_20210519/dide/enterprise-ascm-developer-guide/CreateProjectMember-1-2.html?spm=a2c4g.14484438.10001.561).

## Example Usage

Basic Usage

```terraform
resource "alibabacloudstack_data_works_user" "default" {
  user_id = "5225501456060119238"
  project_id = "10060"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.
* `user_id` - (Required) User ID to be added.
* `role_code` - (Optional) If it is not blank, the user will be added to this role.

