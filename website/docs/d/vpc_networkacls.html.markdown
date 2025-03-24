---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_networkacls"
sidebar_current: "docs-Alibabacloudstack-datasource-vpc-networkacls"
description: |- 
  Provides a list of vpc networkacls owned by an alibabacloudstack account.
---

# alibabacloudstack_vpc_networkacls
-> **NOTE:** Alias name has: `alibabacloudstack_network_acls`

This data source provides a list of vpc network ACLs in an Alibabacloudstack account according to the specified filters.

## Example Usage

Basic Usage:

```terraform
data "alibabacloudstack_vpc_networkacls" "example" {
  ids        = ["example_value"]
  name_regex = "the_resource_name"
  network_acl_name = "example_network_acl"
  status     = "Available"
  vpc_id     = "vpc-1234567890abcdef"
  resource_type = "NETWORKACL"
}

output "first_network_acl_id" {
  value = data.alibabacloudstack_vpc_networkacls.example.acls.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of Network ACL IDs. If specified, the results will be filtered by these IDs.
* `name_regex` - (Optional, ForceNew) A regex string used to filter results by Network ACL name.
* `network_acl_name` - (Optional, ForceNew) The name of the network ACL. The name must be 1 to 128 characters in length and cannot start with `http://` or `https://`.
* `resource_id` - (Optional, ForceNew) The ID of the associated resource. This is required if `resource_type` is specified.
* `resource_type` - (Optional, ForceNew) The type of the associated resource. Valid values: `NETWORKACL`. Both `resource_type` and `resource_id` need to be specified at the same time to take effect.
* `status` - (Optional, ForceNew) The state of the network ACL. Valid values: `Available`, `Modifying`.
* `vpc_id` - (Optional, ForceNew) The ID of the associated VPC.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Network ACL names.
* `acls` - A list of Network ACLs. Each element contains the following attributes:
  * `description` - The description of the network ACL. The description must be 1 to 256 characters in length and cannot start with `http://` or `https://`.
  * `egress_acl_entries` - Outbound direction rule information. Each entry contains:
    * `description` - Description of the outbound direction rule.
    * `destination_cidr_ip` - The destination CIDR block.
    * `network_acl_entry_name` - The name of the outbound direction rule entry.
    * `policy` - The authorization policy (e.g., `Allow`, `Deny`).
    * `port` - Destination port range.
    * `protocol` - Transport layer protocol (e.g., `tcp`, `udp`).
  * `id` - The ID of the Network ACL.
  * `ingress_acl_entries` - Inbound direction rule information. Each entry contains:
    * `source_cidr_ip` - The source CIDR block.
    * `description` - Description of the inbound direction rule.
    * `network_acl_entry_name` - The name of the inbound direction rule entry.
    * `policy` - The authorization policy (e.g., `Allow`, `Deny`).
    * `port` - Source port range.
    * `protocol` - Transport layer protocol (e.g., `tcp`, `udp`).
  * `network_acl_id` - The first ID of the resource.
  * `network_acl_name` - The name of the network ACL.
  * `resources` - The associated resources. Each entry contains:
    * `resource_id` - The ID of the associated resource.
    * `resource_type` - The type of the associated resource.
    * `status` - The state of the associated resource.
  * `status` - The state of the network ACL.
  * `vpc_id` - The ID of the associated VPC.