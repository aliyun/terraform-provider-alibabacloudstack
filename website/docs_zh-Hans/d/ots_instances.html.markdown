---
subcategory: "OTS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ots_instances"
sidebar_current: "docs-Alibabacloudstack-datasource-ots-instances"
description: |- 
  查询表格存储（OTS）实例
---

# alibabacloudstack_ots_instances

根据指定过滤条件列出当前凭证权限可以访问的表格存储（OTS）实例列表。

## 示例用法

```terraform
variable "name" {
  default = "tf-testAcc13773"
}

resource "alibabacloudstack_ots_instance" "default" {
  name          = "${var.name}"
  description   = "${var.name}"
  instance_type = "Capacity"
  tags = {
    Created = "TF-${var.name}"
    For     = "acceptance test"
  }
}

data "alibabacloudstack_ots_instances" "default" {
  ids        = ["${alibabacloudstack_ots_instance.default.id}"]
  name_regex = "^${var.name}-.*$"

  tags = {
    Created = "TF-${var.name}"
    For     = "acceptance test"
  }

  output_file = "instances.txt"
}

output "first_instance_id" {
  value = data.alibabacloudstack_ots_instances.default.instances.0.id
}
```

## 参数参考

以下参数是支持的：

* `ids` - (可选) 实例 ID 列表。如果指定，数据源将仅返回具有这些 ID 的实例。
* `name_regex` - (可选，变更时重建) 用于按实例名称过滤结果的正则表达式字符串。这允许您使用正则表达式匹配实例名称。
* `tags` - (可选) 分配给实例的标签映射。它必须是以下格式：
  ```terraform
  tags = {
    tagKey1 = "tagValue1",
    tagKey2 = "tagValue2"
  }
  ```

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - 匹配指定过滤条件的实例 ID 列表。
* `names` - 匹配指定过滤条件的实例名称列表。
* `instances` - 实例列表。每个元素包含以下属性：
  * `id` - 实例的 ID。
  * `name` - 实例名称。
  * `status` - 实例状态。可能的值：`Running`、`Disabled`、`Deleting`。
  * `write_capacity` - 预留写吞吐量。单位为 CU(Capacity Unit)。只有高性能实例有此返回值。
  * `read_capacity` - 预留读吞吐量。单位为 CU(Capacity Unit)。只有高性能实例有此返回值。
  * `cluster_type` - 实例的集群类型。可能的值：`SSD`、`HYBRID`。
  * `create_time` - 实例的创建时间。
  * `user_id` - 与实例关联的用户 ID。
  * `network` - 实例的网络类型。可能的值：`NORMAL`、`VPC`、`VPC_CONSOLE`。
  * `description` - 实例的描述。
  * `entity_quota` - 实例配额，表示在此实例内可以创建的最大表数。
  * `tags` - 分配给实例的标签。