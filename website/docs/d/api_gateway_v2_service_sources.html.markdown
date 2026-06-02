---
subcategory: "API Gateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_service_sources"
sidebar_current: "docs-Alibabacloudstack-datasource-api-gateway-v2-service-sources"
description: |-
    Provides service sources for API Gateway V2 that can be accessed by an Alibaba Cloud account within the region configured in the provider.
---

# alibabacloudstack_api_gateway_v2_service_sources

This data source provides service sources for API Gateway V2 that can be accessed by an Alibaba Cloud account within the region configured in the provider.

## Example Usage

```hcl
variable "name" {
  default = "tf-testacc-serviceSource19103"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name      = var.name
  node_number        = "1"
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

resource "alibabacloudstack_api_gateway_v2_service_source" "default" {
  source_name      = var.name
  source_type      = "1"
  instance_id      = alibabacloudstack_api_gateway_v2_instance.default.id
  description      = var.name
  check_type       = "1"
  nacos_access_key = "root"
  nacos_secret_key = "12345"
  nacos_registry   = "127.0.0.1:8000"
}

data "alibabacloudstack_api_gateway_v2_service_sources" "default" {
  instance_id = alibabacloudstack_api_gateway_v2_service_source.default.instance_id

  name_regex = "tf-testacc-serviceSource*"
}
```

## Argument Reference

The following arguments support filtering the query results:

* `instance_id` - (String, Required) The ID of the API Gateway instance to which the service sources belong.
* `ids` - (List, Optional) A list of service source IDs to filter specific service sources. The ID format is {instance_id:source_id}.
* `name_regex` - (String, Optional) A regular expression to match service source names for filtering.

## Attributes Reference

The following attributes are exported:

* `id` - (String) The ID of the data source, which is a hash value generated from the list of service source IDs.
* `names` - (List) A list of names of the matched service sources.
* `sources` - (List) A list of matched service sources. Each service source contains the following attributes:
  * `create_time` - (String) The creation time of the service source.
  * `description` - (String) The description of the service source.
  * `id` - (String) The ID of the service source, in the format {instance_id:source_id}.
  * `instance_id` - (String) The ID of the API Gateway instance.
  * `source_id` - (String) The unique identifier ID of the service source.
  * `source_name` - (String) The name of the service source.
  * `source_type` - (String) The type code of the service source (1: NACOS, 2: Microservice Space, 3: Eureka, 5: Database).
  * `source_type_name` - (String) The type name of the service source (e.g., "NACOS").
  * `update_time` - (String) The last update time of the service source.