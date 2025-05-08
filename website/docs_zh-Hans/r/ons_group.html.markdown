---
subcategory: "RocketMQ (ONS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ons_group"
sidebar_current: "docs-alibabacloudstack-resource-ons-group"
description: |-
  编排消息队列（ONS）组资源。

---

# alibabacloudstack_ons_group

使用Provider配置的凭证在指定的资源集编排消息队列（ONS）组资源。


## 示例用法

### 基础用法

```
variable "name" {
  default = "onsInstanceName"
}

variable "group_id" {
  default = "GID-onsGroupDatasourceName"
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

resource "alibabacloudstack_ons_group" "default" {
  group_id = var.group_id
  instance_id = alibabacloudstack_ons_instance.default.id
  remark = "dafault_ons_group_remark"
}
```

## 参数说明

支持以下参数：

* `instance_id` - (必填) 拥有该组的ONS实例的ID。
* `group_id` - (必填) 组的名称。单个实例上的两个组不能具有相同的名称。`group_id`以"GID_"或"GID-"开头，并包含字母、数字、连字符(-)和下划线(_)。
* `remark` - (可选) 此属性是对组的简要描述。长度不得超过256。
* `read_enable` - (可选) 此属性用于设置消息读取是否启用或禁用。只有在组被客户端使用后才能设置。默认值为`true`，表示启用消息读取功能。如果设置为`false`，则禁用消息读取功能。

## 属性说明

导出以下属性：

* `id` - ONS组的唯一标识符，由GroupID和InstanceID组成。格式为`GroupID:InstanceID`，其中`GroupID`是组的ID，`InstanceID`是所属ONS实例的ID。