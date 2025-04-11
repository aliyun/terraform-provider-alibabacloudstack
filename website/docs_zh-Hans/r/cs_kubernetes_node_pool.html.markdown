---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cs_kubernetes_node_pool"
sidebar_current: "docs-alibabacloudstack-resource-cs-kubernetes-node-pool"
description: |-
  编排 Kubernetes 集群中的节点池
---

# alibabacloudstack_cs_kubernetes_node_pool

使用Provider配置的凭证在指定的资源集下编排Kubernetes集群中的节点池。


## 示例用法

托管集群配置，

```terraform
variable "name" {
  default = "tf-test"
}
variable "password" {
}
data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
}
data "alibabacloudstack_instance_types" "default" {
  availability_zone    = data.alibabacloudstack_zones.default.zones.0.id
  cpu_core_count       = 2
  memory_size          = 4
  kubernetes_node_role = "Worker"
}
resource "alibabacloudstack_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.1.0.0/21"
}
resource "alibabacloudstack_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alibabacloudstack_vpc.default.id
  cidr_block   = "10.1.1.0/24"
  zone_id      = data.alibabacloudstack_zones.default.zones.0.id
}
resource "alibabacloudstack_key_pair" "default" {
  key_pair_name = var.name
}
resource "alibabacloudstack_cs_managed_kubernetes" "default" {
  name                         = var.name
  count                        = 1
  cluster_spec                 = "ack.pro.small"
  is_enterprise_security_group = true
  worker_number                = 2
  password                     = var.password
  pod_cidr                     = "172.20.0.0/16"
  service_cidr                 = "172.21.0.0/20"
  worker_vswitch_ids           = [alibabacloudstack_vswitch.default.id]
  worker_instance_types        = [data.alibabacloudstack_instance_types.default.instance_types.0.id]
}
```

创建一个节点池。

```terraform
resource "alibabacloudstack_cs_kubernetes_node_pool" "default" {
  name           = var.name
  cluster_id     = alibabacloudstack_cs_managed_kubernetes.default.0.id
  vswitch_ids    = [alibabacloudstack_vswitch.default.id]
  instance_types = [data.alibabacloudstack_instance_types.default.instance_types.0.id]

  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alibabacloudstack_key_pair.default.key_name

  # 您需要指定节点池中的节点数，可以为 0
  node_count = 1
}
```

创建一个托管节点池。如果需要启用维护窗口，您需要在 `alibabacloudstack_cs_managed_kubernetes` 中设置维护窗口。

```terraform
resource "alibabacloudstack_cs_kubernetes_node_pool" "default" {
  name                 = var.name
  cluster_id           = alibabacloudstack_cs_managed_kubernetes.default.0.id
  vswitch_ids          = [alibabacloudstack_vswitch.default.id]
  instance_types       = [data.alibabacloudstack_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40

  # 只有托管节点池支持 key_name
  key_name = alibabacloudstack_key_pair.default.key_name

  # 托管节点池配置。
  management {
    auto_repair     = true
    auto_upgrade    = true
    surge           = 1
    max_unavailable = 1
  }

}
```

为节点池启用自动扩展。`scaling_config` 是必填的。

```terraform
resource "alibabacloudstack_cs_kubernetes_node_pool" "default" {
  name                 = var.name
  cluster_id           = alibabacloudstack_cs_managed_kubernetes.default.0.id
  vswitch_ids          = [alibabacloudstack_vswitch.default.id]
  instance_types       = [data.alibabacloudstack_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alibabacloudstack_key_pair.default.key_name

  # 自动扩展节点池配置。
  # 启用自动扩展后，节点池中的节点将被标记为 `k8s.aliyun.com=true`，以防止诸如 coredns、metrics-servers 等系统 Pod 被调度到弹性节点上，并防止节点缩减导致业务异常。
  scaling_config {
    min_size = 1
    max_size = 10
  }

}
```

为托管节点池启用自动扩展。

```terraform
resource "alibabacloudstack_cs_kubernetes_node_pool" "default" {
  name                 = var.name
  cluster_id           = alibabacloudstack_cs_managed_kubernetes.default.0.id
  vswitch_ids          = [alibabacloudstack_vswitch.default.id]
  instance_types       = [data.alibabacloudstack_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alibabacloudstack_key_pair.default.key_name
  # 托管节点池配置。
  management {
    auto_repair     = true
    auto_upgrade    = true
    surge           = 1
    max_unavailable = 1
  }
  # 启用自动扩展
  scaling_config {
    min_size = 1
    max_size = 10
    type     = "cpu"
  }
  # 依赖于自动扩展配置，请先通过 alibabacloudstack_cs_autoscaling_config 创建自动扩展配置。
  depends_on = [alibabacloudstack_cs_autoscaling_config.default]
}
```

