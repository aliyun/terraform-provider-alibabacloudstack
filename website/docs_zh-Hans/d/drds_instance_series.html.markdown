---
subcategory: "DRDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_drds_instance_series"
sidebar_current: "docs-alibabacloudstack-drds-instance-series"
description: |-
  查询DRDS实例规格系列
---

# alibabacloudstack_drds_instance_series

根据指定过滤条件列出当前凭证权限可以访问的DRDS实例规格系列列表。

## 示例用法

```
data "alibabacloudstack_drds_instance_series" "default" {
  
}


```

## 参数说明

支持以下参数：

* `ids` - (可选，变更时重建) 指定实例规格系列ID范围。如果未指定，则返回所有可用区的实例类型族。
* `names` - (可选, 变更时重建) 指定实例规格系列名称范围。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `ids` - 包含所有匹配条件的实例规格系列ID的列表。
* `names` - 包含所有匹配条件的实例规格系列名称的列表。
* `series` - 实例规格系列的详细信息列表。每个元素包含以下属性：
  * `id` - 实例规格系列的唯一标识符。
  * `name` - 实例规格系列的名称。
