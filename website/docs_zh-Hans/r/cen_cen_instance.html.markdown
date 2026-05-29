---
subcategory: "云企业网"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_ceninstance"
sidebar_current: "docs-Alibabacloudstack-cen-ceninstance"
description: |-
  提供一个CEN实例资源。
---

# alibabacloudstack\_cen\_ceninstance

提供一个CEN实例资源。

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

* `cen_instance_name` -（可选）CEN实例名称。
* `description` -（可选）CEN实例描述。
* `protection_level` -（可选）保护级别。默认值：`REDUCED`。
* `transit_router_name` -（可选）转发路由器名称。
* `transit_router_description` -（可选）转发路由器描述。
* `transit_router_cidrs` -（可选，v3.18+ 支持）转发路由器的CIDR块。最多可配置5个CIDR块。每个项支持：
  * `cidr` -（必填）CIDR块。
  * `cidr_id` -（计算得出）CIDR块ID。
* `tags` -（可选）标签映射。

## 属性参考

除了上述参数外，还导出以下属性：

* `id` - CEN实例ID。
* `cen_id` - CEN实例ID。
* `cen_bandwidth_package_ids` - 与CEN实例关联的带宽包ID列表。
* `create_time` - CEN实例创建时间。
* `status` - CEN实例状态。
* `transit_router_id` - 转发路由器ID。

## Import

CEN实例可以使用CenId导入，例如：

```
$ terraform import alibabacloudstack_cen_instance.example cen-xxxxxxxxxxx
```