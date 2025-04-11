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

## 参数说明

支持以下参数：
* `cu_name` - (必填，变更时重建，1.110.0+可用) MaxCompute CU 的名称。必须为 3 到 27 个字符。
* `cu_num` - (必填，变更时重建) MaxCompute CU 的 CU 数量。必须至少为 1。
* `cluster_name` - (必填，变更时重建) MaxCompute CU 所属的集群名称。

## 属性说明

导出以下属性：
* `id` - MaxCompute CU 的 ID。

## 导入

MaxCompute 项目可以通过 *name* 或 ID 导入，例如

```bash
$ terraform import alibabacloudstack_maxcompute_project.example tf_maxcompute_project
```