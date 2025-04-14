---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cs_kubernetes_node_pool"
sidebar_current: "docs-alibabacloudstack-resource-cs-kubernetes-node-pool"
description: |-
  Provides a Alibabacloudstack resource to manage container kubernetes node pool.
---

# alibabacloudstack_cs_kubernetes_node_pool

This resource will help you to manage node pool in Kubernetes Cluster. 



-> **NOTE:** From version 1.109.1, support managed node pools, but only for the professional managed clusters.

-> **NOTE:** From version 1.109.1, support remove node pool nodes.

-> **NOTE:** From version 1.111.0, support auto scaling node pool. For more information on how to use auto scaling node pools, see [Use Terraform to create an elastic node pool](https://help.aliyun.com/document_detail/197717.htm). With auto-scaling is enabled, the nodes in the node pool will be labeled with `k8s.aliyun.com=true` to prevent system pods such as coredns, metrics-servers from being scheduled to elastic nodes, and to prevent node shrinkage from causing business abnormalities.

-> **NOTE:** ACK adds a new RamRole (AliyunCSManagedAutoScalerRole) for the permission control of the node pool with auto-scaling enabled. If you are using a node pool with auto scaling, please click [AliyunCSManagedAutoScalerRole](https://ram.console.aliyun.com/role/authorization?request=%7B%22Services%22%3A%5B%7B%22Service%22%3A%22CS%22%2C%22Roles%22%3A%5B%7B%22RoleName%22%3A%22AliyunCSManagedAutoScalerRole%22%2C%22TemplateId%22%3A%22AliyunCSManagedAutoScalerRole%22%7D%5D%7D%5D%2C%22ReturnUrl%22%3A%22https%3A%2F%2Fcs.console.aliyun.com%2F%22%7D) to complete the authorization. 

-> **NOTE:** ACK adds a new RamRole(AliyunCSManagedNlcRole) for the permission control of the management node pool. If you use the management node pool, please click [AliyunCSManagedNlcRole](https://ram.console.aliyun.com/role/authorization?spm=5176.2020520152.0.0.387f16ddEOZxMv&request=%7B%22Services%22%3A%5B%7B%22Service%22%3A%22CS%22%2C%22Roles%22%3A%5B%7B%22RoleName%22%3A%22AliyunCSManagedNlcRole%22%2C%22TemplateId%22%3A%22AliyunCSManagedNlcRole%22%7D%5D%7D%5D%2C%22ReturnUrl%22%3A%22https%3A%2F%2Fcs.console.aliyun.com%2F%22%7D) to complete the authorization.

-> **NOTE:** From version 1.123.1, supports the creation of a node pool of spot instance.

-> **NOTE:** It is recommended to create a cluster with zero worker nodes, and then use a node pool to manage the cluster nodes. 

-> **NOTE:** From version 1.127.0, support for adding existing nodes to the node pool. In order to distinguish automatically created nodes, it is recommended that existing nodes be placed separately in a node pool for management. 

## Example Usage

The managed cluster configuration,

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

Create a node pool.

```terraform
resource "alibabacloudstack_cs_kubernetes_node_pool" "default" {
  name           = var.name
  cluster_id     = alibabacloudstack_cs_managed_kubernetes.default.0.id
  vswitch_ids    = [alibabacloudstack_vswitch.default.id]
  instance_types = [data.alibabacloudstack_instance_types.default.instance_types.0.id]

  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alibabacloudstack_key_pair.default.key_name

  # you need to specify the number of nodes in the node pool, which can be 0
  node_count = 1
}
```

Create a managed node pool. If you need to enable maintenance window, you need to set the maintenance window in `alibabacloudstack_cs_managed_kubernetes`.

```terraform
resource "alibabacloudstack_cs_kubernetes_node_pool" "default" {
  name                 = var.name
  cluster_id           = alibabacloudstack_cs_managed_kubernetes.default.0.id
  vswitch_ids          = [alibabacloudstack_vswitch.default.id]
  instance_types       = [data.alibabacloudstack_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40

  # only key_name is supported in the management node pool
  key_name = alibabacloudstack_key_pair.default.key_name

  # you need to specify the number of nodes in the node pool, which can be zero
  node_count = 1

  # management node pool configuration.
  management {
    auto_repair     = true
    auto_upgrade    = true
    surge           = 1
    max_unavailable = 1
  }

}
```

Enable automatic scaling for the node pool. `scaling_config` is required.

```terraform
resource "alibabacloudstack_cs_kubernetes_node_pool" "default" {
  name                 = var.name
  cluster_id           = alibabacloudstack_cs_managed_kubernetes.default.0.id
  vswitch_ids          = [alibabacloudstack_vswitch.default.id]
  instance_types       = [data.alibabacloudstack_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alibabacloudstack_key_pair.default.key_name

  # automatic scaling node pool configuration.
  # With auto-scaling is enabled, the nodes in the node pool will be labeled with `k8s.aliyun.com=true` to prevent system pods such as coredns, metrics-servers from being scheduled to elastic nodes, and to prevent node shrinkage from causing business abnormalities.
  scaling_config {
    min_size = 1
    max_size = 10
  }

}
```

Enable automatic scaling for managed node pool.

```terraform
resource "alibabacloudstack_cs_kubernetes_node_pool" "default" {
  name                 = var.name
  cluster_id           = alibabacloudstack_cs_managed_kubernetes.default.0.id
  vswitch_ids          = [alibabacloudstack_vswitch.default.id]
  instance_types       = [data.alibabacloudstack_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alibabacloudstack_key_pair.default.key_name
  # management node pool configuration.
  management {
    auto_repair     = true
    auto_upgrade    = true
    surge           = 1
    max_unavailable = 1
  }
  # enable auto-scaling
  scaling_config {
    min_size = 1
    max_size = 10
    type     = "cpu"
  }
  # Rely on auto-scaling configuration, please create auto-scaling configuration through alibabacloudstack_cs_autoscaling_config first.
  depends_on = [alibabacloudstack_cs_autoscaling_config.default]
}
```

Create a `PrePaid` node pool.
```terraform
resource "alibabacloudstack_cs_kubernetes_node_pool" "default" {
  name                 = var.name
  cluster_id           = alibabacloudstack_cs_managed_kubernetes.default.0.id
  vswitch_ids          = [alibabacloudstack_vswitch.default.id]
  instance_types       = [data.alibabacloudstack_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alibabacloudstack_key_pair.default.key_name
  # use PrePaid
  instance_charge_type = "PrePaid"
  period               = 1
  period_unit          = "Month"
  auto_renew           = true
  auto_renew_period    = 1

  # open cloud monitor
  install_cloud_monitor = true

  # enable auto-scaling
  scaling_config {
    min_size = 1
    max_size = 10
    type     = "cpu"
  }
}
```

Create a node pool with spot instance.
```terraform
resource "alibabacloudstack_cs_kubernetes_node_pool" "default" {
  name           = var.name
  cluster_id     = v_cs_managed_kubernetes.default.0.id
  vswitch_ids    = [alibabacloudstack_vswitch.default.id]
  instance_types = [data.alibabacloudstack_instance_types.default.instance_types.0.id]

  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alibabacloudstack_key_pair.default.key_name

  # you need to specify the number of nodes in the node pool, which can be 0
  node_count = 1

  # spot config
  spot_strategy = "SpotWithPriceLimit"
  spot_price_limit {
    instance_type = data.alibabacloudstack_instance_types.default.instance_types.0.id
    # Different instance types have different price caps
    price_limit = "0.70"
  }
}
```

Use Spot instances to create a node pool with auto-scaling enabled 
```terraform
resource "alibabacloudstack_cs_kubernetes_node_pool" "default" {
  name                 = var.name
  cluster_id           = alibabacloudstack_cs_managed_kubernetes.default.0.id
  vswitch_ids          = [alibabacloudstack_vswitch.default.id]
  instance_types       = [data.alibabacloudstack_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alibabacloudstack_key_pair.default.key_name

  # automatic scaling node pool configuration.
  scaling_config {
    min_size = 1
    max_size = 10
    type     = "spot"
  }
  # spot price config
  spot_strategy = "SpotWithPriceLimit"
  spot_price_limit {
    instance_type = data.alibabacloudstack_instance_types.default.instance_types.0.id
    price_limit   = "0.70"
  }
}
```

Create a node pool with platform as Windows 
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

  // if the instance platform is windows, the password is requered.
  password = var.password
  platform = "Windows"
  image_id = "${window_image_id}"
}
```

Add an existing node to the node pool

In order to distinguish automatically created nodes, it is recommended that existing nodes be placed separately in a node pool for management. 

```terraform
resource "alibabacloudstack_cs_kubernetes_node_pool" "default" {
  name                 = "existing-node"
  cluster_id           = alibabacloudstack_cs_managed_kubernetes.default.0.id
  vswitch_ids          = [alibabacloudstack_vswitch.default.id]
  instance_types       = [data.alibabacloudstack_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  instance_charge_type = "PostPaid"

  # add existing node to nodepool
  instances = ["instance_id_01", "instance_id_02", "instance_id_03"]
  # default is false
  format_disk = false
  # default is true
  keep_instance_name = true
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The id of kubernetes cluster.
* `name` - (Required) The name of node pool.
* `vswitch_ids` - (Required) The vswitches used by node pool workers.
* `instance_types` (Required) The instance type of worker node.
* `node_count` (Optional) The worker node number of the node pool. From version 1.111.0, `node_count` is not required.
* `password` - (Required, Sensitive) The password of ssh login cluster node. You have to specify one of `password` `key_name` `kms_encrypted_password` fields.
* `key_name` - (Required) The keypair of ssh login cluster node, you have to create it first. You have to specify one of `password` `key_name` `kms_encrypted_password` fields. Only `key_name` is supported in the management node pool.
* `kms_encrypted_password` - (Required) An KMS encrypts password used to a cs kubernetes. You have to specify one of `password` `key_name` `kms_encrypted_password` fields.
* `system_disk_category` - (Optional) The system disk category of worker node. Its valid value are `cloud_ssd` and `cloud_efficiency`. Default to `cloud_efficiency`.
* `system_disk_size` - (Optional) The system disk category of worker node. Its valid value range [40~500] in GB. Default to `120`.
* `data_disks` - (Optional) The data disk configurations of worker nodes, such as the disk type and disk size. 
  * `category` - The type of the data disks. Valid values:`cloud`, `cloud_efficiency`, `cloud_ssd` and `cloud_essd`.
  * `size` - The size of a data disk, Its valid value range [40~32768] in GB. Default to `40`.
  * `encrypted` - Specifies whether to encrypt data disks. Valid values: true and false. Default to `false`.
* `platform` - (Optional) The platform. One of `AliyunLinux`, `Windows`, `CentOS`, `WindowsCore`. If you select `Windows` or `WindowsCore`, the `passord` is required.
* `image_id` - (Optional) Custom Image support. Must based on CentOS7 or AliyunLinux2.
* `node_name_mode` - (Optional) Each node name consists of a prefix, an IP substring, and a suffix. For example "customized,aliyun.com,5,test", if the node IP address is 192.168.0.55, the prefix is aliyun.com, IP substring length is 5, and the suffix is test, the node name will be aliyun.com00055test.
* `user_data` - (Optional) Windows instances support batch and PowerShell scripts. If your script file is larger than 1 KB, we recommend that you upload the script to Object Storage Service (OSS) and pull it through the internal endpoint of your OSS bucket.
* `tags` - (Optional) A Map of tags to assign to the resource. It will be applied for ECS instances finally.
* `labels` - (Optional) A List of Kubernetes labels to assign to the nodes . Only labels that are applied with the ACK API are managed by this argument.
  * `key` - The label key.
  * `value` - The label value.
* `taints` - (Optional) A List of Kubernetes taints to assign to the nodes.
  * `effect` - (Optional) The scheduling policy.
  * `key` - (Required) The key of a taint.
  * `value` - (Optional) The value of a taint.
* `scaling_policy` - (Optional) The scaling mode. Valid values: `release`, `recycle`, default is `release`. Standard mode(release): Create and release ECS instances based on requests.Swift mode(recycle): Create, stop, and restart ECS instances based on needs. New ECS instances are only created when no stopped ECS instance is avalible. This mode further accelerates the scaling process. Apart from ECS instances that use local storage, when an ECS instance is stopped, you are only chatged for storage space.
* `scaling_config` - (Optional) Auto scaling node pool configuration. With auto-scaling is enabled, the nodes in the node pool will be labeled with `k8s.aliyun.com=true` to prevent system pods such as coredns, metrics-servers from being scheduled to elastic nodes, and to prevent node shrinkage from causing business abnormalities.
  * `min_size` - (Required) Min number of instances in a auto scaling group, its valid value range [0~1000].
  * `max_size` - (Required) Max number of instances in a auto scaling group, its valid value range [0~1000]. `max_size` has to be greater than `min_size`.
  * `type` - (Optional) Instance classification, not required. Vaild value: `cpu`, `gpu`, `gpushare` and `spot`. Default: `cpu`. The actual instance type is determined by `instance_types`.
  * `is_bond_eip` - (Optional) Whether to bind EIP for an instance. Default: `false`.
  * `eip_internet_charge_type` - (Optional) EIP billing type. `PayByBandwidth`: Charged at fixed bandwidth. `PayByTraffic`: Billed as used traffic. Default: `PayByBandwidth`. Conflict with `internet_charge_type`, EIP and public network IP can only choose one. 
  * `eip_bandwidth` - (Optional) Peak EIP bandwidth. Its valid value range [1~500] in Mbps. Default to `5`.
* `system_disk_performance_level` - (Optional) The performance level (PL) of the system disk that you want to use for the node. This parameter takes effect only for ESSDs. Its valid value is one of {"PL0", "PL1", "PL2", "PL3"}.
* `instance_charge_type`- (Optional) Node payment type. Valid values: `PostPaid`, `PrePaid`, default is `PostPaid`. If value is `PrePaid`, the arguments `period`, `period_unit`, `auto_renew` and `auto_renew_period` are required.
* `period`- (Optional) Node payment period. Its valid value is one of {1, 2, 3, 6, 12, 24, 36, 48, 60}.
* `period_unit`- (Optional) Node payment period unit, valid value: `Month`. Default is `Month`.
* `auto_renew`- (Optional) Enable Node payment auto-renew, default is `false`.
* `auto_renew_period`- (Optional) Node payment auto-renew period, one of `1`, `2`, `3`,`6`, `12`.
* `install_cloud_monitor`- (Optional) Install the cloud monitoring plug-in on the node, and you can view the monitoring information of the instance through the cloud monitoring console. Default is `true`.
* `unschedulable`- (Optional) Set the newly added node as unschedulable. If you want to open the scheduling option, you can open it in the node list of the console. If you are using an auto-scaling node pool, the setting will not take effect. Default is `false`.
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group,by default these cloud resources are automatically assigned to the default resource group.
* `internet_charge_type` - (Optional) The billing method for network usage. Valid values `PayByBandwidth` and `PayByTraffic`. Conflict with `eip_internet_charge_type`, EIP and public network IP can only choose one. 
* `internet_max_bandwidth_out` - (Optional) The maximum outbound bandwidth for the public network. Unit: Mbit/s. Valid values: 0 to 100.
* `spot_strategy` - (Optional) The preemption policy for the pay-as-you-go instance. This parameter takes effect only when `instance_charge_type` is set to `PostPaid`. Valid value `SpotWithPriceLimit`.
* `spot_price_limit` - (Optional) The maximum hourly price of the instance. This parameter takes effect only when `spot_strategy` is set to `SpotWithPriceLimit`. A maximum of three decimal places are allowed.
  * `instance_type` - (Optional) Spot instance type.
  * `price_limit` - (Optional) The maximum hourly price of the spot instance.
* `instances` - (Optional) The instance list. Add existing nodes under the same cluster VPC to the node pool. 
* `keep_instance_name` - (Optional) Add an existing instance to the node pool, whether to keep the original instance name. It is recommended to set to `true`.
* `format_disk` - (Optional,) After you select this, if data disks have been attached to the specified ECS instances and the file system of the last data disk is uninitialized, the system automatically formats the last data disk to ext4 and mounts the data disk to /var/lib/docker and /var/lib/kubelet. The original data on the disk will be cleared. Make sure that you back up data in advance. If no data disk is mounted on the ECS instance, no new data disk will be purchased. Default is `false`.
* `security_group_id` - (Optional ) The ID of security group where the current cluster worker node is located.
* `system_disk_size` - (Optional) The system disk size of worker node. Its valid value range [20~32768] in GB. Default to `40`.
* `node_count` - (Optional) The worker node number of the node pool. From version 1.111.0, `node_count` is not required.

#### tags

The tags example：
```
tags {
  "key-a" = "value-a"
  "key-b" = "value-b"
  "env"   = "prod"
}
```

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the node pool, format cluster_id:nodepool_id.
* `cluster_id` - The cluster id.
* `name` - The name of the nodepool.
* `vswitch_ids` - The vswitches used by node pool workers.
* `image_id` - The image used by node pool workers.
* `security_group_id` - The ID of security group where the current cluster worker node is located.
* `scaling_group_id` - (Available in 1.105.0+) Id of the Scaling Group.
* `system_disk_performance_level` - The performance level (PL) of the system disk that you want to use for the node. This parameter takes effect only for ESSDs. Its valid value is one of {"PL0", "PL1", "PL2", "PL3"}.
* `platform` - The platform. One of `AliyunLinux`, `Windows`, `CentOS`, `WindowsCore`.
* `instance_charge_type` - Node payment type. Valid values: `PostPaid`, `PrePaid`.
* `resource_group_id` - The ID of the resource group.
* `internet_charge_type` - The billing method for network usage. Valid values `PayByBandwidth` and `PayByTraffic`.
* `internet_max_bandwidth_out` - The maximum outbound bandwidth for the public network. Unit: Mbit/s. Valid values: 0 to 100.
* `spot_strategy` - The preemption policy for the pay-as-you-go instance. Valid value `SpotWithPriceLimit`.
* `node_count` - The worker node number of the node pool.