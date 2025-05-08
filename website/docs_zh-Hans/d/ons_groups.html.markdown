---
subcategory: "RocketMQ (ONS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ons_groups"
sidebar_current: "docs-alibabacloudstack-datasource-ons-groups"
description: |-
    查询消息队列服务组
---

# alibabacloudstack_ons_groups

根据指定过滤条件列出当前凭证权限可以访问的消息队列服务组列表。


## 示例用法

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
  instance_id = "${alibabacloudstack_ons_instance.default.id}"
  remark = "dafault_ons_group_remark"
}

data "alibabacloudstack_ons_groups" "default" {
  instance_id = alibabacloudstack_ons_group.default.instance_id

}
output "onsgroups" {
  value = data.alibabacloudstack_ons_groups.default.*
}
```

## 参数说明

以下支持以下参数：

* `instance_id` - (必填) 拥有这些组的 ONS 实例 ID。
* `group_id_regex` - (可选) 用于通过组名筛选结果的正则表达式字符串。此参数允许用户根据组名的部分匹配来过滤结果。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - 组名称列表。此列表包含所有符合条件的组名称。
* `groups` - 组列表。每个元素包含以下属性：
  * `id` - 组的名称。这是组的唯一标识符。
  * `group_id` - 组的 ID。此字段与 `id` 类似，但可能包含额外的命名空间信息。
  * `owner` - 组所有者的 ID，即 Apsara Stack Cloud UID。此字段标识了该组的创建者或拥有者。
  * `instance_id` - 命名空间的 Id。此字段表示该组所属的 ONS 实例 ID。
  * `independent_naming` - 表示命名空间是否可用。如果值为 `true`，则表示该组使用独立命名空间；否则，表示共享命名空间。
  * `remark` - 组的备注。此字段提供了关于该组的额外描述信息。
  * `create_time` - 组的创建时间。此字段以时间戳或日期格式表示该组的创建时间。