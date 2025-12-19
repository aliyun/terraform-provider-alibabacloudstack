---
subcategory: "BMS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bms_keypair"
sidebar_current: "docs-Alibabacloudstack-bms-keypair"
description: |-
  管理裸金属服务器密钥对
---

# alibabacloudstack_bms_keypair

管理裸金属服务器密钥对资源，用于创建、读取和删除密钥对。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tf-bms-keypair1641"
}

variable "public_key" {
  default = "ssh-rsa ********* == root@vm010017040011"
}


resource "alibabacloudstack_bms_keypair" "default" {
  name = var.name
}
```

## 参数说明

支持以下参数：

* `name` - (必填, 变更时重建) 密钥对的名称。长度为1-128个字符，支持字母、数字、下划线（_）和短划线（-）。

* `public_key` - (可选, 变更时重建) 公钥内容。如果不提供，系统将自动生成密钥对。公钥格式应为SSH公钥格式，如"ssh-rsa AAAAB3NzaC1yc2EAAA..."。

## 属性说明

以下属性会从API响应中导出：

* `id` - 密钥对的ID，与name相同。
* `key_pair_fingerprint` - 密钥对的指纹信息，格式为"[位数] [算法]:[哈希值] [注释]"，例如"4096 SHA256:+MiS793F6Nxn/ygAGtquHw2e5grziG/AA+0ESwvF3VM root@vm010017040011 (RSA)"。
* `private_key` - 私钥内容（敏感信息），仅在创建时返回，后续读取将为空。格式为PEM格式的RSA私钥。
* `public_key` - 公钥内容，格式为SSH公钥格式。如果创建时提供了公钥，则返回提供的公钥；如果系统自动生成，则返回生成的公钥。