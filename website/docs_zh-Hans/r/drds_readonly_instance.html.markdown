---
subcategory: "分布式关系型数据库"
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
variable "name" {
  default = "tf-testacc-readonlyinstance-19500"
}

data "alibabacloudstack_drds_instance_specifications" "default" {
  sorted_by = "CPU"
}

resource "alibabacloudstack_drds_instance" "default" {
  description   = var.name
  zone_id       = alibabacloudstack_vpc_vswitch.default.availability_zone
  vswitch_id    = alibabacloudstack_vpc_vswitch.default.id
  specification = data.alibabacloudstack_drds_instance_specifications.default.specifications.0.id
}


data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details              = true
}


resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name       = "${var.name}_vsw"
  vpc_id     = alibabacloudstack_vpc_vpc.default.id
  cidr_block = "172.16.1.0/24"
  zone_id    = data.alibabacloudstack_zones.default.zones.0.id
}




resource "alibabacloudstack_drds_readonly_instance" "default" {
  master_instance_id   = alibabacloudstack_drds_instance.default.id
  zone_id              = alibabacloudstack_vpc_vswitch.default.availability_zone
  instance_charge_type = "PostPaid"
  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
  specification        = data.alibabacloudstack_drds_instance_specifications.default.specifications.0.id
  description          = var.name
}
```

## 参数说明

支持以下参数：

* `description` - (必填) 实例描述。该描述可以包含2到256个字符的字符串。
* `zone_id` - (必填，变更时重建) 启动DRDS实例的可用区ID。
* `specification` - (必填，变更时重建) 用户定义的DRDS实例规格。值范围：
* `instance_charge_type` - (可选，变更时重建) 计费类型。有效值为`PrePaid`(预付费)和`PostPaid`(后付费)。默认为`PostPaid`。
* `vswitch_id` - (必填，变更时重建) 要启动的交换机ID。
* `master_instance_id` - (必填，变更时重建) 需要镜像的Drds实力的ID。

### 超时时间

`timeouts` 块允许您指定某些操作的 [超时时间](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts)：

* `create` - (默认为10分钟)用于创建DRDS实例(直到其达到运行状态)。
* `delete` - (默认为10分钟)用于终止DRDS实例。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `id` - DRDS实例ID。