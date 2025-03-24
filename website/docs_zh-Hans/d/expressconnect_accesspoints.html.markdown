---
subcategory: "ExpressConnect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_accesspoints"
sidebar_current: "docs-Alibabacloudstack-datasource-expressconnect-accesspoints"
description: |- 
  查询高速通道接入点
---

# alibabacloudstack_expressconnect_accesspoints
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_express_connect_access_points`

根据指定过滤条件列出当前凭证权限可以访问的高速通道接入点列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testacc-expressConnectAccessPoints19873"
}

data "alibabacloudstack_expressconnect_accesspoints" "default" {
  ids        = ["ap-cn-qingdao-env17-d01-amtest17"]
  name_regex = "^tf-testacc-expressConnectAccessPoints"
  status     = "Available"
  output_file = "access_points_output.txt"
}
```

## 参数参考

以下参数是支持的：
  * `ids` - (可选，变更时重建) - 用于过滤的接入点 ID 列表。此参数可以帮助您快速定位特定的接入点。
  * `name_regex` - (可选，变更时重建) - 按名称过滤接入点的正则表达式模式。通过该参数可以匹配符合特定命名规则的接入点。
  * `status` - (可选，变更时重建) - 资源的状态。有效值包括：`Available`(可用)、`Unavailable`(不可用)等。
  
## 属性参考

除了上述参数外，还导出以下属性：
  * `points` - 接入点列表。列表中的每个元素是一个包含以下键的映射：
    * `id` - 接入点的唯一标识符。
    * `access_point_id` - 接入点 ID，与 `id` 相同。
    * `access_point_name` - 接入点名称，用于标识接入点的逻辑名称。
    * `attached_region_no` - 接入点所在的地域 ID，表示接入点所属的阿里云区域。
    * `description` - 接入点描述信息，提供关于接入点的详细说明。
    * `host_operator` - 接入点所属的运营商，例如中国电信、中国联通等。
    * `location` - 接入点的位置，通常为物理位置或数据中心的地理信息。
    * `status` - 接入点的资源状态，可能的值包括 `Available`(可用)、`Unavailable`(不可用)等。
    * `type` - 物理专线的网络类型，例如 `MPLS` 或其他网络类型。