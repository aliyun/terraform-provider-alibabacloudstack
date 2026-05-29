---
subcategory: "云服务器 ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_dedicatedhostclusters"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-dedicatedhostclusters"
description: |-
  提供阿里云账号下拥有的ecs dedicatedhostclusters列表。
---

# alibabacloudstack\_ecs\_dedicatedhostclusters

此数据源提供根据指定过滤条件列出的阿里云账号下的ecs dedicatedhostclusters资源列表。

## 示例用法
```
data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}


		resource "alibabacloudstack_ecs_dedicated_host_cluster" "default" {
		  dedicated_host_cluster_name = "tf_testAccEcsDedicatedHostsClusterDataSource_5238343"
          zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
		}
	

data "alibabacloudstack_ecs_dedicated_host_cluster" "default" {
  ids = [
          "${alibabacloudstack_ecs_dedicated_host_cluster.default.id}"
        ]
}
```

## 参数参考
以下参数是支持的：
  * `ids` - (选填) - 筛选结果的专有宿主机集群ID列表。
  * `zone_id` - (选填) - 可用区ID
  * `dedicated_host_cluster_name_regex` - (选填) - 用于通过专有宿主机集群名称筛选结果的正则表达式字符串。
  * `dedicated_host_cluster_name` - (选填) - 专有宿主机集群名

## Attributes Reference
除了上述参数外，还导出以下属性：
  * `dedicated_host_clusters` - 专有宿主机集群列表。
    * `id` - 该专有宿主机集群的ID。
    * `dedicated_host_cluster_id` - 专有宿主机集群ID
    * `dedicated_host_cluster_name` - 专有宿主机集群名
    * `description` - 描述
    * `zone_id` - 可用区ID
