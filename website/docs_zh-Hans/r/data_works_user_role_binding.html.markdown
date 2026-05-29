---
subcategory: "一站式大数据开发治理平台"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_data_works_user_role_binding"
sidebar_current: "docs-Alibabacloudstack-resource-data-works-user-role-binding"
description: |-
  提供阿里云专有云 DataWorks 用户角色绑定资源。
---

# alibabacloudstack_data_works_user_role_binding

-> **已弃用:** 该资源已被弃用，因为 `alibabacloudstack_data_works_user` 资源已包含相应功能。该资源计划在 3.21.0 版本中移除。

提供 DataWorks 用户角色绑定资源。

有关 DataWorks 用户角色绑定及其使用方法的信息，
请参阅 [什么是 UserRoleBinding](https://help.aliyun.com/apsara/enterprise/v_3_14_0_20210519/dide/enterprise-ascm-developer-guide/AddProjectMemberToRole-1-2.html?spm=a2c4g.14484438.10001.559)。

## 示例

基本用法

```terraform
resource "alibabacloudstack_data_works_user_role_binding" "default" {
  project_id = "10060"
  user_id    = "5225501456060119238"
  role_code  = "role_project_guest"
}
```

## 参数说明

支持以下参数：

* `project_id` - （必选，ForceNew）DataWorks 项目的 ID。修改此参数将强制创建新资源。
* `user_id` - （必选，ForceNew）要绑定角色的用户 ID。修改此参数将强制创建新资源。
* `role_code` - （必选，ForceNew）DataWorks 项目成员的角色代码。有效值：`role_project_owner`、`role_project_admin`、`role_project_dev`、`role_project_pe`、`role_project_deploy`、`role_project_guest`、`role_project_security`。修改此参数将强制创建新资源。

## 属性说明

导出以下属性：

* `id` - 资源绑定的 ID。格式为 `<role_code>:<project_id>:<user_id>`。

## 导入

DataWorks 用户角色绑定可以使用 `<role_code>:<project_id>:<user_id>` 格式的复合 ID 进行导入，例如：

```
$ terraform import alibabacloudstack_data_works_user_role_binding.example role_project_guest:10060:5225501456060119238
```
