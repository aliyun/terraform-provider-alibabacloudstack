---
subcategory: "裸金属算力平台 BMCP"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bcmp_keypair"
sidebar_current: "docs-Alibabacloudstack-resource-bcmp-keypair"
description: |-
  提供 BCMP 密钥对资源。
---

# alibabacloudstack_bcmp_keypair

提供 BCMP 密钥对资源。

## 示例用法

基本用法

```hcl
resource "alibabacloudstack_bcmp_keypair" "basic" {
  key_pair_name = "terraform-test-key-pair"
}

// 导入现有公钥构建密钥对
resource "alibabacloudstack_bcmp_keypair" "publickey" {
  key_pair_name = "my_public_key"
  public_key    = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq*******<Your Server Certificate String>*****bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
}
```

## 参数说明

支持以下参数：

* `key_pair_name` - (选填, ForceNew) 密钥对的名称。值包含 2 到 128 个英文或中文字符。名称必须以字母或中文字符开头。不能以 `http://` 或 `https://` 开头。值可以包含数字、冒号 (`:`)、下划线 (`_`) 或连字符 (`-`)。如果不指定，Terraform 将生成唯一的名称。名称在区域内必须唯一。
* `public_key` - (选填) 您想要导入并使用 AlibabaCloudStack 密钥对管理的现有公钥。首尾空格会被自动裁剪。修改此参数不会触发资源重建，但由于该资源不支持更新公钥，因此变更不会被持久化。
* `key_file` - (选填, ForceNew) 新创建的密钥对的私钥将保存的文件名。强烈建议在创建密钥对时指定此参数，因为如果不保存，您将无法在之后检索私钥。修改此参数会强制创建新资源。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 密钥对的 ID，与 `key_pair_name` 相同。
* `key_pair_name` - 密钥对的名称。
* `finger_print` - 密钥对的指纹。
* `region` - 密钥对的区域。
* `resource_group_name` - 密钥对的资源组名称。
* `create_time` - 密钥对的创建时间。
* `update_time` - 密钥对的更新时间。

> **注意：** `resource_group` 属性在 Schema 中已定义，但当前 API 读取操作未对其赋值。

## Import

BCMP 密钥对可以使用密钥对名称进行导入，例如：

```
$ terraform import alibabacloudstack_bcmp_keypair.example my-key-pair-name
```
