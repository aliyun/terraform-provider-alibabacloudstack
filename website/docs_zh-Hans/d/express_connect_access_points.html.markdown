---
subcategory: "Express Connect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_express_connect_access_points"
sidebar_current: "docs-alibabacloudstack-datasource-express-connect-access-points"
description: |-
  提供给当前阿里云用户 Express Connect 接入点列表。
---

# alibabacloudstack_express_connect_access_points
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_expressconnect_accesspoints`

此数据源提供当前阿里云用户的 Express Connect 接入点列表。


## 示例用法

### 基础用法

```terraform
data "alibabacloudstack_express_connect_access_points" "ids" {
  ids = ["ap-cn-hangzhou-yh-C"]
}
output "express_connect_access_point_id_1" {
  value = data.alibabacloudstack_express_connect_access_points.ids.points.0.id
}

data "alibabacloudstack_express_connect_access_points" "nameRegex" {
  name_regex = "^杭州-"
}
output "express_connect_access_point_id_2" {
  value = data.alibabacloudstack_express_connect_access_points.nameRegex.points.0.id
}

```

## 参数参考

支持以下参数：

* `ids` - (可选，变更时重建) 接入点 ID 列表。
* `name_regex` - (可选，变更时重建) 按接入点名称使用正则表达式过滤结果。
* `status` - (可选，变更时重建) 接入点所关联的物理连接状态。有效值：`disabled`, `full`, `hot`, `recommended`。
* `names` - (可选，变更时重建) 接入点名称列表。

## 属性参考

除了上述参数外，还导出以下属性：

* `names` - 接入点名称列表。
* `points` - Express Connect 接入点列表。每个元素包含以下属性：
  * `access_point_id` - 接入点 ID。
  * `access_point_name` - 接入点名称。
  * `attached_region_no` - 接入点所在区域 ID。
  * `description` - 接入点描述。
  * `host_operator` - 接入点所属运营商。
  * `id` - 接入点 ID。
  * `location` - 接入点位置。
  * `status` - 接入点所关联的物理连接状态。
  * `type` - 接入点所关联的网络类型。