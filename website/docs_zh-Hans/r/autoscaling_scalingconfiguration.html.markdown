---
subcategory: "Auto Scaling (ESS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_autoscaling_configuration"
sidebar_current: "docs-Alibabacloudstack-autoscaling-configuration"
description: |- 
  编排弹性伸缩配置
---

# alibabacloudstack_ess_scaling_configuration

使用Provider配置的凭证在指定的资源集下编排弹性伸缩配置。

-> **注意：** 某些实例类型在某些地域和可用区已过时，例如 `ecs.t1.*`、`ecs.s2.*`、`ecs.n1.*` 等。如果您想继续使用它们，您应该将 `is_outdated` 设置为 true。有关升级后的实例类型，请参阅 `alibabacloudstack_instance_types` 数据源。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAccEssScCon-694731"
}

data "alibabacloudstack_zones" "default" {
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
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = "${alibabacloudstack_ecs_securitygroup.default.id}"
  cidr_ip           = "172.16.0.0/24"
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_"
  most_recent = true
  owners      = "system"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count    = 2
  memory_size       = 4
}

resource "alibabacloudstack_ecs_deployment_set" "default" {
  strategy            = "Availability"
  domain              = "Default"
  granularity         = "Host"
  deployment_set_name = "example_value"
  description         = "example_value"
}

resource "alibabacloudstack_ess_scaling_group" "default" {
  min_size           = 1
  max_size           = 1
  scaling_group_name = var.name
  removal_policies   = ["OldestInstance", "NewestInstance"]
  vswitch_ids        = [alibabacloudstack_vpc_vswitch.default.id]
}

