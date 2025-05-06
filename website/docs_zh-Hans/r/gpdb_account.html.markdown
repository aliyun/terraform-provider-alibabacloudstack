---
subcategory: "GraphDatabase(GPDB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_gpdb_account"
sidebar_current: "docs-Alibabacloudstack-gpdb-account"
description: |- 
  编排图数据库帐号
---

# alibabacloudstack_gpdb_account

使用Provider配置的凭证在指定的资源集下编排图数据库帐号。

## 示例用法

```terraform
variable "name" {
  default = "tftest1124"
}
variable "password" {
  default = "TestPassword123!"
}

data "alibabacloudstack_gpdb_zones" "default" {}
data "alibabacloudstack_zones" "default" {}
data "alibabacloudstack_vpcs" "default" {
  name_regex = "default-NODELETING"
}

resource "alibabacloudstack_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
  name              = var.name
}

resource "alibabacloudstack_gpdb_instance" "default" {
  availability_zone      = data.alibabacloudstack_zones.default.zones.0.id
  engine                = "gpdb"
  engine_version        = "4.3"
  instance_class        = "gpdb.group.segsdx2"
  instance_group_count  = 2
  description          = "tf-testAccGpdbInstance_new"
  vswitch_id           = alibabacloudstack_vswitch.default.id
}

resource "alibabacloudstack_gpdb_account" "default" {
  account_name         = "tftest1124"
  account_password     = var.password
  account_description  = "tftest1124"
  db_instance_id       = alibabacloudstack_gpdb_instance.default.id
}
```

## 参数说明

支持以下参数：

* `account_description` - (可选，变更时重建) 账户的描述。  
  * 必须以字母开头。
  * 不得以 `http://` 或 `https://` 开头。
  * 只能包含字母、下划线(_)、连字符(-)或数字。
  * 长度必须在2到256个字符之间。

* `account_name` - (必填，变更时重建) 账户名称。账户名称必须唯一，并满足以下要求：
  * 必须以字母开头。
  * 只能包含小写字母、数字或下划线(_)。
  * 长度不得超过16个字符。
  * 不得包含保留关键字。

* `account_password` - (必填) 账户密码。密码必须为8到32个字符长度，并且至少包含以下四种字符类型中的三种：大写字母、小写字母、数字和特殊字符。特殊字符包括 `! @ # $ % ^ & * ( ) _ + - =`。

* `db_instance_id` - (必填，变更时重建) 实例的ID。

## 属性说明

除了上述所有参数外，还导出以下属性：

* `id` - 账户的资源ID。格式为 `<db_instance_id>:<account_name>`。
* `status` - 账户的状态。有效值为：`Active`（活跃）、`Creating`（创建中）和 `Deleting`（删除中）。

### 超时时间

`timeouts` 块允许您为某些操作指定 [超时时间](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts)：

* `create` - (默认为5分钟)用于创建账户时。

## 导入

GPDB账户可以使用ID导入，例如：

```bash
$ terraform import alicloud_gpdb_account.example <db_instance_id>:<account_name>
```