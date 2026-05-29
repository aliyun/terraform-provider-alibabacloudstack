---
subcategory: "硬件安全模块"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_hsm_clusters"
sidebar_current: "docs-Alibabacloudstack-datasource-hsm-clusters"
description: |-
    Provides information about Alibaba Cloud HSM Clusters.
---

# alibabacloudstack_hsm_clusters

> 查询阿里云密码机集群（HSM Cluster）的信息。

## 示例用法

```hcl

variable "name" {
  default = "tf_hsm_cluster6821719"
}

data "alibabacloudstack_zones" "default" {
  provider                    = alibabacloudstack-common
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  provider   = alibabacloudstack-common
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  provider     = alibabacloudstack-common
  vswitch_name = "${var.name}_vsw"
  vpc_id       = alibabacloudstack_vpc.default.id
  cidr_block   = "172.16.1.0/24"
  zone_id      = data.alibabacloudstack_zones.default.zones.0.id
}

data "alibabacloudstack_hsm_vendors" "default" {
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

resource "alibabacloudstack_hsm_instance" "default" {
  product_code = data.alibabacloudstack_hsm_vendors.default.vendors.0.products.0.code
  vendor_code  = data.alibabacloudstack_hsm_vendors.default.vendors.0.code
  vsm_type     = "gvsm"
  zone_no      = data.alibabacloudstack_zones.default.zones.0.id
  remark       = var.name
}

resource "alibabacloudstack_hsm_cluster" "default" {
  cluster_name       = var.name
  master_instance_id = alibabacloudstack_hsm_instance.default.id
  vpc_id             = alibabacloudstack_vpc.default.id
  vswitch_ids        = alibabacloudstack_vswitch.default.id
  zone_nos           = data.alibabacloudstack_zones.default.zones.0.id
  ip_white_list      = "123.12.13.1/16"
  password           = random_password.password.0.result
}

data "alibabacloudstack_hsm_clusters" "default" {
  name_regex = alibabacloudstack_hsm_cluster.default.cluster_name
}

```

## 参数说明
以下参数用于过滤查询结果：

* `vsm_type` (字符串, 变更时重建)：用于过滤特定类型的集群。取值：`evsm`（金融数据密码机）、`gvsm`（通用服务器密码机）、`svsm`（签名验签服务器密码机）。

* `ids` (字符串列表, 可选)：用于过滤特定ID的集群。如果指定，将只返回ID在列表中的集群。

* `name_regex` (字符串, 可选)：用于通过正则表达式过滤集群名称。只有名称匹配正则表达式的集群会被返回。

## 属性说明
以下属性被导出：

* `id` (字符串)：集群ID，作为数据源的唯一标识符。

* `abnormal_type` (字符串)：集群异常类型。只有集群状态异常时，此参数有值，否则为空值。取值：`3`（密码机数据不一致）、`4`（所有密码机不可访问）、`5`（密码机数量不一致）。

* `cluster_id` (字符串)：集群ID，集群唯一标识。

* `cluster_master` (字符串)：集群中主密码机实例ID。

* `cluster_name` (字符串)：用户自定义的集群名称。

* `cluster_size` (字符串)：集群大小，即集群下属密码机实例的数量。

* `cluster_zones` (对象列表)：集群下属的可用区与交换机的数据集。每个交换机只属于一个可用区。
  * `vswitch_id` (字符串)：交换机ID。
  * `zone_no` (字符串)：可用区编号。

* `gmt_create` (整数)：集群的创建时间（Unix时间戳，毫秒）。

* `hsm_cluster_items` (对象列表)：集群中子项数据集，每个子项即为密码机实例。
  * `id` (整数)：集群中子项的主键ID。
  * `instance_id` (字符串)：密码机实例ID。
  * `ip` (字符串)：密码机实例所注册的经典网络IP地址。
  * `is_master` (整数)：当前密码机是否是集群的主密码机（1表示是，0表示否）。
  * `status` (整数)：密码机实例状态。取值：`1`（未初始化）、`3`（已释放）、`4`（生产失败）、`5`（运行中）、`6`（同步中）、`7`（重置中）、`8`（已停用）。
  * `vsm_type` (字符串)：密码机实例的设备类型。
  * `zone_no` (字符串)：密码机实例所在的可用区编号。

* `ip_white_list` (字符串)：集群的白名单IP，多个IP用逗号分隔。

* `master_ip` (字符串)：集群中主密码机实例所注册的经典网络IP地址。

* `product_code` (字符串)：密码机实例所用设备的型号编码。

* `vendor_code` (字符串)：密码机实例所用设备的厂商编码。

* `vpc_id` (字符串)：集群的VPC实例ID。每个集群只有一个VPC实例。

* `vsm_type` (字符串)：集群中主密码机实例的设备类型，此处标识为整个集群的密码机实例的设备类型。取值：`evsm`（金融数据密码机）、`gvsm`（通用服务器密码机）、`svsm`（签名验签服务器密码机）。