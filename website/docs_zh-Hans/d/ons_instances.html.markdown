---
subcategory: "RocketMQ"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ons_instances"
sidebar_current: "docs-alibabacloudstack-datasource-ons-instances"
description: |-
    查询消息队列服务实例
---

# alibabacloudstack_ons_instances

根据指定过滤条件列出当前凭证权限可以访问的消息队列服务实例列表。


## 示例用法

```
variable "name" {
  default = "onsInstanceDatasourceName"
}

resource "alibabacloudstack_ons_instance" "default" {
  tps_receive_max = 500
  tps_send_max = 500
  topic_capacity = 50
  cluster = "cluster1"
  independent_naming = "true"
  name = var.name
  remark = "Ons Instance"
}

data "alibabacloudstack_ons_instances" "instances_ds" {
  name_regex = alibabacloudstack_ons_instance.inst.name
  output_file = "instances.txt"
}

output "first_instance_id" {
  value = data.alibabacloudstack_ons_instances.instances_ds.*
}
```

## 参数说明

支持以下参数：

* `ids` - (可选) 用于过滤结果的实例ID列表。
* `name_regex` - (可选) 用于通过实例名称过滤结果的正则表达式字符串。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - 实例ID列表。
* `names` - 实例名称列表。
* `instances` - 实例列表。每个元素包含以下属性：
  * `id` - 实例ID。
  * `instance_id` - 实例ID。
  * `instance_name` - 实例名称。
  * `instance_type` - 实例类型。
  * `instance_status` - 实例状态。
  * `independent_naming` - 指示是否启用独立命名空间。
  * `tps_receive_max` - 设置主题在某段时间内的最大消息接收吞吐量（TPS）。
  * `tps_send_max` - 设置主题在某段时间内的最大消息发送吞吐量（TPS）。
  * `topic_capacity` - 主题容量，表示该实例支持的主题数量上限。
  * `cluster` - 集群名称，表示该实例所属的集群。
  * `create_time` - 实例创建时间，格式为标准时间戳。
  * `computed_property_example` - 计算属性示例，表示某些动态计算的结果（如果存在）。