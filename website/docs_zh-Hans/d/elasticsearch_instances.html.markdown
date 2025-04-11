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
  * `description` - Elasticsearch实例的描述。
  * `instance_charge_type` - 计费方法。选项：`PostPaid` 表示按量付费，`PrePaid` 表示包年包月。
  * `data_node_amount` - Elasticsearch集群的数据节点数量，范围为2到50。
  * `data_node_spec` - Elasticsearch实例的数据节点规格。
  * `data_node_disk_size` - 单个数据节点的存储空间大小。单位：GB。
  * `data_node_disk_type` - 数据节点磁盘类型。包括值：`cloud_ssd` 和 `cloud_efficiency`。
  * `vswitch_id` - 实例所属的VSwitch ID。
  * `version` - Elasticsearch版本，包括 `5.5.3_with_X-Pack`, `6.3.2_with_X-Pack` 和 `6.7.0_with_X-Pack`。
  * `created_at` - 实例的创建时间。它是GTM格式，例如："2019-01-08T15:50:50.623Z"。
  * `updated_at` - 实例的最后修改时间。它是GMT格式，例如："2019-01-08T15:50:50.623Z"。
  * `status` - 实例的状态。包括 `active`, `activating`, `inactive`。
  * `tags` - 分配给实例的标签映射。
  * `output_file` - 保存数据源结果的文件名(在运行 `terraform plan` 后)。