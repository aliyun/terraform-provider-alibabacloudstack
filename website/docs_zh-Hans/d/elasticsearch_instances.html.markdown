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

```hcl
data "alibabacloudstack_elasticsearch_instances" "default" {
  description_regex = "my-es-instance"
  version           = "6.7.0_with_X-Pack"
  vpc_id            = "vpc-xxxxxxxxxxxxx"
}
```

## 参数说明

以下是支持的参数：

* `description_regex` - (可选) 用于过滤实例描述的正则表达式字符串。
* `ids` - (可选) Elasticsearch 实例 ID 列表。
* `version` - (可选) Elasticsearch 版本。如果不指定，则返回所有版本的实例。
* `vpc_id` - (可选) Elasticsearch 实例所属的 VPC ID。
* `output_file` - (可选, 已弃用) 该字段已弃用，计划在 3.19.0 版本中移除。如需将内容写入文件，请使用 `local_file` provider。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `ids` - Elasticsearch 实例 ID 列表。
* `descriptions` - Elasticsearch 实例描述列表。
* `instances` - Elasticsearch 实例列表。每个元素包含以下属性：
  * `id` - Elasticsearch 实例的 ID。
  * `version` - Elasticsearch 版本。
  * `description` - Elasticsearch 实例的描述。
  * `data_node_amount` - Elasticsearch 集群中的数据节点数量。
  * `data_node_spec` - 数据节点的规格。
  * `data_node_disk_size` - 数据节点的磁盘大小。
  * `data_node_disk_type` - 数据节点的磁盘类型。
  * `kibana_node_spec` - Kibana 节点的规格（仅在启用 Kibana 时可用）。
  * `kibana_slb_address` - Kibana 的 SLB 地址（仅在启用 Kibana 时可用）。
  * `kibana_domain` - Kibana 的域名（仅在启用 Kibana 时可用）。
  * `kibana_protocol` - Kibana 使用的协议（仅在启用 Kibana 时可用）。
  * `kibana_port` - Kibana 的端口号（仅在启用 Kibana 时可用）。
  * `master_node_amount` - 专用主节点数量（仅在启用专用主节点时可用）。
  * `master_node_spec` - 专用主节点的规格。
  * `master_node_disk_size` - 专用主节点的磁盘大小。
  * `master_node_disk_type` - 专用主节点的磁盘类型。
  * `client_node_amount` - 客户端节点数量（仅在启用客户端节点时可用）。
  * `client_node_spec` - 客户端节点的规格。
  * `vswitch_id` - 启动 Elasticsearch 实例的 VSwitch 的 ID。
  * `status` - Elasticsearch 实例的状态。