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

## 参数说明

支持以下参数：

* `ids` - (可选，变更时重建) 接入点 ID 列表。通过指定此参数，可以筛选出特定的接入点。
* `name_regex` - (可选，变更时重建) 使用正则表达式按接入点名称过滤结果。例如，可以通过设置 `^杭州-` 来筛选名称以“杭州-”开头的接入点。
* `status` - (可选，变更时重建) 接入点所关联的物理连接状态。有效值包括：`disabled`（禁用）、`full`（满载）、`hot`（热门）、`recommended`（推荐）。
* `names` - (可选，变更时重建) 接入点名称列表。通过指定此参数，可以筛选出特定名称的接入点。

## 属性说明

除了上述参数外，还导出以下属性：

* `names` - 接入点名称列表。此属性返回所有匹配的接入点名称。
* `points` - Express Connect 接入点列表。每个元素包含以下属性：
  * `access_point_id` - 接入点 ID。用于唯一标识一个接入点。
  * `access_point_name` - 接入点名称。表示接入点的名称。
  * `attached_region_no` - 接入点所在区域 ID。表示接入点所属的地域。
  * `description` - 接入点描述。提供对接入点的详细说明。
  * `host_operator` - 接入点所属运营商。表示接入点由哪个运营商提供支持。
  * `id` - 接入点 ID。与 `access_point_id` 相同，用于唯一标识接入点。
  * `location` - 接入点位置。表示接入点的实际地理位置。
  * `status` - 接入点所关联的物理连接状态。表示接入点当前的状态，可能的值为 `disabled`, `full`, `hot`, `recommended`。
  * `type` - 接入点所关联的网络类型。表示接入点支持的网络类型。
  * `attached_region_no` - 接入点所属区域编号。表示接入点绑定的区域编号。
  * `description` - 接入点描述信息。提供关于接入点的额外说明。
  * `host_operator` - 接入点所属的运营商。表示接入点由哪个运营商提供支持。
  * `location` - 接入点的物理位置。表示接入点所在的地理位置。