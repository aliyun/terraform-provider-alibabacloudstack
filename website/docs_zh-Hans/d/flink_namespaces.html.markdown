---
subcategory: "Realtime Compute for Apache(Flink)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_flink_namespaces"
sidebar_current: "docs-alibabacloudstack-datasource-flink-namespaces"
description: |-
  查询Flink命名空间
---

# alibabacloudstack_flink_namespaces

根据指定过滤条件列出当前凭证权限可以访问的Flink命名空间列表。



## 示例用法

```
# Declare the data source
data "alibabacloudstack_flink_namespaces" "my_namespaces" {
  name_regex  = "my-namespace"
}

```

## 参数说明

支持以下参数：

* `name_regex` - (可选) 用于通过命名空间名称过滤结果的正则表达式字符串。
* `ids` - (可选) 用于过滤命名空间名称的一个字符串列表。
* `owner_id` - (可选) 根据命名空间管理员ID过滤命名空间。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `ids` - 匹配的命名空间列表。其元素是一个命名空间名称。
* `names` - 命名空间名称列表。
* `namespaces` - 匹配的仓库命名空间列表。每个元素包含以下属性：
  * `name` - 仓库命名空间的名称。
  * `cu` -  仓库命名空间保留的资源数量。
  * `cpu_type` - 资源CPU类型，有效值：`intel`。
  * `owner_uid` - 仓库命名空间管理员ID。