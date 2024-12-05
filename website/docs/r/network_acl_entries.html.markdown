---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_network_acl_entries"
sidebar_current: "docs-alibabacloudstack-resource-network-acl-entries"
description: |-
  Provides a Alibabacloudstack Network Acl Entries resource.
---

# alibabacloudstack\_network_acl_entries

Provides a network acl entries resource to create ingress and egress entries.

-> **NOTE:** Currently, the resource are only available in Hongkong(cn-hongkong), India(ap-south-1), and Indonesia(ap-southeast-1) regions.

-> **NOTE:** It doesn't support concurrency and the order of the ingress and egress entries determines the priority.

-> **NOTE:** Using this resource need to open a whitelist.

-> **DEPRECATED:**  This resource  has been deprecated . Replace by `ingress_acl_entries` and `egress_acl_entries` with the resource [alibabacloudstack_network_acl](https://www.terraform.io/docs/providers/alibabacloudstack/r/network_acl.html).

## Example Usage

Basic Usage

```
variable "name" {
  default = "NetworkAclEntries"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_network_acl" "default" {
  vpc_id = alibabacloudstack_vpc.default.id
  name   = var.name
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/21"
  zone_id           = data.alibabacloudstack_zones.default.zones[0].id
  name              = var.name
}

resource "alibabacloudstack_network_acl_attachment" "default" {
  network_acl_id = alibabacloudstack_network_acl.default.id
  resources {
    resource_id   = alibabacloudstack_vswitch.default.id
    resource_type = "VSwitch"
  }
}

resource "alibabacloudstack_network_acl_entries" "default" {
  network_acl_id = alibabacloudstack_network_acl.default.id
  ingress {
    protocol       = "all"
    port           = "-1/-1"
    source_cidr_ip = "0.0.0.0/32"
    name           = var.name
    entry_type     = "custom"
    policy         = "accept"
    description    = var.name
  }
  egress {
    protocol            = "all"
    port                = "-1/-1"
    destination_cidr_ip = "0.0.0.0/32"
    name                = var.name
    entry_type          = "custom"
    policy              = "accept"
    description         = var.name
  }
}
```

## Argument Reference

The following arguments are supported:

* `network_acl_id` - (Required, ForceNew) The id of the network acl, the field can't be changed.
* `ingress` - (Optional) List of the ingress entries of the network acl. The order of the ingress entries determines the priority. The details see Block Ingress. The resources mapping supports the following:
  * `description` - (Optional) The description of the ingress entry.
  * `source_cidr_ip` - (Optional) The source ip of the ingress entry.
  * `entry_type` - (Optional) The entry type of the ingress entry. It must be `custom` or `system`. Default value is `custom`.
  * `name` - (Optional) The name of the ingress entry.
  * `policy` - (Optional) The policy of the ingress entry. It must be `accept` or `drop`.
  * `port` - (Optional) The port of the ingress entry.
  * `protocol` - (Optional) The protocol of the ingress entry.
* `egress` - (Optional) List of the egress entries of the network acl. The order of the egress entries determines the priority. The details see Block Egress. The resources mapping supports the following:
  * `description` - (Optional) The description of the egress entry.
  * `destination_cidr_ip` - (Optional) The destination ip of the egress entry.
  * `entry_type` - (Optional) The entry type of the egress entry. It must be `custom` or `system`. Default value is `custom`.
  * `name` - (Optional) The name of the egress entry.
  * `policy` - (Optional) The policy of the egress entry. It must be `accept` or `drop`.
  * `port` - (Optional) The port of the egress entry.
  * `protocol` - (Optional) The protocol of the egress entry.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the network acl entries. It is formatted as `<network_acl_id>:<a unique id>`.


