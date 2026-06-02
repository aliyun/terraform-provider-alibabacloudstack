---
subcategory: "API Gateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_cascade_instance"
sidebar_current: "docs-Alibabacloudstack-resource-api-gateway-v2-cascade-instance"
description: |-
  Manage API Gateway V2 Cascade Gateway Instances
---

# alibabacloudstack_api_gateway_v2_cascade_instance

Manages API Gateway V2 Cascade Gateway Instances. This resource is used to create and manage cascade API gateway instances, supporting the specification of cascade instance ID and instance name. **Note: This resource cannot be modified after creation; any parameter changes will trigger resource recreation.**

## Example Usage

### Basic Usage

```hcl
variable "name" {
  default = "tf-testAccApiGwV222547"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name      = "${var.name}source"
  node_number        = 1
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

resource "alibabacloudstack_api_gateway_v2_cascade_instance" "default" {
  instance_name       = var.name
  cascade_instance_id = alibabacloudstack_api_gateway_v2_instance.default.id
}
```

## Argument Reference

The following arguments are supported, sorted by type (Required → Recreate on Change → Optional → Deprecated):

* `cascade_instance_id` - (Required, Recreate on Change) The cascade instance ID. Used to associate with the underlying cascade instance, must match an existing cascade instance ID in the API Gateway service. Cannot be modified after creation; changes will trigger resource recreation.
* `instance_name` - (Required, Recreate on Change) The instance name. Length must be 1-128 characters, can contain letters, digits, hyphens (-), and underscores (_). Cannot be modified after creation; changes will trigger resource recreation.
