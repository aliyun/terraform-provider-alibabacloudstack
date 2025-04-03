---
subcategory: "GPDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_gpdb_accounts"
sidebar_current: "docs-Alibabacloudstack-datasource-gpdb-accounts"
description: |- 
  查询图数据库帐号
---

# alibabacloudstack_gpdb_accounts

根据指定过滤条件列出当前凭证权限可以访问的图数据库帐号列表。

## 示例用法

### 基础用法：

```terraform
variable "name" {
  default = "tftestacc3018"
}

data "alibabacloudstack_gpdb_zones" "default" {}

data "alibabacloudstack_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alibabacloudstack_vswitches" "default" {
  vpc_id  = data.alibabacloudstack_vpcs.default.ids.0
  zone_id = data.alibabacloudstack_gpdb_zones.default.zones.2.id
}

resource "alibabacloudstack_vswitch" "default" {
  count             = length(data.alibabacloudstack_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alibabacloudstack_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alibabacloudstack_vpcs.default.vpcs[0].cidr_block, 8, 8)
  availability_zone = data.alibabacloudstack_gpdb_zones.default.zones.3.id
}

resource "alibabacloudstack_gpdb_elastic_instance" "default" {
  engine                   = "gpdb"
  engine_version           = "6.0"
  seg_storage_type         = "cloud_essd"
  seg_node_num             = 4
  storage_size             = 50
  instance_spec            = "2C16G"
  db_instance_description  = var.name
  instance_network_type    = "VPC"
  payment_type             = "PayAsYouGo"
  vswitch_id               = length(data.alibabacloudstack_vswitches.default.ids) > 0 ? data.alibabacloudstack_vswitches.default.ids[0] : concat(alibabacloudstack_vswitch.default.*.id, [""])[0]
}

resource "alibabacloudstack_gpdb_account" "default" {
  account_name        = var.name
  db_instance_id      = alibabacloudstack_gpdb_elastic_instance.default.id
  account_password    = "inputYourCodeHere"
  account_description = var.name
}

data "alibabacloudstack_gpdb_accounts" "default" {  
  db_instance_id = alibabacloudstack_gpdb_elastic_instance.default.id
  ids            = [alibabacloudstack_gpdb_account.default.account_name]
}

output "gpdb_account_id" {
  value = data.alibabacloudstack_gpdb_accounts.default.accounts[0].id
}
```

使用正则表达式筛选账号名称：

```terraform
data "alibabacloudstack_gpdb_accounts" "nameRegex" {
  db_instance_id = "example_value"
  name_regex     = "^my-Account"
}

output "gpdb_account_id_2" {
  value = data.alibabacloudstack_gpdb_accounts.nameRegex.accounts[0].id
}
```

## 参数参考

以下参数是支持的：

* `db_instance_id` - (必填，变更时重建) ：GPDB 实例的 ID。这是查询 GPDB 账号时必须提供的参数。
* `ids` - (选填，变更时重建) ：账号 ID 列表。其元素值与账号名称相同，用于通过账号 ID 筛选结果。
* `name_regex` - (选填，变更时重建) ：用于通过账号名称筛选结果的正则表达式字符串。
* `status` - (选填，变更时重建) ：账号的状态。有效值为 `Active`、`Creating` 和 `Deleting`。

## 属性参考

除了上述参数外，还导出以下属性：

* `names` - 账号名称列表。
* `accounts` - GPDB 账号列表。每个元素包含以下属性：
  * `account_description` - 账号的描述信息。
  * `id` - 账号的 ID。其值与账号名称相同。
  * `account_name` - 账号的名称。
  * `db_instance_id` - GPDB 实例的 ID。
  * `status` - 账号的状态。有效值为 `Active`、`Creating` 和 `Deleting`。