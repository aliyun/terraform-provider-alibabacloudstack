---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_router_interface"
sidebar_current: "docs-alibabacloudstack-resource-router-interface"
description: |-
  Provides a VPC router interface resource to connect two VPCs.
---

# alibabacloudstack\_router\_interface

Provides a VPC router interface resource aim to build a connection between two VPCs.

-> **NOTE:** Only one pair of connected router interfaces can exist between two routers. Up to 5 router interfaces can be created for each router and each account.



## Example Usage

```
resource "alibabacloudstack_vpc" "foo" {
  name       = "tf_test_foo12345"
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_router_interface" "interface" {
  opposite_region = region
  router_type     = "VRouter"
  router_id       = alibabacloudstack_vpc.foo.router_id
  role            = "InitiatingSide"
  specification   = "Large.2"
  name            = "test1"
  description     = "test1"
}
```
## Argument Reference

The following arguments are supported:

* `opposite_region` - (Required, ForceNew) The Region of peer side.
* `router_type` - (Required, ForceNew) Router Type. Optional value: VRouter, VBR. Accepting side router interface type only be VRouter.
* `router_id` - (Required, ForceNew) The Router ID.
* `role` - (Required, ForceNew) The role the router interface plays. Optional value: `InitiatingSide`, `AcceptingSide`.
* `specification` - (Optional) Specification of router interfaces. It is valid when `role` is `InitiatingSide`. Accepting side's role is default to set as 'Negative'.
* `opposite_access_point_id` - (Optional) The ID of the access point of the peer router interface.
* `name` - (Optional) Name of the router interface. Length must be 2-80 characters long. Only Chinese characters, English letters, numbers, period (.), underline (_), or dash (-) are permitted.
                                                    If it is not specified, the default value is interface ID. The name cannot start with http:// and https://.
* `description` - (Optional) Description of the router interface. It can be 2-256 characters long or left blank. It cannot start with http:// and https://.
* `health_check_source_ip` - (Optional) Used as the Packet Source IP of health check for disaster recovery or ECMP. It is only valid when `router_type` is `VBR`. The IP must be an unused IP in the local VPC. It and `health_check_target_ip` must be specified at the same time.
* `health_check_target_ip` - (Optional) Used as the Packet Target IP of health check for disaster recovery or ECMP. It is only valid when `router_type` is `VBR`. The IP must be an unused IP in the local VPC. It and `health_check_source_ip` must be specified at the same time.

## Attributes Reference

The following attributes are exported:

* `id` - Router interface ID.
* `router_id` - Router ID.
* `router_type` - Router type.
* `role` - Router interface role.
* `name` - Router interface name.
* `description` - Router interface description.
* `specification` - Router nterface specification.
* `access_point_id` - Access point of the router interface.
* `opposite_access_point_id` - ID of the access point of the peer                             
* `opposite_router_type` - Peer router type.
* `opposite_router_id` - Peer router ID.
* `opposite_interface_id` - Peer router interface ID.
* `opposite_interface_owner_id` - Peer account ID.
* `health_check_source_ip` - Source IP of Packet of Line HealthCheck.
* `health_check_target_ip` - Target IP of Packet of Line HealthCheck.

