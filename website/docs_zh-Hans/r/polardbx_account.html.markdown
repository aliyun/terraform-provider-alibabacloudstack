---
subcategory: "云原生分布式数据库PolarDB-X 2.0"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_account"
sidebar_current: "docs-Alibabacloudstack-PolarDBX-account"
description: |-
  提供 PolarDBX 账户资源。
---

# alibabacloudstack_polardbx_account

提供 PolarDBX 账户资源。

> **注意：** 该资源也可通过以下别名引用：
> - `apsarastack_polardbx_account`

## 示例用法

```hcl
variable "name" {
  default = "accdbbind90366"
}

variable "password" {
  default = ""
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.1.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_polardbx_instance" "default" {
  description = "testtf1111"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  engine_version = "5.7"
  storage = "50"
  network_type = "vpc"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
  cn_node_class = "polarx.x4.medium.2e"
  cn_node_count = "2"
  dn_node_class = "mysql.n4.medium.25"
  dn_node_count = "2"
}

resource "alibabacloudstack_polardbx_account" "default" {
  instance_id  = alibabacloudstack_polardbx_instance.default.id
  account_name = var.name
  password     = var.password
  description  = "Normal user"
}
```

## 参数参考

支持以下参数：

* `instance_id` - (必填，ForceNew) PolarDBX 实例的 ID。修改此参数会强制重新创建资源。
* `account_name` - (必填，ForceNew) 账户名称。修改此参数会强制重新创建资源。账户名称必须满足以下要求：
  * 以小写字母开头，以字母或数字结尾
  * 由小写字母、数字或下划线组成
  * 长度为 2 到 16 个字符
  * 不能使用保留用户名，如 root 和 admin
* `password` - (必填) 账户密码。该参数为敏感信息，在 Terraform 状态文件中会被标记为敏感。
* `description` - (可选) 账户描述。长度限制为 2 到 256 个字符，不能以 `http://` 或 `https://` 开头。此属性由 API 返回，无法手动设置。

## 属性参考

除上述参数外，还导出以下属性：

* `id` - 资源的唯一标识，格式为 `<instance_id>:<account_name>`
* `description` - 账户描述