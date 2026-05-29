---
subcategory: "Express Connect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_router_interface"
sidebar_current: "docs-alibabacloudstack-resource-router-interface"
description: |-
  Provides a Router Interface resource.
---

# alibabacloudstack_router_interface

Provides a Router Interface resource. Router interfaces are used to establish high-speed connections between Virtual Private Clouds (VPCs) and on-premises data centers.

## Example Usage

```hcl
variable "name" {
  default = "tf-testAccRouterInterface"
}

data "alibabacloudstack_account" "current" {
}

resource "alibabacloudstack_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_router_interface" "accepting_side" {
  opposite_region           = data.alibabacloudstack_account.current.region
  router_type               = "VRouter"
  router_id                 = alibabacloudstack_vpc.default.router_id
  role                      = "AcceptingSide"
  name                      = var.name
  description               = var.name
}

resource "alibabacloudstack_router_interface" "initiating_side" {
  opposite_region           = data.alibabacloudstack_account.current.region
  router_type               = "VRouter"
  router_id                 = alibabacloudstack_vpc.default.router_id
  role                      = "InitiatingSide"
  specification             = "Large.2"
  name                      = "${var.name}-initiating"
  description               = "${var.name}-initiating"
  opposite_interface_id     = alibabacloudstack_router_interface.accepting_side.id
  opposite_router_id        = alibabacloudstack_vpc.default.router_id
  opposite_router_type      = "VRouter"
  opposite_interface_owner_id = data.alibabacloudstack_account.current.id
}
```

## Argument Reference

The following arguments are supported:

### Required

* `opposite_region` - (Required, ForceNew) The region ID of the peer router interface. For cross-region connections, specify the peer region.
* `router_type` - (Required, ForceNew) The type of the router. Valid values: `VRouter` (VPC router), `VBR` (border router).
* `router_id` - (Required, ForceNew) The ID of the router. When router_type is VRouter, this is the VPC ID; when router_type is VBR, this is the VBR ID.
* `role` - (Required, ForceNew) The role of the router interface. Valid values: `InitiatingSide` (initiator), `AcceptingSide` (acceptor).

### Optional

* `specification` - (Optional) The specification of the router interface. Required when role is InitiatingSide, not required when role is AcceptingSide. Valid values vary by region and zone, common values include: `Large.1`, `Large.2`, `Medium.1`, `Medium.2`, `Small.1`, `Small.2`, etc.
* `name` - (Optional) The name of the router interface. The name must be 2 to 128 characters in length. It must start with a letter and cannot start with `http://` or `https://`.
* `description` - (Optional) The description of the router interface. The description must be 2 to 256 characters in length.
* `health_check_source_ip` - (Optional) The source IP address for health checks. Effective when router_type is VRouter. Must be specified together with `health_check_target_ip`.
* `health_check_target_ip` - (Optional) The target IP address for health checks. Effective when router_type is VRouter. Must be specified together with `health_check_source_ip`.
* `opposite_access_point_id` - (Optional) The access point ID on the peer side. Effective when router_type is VBR.
* `opposite_router_type` - (Optional) The router type on the peer side. Valid values: `VRouter`, `VBR`.
* `opposite_router_id` - (Optional) The router ID on the peer side.
* `opposite_interface_id` - (Optional) The router interface ID on the peer side.
* `opposite_interface_owner_id` - (Optional) The account ID of the peer router interface owner. If not specified, defaults to the current account.

-> **Note:** 
- When role is `AcceptingSide`, the `specification` parameter is automatically set to `Negative` and does not need to be manually specified.
- `health_check_source_ip` and `health_check_target_ip` must be specified together or both omitted.
- The value of `opposite_interface_owner_id` must be a primary account ID, not a sub-account.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the router interface.
* `access_point_id` - The access point ID. Effective when router_type is VBR.
* `opposite_access_point_id` - The access point ID on the peer side.
* `opposite_router_type` - The router type on the peer side.
* `opposite_router_id` - The router ID on the peer side.
* `opposite_interface_id` - The router interface ID on the peer side.
* `opposite_interface_owner_id` - The account ID of the peer router interface owner.
* `status` - The status of the router interface. Common statuses: `Idle`, `Waiting`, `Active`, `Inactive`.

## Import

Router Interface can be imported using the router_interface_id, e.g.

```
$ terraform import alibabacloudstack_router_interface.example ri-12345678
```

-> **Note:** This resource can also be referred to by the following alias:
> - `alibabacloudstack_expressconnect_routerinterface`
