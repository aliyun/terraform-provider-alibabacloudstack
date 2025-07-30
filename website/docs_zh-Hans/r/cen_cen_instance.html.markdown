---
subcategory: "CEN"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_ceninstance"
sidebar_current: "docs-Alibabacloudstack-cen-ceninstance"
description: |-
  提供一个cen实例资源。
---

# alibabacloudstack\_cen\_ceninstance

提供一个cen实例资源。

## 示例用法
variable "name" { default = "tf-testaccceninstance48958" }

resource "alibabacloudstack_cen_instance" "default" { description = "tf-testaccceninstance48958" cen_instance_name = "tf-testaccceninstance48958" }


## 参数参考

支持以下参数：
  * `cen_instance_name` - (可选) - cen实例名称。
  * `description` - (可选) - cen实例描述。
  * `protection_level` - (可选) - cen实例保护级别。
  * `status` - (可选) - cen实例状态。

## 属性参考

除上述参数外，还导出以下属性：
  * `cen_id` - cen实例ID。
  * `create_time` - cen实例创建时间。
  * `status` - cen实例状态。