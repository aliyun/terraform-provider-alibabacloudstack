---
subcategory: "MQTT"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mqtt_instance"
sidebar_current: "docs-Alibabacloudstack-datasource-mqtt-instance"
description: |-
    编排 MQTT 实例资源
---
# alibabacloudstack_mqtt_instance

使用Provider配置的凭证在指定的区域创建和管理MQTT实例。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tftestacc729"
}

resource "alibabacloudstack_ons_instance" "default" {
  tps_receive_max    = 500
  tps_send_max       = 500
  topic_capacity     = 50
  cluster            = "cluster1"
  independent_naming = "true"
  name               = "${var.name}MQ"
  remark             = "Ons_instance"
}



resource "alibabacloudstack_mqtt_instance" "default" {
  max_conn          = "1000"
  max_up_tps        = "500"
  max_down_tps      = "1000"
  remark            = "test"
  store_instance_id = alibabacloudstack_ons_instance.default.id
  max_sub           = "10000"
  instance_name     = var.name
  cluster_name      = "mqtt4Private"
}
```

## 参数说明

支持以下参数：

* `instance_name` - (必填) MQTT实例的名称。名称长度为1到128个字符，不能以`http://`或`https://`开头。
* `cluster_name` - (必填, 变更时重建) 集群名称。默认值为"mqtt4Private"，表示使用私有集群。
* `independent_naming` - (可选) 是否使用独立命名空间。默认值为`true`，表示启用独立命名空间。
* `max_conn` - (可选) 最大连接数。取值范围：1000-100000。
* `max_down_tps` - (可选) 最大下行TPS（每秒事务数）。取值范围：1000-100000。
* `max_sub` - (可选) 最大订阅数。默认值为1000000。
* `max_up_tps` - (可选) 最大上行TPS（每秒事务数）。取值范围：500-20000。
* `remark` - (可选) 实例备注信息。描述长度为1到256个字符。
* `store_instance_id` - (可选) 存储实例ID。用于绑定已有的消息队列RocketMQ实例。
* `store_type` - (可选) 存储类型。默认值为1，表示使用标准存储类型。

## 属性说明

以下属性会从API响应中导出：

* `instance_id` - MQTT实例ID。
* `endpoints` - 实例的端点信息，包含内部和外部访问地址，如`tcpInternalEndpoint`、`tcpInternetEndpoint`等。
* `namespace_id` - 命名空间ID，通常与实例ID相同。