---
subcategory: "ACK"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ack_cluster"
sidebar_current: "docs-Alibabacloudstack-ack-cluster"
description: |- 
  Provides a ack Cluster resource.
---

# alibabacloudstack_ack_cluster
-> **NOTE:** Alias name has: `alibabacloudstack_cs_kubernetes`

Provides a ack Cluster resource.

## Example Usage

```hcl
variable "name" {
	default = "tf-testAccCsK8sConfigBasic3740595"
}

data "alibabacloudstack_zones" default {
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
  cidr_block = "172.16.0.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
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
  	cidr_ip = "172.16.0.0/24"
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
  name_regex  = "^anolisos_"
  most_recent = true
  owners      = "system"
}

variable "runtime" {
 default     = [
		{
			name    = "containerd"
			version = "1.6.28"
		}
	]
}

variable "new_nat_gateway" {
  description = "Whether to create a new nat gateway. In this template, a new nat gateway will create a nat gateway, eip and server snat entries."
  default     = "ture"
}

variable "node_cidr_mask" {
  description = "The node cidr block to specific how many pods can run on single node."
  default     = 26
}

variable "enable_ssh" {
  description = "Enable login to the node through SSH."
  default     = true
}

variable "password" {
  description = "The password of ECS instance."
}

variable "worker_number" {
  description = "The number of worker nodes in kubernetes cluster."
  default     = 3
}

variable "pod_cidr" {
  description = "The kubernetes pod cidr block. It cannot be equals to vpc's or vswitch's and cannot be in them."
  default     = "172.24.0.0/16"
}

variable "service_cidr" {
  description = "The kubernetes service cidr block. It cannot be equals to vpc's or vswitch's or pod's and cannot be in them."
  default     = "172.25.0.0/16"
}

resource "alibabacloudstack_cs_kubernetes" "k8s" {
  name                        = var.name
  vpc_id                      = alibabacloudstack_vpc_vpc.default.id
  worker_vswitch_ids         = [alibabacloudstack_vpc_vswitch.default.id]
  master_vswitch_ids         = [alibabacloudstack_vpc_vswitch.default.id, alibabacloudstack_vpc_vswitch.default.id, alibabacloudstack_vpc_vswitch.default.id]
  enable_ssh                 = var.enable_ssh
  is_enterprise_security_group = true
  worker_instance_types      = [local.default_instance_type_id]
  master_instance_types      = [local.default_instance_type_id, local.default_instance_type_id, local.default_instance_type_id]
  master_disk_size           = 40
  service_cidr               = var.service_cidr
  addons {
    name = "flannel"
  }
  addons {
    name = "csi-plugin"
  }
  addons {
    name = "csi-provisioner"
  }
  addons {
    name = "nginx-ingress-controller"
  }
  
  num_of_nodes               = var.worker_number
  version                    = "1.30.1-aliyun.1"
  delete_protection          = false
  master_count               = 3
  worker_disk_category       = data.alibabacloudstack_zones.default.zones[0].available_disk_categories[0]
  timeout_mins               = 60
  pod_cidr                  = var.pod_cidr
  worker_disk_size          = 40
  runtime {
    name = "containerd"
    version = "1.6.28"
  }
  
  proxy_mode                 = "ipvs"
  node_cidr_mask             = var.node_cidr_mask
  image_id                   = data.alibabacloudstack_images.default.images[0].id
  slb_internet_enabled      = false
  password                  = var.password
  os_type                   = "linux"
  platform                  = "AliyunLinux"
  new_nat_gateway           = false
  master_disk_category      = data.alibabacloudstack_zones.default.zones[0].available_disk_categories[0]
}
```

## Argument Reference

The following arguments are supported:

### Global params

