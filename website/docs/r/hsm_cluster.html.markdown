---
subcategory: "Hardware Security Module (HSM)"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_hsm_cluster"
description: |-
  Provides a AlibabacloudStack HSM Cluster resource.
---

# alibabacloudstack\_hsm\_cluster

Provides an HSM cluster resource.


## Example Usage

### Create an HSM cluster

```hcl

variable "name" {
  default = "tf_hsm_cluster6821719"
}

data "alibabacloudstack_zones" "default" {
  provider                    = alibabacloudstack-common
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  provider   = alibabacloudstack-common
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  provider     = alibabacloudstack-common
  vswitch_name = "${var.name}_vsw"
  vpc_id       = alibabacloudstack_vpc.default.id
  cidr_block   = "172.16.1.0/24"
  zone_id      = data.alibabacloudstack_zones.default.zones.0.id
}

data "alibabacloudstack_hsm_vendors" "default" {
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

resource "alibabacloudstack_hsm_instance" "default0" {
  product_code = data.alibabacloudstack_hsm_vendors.default.vendors.0.products.0.code
  vendor_code  = data.alibabacloudstack_hsm_vendors.default.vendors.0.code
  vsm_type     = "gvsm"
  zone_no      = data.alibabacloudstack_zones.default.zones.0.id
  remark       = var.name
}

resource "alibabacloudstack_hsm_instance" "default1" {
  product_code = data.alibabacloudstack_hsm_vendors.default.vendors.0.products.0.code
  vendor_code  = data.alibabacloudstack_hsm_vendors.default.vendors.0.code
  vsm_type     = "gvsm"
  zone_no      = data.alibabacloudstack_zones.default.zones.0.id
  remark       = var.name
}

resource "alibabacloudstack_hsm_cluster" "example" {
  cluster_name        = var.name
  master_instance_id  = alibabacloudstack_hsm_instance.default0.id
  vpc_id              = alibabacloudstack_vpc.default.id
  vswitch_ids         = alibabacloudstack_vswitch.default.id
  zone_nos            = data.alibabacloudstack_zones.default.zones.0.id
  password            = random_password.password.0.result
  ip_white_list       = "10.0.0.0/8"
  
  sub_instance_ids = [
    alibabacloudstack_hsm_instance.default1.id,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required) The name of the HSM cluster. The length is 2 to 128 English or Chinese characters. It must start with a digit or letter and support digits, letters, Chinese characters, underscores (_), and hyphens (-).
* `master_instance_id` - (Required, ForceNew) The ID of the master HSM instance.
* `vpc_id` - (Required, ForceNew) The ID of the VPC where the HSM cluster resides.
* `vswitch_ids` - (Required, ForceNew) The ID of the vSwitch that is associated with the VPC.
* `zone_nos` - (Required, ForceNew) The availability zone where the HSM cluster resides.
* `password` - (Required) The password of the HSM cluster administrator. The password must contain 8 to 30 characters and must include at least three of the following character types: uppercase letters, lowercase letters, digits, and special characters (!@#$%^&*()_+-=).
* `ip_white_list` - (Optional) The IP whitelist of the HSM cluster. Separate multiple IP addresses with commas (,). The whitelist can contain asterisks (*) as wildcards. Default value: 0.0.0.0/0.
* `sub_instance_ids` - (Optional) A list of sub HSM instance IDs that belong to the cluster.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the HSM cluster.

## Import

HSM clusters can be imported using the id, e.g.

```shell
$ terraform import alibabacloudstack_hsm_cluster.example cl-abc123456
```
