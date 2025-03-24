---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_keypairattachment"
sidebar_current: "docs-Alibabacloudstack-ecs-keypairattachment"
description: |- 
  编排绑定云服务器（Ecs）密钥对和实例
---

# alibabacloudstack_ecs_keypairattachment
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_key_pair_attachment`

使用Provider配置的凭证在指定的资源集下编排绑定云服务器（Ecs）密钥对和实例。

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

## 参数参考

支持以下参数：

* `key_name` - (必填，变更时重建) 密钥对的名称。此名称用于标识要绑定到实例的密钥对。
* `instance_ids` - (必填，变更时重建) 需要附加密钥对的 ECS 实例 ID 列表。每个实例都会与指定的密钥对进行绑定。
* `force` - (选填，变更时重建) 如果设置为 `true`，在附加密钥对后，实例将自动重启以确保密钥对生效。默认值为 `false`。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `key_name` - 已成功附加到实例的密钥对名称。
* `instance_ids` - 密钥对已成功附加的 ECS 实例 ID 列表。