---
subcategory: "API Gateway"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_api_gateway_v2_cascade_link"
sidebar_current: "docs-AlibabacloudStack-api_gateway_v2_cascade_link"
description: |-
  Create and manage API Gateway v2 cascade link resources
---

# alibabacloudstack_api_gateway_v2_cascade_link

Manages API Gateway v2 cascade links, used to establish connections between source instances and cascade instances.

## Example Usage

### Basic Configuration

```hcl
variable "name" {
  default = "tf-testAccApiGwV237988"
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
```

## Argument Reference

The following arguments are supported:

* `cascade_instance_id` - (Required, ForceNew) The ID of the cascade instance. Specifies the unique identifier of the target API Gateway instance for establishing cascade relationships.
* `link_name` - (Required, ForceNew) The name of the link. A custom name for the link, used to identify the cascade link.
* `source_instance_address` - (Required) The address of the source instance. The IP address or domain name of the source service instance, used for API request routing.
* `source_instance_id` - (Required, ForceNew) The ID of the source instance. The unique identifier of the source service instance, which must match the source instance address.
* `cascade_instance_name` - (Optional) The name of the cascade instance. The display name of the cascade API Gateway instance, used for identification only.
* `cascade_service_id` - (Optional) The ID of the cascade service. The unique identifier of the associated cascade service, used for service cascade management.
* `source_instance_name` - (Optional) The name of the source instance. The display name of the source service instance, used for identification only.

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID of the cascade link, which is the same as the `link_id` value.
* `link_id` - The system-generated ID of the cascade link, used to uniquely identify the resource.