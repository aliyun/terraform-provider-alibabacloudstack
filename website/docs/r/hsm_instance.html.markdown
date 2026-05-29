---
subcategory: "Hardware Security Module (HSM)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_hsm_instance"
sidebar_current: "docs-Alibabacloudstack-hsm-instance"
description: |-
  Create and manage HSM (Hardware Security Module) instances
---

# alibabacloudstack_hsm_instance

Create and manage Alibaba Cloud HSM (Hardware Security Module) instances.

## Example Usage

### Basic Usage

```hcl

variable "name" {
  default = "test-tf-hsm-instance78983"
}

data "alibabacloudstack_zones" "default" {
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_hsm_instance" "default" {
  product_code = "jnta.SJJ1528"
  vendor_code  = "jnta"
  vsm_type     = "gvsm"
  zone_no      = data.alibabacloudstack_zones.default.zones.0.id
  remark       = var.name
  vpc_id       = "${alibabacloudstack_vpc_vpc.default.id}"
  vswitch_id   ="${alibabacloudstack_vpc_vswitch.default.id}"
  ip           = "172.16.1.100"
  white_list   = "192.168.1.0/24"
}
```

## Argument Reference

The following arguments are supported:

* `product_code` - (Required, Forces new resource) The product code of the HSM instance. You can obtain this parameter through the DescribeProducts API.
* `vendor_code` - (Required, Forces new resource) The vendor code of the HSM instance. You can obtain this parameter through the DescribeVendors API.
* `vsm_type` - (Required, Forces new resource) The type of the HSM instance. Valid values:
  * `evsm`: Financial data HSM.
  * `gvsm`: General server HSM.
  * `svsm`: Signature verification server HSM.
* `zone_no` - (Required, Forces new resource) The zone ID where the HSM instance is located. You can obtain this parameter through the DescribeZones API.
* `hsm_id` - (Optional, Computed) The ID of the HSM device. When this parameter is specified, an HSM instance for the specified device will be created. This attribute is also returned by the API.
* `ip` - (Optional, Computed) The classic network IP address corresponding to the HSM instance. If this parameter is not provided, the system will automatically assign an available IP based on the VPC and VSwitch. This attribute is also returned by the API.
* `remark` - (Optional) The alias of the HSM instance.
* `vpc_id` - (Optional, Computed) The VPC ID configured for the HSM instance. You can obtain this parameter through the DescribeVpc API. This attribute is also returned by the API.
* `vswitch_id` - (Optional, Computed) The VSwitch ID configured for the HSM instance. You can obtain this parameter through the DescribeVpc API. This attribute is also returned by the API.
* `white_list` - (Optional, Computed) The whitelist IP addresses that can access the HSM instance. Multiple values are supported and can be separated by commas. This attribute is also returned by the API.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the HSM instance.
* `cluster_id` - The ID of the cluster where the HSM instance is located.
* `cluster_name` - The name of the cluster where the HSM instance is located.
* `instance_id` - The ID of the HSM instance (same as id).
* `is_master` - Indicates whether the current HSM is the master HSM in the cluster. Valid values:
  * `0`: No.
  * `1`: Yes.
* `product_name` - The product name of the HSM instance.
* `release_protection` - The release protection status.
* `show_create_cluster` - Indicates whether a cluster can be created. Valid values:
  * `true`: Yes.
  * `false`: No.
* `status` - The status of the HSM instance. Valid values:
  * `1`: Not initialized.
  * `3`: Released.
  * `4`: Production failed.
  * `5`: Running.
  * `6`: Synchronizing.
  * `7`: Resetting.
  * `8`: Disabled.
* `vendor_name` - The vendor name of the HSM instance.

## Import

HSM Instance can be imported using the instance ID, e.g.

```
$ terraform import alibabacloudstack_hsm_instance.example hsm-12345678
```