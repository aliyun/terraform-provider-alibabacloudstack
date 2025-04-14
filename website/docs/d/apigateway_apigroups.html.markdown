---
subcategory: "ApiGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_apigateway_apigroups"
sidebar_current: "docs-Alibabacloudstack-datasource-apigateway-apigroups"
description: |- 
  Provides a list of apigateway apigroups owned by an alibabacloudstack account.
---

# alibabacloudstack_apigateway_apigroups
-> **NOTE:** Alias name has: `alibabacloudstack_api_gateway_groups`

This data source provides a list of apigateway apigroups in an AlibabaCloudStack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_apigateway_apigroups" "example" {
  name_regex = "example-group"
  ids        = ["group1", "group2"]
  output_file = "apigroups_output.txt"
}

output "first_group_id" {
  value = data.alibabacloudstack_apigateway_apigroups.example.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string used to filter API Gateway groups by their names. 
* `ids` - (Optional) A list of API Gateway group IDs to filter results.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `names` - A list of API Gateway group names.
* `groups` - A list of API Gateway groups. Each element contains the following attributes:
  * `id` - The unique identifier of the API Gateway group.
  * `region_id` - The ID of the region where the API Gateway group is located.
  * `name` - The name of the API Gateway group.
  * `sub_domain` - The second-level domain name bound to the API Gateway group, which is used for API call testing.
  * `description` - The description of the API Gateway group, which cannot exceed 180 characters.
  * `created_time` - The creation time of the API Gateway group (in Greenwich Mean Time).
  * `modified_time` - The last modification time of the API Gateway group (in Greenwich Mean Time).
  * `traffic_limit` - The maximum QPS limit for the API Gateway group. The default value is 500, but it can be increased by submitting an application.
  * `billing_status` - The billing status of the API Gateway group. Possible values include:
    - `NORMAL`: The API Gateway group is normal.
    - `LOCKED`: Locked due to outstanding payment.
  * `illegal_status` - The illegal status of the API Gateway group. Possible values include:
    - `NORMAL`: The API Gateway group is normal.
    - `LOCKED`: Locked due to illegality. 