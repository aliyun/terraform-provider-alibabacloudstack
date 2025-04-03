---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_router_interface_connection"
sidebar_current: "docs-alibabacloudstack-resource-route-interface-connection"
description: |-
  Provides a Alibabacloudstack Router Interface Connection resource.
---
# alibabacloudstack_router_interface_connection
The Router Interface Connection resource provides a one-side router interface connection resource which connects two router interfaces.

## Argument Reference

The following arguments are supported:

* `interface_id` - (Required, ForceNew) One side router interface ID.
* `opposite_interface_id` - (Required, ForceNew) Another side router interface ID. It must belong the specified "opposite_interface_owner_id" account.
* `opposite_interface_owner_id` - (Optional, ForceNew) Another side router interface account ID. Log on to the AlibabacloudStack Cloud console, select User Info > Account Management to check the account ID. Default to [Provider account_id](https://www.terraform.io/docs/providers/alibabacloudstack/index.html#account_id).
* `opposite_router_id` - (Optional, ForceNew) Another side router ID. It must belong the specified "opposite_interface_owner_id" account. It is valid when field "opposite_interface_owner_id" is specified.
* `opposite_router_type` - (Optional, ForceNew) Another side router Type. Optional value: VRouter, VBR. It is valid when field "opposite_interface_owner_id" is specified. 

-> **NOTE:** The value of "opposite_interface_owner_id" or "account_id" must be main account and not be sub account.

## Attributes Reference

The following attributes are exported:

* `id` - Router interface ID. The value is equal to "interface_id".
* `opposite_router_id` - Another side router ID. It must belong the specified "opposite_interface_owner_id" account. It is valid when field "opposite_interface_owner_id" is specified. 
* `opposite_interface_owner_id` - Another side router interface account ID. 