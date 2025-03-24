---
subcategory: "NAS"
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

## 参数参考

支持以下参数：

* `access_group_name` - (必填，变更时重建) 权限组的名称。一旦设置，无法修改。
* `file_system_type` - (可选，变更时重建) 文件系统类型。有效值：`standard` 和 `extreme`。默认为 `standard`。注意，`extreme` 类型仅支持 `Vpc` 网络。
* `access_group_type` - (必填，变更时重建) 权限组类型。有效值：`Vpc`。不支持新增经典网络类型(Classic)权限组。
* `description` - (可选) 权限组描述信息。这提供了关于权限组的额外详细信息，以便更好地识别。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - 访问组的 ID。
* `access_group_name` - 权限组的名称。
* `access_group_type` - 权限组类型，包括 `Vpc`。不支持新增经典网络类型(Classic)权限组。
* `file_system_type` - 文件系统类型。有效值：`standard` 和 `extreme`。

## 导入

可以使用 ID 导入 NAS 访问组，例如：

```bash
$ terraform import alibabacloudstack_nas_accessgroup.foo tf_testAccNasConfig:standard
```