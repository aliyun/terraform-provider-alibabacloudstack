---
subcategory: "Network Attached Storage (NAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_namespace_group"
sidebar_current: "docs-Alibabacloudstack-datasource-nas-namespace-group"
description: |-
  Query Alibaba Cloud NAS cross-region mount orchestration information
---

# alibabacloudstack_nas_namespace_group

> NAS Cross-Region Mount Orchestration

## Example Usage

```hcl

variable "name" {
  default = "tf-testnasng4165"
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
  nas_namespace_id    = alibabacloudstack_nas_namespace.default.id
  mount_target_domain = alibabacloudstack_nas_namespace_mount_target.default.mount_target_domain
  mapped_path         = var.name
}



data "alibabacloudstack_nas_namespace_groups" "default" {
  ids = [
    "${alibabacloudstack_nas_namespace_group.default.id}"
  ]
}
```

## Argument Reference
The following arguments can be used to filter query results:

* `mount_target_domain` (String, Optional): Used to filter results by mount target domain, which is the domain name for mounting the NAS file system.

* `mapped_path_regex` (String, Optional): Regular expression used to filter results by namespace group `mapped_path`.

* `nas_namespace_id` (String, Optional): Used to filter results by NAS namespace ID.

* `network_type` (String, Optional): Used to filter results by network type, which can be Vpc (Virtual Private Cloud) or Classic (Classic Network).

## Attributes Reference
The following attributes are exported:

* `id` (String): The unique identifier of the data source, generated based on the list of queried namespace IDs.

* `create_time` (String): The creation time of the namespace group, following the ISO 8601 standard format.

* `ids` (List): A list of NAS namespace group IDs.

* `mapped_path` (String): The mapped path, representing the path of the NAS namespace on the mount point.

* `member_id` (String): The member ID, representing the identifier of the member in the namespace group.

* `mount_target_domain` (String): The mount target domain, used for mounting the NAS file system, with trailing dots removed.

* `nas_namespace_id` (String): The NAS namespace ID.

* `names` (List): A list of NAS namespace group names (using NasNamespaceId as the name).

* `network_type` (String): The network type, which can be Vpc (Virtual Private Cloud) or Classic (Classic Network).

* `status` (String): The status of the namespace group, possible values include: Enabled.

* `groups` (List): A list of NAS namespace groups. Each element contains the attributes of a namespace group.

Each object in the `groups` list contains the following attributes:

* `create_time` (String): The creation time of the namespace group.

* `mapped_path` (String): The mapped path.

* `member_id` (String): The member ID.

* `mount_target_domain` (String): The mount target domain.

* `nas_namespace_id` (String): The NAS namespace ID.

* `network_type` (String): The network type.

* `status` (String): The status of the namespace group.