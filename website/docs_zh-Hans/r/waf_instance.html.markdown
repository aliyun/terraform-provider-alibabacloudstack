---
subcategory: "Web应用防火墙"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_waf_instance"
sidebar_current: "docs-alibabacloudstack-resource-waf_instance"
description: |-
  提供Alibabacloudstack waf-onecs交换机资源。
---
# alibabacloudstack\waf-onecs

提供一个 waf-onecs 实例资源

## 示例用法


```
provider "alibabacloudstack" {
  alias = "provider1"
  popgw_domain = "xx"
  access_key   = "xx"
  secret_key   = "xx"
  region                  = "xx"
  proxy                   = "xx"
  protocol                = "xx"
  insecure                = "xx"
  resource_group_set_name = "xx"
  role_arn = "xx"
}

provider "alibabacloudstack" {
  alias = "provider2"
  popgw_domain = "xx"
  # access_key   = "xx"
  # secret_key   = "xx"
  access_key   = "xx"
  secret_key   = "xx"
  region                  = "xx"
  proxy                   = "xx"
  # proxy                   = "xx"
  protocol                = "xx"
  insecure                = "xx"
  resource_group_set_name = "xx"
}

variable "name" {
  default = "terraform_test"
}

# 查询可用域
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  provider = alibabacloudstack.provider2
}

# 创建vpc
resource "alibabacloudstack_vpc" "vpc" {
  vpc_name = var.name
  cidr_block = "192.168.0.0/16" # vpc网段
  provider = alibabacloudstack.provider2
}

# 创建vsw
resource "alibabacloudstack_vswitch" "vsw" {
  provider = alibabacloudstack.provider2
  vpc_id = alibabacloudstack_vpc.vpc.id
  cidr_block = "192.168.0.0/16" # 网段
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id # 可用区
}

resource "alibabacloudstack_waf_instance" "default" {
  provider = alibabacloudstack.provider1
  vswitch_id = alibabacloudstack_vswitch.vsw.id
  name = "terraform_test"
  detector_specs = "exclusive"
  vpc_id = alibabacloudstack_vpc.vpc.id
  detector_version = "basic"
  detector_nodenum = 2
  vpc_vswitch {
    vswitch_name = alibabacloudstack_vswitch.vsw.vswitch_name
    vswitch      = alibabacloudstack_vswitch.vsw.id
    cidr_block   = alibabacloudstack_vswitch.vsw.cidr_block
    available_zone = alibabacloudstack_vswitch.vsw.zone_id
    vpc          = alibabacloudstack_vpc.vpc.id
    vpc_name      = alibabacloudstack_vpc.vpc.vpc_name
  }
}
```
## 参数说明
## 支持以下参数：

* `vpc_id` - (必填，强制更新) VPC ID。
* `vswitch_id` - (必填，强制更新) 配置到堡垒机的 VSwitch ID。
* `detector_specs` - (必填) 检测引擎规格，提供了检测引擎的能力、配置和性能特性的详细信息。
* `detector_version` - (必填) 检测引擎等级。
* `detector_nodenum` - (必填) 单个可用区中检测引擎的数量。
* `name` - (必填，强制更新) WAF 实例的名称。
* `vpc_vswitch` - (必填，强制更新) VPC vswitch 设置的配置块。
## 属性输出
导出以下属性：

* `id` - switch 的 ID。
* `availability_zone` - switch 的可用区。
* `cidr_block` - switch 的 CIDR 块。
* `ipv6_cidr_block` - (可选) switch 的 IPv6 CIDR 块。
* `vpc_id` - VPC ID。
* `name` - switch 的名称。
* `description` - switch 的描述。
* `arch` - WAF 实例关联的架构。
* `cpu_type` - WAF 实例关联的 CPU 类型。
* `wafinstance_id` - WAF 实例的 ID。
* `instance_status` - WAF 实例的状态。
* `instance_make_status` - WAF 实例的生成状态。