* `name` - (Optional) The Kubernetes cluster's name. It must be unique within one AlibabaCloudStack account.
* `vpc_id` - (Required) The ID of the VPC where the current cluster is located.
* `worker_vswitch_ids` - (Required) The vSwitches used by workers. You can specify one or more vSwitches.
* `master_vswitch_ids` - (Required) The vSwitches used by masters. You can specify three or five vSwitches based on the number of master nodes.
* `version` - (Optional) Desired Kubernetes version. If you do not specify a value, the latest available version at resource creation is used and no upgrades will occur except you set a higher version number. Downgrades are not supported by ACK.
* `password` - (Required, Sensitive) The password for SSH login to cluster nodes. You must specify one of `password`, `key_name`, or `kms_encrypted_password`.
* `kms_encrypted_password` - (Optional) An KMS encrypted password used to decrypt the password before creating or updating a CS Kubernetes with `kms_encrypted_password`.
* `enable_ssh` - (Optional) Enable login to the node through SSH. Default is `false`.
* `cpu_policy` - (Optional) kubelet CPU policy. Options: `static` | `none`. Default is `none`.
* `proxy_mode` - (Optional) Proxy mode for kube-proxy. Options: `iptables` | `ipvs`. Default is `ipvs`.
* `user_data` - (Optional) User-defined data for instances. Windows instances support batch and PowerShell scripts. If your script file is larger than 1 KB, we recommend that you upload the script to Object Storage Service (OSS) and pull it through the internal endpoint of your OSS bucket.
* `instances` - (Optional) A list of instances that can be attached as worker nodes in the same VPC.
* `os_type` - (Optional) The operating system of the nodes that run pods, its valid values are either `Linux` or `Windows`. Default is `Linux`.
* `platform` - (Optional) The architecture of the nodes that run pods. Default is `CentOS`.
* `security_group_id` - (Optional) The ID of the security group to which the ECSs in the cluster belong. Conflicts with `is_enterprise_security_group`.
* `is_enterprise_security_group` - (Optional) Specifies whether an advanced security group is automatically created. You must set either the `security_group_id` or `is_enterprise_security_group` parameter.
* `runtime` - (Optional) The platform on which the clusters are going to run.
  * `name` - (Optional) Name of the runtime platform.
  * `version` - (Optional) Version of the runtime platform.
* `tags` - (Optional) A mapping of tags to assign to the resource.
  - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
  - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.
* `keep_instance_name` - (Optional) Add an existing instance to the node pool, whether to keep the original instance name. It is recommended to set to `true`.
* `format_disk` - (Optional) After selecting this option, if data disks have been attached to the specified ECS instances and the file system of the last data disk is uninitialized, the system automatically formats the last data disk to ext4 and mounts the data disk to `/var/lib/docker` and `/var/lib/kubelet`. The original data on the disk will be cleared. Make sure that you back up data in advance. If no data disk is mounted on the ECS instance, no new data disk will be purchased. Default is `false`.
* `image_id` - (Optional) Custom Image support. Must be based on CentOS7 or AliyunLinux2.
* `timeout_mins` - (Optional) Backend service time-out time; unit: minute. Default is 60.
* `delete_protection` - (Optional) Whether the instance should have delete protection.
* `kms_encryption_context` - (Optional) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating a CS Kubernetes with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `addons` - (Optional) The add-ons you want to install in the cluster.
  * `name` - (Optional) Name of the ACK add-on. The name must match one of the names returned by DescribeAddons.
  * `config` - (Optional) The ACK add-on configurations. For more config information, see [cs_kubernetes_addon_metadata](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/data-sources/cs_kubernetes_addon_metadata).
* `cloud_monitor_flags` - (Optional) Whether to install cloud monitoring plugin.
* `node_port_range` - (Optional) Specifies the range of ports used for NodePort services.

### Network

* `pod_cidr` - (Optional) [Flannel Specific] The CIDR block for the pod network when using Flannel.
* `pod_vswitch_ids` - (Optional) [Terway Specific] The vSwitches for the pod network when using Terway. Be careful, the `pod_vswitch_ids` cannot equal to `worker_vswitch_ids` or `master_vswitch_ids` but must be in the same availability zones.
* `new_nat_gateway` - (Optional) Whether to create a new NAT gateway while creating the Kubernetes cluster. Default is `true`.
* `service_cidr` - (Optional) The CIDR block for the service network. It cannot be duplicated with the VPC CIDR and CIDR used by Kubernetes cluster in VPC, cannot be modified after creation.
* `node_cidr_mask` - (Optional) The node CIDR block to specify how many pods can run on a single node. Valid range: 24-28. Default is 24.
* `slb_internet_enabled` - (Optional) Whether to create an internet load balancer for the API Server. Default is `true`.

