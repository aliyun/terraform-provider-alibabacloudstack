---
subcategory: "物联网平台(连接版)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mqtt_clusters"
sidebar_current: "docs-Alibabacloudstack-datasource-mqtt-clusters"
description: |-
  提供MQTT集群列表。
---

# alibabacloudstack_mqtt_clusters

该数据源用于获取专有云中可用的MQTT集群列表。

-> **注意:** 适用于专有云环境。

## 示例

```hcl
data "alibabacloudstack_mqtt_clusters" "example" {
  name_regex = "^cluster-.*"
}

output "mqtt_clusters" {
  value = data.alibabacloudstack_mqtt_clusters.example.clusters
}
```

## 参数说明

以下参数支持配置：

* `ids` - (可选, 变更时强制重建) 用于过滤结果的集群ID列表。
* `name_regex` - (可选, 变更时强制重建) 用于按名称过滤集群的正则表达式字符串。

## 属性参考

以下属性会被导出：

* `ids` - 集群ID列表。
* `clusters` - MQTT集群列表。每个元素包含以下属性：
  * `id` - 集群的ID。
  * `name` - 集群的名称。
