---
subcategory: "Express Connect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_routerinterfaces"
sidebar_current: "docs-Alibabacloudstack-datasource-expressconnect-routerinterfaces"
description: |- 
  Provides a list of expressconnect routerinterfaces owned by an alibabacloudstack account.
---

# alibabacloudstack_expressconnect_routerinterfaces
-> **NOTE:** Alias name has: `alibabacloudstack_router_interfaces`

This data source provides a list of expressconnect router interfaces in an AlibabaCloudStack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_expressconnect_routerinterfaces" "example" {
  name_regex = "^testenv"
  status     = "Active"
}

output "first_router_interface_id" {
  value = "${data.alibabacloudstack_expressconnect_routerinterfaces.example.interfaces.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `status` - (Optional, ForceNew) The status of the resource. Valid values are `Active`, `Inactive`, and `Idle`.
* `name_regex` - (Optional, ForceNew) A regex string used to filter by router interface name.
* `specification` - (Optional, ForceNew) Specification of the link, such as `Small.1` (10Mb), `Middle.1` (100Mb), `Large.2` (2Gb), etc.
* `router_id` - (Optional, ForceNew) ID of the VRouter located in the local region.
* `router_type` - (Optional, ForceNew) Router type in the local region. Valid values are `VRouter` and `VBR` (physical connection).
* `role` - (Optional, ForceNew) Role of the router interface. Valid values are `InitiatingSide` (connection initiator) and `AcceptingSide` (connection receiver). The value of this parameter must be `InitiatingSide` if the `router_type` is set to `VBR`.
* `opposite_interface_id` - (Optional, ForceNew) ID of the peer router interface.
* `opposite_interface_owner_id` - (Optional, ForceNew) Account ID of the owner of the peer router interface.
* `ids` - (Optional) A list of router interface IDs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of router interface IDs.
* `names` - A list of router interface names.
* `interfaces` - A list of router interfaces. Each element contains the following attributes:
  * `id` - Router interface ID.
  * `status` - Router interface status. Possible values: `Active`, `Inactive`, and `Idle`.
  * `name` - Router interface name.
  * `description` - Router interface description.
  * `role` - Router interface role. Possible values: `InitiatingSide` and `AcceptingSide`.
  * `specification` - Router interface specification. Possible values: `Small.1`, `Middle.1`, `Large.2`, etc.
  * `router_id` - ID of the VRouter located in the local region.
  * `router_type` - Router type in the local region. Possible values: `VRouter` and `VBR`.
  * `vpc_id` - ID of the VPC that owns the router in the local region.
  * `access_point_id` - ID of the access point used by the VBR.
  * `creation_time` - Router interface creation time.
  * `opposite_region_id` - Peer router region ID.
  * `opposite_interface_id` - Peer router interface ID.
  * `opposite_router_id` - Peer router ID.
  * `opposite_router_type` - Router type in the peer region. Possible values: `VRouter` and `VBR`.
  * `opposite_interface_owner_id` - Account ID of the owner of the peer router interface.
  * `health_check_source_ip` - Source IP address used to perform health check on the physical connection.
  * `health_check_target_ip` - Destination IP address used to perform health check on the physical connection.
