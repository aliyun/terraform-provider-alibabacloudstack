---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_vpc"
sidebar_current: "docs-Alibabacloudstack-vpc-vpc"
description: |- 
  Provides a vpc Vpc resource.
---

# alibabacloudstack_vpc_vpc
-> **NOTE:** Alias name has: `alibabacloudstack_vpc`

Provides a vpc Vpc resource.

## Example Usage

Basic Usage

```hcl
variable "name" {
    default = "tf-testaccvpcvpc48306"
}

resource "alibabacloudstack_vpc_vpc" "default" {
  cidr_block      = "172.16.0.0/12"
  vpc_name        = var.name
  description     = "RDK更新"
  enable_ipv6     = true
  resource_group_id = "rg-abc123xyz"
  secondary_cidr_blocks = ["192.168.0.0/16"]
  user_cidrs      = ["10.0.0.0/8"]
}
```

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Required, ForceNew) The CIDR block of the VPC. You can specify one of the following CIDR blocks or their subsets as the primary IPv4 CIDR block of the VPC:
  * Standard private CIDR blocks defined by RFC documents: 192.168.0.0/16, 172.16.0.0/12, and 10.0.0.0/8. The subnet mask must be between 8 and 28 bits in length.
  * Custom CIDR blocks other than the following ranges: 100.64.0.0/10, 224.0.0.0/4, 127.0.0.0/8, 169.254.0.0/16, and their subnets.
* `vpc_name` - (Optional) The name of the VPC. The name must be 1 to 128 characters in length and cannot start with `http://` or `https://`. Defaults to null.
* `description` - (Optional) The description of the VPC. The description must be 1 to 256 characters in length and cannot start with `http://` or `https://`. Defaults to null.
* `dry_run` - (Optional, ForceNew) Specifies whether to perform a dry run. Valid values:
  * `true`: Performs a dry run. The system checks the required parameters, request syntax, and limits. If the request fails the dry run, an error message is returned. If the request passes the dry run, the `DryRunOperation` error code is returned.
  * `false` (default): Sends the request and performs the operation. If the request passes the dry run, a 2xx HTTP status code is returned and the operation is performed.
* `enable_ipv6` - (Optional, ForceNew) Specifies whether to enable the IPv6 CIDR block. Valid values:
  * `false` (Default): Disables IPv6 CIDR blocks.
  * `true`: Enables IPv6 CIDR blocks. When this parameter is set to `true`, the system will automatically create a free version of an IPv6 gateway for your private network and assign an IPv6 network segment assigned as /56.
* `resource_group_id` - (Optional) The ID of the resource group to which you want to move the resource. You can use resource groups to facilitate resource grouping and permission management for Alibaba Cloud resources. For more information, see [What is resource management?](https://www.alibabacloud.com/help/en/doc-detail/94475.html)
* `secondary_cidr_blocks` - (Optional) A list of secondary CIDR blocks for the VPC. **Note**: This field has been deprecated from provider version 1.185.0 and will be removed in future versions. Use the new resource `alicloud_vpc_ipv4_cidr_block` instead. `secondary_cidr_blocks` attributes and `alicloud_vpc_ipv4_cidr_block` resource cannot be used at the same time.
* `user_cidrs` - (Optional, ForceNew) A list of user-defined CIDRs.
* `status` - The status of the VPC. Valid values:
  * `Pending`: The VPC is being configured.
  * `Available`: The VPC is available.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the VPC.
* `router_id` - The ID of the router created by default on VPC creation.
* `route_table_id` - The route table ID of the router created by default on VPC creation.
* `ipv6_cidr_block` - The IPv6 CIDR block of the VPC. This attribute is only available when `enable_ipv6` is set to `true`.
* `resource_group_id` - The ID of the resource group to which the VPC belongs.
* `status` - The status of the VPC. Valid values:
  * `Pending`: The VPC is being configured.
  * `Available`: The VPC is available.
* `vpc_name` - The name of the VPC.