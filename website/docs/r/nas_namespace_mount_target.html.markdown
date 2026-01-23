---
subcategory: "NAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_namespace_mount_target"
sidebar_current: "docs-Alibabacloudstack-resource-nas-namespace-mount-target"
description: |-
  Create and manage NAS unified namespace mount targets
---

# alibabacloudstack_nas_namespace_mount_target

> NAS Unified Namespace Mount Settings

## Example Usage

### Basic Usage

```hcl

variable "name" {
  default = "tf-AccNasaccgop86369"
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

resource "alibabacloudstack_nas_accessgroup" "default1" {
  access_group_name = "${var.name}1"
  access_group_type = "Vpc"
}

resource "alibabacloudstack_nas_accessgroup" "default2" {
  access_group_name = "${var.name}2"
  access_group_type = "Vpc"
}



resource "alibabacloudstack_nas_namespace_mount_target" "default" {
  network_type      = "Vpc"
  access_group_name = alibabacloudstack_nas_accessgroup.default1.access_group_name
  nas_namespace_id  = alibabacloudstack_nas_namespace.default.id
  vswitch_id        = alibabacloudstack_vpc_vswitch.default.id
}
```

## Argument Reference

The following arguments are supported:

* `nas_namespace_id` - (Required, Forces new resource when changed) The namespace ID.
* `network_type` - (Required, Forces new resource when changed) The network type of the mount target. Valid values:
  * `Vpc`: VPC network
  * `Classic`: Classic network
* `access_group_name` - (Required) The name of the access group. Constraints:
  * Must be 3 to 64 characters in length
  * Must start with a letter and can contain letters, digits, underscores (_), or hyphens (-)
  * The name of a newly created access group cannot be the same as the two default access groups (DEFAULT_VPC_GROUP_NAME and DEFAULT_CLASSIC_GROUP_NAME)
* `vswitch_id` - (Optional, Forces new resource when changed) The ID of the vSwitch. This field is required and meaningful when the network type is VPC.
* `status` - (Optional) The status of the mount target. Valid values:
  * `Active`: Available
  * `Inactive`: Unavailable

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in the format {NasNamespaceId:MountTargetDomain}.
* `mount_target_domain` - The domain name of the mount target.
* `vpc_id` - The ID of the VPC. This field has a value when the network type is VPC.