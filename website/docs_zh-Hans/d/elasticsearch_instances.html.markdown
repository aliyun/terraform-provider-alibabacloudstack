---
subcategory: "Elasticsearch"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_elasticsearch_instances"
sidebar_current: "docs-alibabacloudstack-datasource-elasticsearch-instances"
description: |-
  查询Elasticsearch。
---

# alibabacloudstack_elasticsearch_instances

根据指定过滤条件列出当前凭证权限可以访问的Elasticsearch实例列表。

## 示例用法

```
data "alibabacloudstack_elasticsearch_instances" "instances" {
  description_regex = "myes"
  version           = "5.5.3_with_X-Pack"
}
```

## 参数说明

以下是支持的参数：

* `description_regex` - (可选) 应用于实例描述的正则表达式字符串。
* `ids` - (可选, 1.52.1+可用) Elasticsearch实例ID列表。
* `version` - (可选) Elasticsearch版本。选项包括 `5.5.3_with_X-Pack`, `6.3.2_with_X-Pack` 和 `6.7.0_with_X-Pack`。如果不指定值，则返回所有版本。
* `tags` - (可选, 1.74.0+可用) 分配给实例的标签映射。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `ids` - Elasticsearch实例ID列表。
* `descriptions` - Elasticsearch实例描述列表。
* `instances` - Elasticsearch实例列表。每个元素包含以下属性：
  * `id` - Elasticsearch实例的ID。
  * `zone_id` - Elasticsearch实例所属的可用区。
  * `cpu_type` - Elasticsearch实例的CPU类型。
  * `version` - 要部署的Elasticsearch版本。
  * `description` - Elasticsearch实例的描述。长度必须在0到30个字符之间，可以包含数字、字母、下划线和连字符。必须以字母、数字或中文字符开头。
  * `scense` - 实例应用场景。
  * `data_node_amount` - Elasticsearch集群中的数据节点数量。
  * `data_node_spec` - 数据节点的规格。
  * `data_node_disk_size` - 数据节点的磁盘大小。
  * `data_node_disk_type` - 数据节点的磁盘类型。
  * `kibana_node_spec` - Kibana节点的规格。
  * `kibana_node_password` - Kibana节点密码。
  * `master_node_amount` - Elasticsearch集群中的主节点数量。
  * `master_node_spec` - 主节点的规格。
  * `master_node_disk_size` - 主节点的磁盘大小。
  * `master_node_disk_type` - 主节点的磁盘类型。
  * `client_node_amount` - Elasticsearch集群中的客户端节点数量。
  * `client_node_spec` - 客户端节点的规格。
  * `vswitch_id` - 启动Elasticsearch实例的VSwitch的ID。