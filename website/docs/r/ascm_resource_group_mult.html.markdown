---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_resource_group_mult"
sidebar_current: "docs-alibabacloudstack-resource-resource_group_mult"
description: |-
Create a new resource set and also create all subsequent resources under this resource set
---

# alibabacloudstack\ascm_resource_group
This resource will help you to  Create different resources in different resource sets 

-> **NOTE:** it is necessary to execute the template twice, and create a new resource set for the first application. The second application executes other resources under the newly created resource set

-> **NOTE:** terraform apply -target=alibabacloudstack_ascm_resource_group.new_rg

-> **NOTE:** terraform apply

-> **NOTE:** terraform apply --destory

## Example Usage

```
terraform {
  required_providers {
    alibabacloudstack = {
      source  = "hashicorp/alibabacloudstack"
      version = "= 0.0.1"
    }
  }
}

provider "alibabacloudstack" {
   domain                  = "<ASAPI_ENDPOINT>"
   access_key              = "<AK>"
   secret_key              = "<SK>"
   region                  = "<REGION_NAME>"
   proxy                   = "HTTP://<IP>:<PORT>"
   protocol                = "HTTP"
   insecure                = "true"
   resource_group_set_name = "<OLD_RG>"
}

resource "alibabacloudstack_ascm_resource_group" "new_rg" {
  name = "<NEW_RG>"
  organization_id = <ORG_ID>
}

provider "alibabacloudstack" {
   alias = "new_rg"
   domain                  = "<ASAPI_ENDPOINT>"
   access_key              = "<AK>"
   secret_key              = "<SK>"
   region                  = "<REGION_NAME>"
   proxy                   = "HTTP://<IP>:<PORT>"
   protocol                = "HTTP"
   insecure                = "true"
   resource_group_set_name = resource.alibabacloudstack_ascm_resource_group.new_rg.name

}

resource "alibabacloudstack_vpc" "default" {
  vpc_name       = "wy_vpc_terraform_test"
  cidr_block = "192.168.0.0/16"
  provider = alibabacloudstack.new_rg
}
```

## Argument Reference

The following arguments are supported:
* `name` - (Required) The name of the resource group. This name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin or end with a hyphen, and must not begin with http:// or https://. Default value is null.
* `organization_id` - (Required) ID of an Organization.

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

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the vpc (until it reaches the initial `Available` status). 
* `delete` - (Defaults to 10 mins) Used when terminating the vpc. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the VPC.
* `cidr_block` - The CIDR block for the VPC.
* `name` - The name of the VPC.
* `description` - The description of the VPC.
* `router_id` - The ID of the router created by default on VPC creation.
* `route_table_id` - The route table ID of the router created by default on VPC creation.
* `ipv6_cidr_block` - The ipv6 cidr block of VPC.