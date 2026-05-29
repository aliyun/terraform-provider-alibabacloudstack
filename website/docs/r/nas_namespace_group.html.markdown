---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_namespace_group"
sidebar_current: "docs-Alibabacloudstack-resource-nas-namespace-group"
description: |-
  Orchestrate NAS cross-domain mount orchestration
---

# alibabacloudstack_nas_namespace_group

Orchestrate NAS cross-domain mount orchestration resources in the specified resource set using credentials configured with the Provider.

## Example Usage

### Basic Usage

```hcl

variable "name" {
  default = "tf-AccNasaccgop35376"
}

data "alibabacloudstack_nas_zones" "default" {
}


data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details              = true
}


resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
  tags = {
    common_test = "terraform"
    filter      = var.name
  }
  enable_ipv6 = true
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
  enable_ipv6  = true
  zone_id      = data.alibabacloudstack_zones.default.zones.0.id
  lifecycle {
    ignore_changes = [
      tags
    ]
  }
}



resource "alibabacloudstack_nas_namespace" "default" {
  zone_id       = data.alibabacloudstack_nas_zones.default.zones.0.zone_id
  cluster_id    = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.cluster_id
  description   = var.name
  storage_type  = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.storage_type
  protocol_type = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.protocol_type
  encrypt_type  = "0"
}

resource "alibabacloudstack_nas_accessgroup" "default" {
  access_group_name = var.name
  access_group_type = "Vpc"
}

resource "alibabacloudstack_nas_namespace_mount_target" "default" {
  network_type      = "Vpc"
  access_group_name = alibabacloudstack_nas_accessgroup.default.access_group_name
  nas_namespace_id  = alibabacloudstack_nas_namespace.default.id
  vswitch_id        = alibabacloudstack_vpc_vswitch.default.id
}



resource "alibabacloudstack_nas_namespace_group" "default" {
  network_type        = "Vpc"
  mapped_path         = var.name
  nas_namespace_id    = alibabacloudstack_nas_namespace.default.id
  mount_target_domain = alibabacloudstack_nas_namespace_mount_target.default.mount_target_domain
}
```

## Argument Reference

The following arguments are supported:

* `mapped_path` - (Required, Forces new resource) The mapped path for cross-domain mount orchestration. Once mapped, the path will be permanently bound to this mount point.
* `mount_target_domain` - (Required, Forces new resource) The domain name of the namespace mount target.
* `nas_namespace_id` - (Required, Forces new resource) The ID of the namespace.
* `network_type` - (Required, Forces new resource) The network type. Valid values: `Vpc` (Virtual Private Cloud) or `Classic` (Classic Network). Once set, it cannot be changed unless the list is cleared.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID, which is the mount target domain name (`mount_target_domain`).
* `create_time` - The creation time of the resource, in ISO 8601 standard format.
* `member_id` - The member ID of the cross-domain mount orchestration, used for deletion operations.
* `status` - The status of the resource, indicating the current enabled state of the mount orchestration.