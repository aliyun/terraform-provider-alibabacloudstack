---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_keypair"
sidebar_current: "docs-Alibabacloudstack-ecs-keypair"
description: |- 
  编排云服务器（Ecs）密钥对
---

# alibabacloudstack_ecs_keypair
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_key_pair`

使用Provider配置的凭证在指定的资源集下编排云服务器（Ecs）密钥对。

## 示例用法

### 基础用法

```hcl
resource "alibabacloudstack_ecs_keypair" "basic" {
  key_pair_name = "terraform-test-key-pair"
}

// 使用名称前缀来创建密钥对
resource "alibabacloudstack_ecs_keypair" "prefix" {
  key_name_prefix = "terraform-test-key-pair-prefix"
}

// 导入现有的公钥以创建密钥对
resource "alibabacloudstack_ecs_keypair" "publickey" {
  key_pair_name = "my_public_key"
  public_key    = "ssh-rsa AB3Napapsod45678qwertyuudsfsg"
}
```

## 参数参考

支持以下参数：

* `key_pair_name` - (必填，变更时重建) 密钥对的名称。长度为2~128个英文或中文字符。必须以字母或中文开头，不能以`http://`或`https://`开头。可以包含数字、半角冒号(`:`)、下划线(`_`)或者短划线(`-`)。
* `key_name_prefix` - (可选，变更时重建) 密钥对名称的前缀。它与`key_pair_name`冲突。如果指定了此参数，Terraform将使用它生成一个唯一的密钥对名称。
* `public_key` - (可选，变更时重建) 您要导入并使用阿里云密钥对管理的现有公钥。
* `key_file` - (可选，变更时重建) 保存新创建密钥对私钥的文件名。强烈建议在创建密钥对时指定此参数，因为如果您没有保存私钥，之后将无法检索它。
* `tags` - (可选) 要分配给资源的标签映射。

> **注意：** 如果既未设置`key_pair_name`也未设置`key_name_prefix`，Terraform将生成一个唯一ID来替换它们。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `key_pair_name` - 密钥对的名称。长度为2~128个英文或中文字符。必须以字母或中文开头，不能以`http://`或`https://`开头。可以包含数字、半角冒号(`:`)、下划线(`_`)或者短划线(`-`)。
* `finger_print` - 密钥对的指纹。根据RFC 4716定义的公钥指纹格式，采用MD5信息摘要算法。更多详情，请参见[RFC 4716](https://tools.ietf.org/html/rfc4716)。