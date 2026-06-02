---
subcategory: "API Gateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_cascade_link"
sidebar_current: "docs-Alibabacloudstack-datasource-api-gateway-v2-cascade-link"
description: |-
  Query cascade links of API Gateway v2
---

# alibabacloudstack_api_gateway_v2_cascade_link

> This data source for API Gateway v2 cascade links, used to query and manage API Gateway cascade link resources

## Example Usage

```hcl

variable "name" {
  default = "tf-testAccApiGwV23033395573387097040"
}

resource "alibabacloudstack_api_gateway_v2_instance" "source" {
  instance_name      = "${var.name}-source"
  node_number        = 1
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

resource "alibabacloudstack_api_gateway_v2_instance" "cascade" {
  instance_name      = "${var.name}-cascade"
  node_number        = 1
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

resource "alibabacloudstack_api_gateway_v2_cascade_instance" "default" {
  instance_name       = var.name
  cascade_instance_id = alibabacloudstack_api_gateway_v2_instance.cascade.id
}

resource "alibabacloudstack_api_gateway_v2_cascade_link" "default" {
  source_instance_id      = alibabacloudstack_api_gateway_v2_instance.source.id
  source_instance_address = "10.17.94.180"
  cascade_instance_id     = alibabacloudstack_api_gateway_v2_cascade_instance.default.id
  link_name               = var.name
}

data "alibabacloudstack_api_gateway_v2_cascade_links" "default" {
  name_regex = alibabacloudstack_api_gateway_v2_cascade_link.default.link_name
}

```

## Argument Reference

The following arguments are supported:

* `cascade_instance_name` (Optional): The name of the cascade instance, used to filter links of a specific cascade instance.

* `ids` (Optional): A list of cascade link IDs, used to precisely match specified cascade links by ID.

* `name_regex` (Optional): A regular expression for the cascade link name, used to filter cascade links by name pattern.

* `source_instance_name` (Optional): The name of the source instance, used to filter links of a specific source instance.

## Attributes Reference

The following attributes are exported:

* `id` (String): The unique identifier of the cascade link, equivalent to link_id.

* `cascade_instance_id` (String): The ID of the cascade instance.

* `cascade_instance_name` (String): The name of the cascade instance.

* `cascade_service_id` (String): The ID of the cascade service.

* `link_id` (String): The ID of the cascade link.

* `link_name` (String): The name of the cascade link.

* `source_instance_address` (String): The address of the source instance.

* `source_instance_id` (String): The ID of the source instance.

* `source_instance_name` (String): The name of the source instance.