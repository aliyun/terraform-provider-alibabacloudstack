---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_ceninstance"
sidebar_current: "docs-Alibabacloudstack-cen-ceninstance"
description: |-
  提供一个cen实例资源。
---

# alibabacloudstack\_cen\_ceninstance

提供一个cen实例资源。

## 示例用法
```
variable "name" {
	default = "tf-testaccceninstance48958"
}


resource "alibabacloudstack_cen_instance" "default" {
  description = "tf-testaccceninstance48958"
  cen_instance_name = "tf-testaccceninstance48958"
}
```
## 参数参考
支持以下参数：

* `cen_instance_name` -（可选）- cen实例名称。
* `description` -（可选）- cen实例描述。
* `protection_level` -（可选）- 保护级别。
* `status` -（计算得出）- cen实例状态。
* `transit_router_name` -（可选）- 转发路由器名称。
* `transit_router_description` -（可选）- 转发路由器描述。
* `transit_router_cidrs` -（可选）- 转发路由器CIDR。
## 属性参考
除了上述参数外，还导出以下属性：

* `cen_id` - cen实例ID。
* `create_time` - 创建时间。
* `status` - cen实例状态。s
* `transit_router_id` - 转发路由器ID。