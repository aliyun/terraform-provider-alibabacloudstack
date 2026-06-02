---
subcategory: "租户侧堡垒机"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bastionhost_instance"
sidebar_current: "docs-Alibabacloudstack-resource-bastionhost-instance"
description: |-
  提供一个 Alibabacloudstack 堡垒机实例资源
---

# alibabacloudstack\bastionhostprivate

提供一个堡垒机实例资源，用于在 Alibaba Cloud Stack 环境中创建和管理阿里云 BastionHost 实例。

示例用法
以下示例通过描述正则表达式检索 Bastionhost 实例，并将结果写入文件：

```hcl
provider "alibabacloudstack" {
  alias = "provider1"
  popgw_domain = "xx"
  access_key   = "xx"
  secret_key   = "xx"
  region       = "xx"
  proxy        = "xx"
  protocol     = "xx"
  insecure     = "xx"
  resource_group_set_name = "xx"
  role_arn     = "xx"
}

provider "alibabacloudstack" {
  alias = "provider2"
  popgw_domain = "xx"
  access_key   = "xx"
  secret_key   = "xx"
  region       = "xx"
  proxy        = "xx"
  protocol     = "xx"
  insecure     = "xx"
  resource_group_set_name = "xx"
}

variable "name" {
  default = "terraform_test"
}

# 查询可用区
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  provider = alibabacloudstack.provider2
}

# 创建 VPC
resource "alibabacloudstack_vpc" "vpc" {
  vpc_name = var.name
  cidr_block = "192.168.0.0/16"
  provider = alibabacloudstack.provider2
}

# 创建 VSwitch
resource "alibabacloudstack_vswitch" "vsw" {
  provider = alibabacloudstack.provider2
  vpc_id = alibabacloudstack_vpc.vpc.id
  cidr_block = "192.168.0.0/16"
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
}

# 创建堡垒机实例
resource "alibabacloudstack_bastionhost_instance" "default" {
  vswitch_id = alibabacloudstack_vswitch.vsw.id
  license_code = "bastionhostah_small_lic"
  vpc_id = alibabacloudstack_vpc.vpc.id
  asset = 50
  highavailability = false
  disasterrecovery = false
  provider = alibabacloudstack.provider1
}
```

## 参数说明
支持以下参数：

* `vpc_id` -（必需，变更时强制重建）VPC 的 ID。
* `vswitch_id` -（必需，变更时强制重建）配置到堡垒机的 VSwitch ID。
* `license_code` -（必需）云堡垒机实例的套餐类型。你可以通过 * `DescribePricingModule` 查询更多支持的套餐类型。
* `asset` -（必需，变更时强制重建）通常指通过堡垒机管理的资源或目标，例如服务器、数据库和网络设备等。堡垒机的主要功能是提供一个安全的访问入口，并对这些资产的访问进行控制和审计。
* `highavailability` -（必需，变更时强制重建）高可用性（HA）是指系统或服务能够在长时间内持续运行而不中断的能力。它通过冗余、故障检测与恢复、负载均衡等技术实现，确保系统在面对各种故障时仍能提供可靠的服务。
* `disasterrecovery` -（必需，变更时强制重建）容灾（DR）是一组策略、工具和流程，用于在遭遇自然灾害或人为灾难后恢复关键的技术基础设施和系统。容灾的主要目标是最小化突发事件的影响，并确保关键业务功能能够尽快恢复。
属性输出
以下属性被导出：

* `id` - 实例的 ID。
* `availability_zone` - 实例所在的可用区。
* `cidr_block` - VSwitch 的 CIDR 网段。
* `ipv6_cidr_block` -（可选）VSwitch 的 IPv6 CIDR 网段。
* `vpc_id` - 所属 VPC 的 ID。
* `name` - 实例的名称。
* `description` - 实例的描述信息。


