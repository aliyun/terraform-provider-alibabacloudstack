---
subcategory: "DRDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_drds_instance"
sidebar_current: "docs-Alibabacloudstack-drds-instance"
description: |- 
  编排云原生分布式数据库（Drds）实例
---

# alibabacloudstack_drds_instance

使用Provider配置的凭证在指定的资源集下编排云原生分布式数据库（Drds）实例。

## 示例用法

```hcl
provider "apsarastack" {
	assume_role {}
}

variable "name" {
	default = "tf-testaccDrdsdatabase-14880"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation = "VSwitch"
}

variable "instance_series" {
	default = "drds.sn2.4c16g"
}

resource "alibabacloudstack_vpc" "default" {
	name       = var.name
	cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
	vpc_id            = alibabacloudstack_vpc.default.id
	cidr_block        = "172.16.0.0/24"
	availability_zone = data.alibabacloudstack_zones.default.zones.0.id
	name              = var.name
}

resource "alibabacloudstack_drds_instance" "default" {
	instance_charge_type = "PostPaid"
	vswitch_id           = alibabacloudstack_vswitch.default.id
	specification        = "drds.sn2.4c16g.8C32G"
	description          = var.name
	zone_id              = alibabacloudstack_vswitch.default.availability_zone
	instance_series      = var.instance_series
}
```

## 参数参考

支持以下参数：

* `description` - (必填) 实例描述。该描述可以包含2到256个字符的字符串。
* `zone_id` - (必填，变更时重建) 启动DRDS实例的可用区ID。
* `specification` - (必填，变更时重建) 用户定义的DRDS实例规格。值范围：
    - 对于`drds.sn1.4c8g`(入门版)：
        - `drds.sn1.4c8g.8c16g`, `drds.sn1.4c8g.16c32g`, `drds.sn1.4c8g.32c64g`, `drds.sn1.4c8g.64c128g`
    - 对于`drds.sn1.8c16g`(标准版)：
        - `drds.sn1.8c16g.16c32g`, `drds.sn1.8c16g.32c64g`, `drds.sn1.8c16g.64c128g`
    - 对于`drds.sn1.16c32g`(企业版)：
        - `drds.sn1.16c32g.32c64g`, `drds.sn1.16c32g.64c128g`
    - 对于`drds.sn1.32c64g`(极致版)：
        - `drds.sn1.32c64g.128c256g`
* `instance_charge_type` - (可选，变更时重建) 计费类型。有效值为`PrePaid`(预付费)和`PostPaid`(后付费)。默认为`PostPaid`。
* `vswitch_id` - (必填，变更时重建) 要启动的交换机ID。
* `instance_series` - (必填，变更时重建) 用户定义的DRDS实例节点规格。值范围：
    - `drds.sn1.4c8g` 用于DRDS实例入门版；
    - `drds.sn1.8c16g` 用于DRDS实例标准版；
    - `drds.sn1.16c32g` 用于DRDS实例企业版；
    - `drds.sn1.32c64g` 用于DRDS实例极致版；

### 超时时间

`timeouts` 块允许您指定某些操作的 [超时时间](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts)：

* `create` - (默认为10分钟)用于创建DRDS实例(直到其达到运行状态)。
* `delete` - (默认为10分钟)用于终止DRDS实例。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `id` - DRDS实例ID。