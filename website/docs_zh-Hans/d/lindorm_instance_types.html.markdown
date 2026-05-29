---
subcategory: "云原生多模数据库 Lindorm"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_lindorm_instance_types"
description: |-
  为用户提供Lindorm实例类型列表
---

# alibabacloudstack\_lindorm\_instance\_types

该数据源为用户提供可用的Lindorm实例类型列表。

## 示例用法

### 基本用法

```terraform
data "alibabacloudstack_lindorm_instance_types" "example" {
}

output "first_instance_type_id" {
  value = data.alibabacloudstack_lindorm_instance_types.example.instance_types.0.id
}
```

### 按引擎类型筛选

```terraform
data "alibabacloudstack_lindorm_instance_types" "example" {
  engine_type = "lindorm"
}

output "lindorm_instance_types" {
  value = data.alibabacloudstack_lindorm_instance_types.example.instance_types
}
```

### 按CPU和内存筛选

```terraform
data "alibabacloudstack_lindorm_instance_types" "example" {
  cpu    = 4
  memory = 8
}

output "filtered_instance_types" {
  value = data.alibabacloudstack_lindorm_instance_types.example.instance_types
}
```

### 按CPU排序

```terraform
data "alibabacloudstack_lindorm_instance_types" "example" {
  sorted_by = "CPU"
}

output "sorted_instance_types" {
  value = data.alibabacloudstack_lindorm_instance_types.example.instance_types
}
```

## 参数参考

以下参数被支持：

* `ids` - (可选) 实例类型ID列表。
* `engine_type` - (可选) 引擎类型。有效值: `lindorm`, `tsdb`, `solr`, `lts`。
* `cpu` - (可选) CPU数量。
* `memory` - (可选) 内存大小，单位GB。
* `sorted_by` - (可选) 排序字段。有效值: `CPU`, `Memory`。

## 属性参考

以下属性会被导出：

* `ids` - 实例类型ID列表。
* `instance_types` - Lindorm实例类型列表。每个元素包含以下属性：
  * `id` - 实例类型ID。
  * `cpu` - CPU数量。
  * `memory` - 内存大小，单位GB。
  * `rate` - 速率。
  * `name` - 实例类型的名称。