---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_clusters"
sidebar_current: "docs-Alibabacloudstack-datasource-edas-clusters"
description: |- 
  查询企业级分布式应用服务集群列表。
---

# alibabacloudstack_edas_clusters

根据指定过滤条件列出当前凭证权限可以访问的企业级分布式应用服务集群列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testacc-edas-clusters7485"
}

resource "alibabacloudstack_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  name       = "${var.name}"
}

resource "alibabacloudstack_edas_cluster" "default" {
  cluster_name = "${var.name}"
  cluster_type = 2
  network_mode = 2
  vpc_id       = "${alibabacloudstack_vpc.default.id}"
  region_id    = "cn-neimeng-env30-d01"
}

data "alibabacloudstack_edas_clusters" "default" {
  logical_region_id = "cn-shenzhen:xxx"
  ids               = ["${alibabacloudstack_edas_cluster.default.id}"]
  name_regex        = "${alibabacloudstack_edas_cluster.default.cluster_name}"
  output_file       = "edas_clusters_output.txt"
}

output "first_cluster_name" {
  value = data.alibabacloudstack_edas_clusters.default.clusters[0].cluster_name
}
```

## 参数参考

以下参数是支持的：

* `logical_region_id` - （必填，变更时重建）EDAS命名空间的ID。这用于指定集群所在的逻辑区域。
* `ids` - （选填）集群ID列表，用于按特定集群ID过滤结果。
* `name_regex` - （选填，变更时重建）用于通过集群名称过滤结果的正则表达式字符串。

## 属性参考

除了上述参数外，还导出以下属性：

* `names` - 集群名称列表。
* `ids` - 集群ID列表。
* `clusters` - 集群列表。每个集群包含以下属性：
  * `cluster_id` - 集群ID。
  * `cluster_name` - 集群名称。
  * `cluster_type` - 集群类型。有效值：
    * `0`: 普通Docker集群。
    * `1`: Swarm集群。
    * `2`: ECS集群。
    * `3`: EDAS自建K8s集群。
    * `4`: Pandora自动注册应用集群类型。
    * `5`: 容器服务K8s集群。
  * `create_time` - 集群创建时间的时间戳。
  * `update_time` - 最后变更时间的时间戳。
  * `cpu` - CPU总核数。
  * `cpu_used` - 已使用CPU核数。
  * `mem` - 内存总量，单位为MB。
  * `mem_used` - 已使用内存量，单位为MB。
  * `network_mode` - 网络类型。有效值：
    * `1`: 经典网络。
    * `2`: VPC。
  * `node_num` - ECS数量。
  * `vpc_id` - VPC网络ID。
  * `region_id` - 地域ID。