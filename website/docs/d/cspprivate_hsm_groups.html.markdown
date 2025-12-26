---
subcategory: "Cspprivate HSM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cspprivate_hsm_group"
sidebar_current: "docs-alibabacloudstack-datasource-cspprivate-hsm-group"
description: |-
    Provides a data source for CSP Private HSM Group to query HSM group information in Alibaba Cloud CSP Private HSM service.
---

# alibabacloudstack_cspprivate_hsm_group

> Data source for CSP Private HSM Group, used to query HSM group information in Alibaba Cloud CSP Private HSM service.

## Example Usage

```hcl

variable "name" {
  default = "tf_hsm_group34787"
}

data "alibabacloudstack_zones" "default" {
  provider                    = alibabacloudstack-common
  available_resource_creation = "VSwitch"
}

data "alibabacloudstack_cspprivate_hsm_vendors" "default" {
}

resource "alibabacloudstack_vpc" "vpc" {
  provider   = alibabacloudstack-common
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16" # VPC CIDR block
}

resource "alibabacloudstack_vswitch" "vsw" {
  provider          = alibabacloudstack-common
  vpc_id            = alibabacloudstack_vpc.vpc.id
  cidr_block        = "192.168.0.0/24" # VSwitch CIDR block
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
}


resource "alibabacloudstack_cspprivate_hsm_instance" "default" {
  product_code   = data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.products.0.code
  vendor_code    = data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.code
  vsm_type       = "gvsm"
  zone_id        = data.alibabacloudstack_zones.default.zones.0.id
  alias_name     = "${var.name}0"
  vpc_id         = alibabacloudstack_vpc.vpc.id
  vpc_cidr_block = alibabacloudstack_vpc.vpc.cidr_block
  vswitch_id     = alibabacloudstack_vswitch.vsw.id
  ip             = "192.168.0.100"
}


resource "random_password" "password" {
  count            = 1
  length           = 12
  special          = true
  override_special = "!@#$^&*()_"
  min_lower        = 1
  min_upper        = 1
  min_numeric      = 1
}

resource "alibabacloudstack_cspprivate_hsm_group" "default" {
  hsm_list  = ["${alibabacloudstack_cspprivate_hsm_instance.default.id}"]
  zone_ids  = ["${data.alibabacloudstack_zones.default.zones.0.id}"]
  password  = random_password.password.0.result
  hsm_count = 1
}

data "alibabacloudstack_cspprivate_hsm_groups" "default" {
  ids = ["${alibabacloudstack_cspprivate_hsm_group.default.id}"]
}

```

## Argument Reference

The following arguments are used to filter and query HSM groups:

* `ids` (List, Optional): A list of HSM group names used to filter results by group name. If specified, only HSM groups with matching names will be returned.

## Attributes Reference

The following attributes are exported:

* `id` (String): The unique identifier of the HSM group, equivalent to the group name.
* `create_time` (String): The creation time of the HSM group, in ISO 8601 standard format.
* `group_name` (String): The name of the HSM group.
* `hsm_count` (Integer): The number of HSMs (Hardware Security Modules) in the HSM group.
* `security_level_tag` (String): The security level tag of the HSM group, such as "public".
* `status` (String): The current status of the HSM group, such as "uninitialized" indicating uninitialized.
* `unique_id` (Integer): The unique ID identifier of the HSM group.
* `update_time` (String): The last update time of the HSM group, in ISO 8601 standard format.
* `vpc_id` (String): The ID of the VPC to which the HSM group belongs.
* `zone_ids` (String): A list of zone IDs where the HSM group is located, represented in JSON array format.