---
subcategory: "Object Storage Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_oss_single_tunnel"
description: |-
  Creates an OSS single tunnel resource to establish a secure connection channel between OSS service and VPC network.
---

# alibabacloudstack_oss_single_tunnel

This resource is used to create an OSS single tunnel, enabling a secure connection between OSS service and a specified VPC network. Once created, the resource cannot be modified; any parameter changes require recreating the resource.

## Example Usage

### Basic Usage

```hcl
variable "name" {
  default = "tf-oss-single-tunnel-25672"
}

data "alibabacloudstack_oss_clusters" "default" {
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details              = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
  lifecycle {
    ignore_changes = [
      secondary_cidr_blocks,
      tags
    ]
  }
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  vswitch_name = "${var.name}_vsw"
  vpc_id       = alibabacloudstack_vpc_vpc.default.id
  cidr_block   = "172.16.1.0/24"
  zone_id      = data.alibabacloudstack_zones.default.zones.0.id
  lifecycle {
    ignore_changes = [
      tags
    ]
  }
}

resource "alibabacloudstack_oss_single_tunnel" "default" {
  shared     = "0"
  label      = var.name
  vpc_id     = alibabacloudstack_vpc_vpc.default.id
  vswitch_id = alibabacloudstack_vpc_vswitch.default.id
  cluster    = data.alibabacloudstack_oss_clusters.default.clusters.0.id
}
```

## Argument Reference

The following arguments are supported, sorted by type (Required → Recreate on Change → Optional → Deprecated):

* `cluster` - (Required, Recreate on Change) OSS cluster name. Must exactly match the target OSS cluster, in string format.
* `label` - (Required, Recreate on Change) Tunnel label identifier. Used for resource classification management, length 1-128 characters, cannot contain `http://` or `https://` prefixes.
* `shared` - (Required, Recreate on Change) Shared identifier. `0` indicates a private tunnel (only available for this account), `1` indicates a shared tunnel (can be accessed across accounts), value passed as a string.
* `vpc_id` - (Required, Recreate on Change) Target VPC network ID. Must be a valid VPC resource ID in the current region, in the format like `vpc-xxx`.
* `vswitch_id` - (Required, Recreate on Change) Target virtual switch ID. Must be a VSwitch resource ID associated with `vpc_id`, in the format like `vsw-xxx`.

## Attributes Reference

The following attributes are exported as resource state (output-only attributes):

* `id` - Unique identifier of the resource, in the format `{cluster}:{vpc_id}:{vip}`.
* `vip` - Assigned VPC intranet IP address. Automatically allocated by the system for OSS service endpoint access, in the format like `172.16.1.146`.