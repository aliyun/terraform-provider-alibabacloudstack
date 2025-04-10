---
subcategory: "waf-onecs"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_waf_instance"
sidebar_current: "docs-alibabacloudstack-resource-waf_instance"
description: |-
  Provides a Alibabacloudstack waf-onecs switch resource.
---

# alibabacloudstack\waf-onecs

Provides a waf-onecs instance resource.

## Example Usage

Basic Usage

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

#创建vpc
resource "alibabacloudstack_vpc" "vpc" {
  vpc_name = var.name
  cidr_block = "192.168.0.0/16" #vpc口段
  provider = alibabacloudstack.provider2
}
#创建vsw
resource "alibabacloudstack_vswitch" "vsw" {
  provider = alibabacloudstack.provider2
  vpc_id = alibabacloudstack_vpc.vpc.id
  cidr_block = "192.168.0.0/16" #⽹段
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id #可⽤区
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

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) The VPC ID.
* `vswitch_id` - (Required, ForceNew) VSwitch ID configured to Bastionhost.
* `detector_specs` - (Required)  The Detection Engine Specifications provide detailed information about the capabilities, configuration, and performance characteristics of a detection engine
* `detector_version` - (Required) The Detection Engine level
* `detector_nodenum` - (Required) Number of Detection Engines in a Single Availability Zone
* `name` - (Required, ForceNew) Name of the WAF instance.
* `vpc_vswitch` - (Required, ForceNew) Configuration block for VPC vswitch settings.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the switch.
* `availability_zone` The AZ for the switch.
* `cidr_block` - The CIDR block for the switch.
* `ipv6_cidr_block` - (Optional) The ipv6 cidr block of switch.
* `vpc_id` - The VPC ID.
* `name` - The name of the switch.
* `description` - The description of the switch.
* `arch` - Architecture associated with the WAF instance.
* `cpu_type` - CPU type associated with the WAF instance.
* `wafinstance_id` - ID of the WAF instance.
* `instance_status` - Status of the WAF instance.
* `instance_make_status` - Make status of the WAF instance.