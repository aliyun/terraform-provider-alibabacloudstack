---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_instance"
sidebar_current: "docs-Alibabacloudstack-ecs-instance"
description: |- 
  使用Provider配置的凭证在指定的资源集下编排云服务器（Ecs）实例
---

# alibabacloudstack_ecs_instance
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_instance`

使用Provider配置的凭证在指定的资源集下编排云服务器（Ecs）实例。

## 示例用法

```hcl
data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details             = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name        = "${var.name}_vsw"
  vpc_id      = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block  = "172.16.0.0/24"
  zone_id     = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_ecs_securitygroup" "default" {
  name   = "${var.name}_sg"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

resource "alibabacloudstack_security_group_rule" "default" {
  type                = "ingress"
  ip_protocol         = "tcp"
  nic_type           = "intranet"
  policy             = "accept"
  port_range         = "22/22"
  priority           = 1
  security_group_id  = "${alibabacloudstack_ecs_securitygroup.default.id}"
  cidr_ip            = "172.16.0.0/24"
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_"
  most_recent = true
  owners      = "system"
}

data "alibabacloudstack_instance_types" "all" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
}

data "alibabacloudstack_instance_types" "any_n4" {
  availability_zone    = data.alibabacloudstack_zones.default.zones[0].id
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count   = 1
  memory_size      = 1
  instance_type_family = "ecs.n4"
  sorted_by        = "Memory"
}

locals {
  default_instance_type_id = try(element(sort(length(data.alibabacloudstack_instance_types.default.instance_types) > 0 ? data.alibabacloudstack_instance_types.default.ids : data.alibabacloudstack_instance_types.any_n4.ids), 0), sort(data.alibabacloudstack_instance_types.all.ids)[0])
}

variable "name" {
  default = "tf-testAccEcsInstanceConfigBasic2648"
}

