---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_networkacl"
sidebar_current: "docs-Alibabacloudstack-vpc-networkacl"
description: |- 
  Provides a vpc Networkacl resource.
---

# alibabacloudstack_vpc_networkacl
-> **NOTE:** Alias name has: `alibabacloudstack_network_acl`

Provides a vpc Networkacl resource.

## Example Usage

Basic Usage

```terraform
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  name       = "VpcConfig"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  name              = "vswitch"
  cidr_block        = cidrsubnet(alibabacloudstack_vpc.default.cidr_block, 4, 4)
  availability_zone = data.alibabacloudstack_zones.default.ids.0
}

resource "alibabacloudstack_network_acl" "default" {
  vpc_id           = alibabacloudstack_vpc.default.id
  network_acl_name = "network_acl"
  description      = "network_acl"

  ingress_acl_entries {
    description            = "tf-testacc-ingress"
    network_acl_entry_name = "tcp23-ingress"
    source_cidr_ip         = "196.168.2.0/21"
    policy                 = "accept"
    port                   = "22/80"
    protocol               = "tcp"
  }

  egress_acl_entries {
    description            = "tf-testacc-egress"
    network_acl_entry_name = "tcp23-egress"
    destination_cidr_ip    = "0.0.0.0/0"
    policy                 = "accept"
    port                   = "-1/-1"
    protocol               = "all"
  }

  resources {
    resource_id   = alibabacloudstack_vswitch.default.id
    resource_type = "VSwitch"
  }
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) The ID of the associated VPC. This field cannot be changed after creation.
* `network_acl_name` - (Optional) The name of the network ACL. The name must be 1 to 128 characters in length and cannot start with `http://` or `https://`.
* `description` - (Optional) The description of the network ACL. The description must be 1 to 256 characters in length and cannot start with `http://` or `https://`.

### Ingress Acl Entries

* `ingress_acl_entries` - (Optional) List of the ingress entries of the network ACL. The order of the ingress entries determines the priority. Each entry supports the following:
  
  * `description` - (Optional) The description of the ingress entry.
  * `network_acl_entry_name` - (Optional) The name of the ingress entry.
  * `policy` - (Optional) The policy of the ingress entry. Valid values: `accept`, `drop`.
  * `port` - (Optional) The port range for the ingress entry.
  * `protocol` - (Optional) The protocol for the ingress entry. Valid values: `icmp`, `gre`, `tcp`, `udp`, `all`.
  * `source_cidr_ip` - (Optional) The source CIDR IP for the ingress entry.

### Egress Acl Entries

* `egress_acl_entries` - (Optional) List of the egress entries of the network ACL. The order of the egress entries determines the priority. Each entry supports the following:
  
  * `description` - (Optional) The description of the egress entry.
  * `network_acl_entry_name` - (Optional) The name of the egress entry.
  * `policy` - (Optional) The policy of the egress entry. Valid values: `accept`, `drop`.
  * `port` - (Optional) The port range for the egress entry.
  * `protocol` - (Optional) The protocol for the egress entry. Valid values: `icmp`, `gre`, `tcp`, `udp`, `all`.
  * `destination_cidr_ip` - (Optional) The destination CIDR IP for the egress entry.

### Resources

* `resources` - (Optional) The associated resources that this network ACL applies to. Each resource supports the following:
  
  * `resource_id` - (Optional) The ID of the associated resource.
  * `resource_type` - (Optional) The type of the associated resource. Valid values: `VSwitch`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the network ACL instance.
* `status` - The status of the network ACL.
* `egress_acl_entries` - Out direction rule information.
* `ingress_acl_entries` - Inward direction rule information.
* `network_acl_name` - The name of the network ACL.
* `name` - Deprecated field, use `network_acl_name` instead.
* `vpc_id` - (Computed) The ID of the associated VPC.