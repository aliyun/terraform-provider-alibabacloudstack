---
subcategory: "RocketMQ (ONS)"
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

## 参数参考

支持以下参数：

* `ids` - (可选) 用于过滤结果的实例ID列表。
* `name_regex` - (可选) 用于通过实例名称过滤结果的正则表达式字符串。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - 实例ID列表。
* `names` - 实例名称列表。
* `instances` - 实例列表。每个元素包含以下属性：
  * `id` - 实例ID。
  * `instance_id` - 实例ID。
  * `instance_name` - 实例名称。
  * `instance_type` - 实例类型。
  * `instance_status` - 实例状态。
  * `independent_naming` - 指示是否可用命名空间。
  * `tps_receive_max` - 此属性用于设置主题在某段时间内的消息接收每秒事务数(TPS)。
  * `tps_send_max` - 此属性用于设置主题在某段时间内的消息发送每秒事务数(TPS)。
  * `topic_capacity` - 此属性用于设置主题容量。
  * `cluster` - 此属性用于添加集群名称。
  * `create_time` - 实例创建时间。
  * `computed_property_example` - 计算属性示例。