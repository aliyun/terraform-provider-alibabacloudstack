---
subcategory: "Hologram"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_hologram_clusters"
sidebar_current: "docs-alibabacloudstack-datasource-hologram-clusters"
description: |-
  提供Hologram集群列表
---

# alibabacloudstack_hologram_clusters

该数据源根据指定的过滤条件提供Hologram集群列表。

-> **注意：** 该数据源在v1.187.0及以上版本中可用。

## 示例用法

```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

data "alibabacloudstack_hologram_clusters" "default" {
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  compute_type = "Standard"
  cpu = "intel"
}
```

## 参数说明

以下参数是支持的：

* `zone_id` - （必选）集群所属的可用区ID。
* `compute_type` - （可选）集群的计算类型。有效值：`Standard`，`Follower`。默认值为`Standard`。
* `cpu` - （可选）CPU品牌。默认值为`intel`。
* `name_regex` - （可选）用于按集群名称过滤结果的正则表达式。
* `ids` - （可选）集群ID列表。

## 属性导出

以下属性会被导出：

* `ids` - 集群ID列表。
* `clusters` - Hologram集群列表。每个元素包含以下属性：
  * `id` - 集群ID。
  * `type` - 集群类型。
  * `cpu` - 集群的CPU品牌。
  * `cpu_arch` - 集群的CPU架构。
  * `performance` - 集群的性能等级。
  * `support_replica` - 集群是否支持副本。
  * `cluster` - 集群名称。