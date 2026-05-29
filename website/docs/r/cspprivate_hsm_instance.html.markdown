---
subcategory: "Cspprivate HSM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cspprivate_hsm_instance"
sidebar_current: "docs-alibabacloudstack-resource-cspprivate-hsm-instance"
description: |-
  Provides a Alibaba Cloud Hardware Security Module (HSM) Instance resource.
---

# alibabacloudstack_cspprivate_hsm_instance

Provides a Hardware Security Module (HSM) Instance resource.

## Example Usage

### Basic Usage

```hcl
variable "name" {
  default = "tf-example"
}

data "alibabacloudstack_zones" "default" {
  enable_details = true
}

data "alibabacloudstack_cspprivate_hsm_vendors" "default" {
}

resource "alibabacloudstack_vpc" "vpc" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alibabacloudstack_vswitch" "vsw" {
  vpc_id       = alibabacloudstack_vpc.vpc.id
  cidr_block   = "192.168.0.0/24"
  zone_id      = data.alibabacloudstack_zones.default.zones.0.id
}

resource "alibabacloudstack_cspprivate_hsm_instance" "default" {
  product_code   = data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.products.0.code
  vendor_code    = data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.code
  vsm_type       = "gvsm"
  zone_id        = data.alibabacloudstack_zones.default.zones.0.id
  alias_name     = var.name
  vpc_id         = alibabacloudstack_vpc.vpc.id
  vpc_cidr_block = alibabacloudstack_vpc.vpc.cidr_block
  vswitch_id     = alibabacloudstack_vswitch.vsw.id
  ip             = "192.168.0.100"
}
```

### Usage with specific device

```hcl
variable "name" {
  default = "tf-example"
}

data "alibabacloudstack_zones" "default" {
  enable_details = true
}

data "alibabacloudstack_cspprivate_hsm_vendors" "default" {
}

data "alibabacloudstack_cspprivate_hsms" "default" {
  vendor_code  = data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.code
  product_code = data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.products.0.code
  zone_id      = data.alibabacloudstack_zones.default.zones.0.id
}

resource "alibabacloudstack_cspprivate_hsm_instance" "default" {
  product_code = data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.products.0.code
  vendor_code  = data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.code
  vsm_type     = "gvsm"
  zone_id      = data.alibabacloudstack_zones.default.zones.0.id
  alias_name   = var.name
  device_id    = data.alibabacloudstack_cspprivate_hsms.default.hsms.0.id
}
```

## Argument Reference

The following arguments are supported:

* `product_code` - (Required, ForceNew) The product model code of the HSM device.
* `vendor_code` - (Required, ForceNew) The vendor code of the HSM device.
* `vsm_type` - (Required, ForceNew) The type of the HSM instance. Valid values: `evsm`, `gvsm`, `svsm`.
* `zone_id` - (Required, ForceNew) The zone ID where the HSM instance is located.
* `device_id` - (Optional, Computed) The ID of the HSM device.
* `vpc_id` - (Optional, Computed) The ID of the VPC where the HSM instance is located.
* `vpc_cidr_block` - (Optional, Computed) The CIDR block of the VPC.
* `vswitch_id` - (Optional, Computed) The ID of the VSwitch where the HSM instance is located.
* `ip` - (Optional, Computed) The IP address of the HSM instance.
* `alias_name` - (Optional) The alias name of the HSM instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the HSM instance.
* `port` - The port of the HSM instance.
* `status` - The status of the HSM instance.

## Import

HSM Instance can be imported using the id, e.g.

```shell
$ terraform import alibabacloudstack_cspprivate_hsm_instance.example <id>
```