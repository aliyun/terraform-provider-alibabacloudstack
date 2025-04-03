---
subcategory: "RocketMQ"
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

## 参数参考

以下支持以下参数：

* `instance_id` - (必填) 拥有这些组的 ONS 实例 ID。
* `group_id_regex` - (可选) 用于通过组名筛选结果的正则表达式字符串。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - 组名称列表。
* `groups` - 组列表。每个元素包含以下属性：
  * `id` - 组的名称。
  * `group_id` - 组的 ID。
  * `owner` - 组所有者的 ID，即 Apsara Stack Cloud UID。
  * `instance_id` - 命名空间的 Id。
  * `independent_naming` - 表示命名空间是否可用。
  * `remark` - 组的备注。
  * `create_time` - 组的创建时间。