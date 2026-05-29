---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cs_kubernetes_node_pool"
sidebar_current: "docs-Alibabacloudstack-resource-cs-kubernetes-node-pool"
description: |-
  Provides a Alibabacloudstack resource to manage container kubernetes node pool.
---

# alibabacloudstack_cs_kubernetes_node_pool

This resource will help you to manage node pool in Kubernetes Cluster. 


## Example Usage

The managed cluster configuration,

```terraform

variable "name" {
	default = "tf-testAccNodePool-9633174"
}


data "alibabacloudstack_images" "default" {
  name_regex  = "^aliyun_.*"
  most_recent = true
  owners      = "system"
}

variable "existed_k8s_cluster_id" {
	default = ""
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
  tags = {
    common_test = "terraform"
	filter = var.name
  }
  lifecycle {
      ignore_changes = [
		secondary_cidr_blocks,
        tags
      ]
  }
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  vswitch_name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.1.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  lifecycle {
      ignore_changes = [
        tags
      ]
  }
}


resource "alibabacloudstack_ecs_securitygroup" "default" {
  name   = "${var.name}_sg"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

resource "alibabacloudstack_security_group_rule" "default" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alibabacloudstack_ecs_securitygroup.default.id}"
  	cidr_ip = "192.168.0.0/16"
}

resource "random_password" "password" {
	count            = 1
	length           = 12
	special          = true
	override_special = "!@#$^&*()_"
	min_lower        = 1
	min_upper        = 1
	min_numeric      = 1
}

data "alibabacloudstack_instance_types" "all" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  sorted_by         = "CPU"
}

data "alibabacloudstack_instance_types" "default" {
  count = 8  # Traverse 1-8 core CPU configurations

  availability_zone    = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count       = count.index + 1  # 1-8
  sorted_by            = "Memory"
}

locals {
  filtered_default = [for d in data.alibabacloudstack_instance_types.default : d if length(d.ids) > 0]
  fallback_all     = length(data.alibabacloudstack_instance_types.all.ids) > 0 ? data.alibabacloudstack_instance_types.all.ids : []
  
  default_instance_type_id = coalesce(
    try(local.filtered_default[0].ids[0], null),
    try(local.fallback_all[0], null),
    "no-available-instance-type"
  )
}

data "alibabacloudstack_cs_kubernetes_clusters" "default" {
	ids = var.existed_k8s_cluster_id == "" ? [] : [var.existed_k8s_cluster_id]
}

locals {
	create_count = length(data.alibabacloudstack_cs_kubernetes_clusters.default.ids) > 0 ? 0 : 1
}

resource "alibabacloudstack_cs_kubernetes" "default" {
	count						= local.create_count
	name						= var.name
	version						= "1.34.1-aliyun.1"
	os_type						= "linux"
	platform					= "AliyunLinux"
	num_of_nodes				= "1"
	master_count				= "3"
	master_vswitch_ids			= ["${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}"]
	master_instance_types		= ["${local.default_instance_type_id}","${local.default_instance_type_id}","${local.default_instance_type_id}"]
	master_disk_category		= "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
	vpc_id						= "${alibabacloudstack_vpc_vpc.default.id}"
	worker_instance_types		= ["${local.default_instance_type_id}"]
	worker_vswitch_ids			= ["${alibabacloudstack_vpc_vswitch.default.id}"]
	worker_disk_category		= "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
	password					= random_password.password.0.result
	pod_cidr					= "172.20.0.0/16"
	service_cidr				= "172.21.0.0/20"
	worker_disk_size			= "40"
	master_disk_size			= "40"
	slb_internet_enabled		= "true"
	security_group_id			= alibabacloudstack_ecs_securitygroup.default.id
	runtime	 {
		name	= "containerd"
		version	= "2.1.5"
	}
}

locals {
	k8s_cluster_id = length(data.alibabacloudstack_cs_kubernetes_clusters.default.ids) > 0 ? data.alibabacloudstack_cs_kubernetes_clusters.default.ids.0 : alibabacloudstack_cs_kubernetes.default.0.id
	k8s_cluster_name = length(data.alibabacloudstack_cs_kubernetes_clusters.default.ids) > 0 ? data.alibabacloudstack_cs_kubernetes_clusters.default.names.0 : alibabacloudstack_cs_kubernetes.default.0.name
}

resource "alibabacloudstack_ecs_keypair" "default" {
  key_name = var.name
}


resource "alibabacloudstack_cs_kubernetes_node_pool" "default" {
  name = "tf-testAccNodePool-9633174"
  cluster_id = "${local.k8s_cluster_id}"
  instance_types = [
                     "${local.default_instance_type_id}"
                   ]
  password = "${random_password.password.0.result}"
  tags = {
           Created = "TF"
           Foo = "Bar"
         }
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  system_disk_size = "40"
  data_disks {
    size = "100"
    category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  }
  
  node_count = "1"
  install_cloud_monitor = "false"
  vswitch_ids = [
                  "${alibabacloudstack_vpc_vswitch.default.id}"
                ]
}

Enable automatic scaling for the node pool. `scaling_config` is required.

```terraform

resource "alibabacloudstack_cs_kubernetes_node_pool" "autoscaling" {
  name = "tf-testAccNodePoolAuto-9719402"
  cluster_id = "${local.k8s_cluster_id}"
  vswitch_ids = [
                  "${alibabacloudstack_vpc_vswitch.default.id}"
                ]
  key_name = "${alibabacloudstack_ecs_keypair.default.key_name}"
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  install_cloud_monitor = "false"
  scaling_config {
    min_size = "1"
    max_size = "10"
    type = "cpu"
    is_bond_eip = "true"
    eip_bandwidth = "5"
  }
  
  instance_types = [
                     "${local.default_instance_type_id}"
                   ]
  platform = "Custom"
  image_id = "${data.alibabacloudstack_images.default.images.0.id}"
  system_disk_size = "40"
  scaling_policy = "release"
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
* `platform` - (Optional) The platform. One of `AliyunLinux`, `Windows`, `CentOS`, `WindowsCore`, `Custom`. If you select `Windows` or `WindowsCore`, the `passord` is required.
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

The tags example:
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
* `platform` - The platform. One of `AliyunLinux`, `Windows`, `CentOS`, `WindowsCore`, `Custom`.
* `instance_charge_type` - Node payment type. Valid values: `PostPaid`, `PrePaid`.
* `resource_group_id` - The ID of the resource group.
* `internet_charge_type` - The billing method for network usage. Valid values `PayByBandwidth` and `PayByTraffic`.
* `internet_max_bandwidth_out` - The maximum outbound bandwidth for the public network. Unit: Mbit/s. Valid values: 0 to 100.
* `spot_strategy` - The preemption policy for the pay-as-you-go instance. Valid value `SpotWithPriceLimit`.
* `node_count` - The worker node number of the node pool.