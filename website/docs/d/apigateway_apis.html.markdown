---
subcategory: "ApiGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_apigateway_apis"
sidebar_current: "docs-Alibabacloudstack-datasource-apigateway-apis"
description: |- 
  Provides a list of apigateway apis owned by an Alibabacloudstack account.
---

# alibabacloudstack_apigateway_apis
-> **NOTE:** Alias name has: `alibabacloudstack_api_gateway_apis`

This data source provides a list of API Gateway APIs in an Alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_apigateway_apis" "example" {
  group_id   = "your_group_id"
  api_id     = "your_api_id"
  name_regex = "your_api_name_pattern"
  output_file = "output_api_list.txt"
}

output "first_api_id" {
  value = data.alibabacloudstack_apigateway_apis.example.apis[0].id
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Optional) The ID of the API group.
* `api_id` - (Optional) The ID of the specific API. This field is deprecated from version 1.52.2 and replaced with the `ids` field.
* `ids` - (Optional, Available from version 1.52.2) A list of API IDs to filter the results.
* `name_regex` - (Optional) A regex string to filter API Gateway APIs by name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `names` - A list of API names.
* `apis` - A list of APIs. Each element contains the following attributes:
  * `id` - The unique identifier of the API.
  * `name` - The name of the API.
  * `region_id` - The region where the API is located.
  * `group_id` - The ID of the group that the API belongs to.
  * `group_name` - The name of the group that the API belongs to.
  * `description` - The description of the API.
