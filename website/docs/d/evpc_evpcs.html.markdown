---
subcategory: "EasyAI"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_evpc_evpcs"
sidebar_current: "docs-Alibabacloudstack-datasource-evpc-evpcs"
description: |- 
  Provides a list of EVPCs owned by an AlibabacloudStack account.
---

# alibabacloudstack_evpc_evpcs

This data source provides a list of EVPCs (Elastic Virtual Private Clouds) in an AlibabacloudStack account according to the specified filters.

## Example Usage

```hcl
# Declare the data source
data "alibabacloudstack_evpc_evpcs" "example" {
  evpc_name = "my-evpc"
}

output "evpc_ids" {
  value = data.alibabacloudstack_evpc_evpcs.example.ids
}

output "evpc_names" {
  value = data.alibabacloudstack_evpc_evpcs.example.names
}
```

## Argument Reference

The following arguments are supported:

* `evpc_name` - (Optional) The name of the EVPC to filter results.
* `status` - (Optional) The status of the EVPC to filter results.
* `name_regex` - (Optional) A regex string to filter results by EVPC name.
* `ids` - (Optional) A list of EVPC IDs to filter results. If not specified, all EVPCs will be considered.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `names` - A list of names of matched EVPCs.
* `ids` - A list of IDs of matched EVPCs.
* `evpcs` - A list of matched EVPCs. Each element contains the following attributes:
  * `evpc_id` - The ID of the EVPC.
  * `evpc_name` - The name of the EVPC.
  * `status` - The status of the EVPC.
  * `description` - The description of the EVPC.
  * `cidr` - The CIDR block of the EVPC.
  * `tenant_id` - The tenant ID of the EVPC.
  * `department` - The department ID of the EVPC.
  * `department_name` - The department name of the EVPC.
  * `region_id` - The region ID of the EVPC.
  * `resource_group` - The resource group ID of the EVPC.
  * `resource_group_name` - The resource group name of the EVPC.
  * `cluster_id` - The cluster ID of the EVPC.
  * `ascm_create_user` - The user who created the EVPC.
  * `create_time` - The creation time of the EVPC.
  * `update_time` - The last update time of the EVPC.
