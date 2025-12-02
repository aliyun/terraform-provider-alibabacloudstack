---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_ceninstances"
sidebar_current: "docs-Alibabacloudstack-datasource-cen-ceninstances"
description: |-
  提供阿里云账户拥有的cen实例列表。
---

# alibabacloudstack\_cen\_ceninstances

本数据源根据指定的过滤条件，提供阿里云账户中的cen实例列表。

## 示例用法
```
variable "name" {
  default = "tf-testAccCenInstancesDatasource15452"
}

resource "alibabacloudstack_cen_instance" "default" {
    cen_instance_name = "${var.name}"
	description = "${var.name}"
}

data "alibabacloudstack_cen_instances" "default" {
	name_regex = "${alibabacloudstack_cen_instance.default.cen_instance_name}"
}
```


## 参数参考
支持以下参数：
* `ids` - （可选）- cen实例ID。
* `name_regex` - （可选）- cen实例名称正则表达式。
* `description_regex` - （可选）- cen实例描述正则表达式。
* `transit_router_name_regex` - （可选）- cen实例转发路由器名称正则表达式。
* `transit_router_description_regex` - （可选）- cen实例转发路由器描述正则表达式。
* `cidr` - （可选）- cen实例路由CIDR。
## 属性参考
除了上述参数外，还导出以下属性：

* `cens` - cen实例列表。