If you want to use `Terway` as CNI network plugin, you need to specify the `pod_vswitch_ids` field and add-ons such as `csi-plugin`, `csi-provisioner`, `logtail-ds`, and `nginx-ingress-controller`.  
If you want to use `Flannel` as CNI network plugin, you need to specify the `pod_cidr` field and add-ons such as `flannel`.

### Master Params

* `master_count` - (Optional) Number of master nodes. Default is 3.
* `master_disk_category` - (Optional) The system disk category of master nodes. Valid values are `cloud_ssd` and `cloud_efficiency`. Default is `cloud_efficiency`.
* `master_disk_size` - (Optional) The system disk size of master nodes. Valid range: [20~500] in GB. Default is 20.
* `master_instance_types` - (Required) The instance type of master nodes. Specify one type for single AZ Cluster, three types for MultiAZ Cluster.
* `master_system_disk_performance_level` - (Optional) Performance level of the master node's system disk.

### Worker Params

* `num_of_nodes` - (Required) The worker node number of the Kubernetes cluster. Default is 3. It is limited up to 50, and if you want to enlarge it, please apply for a whitelist or contact us.
* `worker_disk_size` - (Optional) The system disk size of worker nodes. Valid range: [20~32768] in GB.
* `worker_disk_category` - (Optional) The system disk category of worker nodes. Valid values are `cloud`, `cloud_ssd`, `cloud_essd`, and `cloud_efficiency`.
* `worker_data_disks` - (Optional) The configurations of the data disks that are mounted to worker nodes.
  * `category` - (Optional) The data disk category of worker nodes. Valid values are `cloud`, `cloud_ssd`, `cloud_essd`, and `cloud_efficiency`.
  * `size` - (Optional) The data disk size of worker nodes. Valid range: [40~500] in GB.
  * `encrypted` - (Optional) Specifies whether disk encryption is enabled.
  * `auto_snapshot_policy_id` - (Optional) The ID of the policy that is used to back up the data disk.
  * `performance_level` - (Optional) The performance level of the data disk.
* `worker_instance_types` - (Required) The instance type of worker nodes. Specify one type for single AZ Cluster, three types for MultiAZ Cluster.
* `worker_system_disk_performance_level` - (Optional) Performance level of the worker node's system disk.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the container cluster.
* `name` - The name of the container cluster.
* `availability_zone` - The ID of the availability zone.
* `vpc_id` - The ID of VPC where the current cluster is located.
* `slb_intranet` - The ID of private load balancer where the current cluster master node is located.
* `security_group_id` - The ID of security group where the current cluster worker node is located.
* `nat_gateway_id` - The ID of NAT gateway used to launch the Kubernetes cluster.
* `master_nodes` - List of cluster master nodes. It contains several attributes to `Block Nodes`.
  * `id` - Id of the Node.
  * `name` - Name of the Node.
  * `private_ip` - The IP address of the node.
* `worker_nodes` - List of cluster worker nodes. It contains several attributes to `Block Nodes`.
  * `id` - Id of the Node.
  * `name` - Name of the Node.
  * `private_ip` - The private IP address of the node.
* `version` - The Kubernetes server version for the cluster.
* `worker_ram_role_name` - The RAM Role Name attached to worker nodes.
* `nodepool_id` - Id of the NodePool.
* `kube_config` - The path of kube config, like `~/.kube/config`.
* `client_cert` - The path of client certificate, like `~/.kube/client-cert.pem`.
* `client_key` - The path of client key, like `~/.kube/client-key.pem`.
* `cluster_ca_cert` - The path of cluster CA certificate, like `~/.kube/cluster-ca-cert.pem`.
* `master_system_disk_performance_level` - The performance level of the master node's system disk.
* `worker_system_disk_performance_level` - The performance level of the worker node's system disk.
* `is_enterprise_security_group` - Indicates whether an enterprise-level security group is used.
* `cloud_monitor_flags` - Indicates whether cloud monitoring is enabled.