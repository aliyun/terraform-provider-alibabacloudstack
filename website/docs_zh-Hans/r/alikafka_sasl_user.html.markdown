---
subcategory: "Alikafka"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_alikafka_sasl_user"
sidebar_current: "docs-alibabacloudstack-resource-alikafka-sasl_user"
description: |-
  编排Alikafka SASL用户
---

# alibabacloudstack_alikafka_sasl_user

使用Provider配置的凭证在指定的资源集下编排Alikafka SASL用户资源。

## 示例用法

### 基础用法

```
variable "username" {
  default = "testusername"
}

variable "password" {
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  zone_id           = data.alibabacloudstack_zones.default.zones[0].id
}

resource "alibabacloudstack_alikafka_instance" "default" {
  name        = "tf-testacc-alikafkainstance"
  topic_quota = "50"
  disk_type   = "1"
  disk_size   = "500"
  deploy_type = "5"
  io_max      = "20"
  vswitch_id  = alibabacloudstack_vswitch.default.id
}

resource "alibabacloudstack_alikafka_sasl_user" "default" {
  instance_id = alibabacloudstack_alikafka_instance.default.id
  username    = var.username
  password    = var.password
}
```

## 参数说明

支持以下参数：

* `instance_id` - (必填，变更时重建) 拥有该组的ALIKAFKA实例的ID。
* `username` - (必填，变更时重建) SASL用户的用户名。长度应在1到64个字符之间。字符只能包含'a'-'z'、'A'-'Z'、'0'-'9'、'_'和'-'。
* `password` - (可选，敏感) 操作密码。它可能由字母、数字或下划线组成，长度为1到64个字符。您必须指定`password`和`kms_encrypted_password`字段之一。
* `kms_encrypted_password` - (可选) 用于数据库账户的KMS加密密码。您必须指定`password`和`kms_encrypted_password`字段之一。
* `kms_encryption_context` - (可选，MapString) 用于在使用`kms_encrypted_password`创建或更新用户之前解密`kms_encrypted_password`的KMS加密上下文。参见[Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm)。当设置了`kms_encrypted_password`时有效。
* `type` - (可选，变更时重建，1.159.0版本及以上可用) 认证机制。有效值：`plain`，`scram`。默认值：`plain`。

## 属性说明

导出以下属性：

* `id` - 资源的ID。格式为 `<instance_id>:<username>`。
* `type` - (计算得出，1.159.0版本及以上可用) 认证机制。有效值：`plain`，`scram`。

## 导入

Alikafka SASL 用户可以通过ID导入，例如

```bash
$ terraform import alibabacloudstack_alikafka_sasl_user.example <instance_id>:<username>
```