resource "alibabacloudstack_ess_scaling_configuration" "default" {
  scaling_group_id  = alibabacloudstack_ess_scaling_group.default.id
  image_id          = data.alibabacloudstack_images.default.images[0].id
  instance_type     = data.alibabacloudstack_instance_types.default.instance_types[0].id
  security_group_ids = [alibabacloudstack_ecs_securitygroup.default.id]
  deployment_set_id = alibabacloudstack_ecs_deployment_set.default.id
  force_delete      = true
  active            = true
  internet_max_bandwidth_in = 50
  system_disk_category = "cloud_efficiency"
  system_disk_size = 40
  user_data = base64encode("echo 'Hello World' > /tmp/hello.txt")
  key_pair_name = "my-key-pair"
  tags = {
    Environment = "Test"
    Owner      = "Terraform"
  }
}
```

## 参数说明

支持以下参数：

* `status` - (可选) 伸缩配置在伸缩组中的状态。可能的值：
  * `Active`: 表示该配置处于活动状态，当前伸缩组将使用此配置自动创建ECS实例。
  * `Inactive`: 表示该配置处于非活动状态，当前伸缩组不会使用此配置自动创建ECS实例。

* `active` - (可选) 是否在指定的伸缩组中激活当前伸缩配置。默认值为 `false`。

* `enable` - (可选) 是否启用指定的伸缩组(使其处于活动状态)，当前伸缩配置属于该伸缩组。

* `scaling_group_id` - (必填, 强制更新) 伸缩配置所属的伸缩组ID。

* `image_id` - (必填) 创建ECS实例时使用的镜像文件ID。

* `instance_type` - (必填) ECS实例的规格。

* `security_group_ids` - (必填) ECS实例所属的安全组ID列表。

* `deployment_set_id` - (必填) ECS实例所属的部署集ID。

* `zone_id` - (可选) ECS实例所属的可用区ID。

* `scaling_configuration_name` - (可选) 伸缩配置名称。必须是2~64个字符，并且可以包含字母、数字、下划线 (`_`)、连字符 (`-`) 或点 (`.`)。如果未指定，默认值为伸缩配置ID。

* `internet_max_bandwidth_in` - (可选) 公网入方向的最大带宽，单位为Mbps(兆比特每秒)。取值范围为 [1, 200]。如果未指定，默认值为200 Mbps。

* `system_disk_category` - (可选) 系统盘的类别。有效选项包括 `ephemeral_ssd`、`cloud_efficiency`、`cloud_ssd`、`cloud_essd` 和 `cloud`。默认值为 `cloud_efficiency`。

* `system_disk_size` - (可选) 系统盘大小，单位为GiB。有效范围取决于磁盘类别：
  * `cloud`: [20, 500]
  * `cloud_efficiency`, `cloud_ssd`, `cloud_essd`: [20, 500]
  * `ephemeral_ssd`: [20, 500]

* `data_disk` - (可选) 要附加到ECS实例的数据盘列表。详见下方 [Block datadisk](#block-datadisk)。

* `substitute` - (可选) 当前配置设置为 `false` 时，自动成为活动状态并替换当前配置的伸缩配置。当 `active` 为 `true` 时，此参数无效。

* `system_disk_auto_snapshot_policy_id` - (可选) 系统盘的自动快照策略ID。

* `is_outdated` - (可选) 是否使用过时的实例类型。默认值为 `false`。

* `user_data` - (可选) 自定义数据用于ECS实例。它必须经过Base64编码，原始数据不得超过16KB。

* `ram_role_name` - (可选) ECS实例的RAM角色名称。您可以使用 [ListRoles](~~ 28713 ~~) API 查询可用的RAM角色。要创建RAM角色，请参见 [CreateRole](~~ 28710 ~~)。

* `key_pair_name` - (可选) 用于登录ECS实例的密钥对名称。对于Windows实例，此参数将被忽略。对于Linux实例，如果指定了此参数，则密码登录将被禁用。

* `force_delete` - (可选) 是否强制删除最后一个伸缩配置及其所属的伸缩组。默认值为 `false`。

* `tags` - (可选) 分配给资源的标签映射。这些标签将应用于伸缩组创建的ECS实例。

* `instance_name` - (可选) ECS实例名称。

* `override` - (可选) 是否覆盖现有数据。默认值为 `false`。

* `host_name` - (可选) 服务器主机名。限制因操作系统而异：
  * Windows: 主机名长度为2~15个字符，可以包含大写字母、数字和连字符 (`-`)。不能以句点 (`.`) 或连字符 (`-`) 开头或结尾。
  * Linux: 主机名长度为2~64个字符，可以包含多个句点 (`.`)。每个段之间的句点可以包含大写字母、数字和连字符 (`-`)。

### Block datadisk

`data_disk` 块支持以下参数：

* `size` - (可选) 数据盘大小，单位为GB。有效范围取决于磁盘类别：
  * `cloud`: [5, 2000]
  * `ephemeral`: [5, 1024]
  * `ephemeral_ssd`: [5, 800]
  * `cloud_efficiency`, `cloud_ssd`, `cloud_essd`: [20, 32768]

* `device` - (可选) 数据盘挂载点。有效值为 `/dev/xvdb` 到 `/dev/xvdz`。

* `category` - (可选) 数据盘类别。有效选项包括 `ephemeral_ssd`、`cloud_efficiency`、`cloud_ssd` 和 `cloud`。

* `snapshot_id` - (可选) 用于创建数据盘的快照ID。如果指定，则忽略 `size` 参数。

* `delete_with_instance` - (可选) 是否在释放ECS实例时删除数据盘。有效值为 `true` 或 `false`。默认值为 `true`。

* `encrypted` - (可选) 是否加密数据盘。有效值为 `true` 或 `false`。默认值为 `false`。

* `kms_key_id` - (可选) 用于加密数据盘的CMK ID。

* `name` - (可选) 数据盘名称。必须是2~128个字符，并且不能以 `http://` 或 `https://` 开头。

* `description` - (可选) 数据盘描述。必须是2~256个字符，并且不能以 `http://` 或 `https://` 开头。

* `auto_snapshot_policy_id` - (可选) 数据盘的自动快照策略ID。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 伸缩配置的ID。

* `status` - 伸缩组内伸缩配置的状态。可能的值：
  * `Active`: 表示该配置处于活动状态。
  * `Inactive`: 表示该配置处于非活动状态。

* `active` - 当前伸缩配置是否在指定的伸缩组中处于活动状态。

* `zone_id` - ECS实例所属的可用区ID。

* `scaling_configuration_name` - 伸缩配置名称。

* `internet_max_bandwidth_in` - 公网入方向的最大带宽，单位为Mbps。

* `substitute` - 当前配置设置为 `false` 时，自动成为活动状态并替换当前配置的伸缩配置。

* `ram_role_name` - ECS实例的RAM角色名称。

* `role_name` - ECS实例的RAM角色名称。

* `key_pair_name` - 用于登录ECS实例的密钥对名称。

* `is_outdated` - 是否使用过时的实例类型。

* `key_name` - 用于登录ECS实例的密钥对名称。