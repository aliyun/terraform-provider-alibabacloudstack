---
subcategory: "Prometheus"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_prometheus_v2_instance"
sidebar_current: "docs-Alibabacloudstack-resource-prometheus-v2-instance"
description: |-
  Manages Prometheus V2 instance resources
---

# alibabacloudstack_prometheus_v2_instance

Manages Alibaba Cloud Prometheus V2 instance resources for creating, reading, updating, and deleting Prometheus monitoring instances.

## Example Usage

### Basic Usage

```hcl
variable "name" {
  default = "tfacc_prometheus15993"
}

resource "alibabacloudstack_prometheus_v2_instance" "default" {
  cluster_name = var.name
  tags = [
    "test1",
    "test2"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required, Forces new resource when changed) The name of the Prometheus instance. The name length must comply with API specifications. It cannot be directly modified after creation; modifying it will trigger resource recreation.
* `tags` - (Optional) A set of tags for the instance. Supports managing resources through tags; updates only affect tag configuration.

## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier (ClusterId) of the instance.
* `cluster_id` - The cluster ID of the Prometheus instance, same as the `id` attribute.
* `http_api` - The access address for Prometheus HTTP API, used for querying and managing monitoring data.
* `objid` - The internal object ID (integer type) of the instance, used for API operation identification.
* `push_gateway_url` - The access address for Push Gateway, used for receiving metric pushes.
* `remote_write_url` - The access address for Remote Write, used for remote writing of monitoring data.