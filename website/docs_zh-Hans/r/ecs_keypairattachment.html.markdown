---
subcategory: "云服务器 ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_keypairattachment"
sidebar_current: "docs-Alibabacloudstack-resource-ecs-keypairattachment"
description: |-
  编排绑定云服务器（ECS）密钥对和实例
---

# alibabacloudstack_ecs_keypairattachment

-> **NOTE:** 该资源等效别名有：`alibabacloudstack_key_pair_attachment`。

使用 Provider 配置的凭证在指定的资源集下编排绑定 SSH 密钥对到一个或多个 Linux 实例。

## 示例用法

以下是一个### 基础用法示例：

```hcl
data "alibabacloudstack_zones" "default" {
  available_disk_category     = "cloud_ssd"
  available_resource_creation = "VSwitch"
}

data "alibabacloudstack_instance_types" "type" {
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  cpu_core_count    = 1
  memory_size       = 2
}

data "alibabacloudstack_images" "images" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

variable "name" {
  default = "keyPairAttachmentName"
}

variable "password" {}

resource "alibabacloudstack_vpc" "vpc" {
  name       = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alibabacloudstack_vswitch" "vswitch" {
  vpc_id            = "${alibabacloudstack_vpc.vpc.id}"
  cidr_block        = "10.1.1.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alibabacloudstack_security_group" "group" {
  name        = "${var.name}"
  description = "New security group"
  vpc_id      = "${alibabacloudstack_vpc.vpc.id}"
}

resource "alibabacloudstack_instance" "instance" {
  instance_name   = "${var.name}-${count.index + 1}"
  image_id        = "${data.alibabacloudstack_images.images.images.0.id}"
  instance_type   = "${data.alibabacloudstack_instance_types.type.instance_types.0.id}"
  count           = 2
  security_groups = ["${alibabacloudstack_security_group.group.id}"]
  vswitch_id      = "${alibabacloudstack_vswitch.vswitch.id}"
  internet_max_bandwidth_out = 5
  password                   = var.password
  system_disk_category = "cloud_ssd"
}

resource "alibabacloudstack_key_pair" "pair" {
  key_name = "${var.name}"
}

resource "alibabacloudstack_ecs_keypairattachment" "attachment" {
  key_name     = "${alibabacloudstack_key_pair.pair.key_name}"
  instance_ids = ["${alibabacloudstack_instance.instance.*.id}"]
  force        = true
}
```

## 参数说明

支持以下参数：

* `key_name` - (必填，变更时重建) SSH 密钥对的名称，用于标识要绑定到实例的密钥对。名称长度必须在 2 到 128 个字符之间。变更此参数会强制重新创建资源。
* `instance_ids` - (必填，变更时重建) 需要绑定 SSH 密钥对的 ECS 实例 ID 列表。最多可指定 50 个实例 ID，以 JSON 数组格式表示。变更此参数会强制重新创建资源。
* `force` - (选填，变更时重建) 如果设置为 `true`，在绑定密钥对后，实例将自动重启以确保密钥对立即生效。默认值为 `false`。变更此参数会强制重新创建资源。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 资源 ID，格式为 `<key_name>:<instance_ids>`，其中 `instance_ids` 是 JSON 数组字符串。
* `key_name` - 已成功绑定到实例的 SSH 密钥对名称。
* `instance_ids` - SSH 密钥对已成功绑定的 ECS 实例 ID 列表。

## Import

ECS 密钥对绑定可以通过密钥对名称和实例 ID 组合导入，格式为 `<key_name>:<instance_ids>`，例如：

```
$ terraform import alibabacloudstack_ecs_keypairattachment.example test-key-pair:["i-bp1d6tsvznfghy7y****","i-bp1ippxbaql9zet7****"]
```