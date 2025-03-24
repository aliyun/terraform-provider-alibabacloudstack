---
subcategory: "RocketMQ"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ons_instance"
sidebar_current: "docs-alibabacloudstack-resource-ons-instance"
description: |-
  编排消息队列（ONS）实例
---

# alibabacloudstack_ons_instance

使用Provider配置的凭证在指定的资源集编排消息队列（ONS）实例。

## 示例用法

### 基础用法

```
resource "alibabacloudstack_ons_instance" "default" {
  tps_receive_max = 500
  tps_send_max = 500
  topic_capacity = 50
  cluster = "cluster1"
  independent_naming = "true"
  name = "Ons_Apsara_instance"
  remark = "Ons Instance"
}

output "inst" {
  value = alibabacloudstack_ons_instance.default.*
}
```

## 参数参考

支持以下参数：


* `name` - (必填) 同一区域内的单个账户的两个实例不能具有相同名称。长度必须为3到64个字符。允许使用中文字符、英文字母、数字和连字符。
* `tps_receive_max` - (必填) 此属性用于设置某个时间段内主题的消息接收每秒事务数(TPS)。
* `tps_send_max` - (必填) 此属性用于设置某个时间段内主题的消息发送每秒事务数(TPS)。
* `topic_capacity` - (必填) 此属性用于设置主题容量。
* `independent_naming` - (必填) 此属性用于定义是否具有独立命名。它只接受布尔值。
* `cluster` - (必填) 此属性用于添加集群名称。
* `remark` - (可选) 此属性是对实例的简要描述。长度不得超过128。
* `instance_type` - (必填) 此属性指定实例的类型。
* `instance_status` - (必填) 此属性指定实例的状态。

## 属性参考

导出以下属性：

* `id` - 上述提供的资源的`key`。
* `instance_type` - 实例的版本。1 表示后付费版本，2 表示铂金版本。
* `instance_status` - 实例的状态。1 表示铂金版本实例正在部署中。2 表示后付费版本实例已过期。5 表示后付费或铂金版本实例正在服务中。7 表示铂金版本实例正在升级且服务可用。
* `create_time` - 资源的创建时间。
* `name` - 此属性表示实例的名称。
* `tps_receive_max` - 此属性表示实例的最大接收 TPS。
* `tps_send_max` - 此属性表示实例的最大发送 TPS。
* `topic_capacity` - 此属性表示实例的主题容量。
* `independent_naming` - 此属性表示实例是否有独立命名。
* `cluster` - 此属性表示与实例关联的集群。
* `remark` - 此属性表示实例的备注。