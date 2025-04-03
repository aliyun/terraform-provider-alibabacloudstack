---
subcategory: "MaxCompute"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_maxcompute_cu"
sidebar_current: "docs-alibabacloudstack-resource-maxcompute-cu"
description: |-
  编排Max Compute Cu
---

# alibabacloudstack_maxcompute_cu

使用Provider配置的凭证在指定的资源集编排Max Compute Cu。


## 示例用法

### 基础用法

```terraform
resource "alibabacloudstack_maxcompute_cu" "example" {
   cu_name      = "tf_testAccAlibabacloudStack7898"
   cu_num       = "1"
   cluster_name = "HYBRIDODPSCLUSTER-A-20210520-07B0"
}
```
## 参数参考

支持以下参数：
* `id` - (必填，变更时重建) MaxCompute CU 的 ID。
* `cu_name` - (必填，变更时重建) MaxCompute CU 的名称。 
* `cu_num` - (必填，变更时重建) MaxCompute CU 的 CU 数量。必须至少为 1。 
* `cluster_name` - (必填，变更时重建) MaxCompute CU 的集群名称。

## 属性参考

导出以下属性：
* `id` - MaxCompute CU 的 ID。 

## 导入

MaxCompute 项目可以通过 *name* 或 ID 导入，例如

```bash
$ terraform import alibabacloudstack_maxcompute_project.example tf_maxcompute_project
```