创建一个 `PrePaid` 节点池。
```terraform
resource "alibabacloudstack_cs_kubernetes_node_pool" "default" {
  name                 = var.name
  cluster_id           = alibabacloudstack_cs_managed_kubernetes.default.0.id
  vswitch_ids          = [alibabacloudstack_vswitch.default.id]
  instance_types       = [data.alibabacloudstack_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alibabacloudstack_key_pair.default.key_name
  # 使用 PrePaid
  instance_charge_type = "PrePaid"
  period               = 1
  period_unit          = "Month"
  auto_renew           = true
  auto_renew_period    = 1

  # 开启云监控
  install_cloud_monitor = true

  # 启用自动扩展
  scaling_config {
    min_size = 1
    max_size = 10
    type     = "cpu"
  }
}
```

创建一个抢占式实例节点池。
```terraform
resource "alibabacloudstack_cs_kubernetes_node_pool" "default" {
  name           = var.name
  cluster_id     = v_cs_managed_kubernetes.default.0.id
  vswitch_ids    = [alibabacloudstack_vswitch.default.id]
  instance_types = [data.alibabacloudstack_instance_types.default.instance_types.0.id]

  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alibabacloudstack_key_pair.default.key_name

  # 您需要指定节点池中的节点数，可以为 0
  node_count = 1

  # 抢占式实例配置
  spot_strategy = "SpotWithPriceLimit"
  spot_price_limit {
    instance_type = data.alibabacloudstack_instance_types.default.instance_types.0.id
    # 不同实例类型有不同的价格上限
    price_limit = "0.70"
  }
}
```

使用抢占式实例创建启用自动扩展的节点池 
```terraform
resource "alibabacloudstack_cs_kubernetes_node_pool" "default" {
  name                 = var.name
  cluster_id           = alibabacloudstack_cs_managed_kubernetes.default.0.id
  vswitch_ids          = [alibabacloudstack_vswitch.default.id]
  instance_types       = [data.alibabacloudstack_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alibabacloudstack_key_pair.default.key_name

  # 自动扩展节点池配置。
  scaling_config {
    min_size = 1
    max_size = 10
    type     = "spot"
  }
  # 抢占式实例价格配置
  spot_strategy = "SpotWithPriceLimit"
  spot_price_limit {
    instance_type = data.alibabacloudstack_instance_types.default.instance_types.0.id
    price_limit   = "0.70"
  }
}
```

创建一个平台为 Windows 的节点池
```terraform

variable "password" {
}

resource "alibabacloudstack_cs_kubernetes_node_pool" "default" {
  name                 = "windows-np"
  cluster_id           = alibabacloudstack_cs_managed_kubernetes.default.0.id
  vswitch_ids          = [alibabacloudstack_vswitch.default.id]
  instance_types       = [data.alibabacloudstack_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  instance_charge_type = "PostPaid"
  node_count           = 1

  // 如果实例平台为 windows，则密码是必填项。
  password = var.password
  platform = "Windows"
  image_id = "${window_image_id}"
}
```

将现有节点添加到节点池

为了区分自动创建的节点，建议将现有节点单独放在一个节点池中进行管理。 

```terraform
resource "alibabacloudstack_cs_kubernetes_node_pool" "default" {
  name                 = "existing-node"
  cluster_id           = alibabacloudstack_cs_managed_kubernetes.default.0.id
  vswitch_ids          = [alibabacloudstack_vswitch.default.id]
  instance_types       = [data.alibabacloudstack_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  instance_charge_type = "PostPaid"

  # 将现有节点添加到节点池
  instances = ["instance_id_01", "instance_id_02", "instance_id_03"]
  # 默认值为 false
  format_disk = false
  # 默认值为 true
  keep_instance_name = true
}
```

## 参数说明

支持以下参数：