resource "alibabacloudstack_ecs_instance" "default" {
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  instance_name        = "${var.name}"
  user_data            = "I_am_user_data"
  security_groups      = ["${alibabacloudstack_ecs_securitygroup.default.id}"]
  vswitch_id          = "${alibabacloudstack_vpc_vswitch.default.id}"
  tags = {
    Bar = "Bar"
    foo = "foo"
  }
  image_id            = "${data.alibabacloudstack_images.default.images.0.id}"
  security_enhancement_strategy = "Active"
  instance_type       = "${local.default_instance_type_id}"
  availability_zone   = "${data.alibabacloudstack_zones.default.zones[0].id}"

  # IPv6 配置
  enable_ipv6         = true
  ipv6_cidr_block     = "fd00::/64"
  ipv6_address_count  = 3

  # 数据盘
  data_disks = [
    {
      category         = "cloud_efficiency"
      size             = 50
      delete_with_instance = true
    },
    {
      category         = "cloud_ssd"
      size             = 100
      snapshot_id      = "snap-12345678"
      delete_with_instance = false
    }
  ]
}
```

## 参数参考

支持以下参数：
  * `availability_zone` - (选填, 变更时重建) - 启动实例所在的可用区。必须与指定的交换机的区域匹配。
  * `zone_id` - (选填, 变更时重建) - 实例所属可用区。
  * `image_id` - (必填) - 实例运行的镜像ID。更改此值将强制创建新资源。
  * `instance_type` - (必填) - 实例规格 ID。更改此值将强制创建新资源。
  * `instance_name` - (选填) - 实例名称，长度为 2 到 128 个字符，不能以 http:// 或 https:// 开头。
  * `description` - (选填) - 实例描述，长度为 2 到 256 个字符，不能以 http:// 或 https:// 开头。
  * `internet_max_bandwidth_in` - (选填) - 公网入带宽最大值，单位 Mbps(兆比特每秒)。取值范围：[1, 200]。默认值为 200 Mbps。
  * `internet_max_bandwidth_out` - (选填) - 公网出带宽最大值，单位 Mbps(兆比特每秒)。取值范围：[0, 100]。默认值为 0 Mbps。
  * `host_name` - (选填) - 实例主机名。更改此值会导致实例重启。
  * `password` - (选填, 敏感信息) - 实例的密码。长度为8至30个字符，必须同时包含大小写英文字母、数字和特殊符号中的三类字符。特殊符号可以是：```()`~!@#$%^&*-_+=|{}[]:;'<>,.?/```。您需要注意：- 如果传入Password参数，建议您使用HTTPS协议发送请求，避免密码泄露。- Windows实例不能以正斜线(/)为密码首字符。- 部分操作系统的实例不支持配置密码，仅支持配置密钥对。例如：Others Linux、Fedora CoreOS。
  * `kms_encrypted_password` - (选填) - 使用 KMS 加密的实例密码。如果提供了 `password`，此字段将被忽略。更改此值会导致实例重启。
  * `kms_encryption_context` - (选填) - 用于解密 `kms_encrypted_password` 的 KMS 加密上下文。当设置 `kms_encrypted_password` 时有效。更改此值会导致实例重启。
  * `is_outdated` - (选填) - 是否使用过时的实例类型。默认值为 `false`。
  * `system_disk_category` - (选填, 变更时重建) - 系统盘的类别。有效值：`ephemeral_ssd`, `cloud_efficiency`, `cloud_ssd`, `cloud_essd`, `cloud`。默认值为 `cloud_efficiency`。
  * `system_disk_size` - (选填) - 系统盘大小，单位 GiB。取值范围：[20, 500]。默认值为 {40, ImageSize} 中的最大值。
  * `system_disk_name` - (选填) - 系统盘的名称。更改此值会导致实例重启。
  * `system_disk_description` - (选填) - 系统盘的描述。更改此值会导致实例重启。
  * `data_disks` - (选填, 变更时重建) - 随实例一起创建的数据盘列表。每个数据盘支持以下属性：
    * `category` - (选填, 变更时重建) - 数据盘的类别。有效值：`cloud`, `cloud_efficiency`, `cloud_ssd`, `ephemeral_ssd`。默认值为 `cloud_efficiency`。
    * `size` - (必填, 变更时重建) - 数据盘的大小，单位 GiB。
    * `snapshot_id` - (选填, 变更时重建) - 用于初始化数据盘的快照 ID。
    * `delete_with_instance` - (选填, 变更时重建) - 是否在销毁实例时删除数据盘。默认值为 `true`。
    * `encrypted` - (选填, 布尔值, 变更时重建) - 是否加密数据盘。默认值为 `false`。
    * `kms_key_id` - (选填) - 数据盘对应的 KMS 密钥 ID。
    * `name` - (选填, 变更时重建) - 数据盘的名称。
    * `description` - (选填, 变更时重建) - 数据盘的描述。
  * `subnet_id` - (自 v1.210.0 起已移除) - 子网 ID。与 `vswitch_id` 冲突。
  * `vswitch_id` - (选填) - 在 VPC 中启动的虚拟交换机 ID。除非可以创建经典网络实例，否则必须设置此参数。
  * `private_ip` - (选填) - 分配给实例的私有 IP 地址。当指定了 `vswitch_id` 时有效。
  * `hpc_cluster_id` - (选填, 变更时重建) - 实例所属的弹性高性能计算(E-HPC)集群 ID。
  * `user_data` - (选填) - 自定义用户数据以定制 ECS 实例的启动行为。更改此值会导致实例重启。
  * `role_name` - (选填, 变更时重建) - 关联到实例的 RAM 角色名称。
  * `key_name` - (选填, 变更时重建) - 实例使用的密钥对名称。
  * `storage_set_id` - (选填, 变更时重建) - 存储集 ID。
  * `storage_set_partition_number` - (选填, 变更时重建) - 存储集中的分区数量。
  * `security_enhancement_strategy` - (选填, 变更时重建) - 安全增强策略。有效值：`Active`(启用安全增强策略)，`Deactive`(禁用安全增强策略)。
  * `enable_ipv6` - (选填, 变更时重建) - 是否启用 IPv6。有效值：`false`(禁用)，`true`(启用)。
  * `ipv6_cidr_block` - (选填) - VPC 的 IPv6 CIDR 段。
  * `ipv6_address_count` - (选填) - 请求分配的 IPv6 地址数量。如果 `enable_ipv6` 为 `true`，则此值必须大于 0。
  * `ipv6_address_list` - (选填, 变更时重建) - 要分配给主 ENI 的 IPv6 地址列表。支持最多 10 个地址。
  * `tags` - (选填) - 要分配给资源的标签映射。
  * `system_disk_tags` - (选填) - 要分配给系统盘的标签映射。
  * `data_disk_tags` - (选填) - 要分配给数据盘的标签映射。

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `availability_zone` - 实例所在的可用区。
  * `zone_id` - 实例所属可用区。
  * `internet_max_bandwidth_in` - 公网入带宽最大值。
  * `host_name` - 实例主机名。
  * `system_disk_id` - 系统盘的 ID。
  * `subnet_id` - 子网的 ID。
  * `private_ip` - 实例的私有 IP 地址。
  * `hpc_cluster_id` - 实例所属的 HPC 集群 ID。
  * `status` - 实例的状态。
  * `role_name` - 关联到实例的 RAM 角色名称。
  * `key_name` - 实例使用的密钥对名称。
  * `storage_set_id` - 存储集 ID。
  * `storage_set_partition_number` - 存储集中的分区数量。