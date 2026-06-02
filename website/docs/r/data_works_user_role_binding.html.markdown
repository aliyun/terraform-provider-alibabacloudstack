---
subcategory: "One-stop Big Data Development and Governance Platform"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_data_works_user_role_binding"
description: |-
  Provides a AlibabacloudStack Data Works UserRoleBinding resource.
---

# alibabacloudstack_data_works_user_role_binding

-> **Deprecated:** This resource is deprecated because `alibabacloudstack_data_works_user` already includes corresponding functions. It is scheduled for removal in version 3.21.0.

Provides a DataWorks User Role Binding resource.

For information about DataWorks User Role Binding and how to use it,
see [What is UserRoleBinding](https://help.aliyun.com/apsara/enterprise/v_3_14_0_20210519/dide/enterprise-ascm-developer-guide/AddProjectMemberToRole-1-2.html?spm=a2c4g.14484438.10001.559).

## Example Usage

Basic Usage

```terraform
resource "alibabacloudstack_data_works_user_role_binding" "default" {
  project_id = "10060"
  user_id    = "5225501456060119238"
  role_code  = "role_project_guest"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, ForceNew) The ID of the DataWorks project. Changing this parameter will force a new resource to be created.
* `user_id` - (Required, ForceNew) The ID of the user to bind the role to. Changing this parameter will force a new resource to be created.
* `role_code` - (Required, ForceNew) The role code for the DataWorks project member. Valid values: `role_project_owner`, `role_project_admin`, `role_project_dev`, `role_project_pe`, `role_project_deploy`, `role_project_guest`, `role_project_security`. Changing this parameter will force a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource binding. The format is `<role_code>:<project_id>:<user_id>`.

## Import

DataWorks User Role Binding can be imported using the composite ID in the format `<role_code>:<project_id>:<user_id>`, e.g.

```
$ terraform import alibabacloudstack_data_works_user_role_binding.example role_project_guest:10060:5225501456060119238
```
