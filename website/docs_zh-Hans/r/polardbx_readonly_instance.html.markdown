---
subcategory: "云原生分布式数据库PolarDB-X 2.0"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_readonly_instance"
sidebar_current: "docs-Alibabacloudstack-resource-polardbx-readonly-instance"
description: |-
  Provides a PolarDB-X read-only instance resource.
---

# alibabacloudstack\_polardbx\_readonly\_instance

提供 PolarDB-X 只读实例资源。


## 示例用法

```hcl
variable "name" {
  default = "tf_acc_drds_polardb_readonly"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details              = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name       = "${var.name}_vsw"
  vpc_id     = alibabacloudstack_vpc_vpc.default.id
  cidr_block = "172.16.1.0/24"
  zone_id    = data.alibabacloudstack_zones.default.zones.0.id
}

# 先创建主实例
resource "alibabacloudstack_polardbx_instance" "primary" {
  vswitch_id     = alibabacloudstack_vpc_vswitch.default.id
  cn_node_count  = 2
  dn_node_class  = "mysql.n4.medium.25"
  dn_node_count  = 2
  description    = "primary instance"
  storage        = 50
  cn_node_class  = "polarx.x4.medium.2e"
  engine_version = "5.7"
  zone_id        = data.alibabacloudstack_zones.default.zones.0.id
}

# 创建只读实例
resource "alibabacloudstack_polardbx_readonly_instance" "readonly" {
  primary_db_instance_id = alibabacloudstack_polardbx_instance.primary.id
  vswitch_id             = alibabacloudstack_vpc_vswitch.default.id
  cn_node_count          = 2
  dn_node_class          = "mysql.n4.medium.25"
  dn_node_count          = 2
  description            = "readonly instance"
  storage                = 50
  cn_node_class          = "polarx.x4.medium.2e"
  engine_version         = "5.7"
  zone_id                = data.alibabacloudstack_zones.default.zones.0.id
}
```

## 参数参考

支持以下参数：

### 必填参数

* `primary_db_instance_id` - (必填，变更时强制重建) 主实例的 ID。修改此参数会强制重新创建只读实例。
* `storage` - (必填，变更时强制重建) PolarDB-X 只读实例的存储容量，单位为 GB。
* `cn_node_class` - (必填) 计算节点（CN）的规格。示例值：`polarx.x4.medium.2e`、`polarx.x4.large.2e`、`polarx.x8.large.2e`、`polarx.x4.xlarge.2e`、`polarx.x8.xlarge.2e`、`polarx.x4.2xlarge.2e`、`polarx.x8.2xlarge.2e`、`polarx.x4.4xlarge.2e`、`polarx.x8.4xlarge.2e`。
* `cn_node_count` - (必填) 计算节点（CN）的数量。
* `dn_node_class` - (必填) 数据节点（DN）的规格。示例值：`mysql.n4.medium.25`、`mysql.n4.large.25`、`mysql.n4.xlarge.25`、`mysql.n4.2xlarge.25`、`mysql.x4.medium.25`。
* `dn_node_count` - (必填) 数据节点（DN）的数量。
* `vswitch_id` - (必填，变更时强制重建) 交换机 ID。

### 选填参数

* `cpu_type` - (选填) PolarDB-X 只读实例的 CPU 架构类型。
* `description` - (选填) PolarDB-X 只读实例的描述信息。
* `engine_version` - (选填，变更时强制重建) PolarDB-X 只读实例的引擎版本。取值：`5.7`、`8.0`。
* `gms_node_class` - (选填) GMS（全局元数据服务）节点的规格。
* `polardbx_instance_id` - (选填) PolarDB-X 只读实例 ID。
* `primary_zone` - (选填，变更时强制重建) 主可用区。
* `resource_type` - (选填，变更时强制重建) 资源类型。目前仅支持 PolarDB-X 2.0 实例。
* `secondary_zone` - (选填，变更时强制重建) 次可用区。
* `tertiary_zone` - (选填，变更时强制重建) 第三可用区。
* `topology_type` - (选填，变更时强制重建) 实例的拓扑类型。取值：`1azone`（单可用区）、`3azones`（三可用区）。默认值：`1azone`。
* `enable_tde` - (选填) 是否启用 TDE（透明数据加密）。取值：`true`、`false`。默认值：`false`。启用后无法禁用。
* `enable_ssl` - (选填) 是否启用 SSL 加密。取值：`true`、`false`。默认值：`false`。
* `zone_id` - (选填，变更时强制重建) 实例所属的可用区 ID。
* `compute_parameters` - (选填) 实例的计算资源配置。这是一个键值对集合。
    * `name` - (必填) 计算参数的名称。
    * `value` - (必填) 计算参数的值。
* `storage_parameters` - (选填) 实例的存储配置。这是一个键值对集合。
    * `name` - (必填) 存储参数的名称。
    * `value` - (必填) 存储参数的值。
* `security_groups` - (选填) 实例的安全 IP 组。
    * `group_name` - (必填) 安全 IP 组的名称。
    * `ips` - (必填) 安全 IP 地址。多个 IP 需要用 "," 分隔，示例：`192.168.0.1,192.168.0.0/24`。
* `private_connection_string_prefix` - (选填) 私有连接串的前缀。
* `private_connection_port` - (选填) 私有连接的端口。取值范围：3000-6000。默认值：`3306`。
* `enable_public_connection` - (选填) 是否启用公网连接。取值：`true`、`false`。默认值：`false`。
* `public_connection_string_prefix` - (选填) 公网连接串的前缀。当 `enable_public_connection` 为 `true` 时必填。
* `public_connection_port` - (选填) 公网连接的端口。取值范围：3000-6000。默认值：`3306`。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - PolarDB-X 只读实例 ID。
* `create_time` - 实例的创建时间。
* `network_type` - 实例的网络类型。
* `status` - 实例的状态。
* `connection_string` - 实例的私有连接串。
* `private_connection_string` - 私有连接串。
* `private_connection_port` - 私有连接端口。
* `public_connection_string` - 公网连接串（如果已启用）。
* `public_connection_port` - 公网连接端口（如果已启用）。

## Import

PolarDB-X 只读实例可以使用实例 ID 进行导入，例如：

```
$ terraform import alibabacloudstack_polardbx_readonly_instance.example px-12345678
```
