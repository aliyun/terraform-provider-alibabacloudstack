---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_adb_zones"
sidebar_current: "docs-alibabacloudstack-datasource-adb-zones"
description: |-
    查询ADB可用区
---

# alibabacloudstack_adb_zones

根据指定过滤条件列出当前凭证权限可以访问的ADB可用区列表。

## 示例用法

```
# 声明数据源
data "alibabacloudstack_adb_zones" "zones_ids" {}
```

## 参数参考

支持以下参数：

* `multi` - (可选) 指示这些可用区是否可以用于多AZ配置。默认为`false`。多AZ通常用于启动ADB实例。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - 区域ID列表。
* `zones` - 可用区列表。每个元素包含以下属性：
  * `id` - 区域的ID。
  * `multi_zone_ids` - 多区域中的区域ID列表。