---
subcategory: "VPC"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_router_interface_connection"
sidebar_current: "docs-apsarastack-resource-router-interface-connection"
description: |-
  Provides a VPC router interface connection resource to connect two VPCs.
---

# apsarastack\_router\_interface\_connection

Provides a VPC router interface connection resource to connect two router interfaces which are in two different VPCs.
After that, all of the two router interfaces will be active.

-> **NOTE:** At present, Router interface does not support changing opposite router interface, the connection delete action is only deactivating it to inactive, not modifying the connection to empty.

-> **NOTE:** If you want to changing opposite router interface, you can delete router interface and re-build them.

-> **NOTE:** A integrated router interface connection tunnel requires both InitiatingSide and AcceptingSide configuring opposite router interface.

-> **NOTE:** Please remember to add a `depends_on` clause in the router interface connection from the InitiatingSide to the AcceptingSide, because the connection from the AcceptingSide to the InitiatingSide must be done first.

## Example Usage

```
provider "apsarastack" {
  region = var.region
}

variable "region" {
  default = "region"
}

variable "name" {
  default = "apsarastackRouterInterfaceConnectionBasic"
}

resource "apsarastack_vpc" "foo" {
  name       = var.name
  cidr_block = "172.16.0.0/12"
}

resource "apsarastack_vpc" "bar" {
  provider   = apsarastack
  name       = var.name
  cidr_block = "192.168.0.0/16"
}

resource "apsarastack_router_interface" "initiate" {
  opposite_region      = var.region
  router_type          = "VRouter"
  router_id            = apsarastack_vpc.foo.router_id
  role                 = "InitiatingSide"
  specification        = "Large.2"
  name                 = var.name
  description          = var.name
  instance_charge_type = "PostPaid"
}

resource "apsarastack_router_interface" "opposite" {
  provider        = apsarastack
  opposite_region = var.region
  router_type     = "VRouter"
  router_id       = apsarastack_vpc.bar.router_id
  role            = "AcceptingSide"
  specification   = "Large.1"
  name            = "${var.name}-opposite"
  description     = "${var.name}-opposite"
}

// A integrated router interface connection tunnel requires both InitiatingSide and AcceptingSide configuring opposite router interface.
resource "apsarastack_router_interface_connection" "foo" {
  interface_id          = apsarastack_router_interface.initiate.id
  opposite_interface_id = apsarastack_router_interface.opposite.id
  depends_on            = [apsarastack_router_interface_connection.bar] // The connection must start from the accepting side.
}

resource "apsarastack_router_interface_connection" "bar" {
  provider              = apsarastack
  interface_id          = apsarastack_router_interface.opposite.id
  opposite_interface_id = apsarastack_router_interface.initiate.id
}
```
## Argument Reference

The following arguments are supported:

* `interface_id` - (Required, ForceNew) One side router interface ID.
* `opposite_interface_id` - (Required, ForceNew) Another side router interface ID. It must belong the specified "opposite_interface_owner_id" account.
* `opposite_interface_owner_id` - (Optional, ForceNew) Another side router interface account ID. Log on to the ApsaraStack Cloud console, select User Info > Account Management to check the account ID. Default to [Provider account_id](https://www.terraform.io/docs/providers/apsarastack/index.html#account_id).
* `opposite_router_id` - (Optional, ForceNew) Another side router ID. It must belong the specified "opposite_interface_owner_id" account. It is valid when field "opposite_interface_owner_id" is specified.
* `opposite_router_type` - (Optional, ForceNew) Another side router Type. Optional value: VRouter, VBR. It is valid when field "opposite_interface_owner_id" is specified.

-> **NOTE:** The value of "opposite_interface_owner_id" or "account_id" must be main account and not be sub account.

## Attributes Reference

The following attributes are exported:

* `id` - Router interface ID. The value is equal to "interface_id".

