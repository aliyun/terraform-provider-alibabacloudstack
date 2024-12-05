---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc"
sidebar_current: "docs-alibabacloudstack-resource-vpc"
description: |-
  Provides a Alibabacloudstack VPC resource.
---

# alibabacloudstack\_vpc

Provides a VPC resource.

-> **NOTE:** Terraform will auto build a router and a route table while it uses `alibabacloudstack_vpc` to build a vpc resource.

## Example Usage

Basic Usage

```
resource "alibabacloudstack_vpc" "vpc" {
  vpc_name       = "tf_test_foo"
  cidr_block     = "172.16.0.0/12"
}
```


## Argument Reference

The following arguments are supported:

* `cidr_block` - (Required, ForceNew) The CIDR block for the VPC. The `cidr_block` is Optional and default value is `172.16.0.0/12`.
* `vpc_name` - (Optional) The name of the VPC. Defaults to null.
* `name` - (Optional) Field `name` has been deprecated from provider. New field `vpc_name` instead.
* `description` - (Optional) The VPC description. Defaults to null.
* `resource_group_id` - (Optional) The Id of resource group which the VPC belongs.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `secondary_cidr_blocks` - (Optional) The secondary CIDR blocks for the VPC.
* `dry_run` - (Optional, ForceNew) Specifies whether to precheck this request only. Valid values: `true` and `false`.
* `user_cidrs` - (Optional, ForceNew) The user cidrs of the VPC.
* `enable_ipv6` - (Optional, ForceNew) Specifies whether to enable the IPv6 CIDR block. Valid values: `false` (Default): disables IPv6 CIDR blocks. `true`: enables IPv6 CIDR blocks. If the `enable_ipv6` is `true`, the system will automatically create a free version of an IPv6 gateway for your private network and assign an IPv6 network segment assigned as /56.
* `router_table_id` - (Deprecated). Field 'router_table_id' has been deprecated. New field 'route_table_id' instead.
* `status` - The status of the VPC. Pending: The VPC is being configured. Available: The VPC is available.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the VPC.
* `cidr_block` - The CIDR block for the VPC.
* `name` - The name of the VPC.
* `description` - The description of the VPC.
* `router_id` - The ID of the router created by default on VPC creation.
* `route_table_id` - The route table ID of the router created by default on VPC creation.
* `ipv6_cidr_block` - The ipv6 cidr block of VPC.