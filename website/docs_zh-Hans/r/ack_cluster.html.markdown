---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ack_cluster"
sidebar_current: "docs-Alibabacloudstack-ack-cluster"
description: |- 
  编排ack集群
---

# alibabacloudstack_ack_cluster
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_cs_kubernetes`

使用Provider配置的凭证在指定的资源集下编排ack集群。

## 示例用法

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
  default     = "true"
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
  default     = "Alibaba@1688"
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

## 参数参考

支持以下参数：

* `name` - (选填) Kubernetes 集群的名称。在同一 AlibabaCloudStack 账户内必须唯一。
* `vpc_id` - (必填) 当前集群所在的 VPC 的 ID。
* `worker_vswitch_ids` - (必填) 工作节点使用的交换机。可以指定一个或多个交换机。
* `master_vswitch_ids` - (必填) 主节点使用的交换机。可以基于主节点数量指定三个或五个交换机。
* `version` - (选填) 所需的 Kubernetes 版本。如果您不指定值，将使用资源创建时的最新可用版本，并且除非您设置更高的版本号，否则不会进行升级。降级不受 ACK 支持。
* `password` - (必填，敏感) 用于通过 SSH 登录到集群节点的密码。必须指定 `password`、`key_name` 或 `kms_encrypted_password` 中的一个。
* `kms_encrypted_password` - (选填) 用于在创建或更新 CS Kubernetes 前解密密码的 KMS 加密密码。
* `enable_ssh` - (选填) 是否启用通过 SSH 登录到节点。默认为 `false`。
* `cpu_policy` - (选填) kubelet CPU 策略。选项：`static` | `none`。默认为 `none`。
* `proxy_mode` - (选填) kube-proxy 的代理模式。选项：`iptables` | `ipvs`。默认为 `ipvs`。
* `user_data` - (选填) 实例的用户定义数据。Windows 实例支持批处理和 PowerShell 脚本。如果您的脚本文件大于 1 KB，建议将其上传到对象存储服务 (OSS)，并通过 OSS 存储桶的内部端点拉取它。
* `instances` - (选填) 可以作为工作节点附加到同一 VPC 的实例列表。
* `os_type` - (选填) 运行 pod 的节点的操作系统类型，其有效值为 `Linux` 或 `Windows`。默认为 `Linux`。
* `platform` - (选填) 运行 pod 的节点架构。默认为 `CentOS`。
* `security_group_id` - (选填) 集群所属的 ECS 的安全组 ID。与 `is_enterprise_security_group` 冲突。
* `is_enterprise_security_group` - (选填) 是否自动创建高级安全组。必须设置 `security_group_id` 或 `is_enterprise_security_group` 参数之一。
* `runtime` - (选填) 集群运行的平台。
  * `name` - (选填) 运行平台的名称。
  * `version` - (选填) 运行平台的版本。
* `tags` - (选填) 分配给资源的标签映射。
  - Key: 最大长度为 64 个字符。不能以 "aliyun"、"acs:"、"http://" 或 "https://" 开头。不能是空字符串。
  - Value: 最大长度为 128 个字符。不能以 "aliyun"、"acs:"、"http://" 或 "https://" 开头。可以是空字符串。
* `keep_instance_name` - (选填) 将现有实例添加到节点池时，是否保留原始实例名称。建议设置为 `true`。
* `format_disk` - (选填) 选择此选项后，如果已将数据盘附加到指定的 ECS 实例并且最后一个数据盘的文件系统未初始化，系统将自动将最后一个数据盘格式化为 ext4 并挂载到 `/var/lib/docker` 和 `/var/lib/kubelet`。磁盘上的原始数据将被清除。请确保提前备份数据。如果没有在 ECS 实例上挂载数据盘，则不会购买新的数据盘。默认为 `false`。
* `image_id` - (选填) 自定义镜像支持。必须基于 CentOS7 或 AliyunLinux2。
* `timeout_mins` - (选填) 后端服务超时时间；单位：分钟。默认为 60。
* `delete_protection` - (选填) 实例是否应具有删除保护。
* `kms_encryption_context` - (选填) 用于在创建或更新 CS Kubernetes 前解密 `kms_encrypted_password` 的 KMS 加密上下文。有关更多信息，请参阅 [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm)。当设置了 `kms_encrypted_password` 时有效。
* `addons` - (选填) 您要在集群中安装的插件。
  * `name` - (选填) ACK 插件的名称。名称必须与 DescribeAddons 返回的名称之一匹配。
  * `config` - (选填) ACK 插件配置。有关更多配置信息，请参阅 [cs_kubernetes_addon_metadata](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/data-sources/cs_kubernetes_addon_metadata)。
* `cloud_monitor_flags` - (选填) 是否安装云监控插件。

### 网络

* `pod_cidr` - (选填) [Flannel 特定] 使用 Flannel 时的 Pod 网络 CIDR 块。
* `pod_vswitch_ids` - (选填) [Terway 特定] 使用 Terway 时的 Pod 网络交换机。请注意，`pod_vswitch_ids` 不能等于 `worker_vswitch_ids` 或 `master_vswitch_ids`，但必须在相同的可用区中。
* `new_nat_gateway` - (选填) 创建 Kubernetes 集群时是否创建新的 NAT 网关。默认为 `true`。
* `service_cidr` - (选填) 服务网络的 CIDR 块。它不能与 VPC CIDR 和 VPC 中 Kubernetes 集群使用的 CIDR 重复，创建后无法修改。
* `node_cidr_mask` - (选填) 节点 CIDR 块，用于指定单个节点上可以运行多少个 Pod。有效范围：24-28。默认为 24。
* `slb_internet_enabled` - (选填) 是否为 API Server 创建互联网负载均衡器。默认为 `true`。

### Master 参数

* `master_count` - (选填) 主节点的数量。默认为 3。
* `master_disk_category` - (选填) 主节点的系统盘类别。有效值为 `cloud_ssd` 和 `cloud_efficiency`。默认为 `cloud_efficiency`。
* `master_disk_size` - (选填) 主节点的系统盘大小。有效范围：[20~500] GB。默认为 20。
* `master_instance_types` - (必填) 主节点的实例类型。对于单 AZ 集群，指定一种类型；对于多 AZ 集群，指定三种类型。
* `master_system_disk_performance_level` - (选填) 主节点系统盘的性能级别。

### Worker 参数

* `num_of_nodes` - (必填) Kubernetes 集群的工作节点数量。默认为 3。限制最多为 50，如果需要扩大，请申请白名单或联系我们。
* `worker_disk_size` - (选填) 工作节点的系统盘大小。有效范围：[20~32768] GB。
* `worker_disk_category` - (选填) 工作节点的系统盘类别。有效值为 `cloud`、`cloud_ssd`、`cloud_essd` 和 `cloud_efficiency`。
* `worker_data_disks` - (选填) 挂载到工作节点的数据盘配置。
  * `category` - (选填) 工作节点的数据盘类别。有效值为 `cloud`、`cloud_ssd`、`cloud_essd` 和 `cloud_efficiency`。
  * `size` - (选填) 工作节点的数据盘大小。有效范围：[40~500] GB。
  * `encrypted` - (选填) 是否启用磁盘加密。
  * `auto_snapshot_policy_id` - (选填) 用于备份数据盘的策略 ID。
  * `performance_level` - (选填) 数据盘的性能级别。
* `worker_instance_types` - (必填) 工作节点的实例类型。对于单 AZ 集群，指定一种类型；对于多 AZ 集群，指定三种类型。
* `worker_system_disk_performance_level` - (选填) 工作节点系统盘的性能级别。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - 容器集群的 ID。
* `name` - 容器集群的名称。
* `availability_zone` - 可用区的 ID。
* `vpc_id` - 当前集群所在的 VPC 的 ID。
* `slb_intranet` - 当前集群主节点所在的私有负载均衡器的 ID。
* `security_group_id` - 当前集群工作节点所在的安全组的 ID。
* `nat_gateway_id` - 用于启动 Kubernetes 集群的 NAT 网关的 ID。
* `master_nodes` - 集群主节点列表。它包含多个 `Block Nodes` 属性。
  * `id` - 节点的 ID。
  * `name` - 节点的名称。
  * `private_ip` - 节点的 IP 地址。
* `worker_nodes` - 集群工作节点列表。它包含多个 `Block Nodes` 属性。
  * `id` - 节点的 ID。
  * `name` - 节点的名称。
  * `private_ip` - 节点的私有 IP 地址。
* `version` - 集群的 Kubernetes 服务器版本。
* `worker_ram_role_name` - 附加到工作节点的 RAM 角色名称。
* `nodepool_id` - 节点池的 ID。
* `kube_config` - kube config 的路径，例如 `~/.kube/config`。
* `client_cert` - 客户端证书的路径，例如 `~/.kube/client-cert.pem`。
* `client_key` - 客户端密钥的路径，例如 `~/.kube/client-key.pem`。
* `cluster_ca_cert` - 集群 CA 证书的路径，例如 `~/.kube/cluster-ca-cert.pem`。