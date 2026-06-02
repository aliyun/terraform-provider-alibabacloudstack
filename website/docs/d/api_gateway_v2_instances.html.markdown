---
subcategory: "API Gateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_instances"
sidebar_current: "docs-Alibabacloudstack-datasource-api-gateway-v2-instances"
description: |-
    Provides a datasource to query the API Gateway V2 instances.
---

The API Gateway V2 Instances data source provides a list of API Gateway V2 instances that match the request parameters.
## Example Usage

### Basic Usage

```hcl

data "alibabacloudstack_api_gateway_v2_instances" "example" {
  instance_id = "i-example123"
}

output "instance_ids" {
  value = data.alibabacloudstack_api_gateway_v2_instances.example.ids
}
```

### Filter by Deploy Mode

```hcl
data "alibabacloudstack_api_gateway_v2_instances" "edas_instances" {
  deploy_mode = "edas"
}

output "edas_instance_details" {
  value = data.alibabacloudstack_api_gateway_v2_instances.edas_instances.instances
}
```

### Filter by Broker Engine Type

```hcl
data "alibabacloudstack_api_gateway_v2_instances" "scg_instances" {
  broker_engine_type = "SCG"
}
```

### Using Regular Expression Filter

```hcl
data "alibabacloudstack_api_gateway_v2_instances" "name_filtered" {
  description_regex = "^test-.*"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of instance IDs to filter the results.
* `instance_id` - (Optional) The ID of the specific instance to retrieve.
* `broker_engine_type` - (Optional) The broker engine type of the instances. Valid values: `HIGRESS`, `SCG`.
* `deploy_mode` - (Optional) The deployment mode of the instances. Valid values: `k8s`, `edas`, `custom`.
* `name_regex` - (Optional, Deprecated) A regex string to filter instances by name. This field is deprecated and will be removed in a future release. Please use `description_regex` instead.
* `description_regex` - (Optional) A regex string to filter instances by description (name).

-> **NOTE:** `name_regex` and `description_regex` are mutually exclusive. Only one can be specified.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of instance IDs.
* `instances` - A list of instances. Each element contains the following attributes:
  * `id` - The ID of the instance.
  * `instance_id` - The ID of the instance.
  * `instance_name` - The name of the instance.
  * `broker_engine_type` - The broker engine type of the instance.
  * `deploy_mode` - The deployment mode of the instance.
  * `instance_class` - The class/specification of the instance.
  * `status` - The status of the instance.
  * `k8s_cluster_id` - The ID of the Kubernetes cluster where the instance is deployed.
  * `k8s_namespace` - The Kubernetes namespace where the instance is deployed.
  * `edas_app_id` - The EDAS application ID associated with the instance.
  * `edas_namespace_id` - The EDAS namespace ID associated with the instance.
  * `shared_instance` - Whether the instance is shared.
  * `node_number` - The number of nodes in the instance.
  * `edas_app_infos` - A set of EDAS application information objects, each containing:
    * `edas_namespace` - The EDAS namespace ID.
    * `app_id` - The application ID.
    * `k8s_cluster_id` - The Kubernetes cluster ID.
    * `k8s_namespace` - The Kubernetes namespace.
  * `custom_deploy_config` - A map of custom deployment configurations (only available when `deploy_mode` is `custom`).