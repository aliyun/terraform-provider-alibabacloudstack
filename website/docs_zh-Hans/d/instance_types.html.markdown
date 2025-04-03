---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_instance_types"
sidebar_current: "docs-alibabacloudstack-datasource-instance-types"
description: |-
    查询云服务器实例类型
---

# alibabacloudstack_instance_types

根据指定过滤条件列出当前凭证权限可以访问的云服务器实例类型列表。

~> **注意：** 默认情况下，仅返回已升级的实例类型。如果需要获取过时的实例类型，必须将`is_outdated`设置为true。

## 示例用法

```
# 声明数据源
data "alibabacloudstack_instance_types" "types_ds" {
  cpu_core_count = 1
  memory_size    = 2
}

output "instance_types"{
  value=data.alibabacloudstack_instance_types.types_ds.*
}
```

## 参数参考

支持以下参数：

* `availability_zone` - (可选) 支持实例类型的可用区。
* `cpu_core_count` - (可选) 筛选特定数量CPU核心的结果。
* `cpu_type` - (可选) 筛选特定类型的CPU结果。可选值：`intel`, `hg`, `kp`, `ft`。
* `memory_size` - (可选) 筛选特定内存大小(GB)的结果。
* `sorted_by` - (可选，强制更新) 排序模式，有效值：`CPU`, `Memory`, `Price`。
* `instance_type_family` - (可选) 根据系列名称筛选结果。例如：'ecs.n4'。
* `eni_amount` - (可选) 筛选网络接口数不超过`eni_amount`的结果。
* `kubernetes_node_role` - (可选) 筛选用于创建Kubernetes集群的结果。可选值：`Master` 和 `Worker`。
* `ids` - (可选) 实例类型ID列表。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - 实例类型ID列表。
* `instance_types` - 实例类型列表。每个元素包含以下属性：
  * `id` - 实例类型的ID。
  * `price` - 实例类型的价格。
  * `cpu_core_count` - CPU核心数量。
  * `memory_size` - 内存大小，以GB为单位。
  * `family` - 实例类型家族。
  * `availability_zones` - 支持该实例类型的可用区列表。
  * `burstable_instance` - 突发性能实例属性：
    * `initial_credit` - 突发性能实例的初始CPU信用。
    * `baseline_credit` - 突发性能实例的基准计算性能CPU信用。
  * `eni_amount` - 实例类型可以附加的最大网络接口数。
  * `local_storage` - 实例类型的本地存储：
    * `capacity` - 本地存储容量，以GB为单位。
    * `amount` - 实例已附加的本地存储设备数量。
    * `category` - 实例已附加的本地存储类别。
  * `cpu_type` - CPU类型。