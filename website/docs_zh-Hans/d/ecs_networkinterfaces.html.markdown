---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_networkinterfaces"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-networkinterfaces"
description: |- 
  查询云服务器弹性网卡
---

# alibabacloudstack_ecs_networkinterfaces
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_network_interfaces`

根据指定过滤条件列出当前凭证权限可以访问的云服务器弹性网卡列表。

## 示例用法

```hcl
variable "name" {
  default = "networkInterfaceDatasource"
}

resource "alibabacloudstack_vpc" "vpc" {
  name       = var.name
  cidr_block = "192.168.0.0/24"
}

resource "alibabacloudstack_vswitch" "vswitch" {
  name              = var.name
  cidr_block        = "192.168.0.0/24"
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  vpc_id            = alibabacloudstack_vpc.vpc.id
}

resource "alibabacloudstack_security_group" "group" {
  name   = var.name
  vpc_id = alibabacloudstack_vpc.vpc.id
}

data "alibabacloudstack_instance_types" "instance_type" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  eni_amount        = 2
  sorted_by         = "Memory"
}

resource "alibabacloudstack_instance" "instance" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  security_groups   = [alibabacloudstack_security_group.group.id]
  instance_type     = data.alibabacloudstack_instance_types.instance_type.instance_types[0].id
  system_disk_category = "cloud_efficiency"
  image_id          = data.alibabacloudstack_images.default.images[0].id
  instance_name     = var.name
  vswitch_id        = alibabacloudstack_vswitch.vswitch.id
}

resource "alibabacloudstack_network_interface" "interface" {
  name            = var.name
  vswitch_id      = alibabacloudstack_vswitch.vswitch.id
  security_groups = [alibabacloudstack_security_group.group.id]
}

resource "alibabacloudstack_network_interface_attachment" "attachment" {
  instance_id          = alibabacloudstack_instance.instance.id
  network_interface_id = alibabacloudstack_network_interface.interface.id
}

data "alibabacloudstack_ecs_networkinterfaces" "enis" {
  ids            = [alibabacloudstack_network_interface_attachment.attachment.network_interface_id]
  name_regex     = var.name
  vpc_id         = alibabacloudstack_vpc.vpc.id
  vswitch_id     = alibabacloudstack_vswitch.vswitch.id
  private_ip     = "192.168.0.2"
  security_group_id = alibabacloudstack_security_group.group.id
  type           = "Secondary"
  instance_id    = alibabacloudstack_instance.instance.id
  output_file    = "eni_list.txt"
}

output "eni_name" {
  value = data.alibabacloudstack_ecs_networkinterfaces.enis.interfaces.0.name
}
```

## 参数参考

以下参数是支持的：

* `ids` - (可选) ENI ID 列表，用于精确匹配需要查询的弹性网卡。
* `name_regex` - (可选) 用于通过 ENI 名称过滤结果的正则表达式字符串。
* `vpc_id` - (可选) 与 ENI 链接的 VPC ID。
* `vswitch_id` - (可选) 与 ENI 链接的 VSwitch ID。
* `private_ip` - (可选) ENI 的主私有 IP 地址。
* `security_group_id` - (可选) 与 ENI 链接的安全组 ID。
* `type` - (可选) ENI 类型，仅支持 "Primary" 或 "Secondary"。
* `instance_id` - (可选) ENI 所附加的 ECS 实例 ID。

## 属性参考

除了上述参数外，还导出以下属性：

* `names` - ENI 名称列表。
* `interfaces` - ENI 列表。每个元素包含以下属性：
  * `id` - ENI 的 ID。
  * `status` - ENI 当前状态。
  * `vpc_id` - ENI 所属的 VPC 的 ID。
  * `vswitch_id` - ENI 链接的 VSwitch 的 ID。
  * `zone_id` - ENI 所属的可用区 ID。
  * `public_ip` - ENI 的公网 IP。
  * `private_ip` - ENI 的主私有 IP。
  * `private_ips` - 分配给 ENI 的次级私有 IP 地址列表。
  * `mac` - ENI 的 MAC 地址。
  * `security_groups` - ENI 所属的安全组列表。
  * `name` - ENI 的名称。
  * `description` - ENI 的描述。
  * `instance_id` - ENI 所附加的实例 ID。
  * `creation_time` - ENI 的创建时间。
  * `tags` - 分配给 ENI 的标签映射。