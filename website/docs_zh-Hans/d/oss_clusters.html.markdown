---
subcategory: "对象存储 OSS"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_oss_clusters"
sidebar_current: "docs-Alibabacloudstack-datasource-oss-clusters"
description: |-
  Provides a list of OSS Clusters to the user.
---

# alibabacloudstack_oss_clusters

This data source provides the OSS Clusters of the current Alibaba Cloud user.

## Example Usage

```hcl
data "alibabacloudstack_oss_clusters" "example" {
  ids = ["cluster-id-1"]
}

output "first_cluster_id" {
  value = data.alibabacloudstack_oss_clusters.example.clusters.0.id
}
```

```hcl
data "alibabacloudstack_oss_clusters" "filtered" {
  name_regex = "^master-"
}

output "filtered_clusters" {
  value = data.alibabacloudstack_oss_clusters.filtered.clusters
}
```

## 参数说明

* `ids` - (可选, 变更后新建) 指定集群ID列表,用于精确匹配特定集群的信息。(*可选*)
* `name_regex` - (可选, 变更后新建) 通过正则表达式过滤集群名称,用于筛选符合条件的集群。(*可选*)

## 属性导出

* `ids` - 匹配到的集群ID列表。
* `clusters` - 匹配到的集群信息列表,每个元素包含以下属性:
  * `id` - 集群ID,对应集群标识。
  * `cluster` - 集群标识。
  * `ha_apsara_stack` - 是否为高可用ApsaraStack。
  * `api_zonelocal_endpoint` - Zone本地API端点。
  * `api_zonelocal_public_endpoint` - Zone本地公共API端点。
  * `oss_public_endpoint` - OSS公共端点。
  * `real_zone` - 实际Zone信息。
  * `oss_ha_enable_single_cluster_access` - 是否启用单集群访问的高可用OSS。
  * `oss_cs_public_endpoint` - OSS容器服务公共端点。
  * `oss_unique_domain` - 是否使用唯一域名。
  * `cluster_name` - 集群名称。
  * `is_master_zone` - 是否为主Zone。
  * `location` - 地理位置信息。
  * `oss_endpoint` - OSS端点地址。
  * `oss_suffix` - OSS后缀信息。