* `cluster_id` - (必填) Kubernetes 集群的 ID。
* `name` - (必填) 节点池的名称。
* `vswitch_ids` - (必填) 节点池工作者使用的交换机。
* `instance_types` (必填) 工作节点的实例类型。
* `node_count` (可选) 节点池的工作节点数量。
* `password` - (必填, 敏感) SSH 登录集群节点的密码。您必须指定 `password`、`key_name` 或 `kms_encrypted_password` 字段之一。
* `key_name` - (必填) SSH 登录集群节点的密钥对，您必须先创建它。您必须指定 `password`、`key_name` 或 `kms_encrypted_password` 字段之一。只有托管节点池支持 `key_name`。
* `kms_encrypted_password` - (必填) 用于加密 cs kubernetes 密码的 KMS。您必须指定 `password`、`key_name` 或 `kms_encrypted_password` 字段之一。
* `system_disk_category` - (可选) 工作节点的系统磁盘类别。其有效值为 `cloud_ssd` 和 `cloud_efficiency`。默认为 `cloud_efficiency`。
* `system_disk_size` - (可选) 工作节点的系统磁盘大小。其有效值范围 [40~500] GB。默认为 `120`。
* `data_disks` - (可选) 工作节点的数据磁盘配置，例如磁盘类型和磁盘大小。 
  * `category` - 数据磁盘的类型。有效值：`cloud`, `cloud_efficiency`, `cloud_ssd` 和 `cloud_essd`。
  * `size` - 数据磁盘的大小，其有效值范围 [40~32768] GB。默认为 `40`。
  * `encrypted` - 指定是否加密数据磁盘。有效值：true 和 false。默认为 `false`。
* `platform` - (可选) 平台。其中之一 `AliyunLinux`, `Windows`, `CentOS`, `WindowsCore`。如果您选择 `Windows` 或 `WindowsCore`，则需要提供 `password`。
* `image_id` - (可选) 自定义镜像支持。必须基于 CentOS7 或 AliyunLinux2。
* `node_name_mode` - (可选) 每个节点名由前缀、IP 子串和后缀组成。例如 "customized,aliyun.com,5,test"，如果节点 IP 地址是 192.168.0.55，前缀是 aliyun.com，IP 子串长度是 5，后缀是 test，那么节点名将是 aliyun.com00055test。
* `user_data` - (可选) Windows 实例支持批处理和 PowerShell 脚本。如果您的脚本文件大于 1 KB，我们建议您将脚本上传到对象存储服务 (OSS)，并通过您的 OSS 存储桶的内部端点拉取它。
* `tags` - (可选) 分配给资源的标签映射。最终将应用于 ECS 实例。
* `labels` - (可选) 分配给节点的 Kubernetes 标签列表。仅通过 ACK API 庻用的标签由该参数管理。
  * `key` - 标签键。
  * `value` - 标签值。
* `taints` - (可选) 分配给节点的 Kubernetes 污点列表。
  * `effect` - (可选) 调度策略。
  * `key` - (必填) 污点的键。
  * `value` - (可选) 污点的值。
* `scaling_policy` - (可选) 缩放模式。有效值：`release`, `recycle`，默认为 `release`。标准模式(release): 基于请求创建和释放 ECS 实例。Swift 模式(recycle): 基于需求创建、停止和重新启动 ECS 实例。当没有可用的已停止 ECS 实例时，才会创建新的 ECS 实例。此模式进一步加速了缩放过程。除了使用本地存储的 ECS 实例外，当 ECS 实例停止时，您只需支付存储空间费用。
* `scaling_config` - (可选) 自动扩展节点池配置。启用自动扩展后，节点池中的节点将被标记为 `k8s.aliyun.com=true`，以防止诸如 coredns、metrics-servers 等系统 Pod 被调度到弹性节点上，并防止节点缩减导致业务异常。
  * `min_size` - (必填) 自动伸缩组中的最小实例数，其有效值范围 [0~1000]。
  * `max_size` - (必填) 自动伸缩组中的最大实例数，其有效值范围 [0~1000]。`max_size` 必须大于 `min_size`。
  * `type` - (可选) 实例分类，非必填。有效值：`cpu`, `gpu`, `gpushare` 和 `spot`。默认值为 `cpu`。实际实例类型由 `instance_types` 决定。
  * `is_bond_eip` - (可选) 是否为实例绑定 EIP，默认为 `false`。
  * `eip_internet_charge_type` - (可选) EIP 计费类型。`PayByBandwidth`: 按固定带宽计费。`PayByTraffic`: 按使用流量计费。默认为 `PayByBandwidth`。与 `internet_charge_type` 冲突，EIP 和公网 IP 只能选择一个。 
  * `eip_bandwidth` - (可选) 峰值 EIP 带宽。其有效值范围 [1~500] Mbps。默认为 `5`。
