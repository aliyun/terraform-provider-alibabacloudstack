---
subcategory: "云原生多模数据库 Lindorm"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_lindorm_lts_instance"
sidebar_current: "docs-alibabacloudstack-resource-lindorm-lts-instance"
description: |-
  提供阿里云 Lindorm LTS 实例资源
---

# alibabacloudstack_lindorm_lts_instance

提供 Lindorm LTS 实例资源。

## 示例用法

基础用法

```terraform

variable "name" {
    default = "tf-testacc97984"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

data "alibabacloudstack_lindorm_instance_types" "sortbycpu" {
	sorted_by = "CPU"
	engine_type = "lts"
}

resource "alibabacloudstack_lindorm_lts_instance" "example" {
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  instance_alias = "${var.name}"
  cpu_brand      = "Intel"
  instance_type  = "${data.alibabacloudstack_lindorm_instance_types.sortbycpu.instance_types[0].name}"
  lts_num        = 2
}

```

## 参数说明

以下参数是可支持的：

* `zone_id` - （必选，ForceNew）实例的可用区 ID。
* `instance_alias` - （必选）实例的别名。
* `cpu_brand` - （必选，ForceNew）CPU 品牌。
* `instance_type` - （必选）实例规格
* `lts_num` - （必选）LTS 节点数量。
* `deletion_protection` - （可选，计算得出）是否为实例开启删除保护。默认值：`false`。

## 属性说明

以下属性会被导出：

* `id` - 实例的 ID。
* `instance_id` - 实例的 ID。
* `instance_status` - 实例的状态。
* `create_time` - 实例的创建时间。
* `service_type` - 实例的服务类型。
* `network_type` - 实例的网络类型。

## 导入资源

Lindorm LTS 实例可以使用实例 ID 进行导入，例如：

```bash
$ terraform import alibabacloudstack_lindorm_lts_instance.example <id>
```