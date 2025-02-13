---
subcategory: "bastionhostprivate"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bastionhost_instance"
sidebar_current: "docs-alibabacloudstack-resource-bastionhost_instance"
description: |-
  Provides a Alibabacloudstack bastionhostprivate switch resource.
---

# alibabacloudstack\bastionhostprivate

Provides a bastionhost instance resource.

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


resource "alibabacloudstack_bastionhost_instance" "default" {
  vswitch_id = alibabacloudstack_vswitch.vsw.id
  # license_code = "bastionhost4sfd_small_lic" 
  license_code = "bastionhostah_small_lic"
  vpc_id = alibabacloudstack_vpc.vpc.id
  asset = "50"
  highavailability = "false"
  disasterrecovery = "false"
  provider = alibabacloudstack.provider1
}

```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) The VPC ID.
* `vswitch_id` - (Required, ForceNew) VSwitch ID configured to Bastionhost.
* `license_code` - (Required)  The package type of Cloud Bastionhost instance. You can query more supported types through the [DescribePricingModule](https://help.aliyun.com/document_detail/96469.html).
* `asset` - (Required, ForceNew) These typically refer to managed resources or targets, such as servers, databases, and network devices, which are accessed and managed through the bastion host. The primary function of a bastion host is to provide a secure access point, control, and audit access to these assets.
* `highavailability` - (Required, ForceNew) High availability (HA) refers to the ability of a system or service to operate continuously over a long period without interruption. It is achieved through techniques such as redundancy, fault detection and recovery, and load balancing, ensuring that the system can provide reliable service even in the face of various failures
* `disasterrecovery` - (Required, ForceNew) Disaster recovery (DR) is a set of policies, tools, and procedures that enable the recovery or continuation of vital technology infrastructure and systems following a natural or human-induced disaster. The primary goal of disaster recovery is to minimize the impact of a disruptive event and ensure that critical business functions can be restored as quickly as possible


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the switch.
* `availability_zone` The AZ for the switch.
* `cidr_block` - The CIDR block for the switch.
* `ipv6_cidr_block` - (Optional) The ipv6 cidr block of switch.
* `vpc_id` - The VPC ID.
* `name` - The name of the switch.
* `description` - The description of the switch.


