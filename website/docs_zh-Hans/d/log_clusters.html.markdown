---
layout: "alicloud-doc"
page_title: "数据源: alibabacloudstack_log_clusters"
subcategory: "日志服务 SLS"
---

# alibabacloudstack_log_clusters

此数据源提供可用的 SLS 集群列表。

## 示例用法

```hcl
data "alibabacloudstack_log_clusters" "default" {
  name_regex = "^my-cluster.*"
}

output "first_cluster_name" {
  value = data.alibabacloudstack_log_clusters.default.clusters.0.name
}
```

## 参数说明

支持以下参数：

- `ids` -（可选，变更后重建）集群 ID 列表。
- `name_regex` -（可选，变更后重建）用于按集群名称过滤结果的正则表达式字符串。

## 属性说明

除上述参数外，还导出以下属性：

- `ids` - 集群 ID 列表。
- `clusters` - 集群列表。每个元素包含以下属性：
  - `name` - 集群名称。
  - `data_server` - 集群的数据服务器地址。
  - `status` - 集群状态。
  - `zone` - 集群所在的可用区。
  - `description` - 集群描述。
