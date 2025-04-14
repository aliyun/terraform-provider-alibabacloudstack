---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_metering_query_ecs"
sidebar_current: "docs-alibabacloudstack-ascm-metering-query-ecs"
description: |-
  查询ecs计量数据
---
# alibabacloudstack_ascm_metering_query_ecs

根据指定过滤条件列出当前凭证权限可以访问的ECS实例的计量数据列表。

## 示例用法

```hcl
data "alibabacloudstack_ascm_metering_query_ecs" "example" {
  start_time = "2023-01-01T00:00:00Z"
  end_time   = "2023-01-31T23:59:59Z"
  product_name = "ECS"
}

output "ecs_metering_data" {
  value = data.alibabacloudstack_ascm_metering_query_ecs.example.data
}
```

## 参数说明
支持以下参数：

* `start_time` - (必填) 计量查询的开始时间，ISO 8601格式(例如："2023-01-01T00:00:00Z")。
* `end_time` - (必填) 计量查询的结束时间，ISO 8601格式(例如："2023-01-31T23:59:59Z")。
* `org_id` - (可选) 要检索计量数据的组织ID。如果未指定，则默认使用当前用户所属的组织。
* `product_name` - (必填) 要检索计量数据的产品名称(例如："ECS")。
* `is_parent_id` - (可选) 指示组织ID是否为父ID。
* `ins_id` - (可选) 要检索计量数据的实例ID。
* `region` - (可选) 要检索计量数据的区域。
* `resource_group_id` - (可选) 要检索计量数据的资源组ID。
* `name_regex` - (可选) 用于按实例名称过滤结果的正则表达式模式。

## 属性说明
导出以下属性：

* `data` - ECS计量数据列表。每个元素包含以下属性：
    * `private_ip_address` - ECS实例的私有IP地址。
    * `instance_type_family` - ECS实例的实例类型族。
    * `memory` - ECS实例的内存大小（以GB为单位）。
    * `cpu` - ECS实例中的CPU数量。
    * `os_name` - ECS实例的操作系统名称。
    * `org_name` - 组织名称。
    * `instance_network_type` - ECS实例的网络类型（例如：经典网络或VPC）。
    * `eip_address` - ECS实例的弹性公网IP地址。
    * `resource_g_name` - 资源组名称。
    * `instance_type` - ECS实例的实例类型。
    * `status` - ECS实例的状态（例如：运行中、停止等）。
    * `sys_disk_size` - ECS实例的系统盘大小（以GB为单位）。
    * `gpu_amount` - ECS实例中的GPU数量。
    * `instance_name` - ECS实例的名称。
    * `vpc_id` - ECS实例的VPC ID。
    * `start_time` - 计量数据的开始时间，ISO 8601格式。
    * `end_time` - 计量数据的结束时间，ISO 8601格式。
    * `create_time` - ECS实例的创建时间，ISO 8601格式。
    * `data_disk_size` - ECS实例的数据盘大小（以GB为单位）。
    * `is_parent_id` - 指示组织ID是否为父ID。
    * `ins_id` - 实例ID。
    * `region` - 实例所在的区域。
    * `resource_group_id` - 实例所属的资源组ID。