* `system_disk_performance_level` - (可选) 节点要使用的系统磁盘的性能级别 (PL)。此参数仅对 ESSD 生效。其有效值为 {"PL0", "PL1", "PL2", "PL3"}。
* `instance_charge_type`- (可选) 节点支付类型。有效值：`PostPaid`, `PrePaid`，默认为 `PostPaid`。如果值为 `PrePaid`，参数 `period`, `period_unit`, `auto_renew` 和 `auto_renew_period` 是必填的。
* `period`- (可选) 节点支付周期。其有效值为 {1, 2, 3, 6, 12, 24, 36, 48, 60}。
* `period_unit`- (可选) 节点支付周期单位，有效值：`Month`。默认为 `Month`。
* `auto_renew`- (可选) 启用节点支付自动续订，默认为 `false`。
* `auto_renew_period`- (可选) 节点支付自动续订周期，其有效值为 `1`, `2`, `3`,`6`, `12`。
* `install_cloud_monitor`- (可选) 在节点上安装云监控插件，可以通过云监控控制台查看实例的监控信息。默认为 `true`。
* `unschedulable`- (可选) 将新添加的节点设置为不可调度。如果要在控制台的节点列表中打开调度选项，可以打开它。如果使用的是自动扩展节点池，该设置不会生效。默认为 `false`。
* `resource_group_id` - (可选, 强制更改) 资源组 ID，默认情况下这些云资源会被自动分配到默认资源组。
* `internet_charge_type` - (可选) 网络使用计费方式。有效值 `PayByBandwidth` 和 `PayByTraffic`。与 `eip_internet_charge_type` 冲突，EIP 和公网 IP 只能选择一个。 
* `internet_max_bandwidth_out` - (可选) 公网最大出带宽。单位：Mbit/s。有效值：0 到 100。
* `spot_strategy` - (可选) 按量付费实例的抢占策略。此参数仅在 `instance_charge_type` 设置为 `PostPaid` 时生效。有效值 `SpotWithPriceLimit`。
* `spot_price_limit` - (可选) 实例的最大每小时价格。此参数仅在 `spot_strategy` 设置为 `SpotWithPriceLimit` 时生效。最多允许三位小数。
  * `instance_type` - (可选) 抢占式实例类型。
  * `price_limit` - (可选) 抢占式实例的最大每小时价格。
* `instances` - (可选) 实例列表。将同一集群 VPC 下的现有节点添加到节点池中。 
* `keep_instance_name` - (可选) 将现有实例添加到节点池时，是否保留原始实例名称。建议设置为 `true`。
* `format_disk` - (可选,) 选择此项后，如果指定的 ECS 实例已挂载数据盘且最后一个数据盘的文件系统未初始化，系统将自动将最后一个数据盘格式化为 ext4 并挂载到 /var/lib/docker 和 /var/lib/kubelet。磁盘上的原始数据将被清除。请确保提前备份数据。如果 ECS 实例未挂载数据盘，则不会购买新的数据盘。默认为 `false`。
* `security_group_id` - (可选 ) 当前集群工作节点所在的安全组 ID。
* `system_disk_size` - (可选) 工作节点的系统磁盘大小。其有效值范围 [20~32768] GB。默认为 `40`。

#### tags

标签示例：
```
tags {
  "key-a" = "value-a"
  "key-b" = "value-b"
  "env"   = "prod"
}
```

## 属性说明

导出以下属性：

* `id` - 节点池的 ID，格式为 cluster_id:nodepool_id。
* `cluster_id` - 集群 id。
* `name` - 节点池的名称。
* `vswitch_ids` - 节点池工作者使用的交换机。
* `image_id` - 节点池工作者使用的镜像。
* `security_group_id` - 当前集群工作节点所在的安全组 ID。
* `scaling_group_id` - 伸缩组 ID。
* `system_disk_performance_level` - 节点要使用的系统磁盘的性能级别 (PL)。此参数仅对 ESSD 生效。其有效值为 {"PL0", "PL1", "PL2", "PL3"}。
* `platform` - 平台。其中之一 `AliyunLinux`, `Windows`, `CentOS`, `WindowsCore`。
* `instance_charge_type` - 节点支付类型。有效值：`PostPaid`, `PrePaid`。
* `resource_group_id` - 资源组 ID。
* `internet_charge_type` - 网络使用计费方式。有效值 `PayByBandwidth` 和 `PayByTraffic`。
* `internet_max_bandwidth_out` - 公网最大出带宽。单位：Mbit/s。有效值：0 到 100。
* `spot_strategy` - 按量付费实例的抢占策略。有效值 `SpotWithPriceLimit`。
* `node_count` - 节点池的工作节点数量。