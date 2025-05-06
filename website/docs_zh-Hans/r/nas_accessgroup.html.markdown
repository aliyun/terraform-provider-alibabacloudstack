---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_accessgroup"
sidebar_current: "docs-Alibabacloudstack-nas-accessgroup"
description: |- 
  编排文件存储（NAS）访问权限组
---

# alibabacloudstack_nas_accessgroup
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_nas_access_group`

使用Provider配置的凭证在指定的资源集编排文件存储（NAS）访问权限组。

## 示例用法

### 基础用法

```terraform
resource "alibabacloudstack_nas_accessgroup" "foo" {
  access_group_name = "CreateAccessGroup"
  access_group_type = "Vpc"
  description       = "test_AccessG"
  file_system_type = "extreme"
}
```

包含所有参数的示例

```terraform
variable "name" {
  default = "tf-testaccnasaccess_group31001"
}

resource "alibabacloudstack_nas_accessgroup" "default" {
  access_group_name = "accssGroupExtremeVpcTest"
  file_system_type  = "extreme"
  access_group_type = "Vpc"
  description       = "test"
}
```

## 参数说明

支持以下参数：

* `access_group_name` - (必填，变更时重建) 权限组的名称。一旦设置，无法修改。
* `file_system_type` - (可选，变更时重建) 文件系统类型。有效值：`standard` 和 `extreme`。默认为 `standard`。注意，`extreme` 类型仅支持 `Vpc` 网络。
* `access_group_type` - (必填，变更时重建) 权限组类型。有效值：`Vpc` 和 `Classic`。
* `description` - (可选) 权限组描述信息。这提供了关于权限组的额外详细信息，以便更好地识别。
* `name` - (可选，变更时重建) 已废弃字段，用于权限组的名称。将在未来的版本中移除，请改用 `access_group_name`。
* `type` - (可选，变更时重建) 已废弃字段，用于权限组的类型。将在未来的版本中移除，请改用 `access_group_type`。

### 已废弃参数

以下参数已被废弃并在版本 1.92.0 中替换：

* `name` - (已废弃) 在版本 1.92.0 后被 `access_group_name` 替代。
* `type` - (已废弃) 在版本 1.92.0 后被 `access_group_type` 替代。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 访问组的 ID。从版本 1.92.0 开始，格式化为 `<access_group_name>:<file_system_type>`。
* `access_group_name` - 权限组的名称。
* `access_group_type` - 权限组类型，包括 `Vpc` 和 `Classic`。
* `file_system_type` - 文件系统类型。有效值：`standard` 和 `extreme`。
* `name` - 已废弃字段，用于权限组的名称。将在未来的版本中移除，请改用 `access_group_name`。
* `type` - 已废弃字段，用于权限组的类型。将在未来的版本中移除，请改用 `access_group_type`。

## 导入

可以使用 ID 导入 NAS 访问组，例如：

```bash
$ terraform import alibabacloudstack_nas_accessgroup.foo tf_testAccNasConfig:standard
```