---
subcategory: "Express Connect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_routerinterface"
sidebar_current: "docs-Alibabacloudstack-expressconnect-routerinterface"
description: |- 
  Provides a expressconnect Routerinterface resource.
---

# alibabacloudstack_expressconnect_routerinterface
-> **NOTE:** Alias name has: `alibabacloudstack_router_interface`

Provides a expressconnect Routerinterface resource.

## Example Usage

```hcl
variable "name" {
  default = "tf-testAccRouterInterfaceConfig18787"
}

resource "alibabacloudstack_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/12"
}

variable "region" {
  default = "cn-hangzhou"
}

resource "alibabacloudstack_router_interface" "default" {
  opposite_region      = var.region
  router_type          = "VRouter"
  router_id            = alibabacloudstack_vpc.default.router_id
  role                 = "AcceptingSide"
  specification        = "Large.2"
  name                = var.name
  description         = "This is a test router interface"
  health_check_source_ip = "172.16.0.10"
  health_check_target_ip = "172.16.0.11"
  opposite_access_point_id = "ap-abc123456"
}
```

## Argument Reference

The following arguments are supported:

* `opposite_region` - (Required, ForceNew) The Region of the peer side.
* `router_type` - (Required, ForceNew) The type of the router. Valid values: `VRouter`, `VBR`.
* `router_id` - (Required, ForceNew) The ID of the router to which the router interface belongs.
* `role` - (Required, ForceNew) The role of the router interface. Valid values: `InitiatingSide`, `AcceptingSide`.
* `specification` - (Optional) The specification of the router interface. It is valid when `role` is set to `InitiatingSide`. Possible values include: `Small.1`, `Middle.1`, `Large.2`.
* `name` - (Optional) The name of the router interface. The length must be between 2 and 80 characters. Only Chinese characters, English letters, numbers, periods (`.`), underscores (`_`), or hyphens (`-`) are allowed. If not specified, it defaults to the router interface ID. The name cannot start with `http://` or `https://`.
* `description` - (Optional) The description of the router interface. The length must be between 2 and 256 characters or left blank. It cannot start with `http://` or `https://`.
* `health_check_source_ip` - (Optional) The source IP address for the health check packet. This is only valid when `router_type` is set to `VBR`. The IP address must be an unused IP within the local VPC subnet. It must be specified together with `health_check_target_ip`.
* `health_check_target_ip` - (Optional) The target IP address for the health check packet. This is only valid when `router_type` is set to `VBR`. The IP address must be an unused IP within the local VPC subnet. It must be specified together with `health_check_source_ip`.
* `opposite_access_point_id` - (Optional) The ID of the access point of the peer router interface.
* `opposite_router_type` - (Optional) The type of the peer router. Valid values: `VRouter`, `VBR`.
* `opposite_router_id` - (Optional) The ID of the peer router.
* `opposite_interface_id` - (Optional) The ID of the peer router interface.
* `opposite_interface_owner_id` - (Optional) The ID of the account that owns the peer router interface.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the router interface.
* `access_point_id` - The ID of the access point of the router interface.
* `opposite_router_type` - The type of the peer router.
* `opposite_router_id` - The ID of the peer router.
* `opposite_interface_id` - The ID of the peer router interface.
* `opposite_interface_owner_id` - The ID of the account that owns the peer router interface.
* `health_check_source_ip` - The source IP address for the health check packet.
* `health_check_target_ip` - The target IP address for the health check packet.