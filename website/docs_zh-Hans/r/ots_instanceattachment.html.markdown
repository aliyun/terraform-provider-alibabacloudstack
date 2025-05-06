---
subcategory: "Table Store (OTS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ots_instanceattachment"
sidebar_current: "docs-Alibabacloudstack-ots-instanceattachment"
description: |- 
  编排绑定表格存储服务(OTS）实例和网络
---

# alibabacloudstack_ots_instanceattachment
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_ots_instance_attachment`

使用Provider配置的凭证在指定的资源集编排绑定表格存储服务(OTS）实例和网络。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAcc84399"
}

resource "alibabacloudstack_ots_instance" "default" {
  name        = "${var.name}"
  description = "${var.name}"
  accessed_by = "Vpc"
  instance_type = "Capacity"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  cidr_block = "172.16.0.0/16"
  name       = "${var.name}"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  vswitch_name      = "${var.name}"
  cidr_block        = "172.16.1.0/24"
  zone_id           = data.alibabacloudstack_zones.default.zones[0].id
}

resource "alibabacloudstack_ots_instance_attachment" "default" {
  instance_name = alibabacloudstack_ots_instance.default.name
  vpc_name      = "test-attachment"
  vswitch_id    = alibabacloudstack_vswitch.default.id
}
```

## 参数说明

支持以下参数：
* `instance_name` - (必填, 变更时重建) OTS实例的名称。这必须与现有的OTS实例名称匹配。
* `vpc_name` - (必填, 变更时重建) 要附加到OTS实例的VPC的名称。这用于标识目的。
* `vswitch_id` - (必填, 变更时重建) 要附加到OTS实例的交换机的ID。这必须在与`vpc_name`指定的相同VPC内。
* `vpc_id` - (可选, 变更时重建) 要附加到OTS实例的VPC的ID。

## 属性说明

除了上述所有参数外，还导出了以下属性：
* `id` - 资源ID。其值与`instance_name`相同。
* `instance_name` - OTS实例的名称。
* `vpc_name` - 附加到OTS实例的VPC的名称。
* `vswitch_id` - 附加到OTS实例的交换机的ID。
* `vpc_id` - 附加到OTS实例的VPC的ID。此属性是从`vswitch_id`自动派生的。