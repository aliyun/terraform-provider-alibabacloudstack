---
subcategory: "Elasticsearch"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_elasticsearch_instance"
sidebar_current: "docs-Alibabacloudstack-elasticsearch-instance"
description: |-
  提供一个在阿里云上的Elasticsearch实例资源。
---

# alibabacloudstack_elasticsearch_instance

提供一个在阿里云上的Elasticsearch实例资源。

## 概述

`alibabacloudstack_elasticsearch_instance` 资源允许您在阿里云上管理Elasticsearch实例。

## 示例用法
```hcl
variable "name" {
  default = "tf-testAccES"
}

data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "192.168.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "192.168.0.0/16"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_elasticsearch_instance" "default" {
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  description = "tf-testAccES16373"
  data_node_spec = "1C 2Gi"
  client_node_amount = "2"
  monitor_password = "Pz*Pu3NIXuk("
  cpu_type = "Intel"
  data_node_disk_type = "yoda-lvm"
  master_node_disk_size = "100"
  master_node_disk_type = "yoda-lvm"
  data_node_disk_size = "500"
  password = "MSPFu4fXosn^"
  vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
  master_node_spec = "1C 2Gi"
  master_node_amount = "3"
  client_node_spec = "1C 2Gi"
  data_node_amount = "3"
  kibana_node_spec = "1C 2Gi"
  kibana_password = "W7mI##At(q*M"
  version = "7.10.0_ali1.6.0"
  scene = "normal"
}
```

## 参数说明
以下参数是支持的：

* `zone_id` - (必选, 变更时重建) Elasticsearch实例所属的可用区。
* `cpu_type` - (必选, 变更时重建) Elasticsearch实例的CPU类型，有效值为`intel`。
* `version` - (必选) 要部署的Elasticsearch版本。
* `description` - (可选) Elasticsearch实例的描述。长度必须在0到30个字符之间，可以包含数字、字母、下划线和连字符。必须以字母、数字或中文字符开头。
* `scense` - (必选, 变更时重建) 实例应用场景。 有效值为 "high", "normal", "log"。
* `data_node_amount` - (必选)Elasticsearch集群中的数据节点数量。有效范围是从 2 到 50。
* `data_node_spec` - (必选) 数据节点的规格。
* `data_node_disk_size` - (必选, 变更时重建) 数据节点的磁盘大小。
* `data_node_disk_type` - (必选, 变更时重建) 数据节点的磁盘类型。
* `kibana_node_spec` - (可选) Kibana节点的规格。
* `kibana_node_password` - (可选)  Kibana节点密码。
* `master_node_amount` - (必选) Elasticsearch集群中的主节点数量。有效范围是从 3 到 50。
* `master_node_spec` - (可选) 主节点的规格。
* `master_node_disk_size` - (可选, 变更时重建) 主节点的磁盘大小。
* `master_node_disk_type` - (可选, 变更时重建) 主节点的磁盘类型。
* `client_node_amount` - (可选) Elasticsearch集群中的客户端节点数量。有效范围是从 2 到 25。
* `client_node_spec` - (可选) 客户端节点的规格。
* `vswitch_id` - (可选) 启动Elasticsearch实例的VSwitch的ID。
* `password` - (必选选) Elasticsearch实例的密码。
* `monitor_password` - (必选) Elasticsearch实例监控的密码。
* `protocol` - (可选) Elasticsearch实例使用的协议。有效值为 HTTP 和 HTTPS。默认值为 HTTP。
* `vpc_whitelist` - (可选) VPC私网访问白名单。
* `apack_accesslog_enabled` - (可选) 访问日志审计功能，仅7.10.0支持。
* `apack_accesslog_search_enabled` - (可选) 搜索请求字段打印功能，仅7.10.0支持。
* `thread_pool_write_queue_size` - (可选) 文档写入队列大小。
* `thread_pool_search_queue_size` - (可选) 文档搜索队列大小。
* `cluster_routing_allocation_disk_watermark_low` - (可选) 当磁盘使用率达到此阈值时，Elasticsearch 会尝试不再向这个节点分配新的分片。
* `cluster_routing_allocation_disk_watermark_high` - (可选) 当磁盘使用率达到此阈值时，Elasticsearch 会尝试将分片重新分配到其他使用率较低的节点。
* `cluster_routing_allocation_disk_watermark_flood_stage` - (可选) 当磁盘使用率达到此阈值时，Elasticsearch会将节点标记为不可分配，并尝试将所有分片从该节点移动到其他节点,节点上的分片设置为只读索引。
* `action_auto_create_index` - (可选) 接收到新文档后，如果没有对应索引，是否允许系统自动创建索引。
* `action_destructive_requires_name` - (可选) 在删除索引时是否需要明确指定索引名称。
* `setting_config` - (可选) Elasticsearch实例配置参数，修改会导致实例重启。

## 属性说明
以下属性会被导出：

* `slb_address` - Elasticsearch实例的IP地址。
* `domain` - Elasticsearch实例的域名。
* `port` - Elasticsearch实例的端口号。
* `status` - Elasticsearch实例的状态。
* `kibana_slb_address` - Kibana实例的IP地址.
* `kibana_domain` - Kibana实例的域名。
* `kibana_port` - Kibana实例的端口号。

## 导入
Elasticsearch实例可以使用ID进行导入，例如：

```bash
$ terraform import alibabacloudstack_elasticsearch_instance.example i-1234567890abcdef0
```