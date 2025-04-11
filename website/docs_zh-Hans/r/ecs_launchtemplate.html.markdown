---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_launchtemplate"
sidebar_current: "docs-Alibabacloudstack-ecs-launchtemplate"
description: |- 
  编排云服务器（Ecs）启动模板资源
---

# alibabacloudstack_ecs_launchtemplate
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_launch_template`

使用Provider配置的凭证在指定的资源集下编排云服务器（Ecs）启动模板资源。

## 示例用法

```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

data "alibabacloudstack_instance_types" "all" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
}

data "alibabacloudstack_instance_types" "any_n4" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count       = 1
  memory_size          = 1
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

locals {
  default_instance_type_id = try(element(sort(length(data.alibabacloudstack_instance_types.default.instance_types) > 0 ? data.alibabacloudstack_instance_types.default.ids : data.alibabacloudstack_instance_types.any_n4.ids), 0), sort(data.alibabacloudstack_instance_types.all.ids)[0])
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_"
  most_recent = true
  owners      = "system"
}

variable "name" {
  default = "tf-testaccLaunchTemplateBasic12183"
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

resource "alibabacloudstack_security_group" "default" {
  name   = var.name
  vpc_id = alibabacloudstack_vpc.default.id
}

resource "alibabacloudstack_security_group_rule" "default" {
  type                = "ingress"
  ip_protocol         = "tcp"
  nic_type           = "intranet"
  policy             = "accept"
  port_range         = "22/22"
  priority           = 1
  security_group_id  = alibabacloudstack_security_group.default.id
  cidr_ip            = "172.16.0.0/24"
}

resource "alibabacloudstack_launch_template" "default" {
  name                          = var.name
  description                   = "Test launch template"
  host_name                     = var.name
  image_id                      = data.alibabacloudstack_images.default.images.0.id
  instance_name                 = var.name
  instance_type                 = local.default_instance_type_id
  internet_max_bandwidth_in     = 5
  internet_max_bandwidth_out    = 0
  io_optimized                  = "none"
  key_pair_name                 = "test-key-pair"
  ram_role_name                 = "xxxxx"
  network_type                  = "vpc"
  security_enhancement_strategy = "Active"
  spot_price_limit              = 5
  spot_strategy                 = "SpotWithPriceLimit"
  security_group_id             = alibabacloudstack_security_group.default.id
  system_disk_category          = "cloud_ssd"
  system_disk_description       = "Test disk"
  system_disk_name              = "hello"
  system_disk_size              = 40
  resource_group_id             = "rg-zkdfjahg9zxncv0"
  userdata                      = "xxxxxxxxxxxxxx"
  vswitch_id                    = alibabacloudstack_vswitch.default.id
  vpc_id                        = alibabacloudstack_vpc.default.id
  zone_id                       = data.alibabacloudstack_zones.default.zones.0.id

  tags = {
    tag1 = "hello"
    tag2 = "world"
  }

  network_interfaces {
    name              = "eth0"
    description       = "NI"
    primary_ip        = "10.0.0.2"
    security_group_id = "xxxx"
    vswitch_id        = "xxxxxxx"
  }

  data_disks {
    name        = "disk1"
    size        = 20
    category    = "cloud_efficiency"
    description = "test1"
  }

  data_disks {
    name        = "disk2"
    size        = 30
    category    = "cloud_ssd"
    description = "test2"
  }
}
```

## 参数说明

支持以下参数：

* `name` - (选填, 变更时重建) - 模板的名称。必须以英文字母(大写或小写)开头，并可以包含数字、句点 (.)、冒号 (:)、下划线 (_) 和连字符 (-)。长度应在 2 到 128 个字符之间。不能以 "http://" 或 "https://" 开头。
* `launch_template_name` - (选填, 变更时重建) - 启动模板的名称。
* `description` - (选填) - 启动模板版本1的描述。长度为2~256个英文或中文字符，不能以`http://`或`https://`开头。默认值为空。
* `host_name` - (选填) - 实例的主机名。
  - 半角句号 (.) 和短划线 (-) 不能作为首尾字符，更不能连续使用。
  - Windows 实例：字符长度为 2~15，不支持半角句号 (.)，不能全是数字。允许大小写英文字母、数字和短划线 (-)。
  - 其他类型实例 (Linux 等)：字符长度为 2~64，支持多个半角句号 (.)，半角句号之间为一段，每段允许大小写英文字母、数字和短划线 (-)。
* `image_id` - (选填) - 创建实例所使用的镜像 ID。
* `image_owner_alias` - (选填) - 镜像来源。有效值：
  - `system`: 阿里云提供的公共镜像。
  - `self`: 您创建的自定义镜像。
  - `others`: 来自另一个阿里云账户的共享镜像。
  - `marketplace`: 市场镜像。
* `instance_charge_type` - (选填) - 实例的计费方式。有效值：
  - `PrePaid`: 包年包月。
  - `PostPaid`: 按量付费。
* `instance_name` - (选填) - 实例名称。长度为 2~128 个英文或中文字符。必须以大小写字母或中文开头，不能以 `http://` 或 `https://` 开头。可以包含数字、半角冒号 (:)、下划线 (_) 或者短划线 (-)。
* `instance_type` - (选填) - 实例规格。您可以使用 `data alibabacloudstack_instance_types` 数据源来获取最新的实例类型列表。
* `internet_charge_type` - (选填) - 网络计费方式。有效值：
  - `PayByBandwidth`: 按固定带宽计费。
  - `PayByTraffic`: 按使用流量计费。
* `internet_max_bandwidth_in` - (选填) - 最大入网带宽，单位为 Mbit/s。取值范围为 [1, 200]。
* `internet_max_bandwidth_out` - (选填) - 最大出网带宽，单位为 Mbit/s。取值范围为 [0, 100]。
* `io_optimized` - (选填) - 是否为 I/O 优化实例。有效值：
  - `none`: 非 I/O 优化。
  - `optimized`: I/O 优化。
* `key_pair_name` - (选填) - 密钥对名称。对于 Windows 实例将被忽略。
* `network_type` - (选填) - 实例的网络类型。有效值：`Classic`, `VPC`。
* `ram_role_name` - (选填) - 分配给实例的 RAM 角色名称。
* `resource_group_id` - (选填) - 实例所属的资源组 ID。
* `security_enhancement_strategy` - (选填) - 是否激活安全增强功能。有效值：`Active`, `Deactive`。
* `security_group_id` - (选填) - 安全组 ID。
* `spot_price_limit` - (选填) - 竞价实例的最大小时价格。支持最多三位小数。
* `spot_strategy` - (选填) - 按量付费实例的竞价策略。有效值：
  - `NoSpot`: 普通按量付费实例。
  - `SpotWithPriceLimit`: 带最高价格限制的竞价实例。
  - `SpotAsPriceGo`: 系统自动计算价格。
* `system_disk_category` - (选填) - 系统盘的类别。有效值：
  - `cloud`: 普通云盘。
  - `cloud_efficiency`: 高效云盘。
  - `cloud_ssd`: SSD 云盘。
  - `ephemeral_ssd`: 本地 SSD 盘。
  - `cloud_essd`: ESSD 云盘。
* `system_disk_description` - (选填) - 系统盘的描述。
* `system_disk_name` - (选填) - 系统盘的名称。
* `system_disk_size` - (选填) - 系统盘大小，单位为 GB。取值范围为 [20, 500]。
* `userdata` - (选填) - 实例的用户数据，Base64 编码。原始数据大小不能超过 16 KB。
* `vswitch_id` - (选填) - 创建 VPC 类型实例时的交换机 ID。
* `vpc_id` - (选填) - VPC ID。
* `zone_id` - (选填) - 实例所在的可用区 ID。
* `network_interfaces` - (选填) - 随实例一起创建的网络接口列表。
  * `name` - (选填) - 网络接口的名称。
  * `description` - (选填) - 网络接口的描述。
  * `primary_ip` - (选填) - 网络接口的主要私有 IP 地址。
  * `security_group_id` - (选填) - 网络接口的安全组 ID。
  * `vswitch_id` - (选填) - 网络接口的交换机 ID。
* `data_disks` - (选填) - 随实例一起创建的数据盘列表。
  * `name` - (选填) - 数据盘的名称。
  * `size` - (必填) - 数据盘的大小，单位为 GB。
    - `cloud`: [5, 2000]
    - `cloud_efficiency`: [20, 32768]
    - `cloud_ssd`: [20, 32768]
    - `cloud_essd`: [20, 32768]
    - `ephemeral_ssd`: [5, 800]
  * `category` - (选填) - 数据盘的类别。默认值为 `cloud_efficiency`。
  * `encrypted` - (选填, Bool) - 数据盘是否加密。默认值为 `false`。
  * `snapshot_id` - (选填) - 用于初始化数据盘的快照 ID。
  * `delete_with_instance` - (选填) - 实例销毁时是否删除数据盘。默认值为 `true`。
  * `description` - (选填) - 数据盘的描述。
* `tags` - (选填) - 要分配给资源的标签映射。
  - Key: 最多 64 个字符长度。不能以 "aliyun", "acs:", "http://", 或 "https://" 开头。
  - Value: 最多 128 个字符长度。可以是空字符串。
* `auto_release_time` - (选填) - 实例的定时释放时间。
* `user_data` - (选填) - 实例的用户数据，Base64 编码。原始数据大小不能超过 16 KB。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 启动模板的 ID。
* `launch_template_name` - 启动模板的名称。
* `internet_max_bandwidth_in` - 最大公网入方向带宽，单位为 Mbit/s。
* `internet_max_bandwidth_out` - 最大公网出方向带宽，单位为 Mbit/s。
* `name` - (计算后返回) 启动模板的名称。