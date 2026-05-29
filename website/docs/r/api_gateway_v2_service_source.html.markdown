---
subcategory: "API Gateway V2"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_service_source"
sidebar_current: "docs-Alibabacloudstack-resource-api-gateway-v2-service-source"
description: |-
  Create and manage service sources for API Gateway v2
---

# alibabacloudstack_api_gateway_v2_service_source

Create and manage service sources in the specified API Gateway instance using credentials configured with the Provider.

Service sources for API Gateway v2 are used to define backend service addresses for APIs, supporting various types including NACOS, Microservice Spaces, Eureka, and Database.

## Example Usage

### Basic Usage

```hcl
variable "name" {
  default = "testtf-apigw-1739"
}

data "alibabacloudstack_api_gateway_v2_instance_types" "default" {
  sorted_by = "CPU"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name      = var.name
  node_number        = "1"
  instance_class     = data.alibabacloudstack_api_gateway_v2_instance_types.default.instance_types[0].id
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

// NACOS
resource "alibabacloudstack_api_gateway_v2_service_source" "default" {
  nacos_access_key = "root"
  nacos_secret_key = "12345"
  nacos_registry   = "127.0.0.1:8000"
  source_name      = var.name
  source_type      = "1"
  instance_id      = alibabacloudstack_api_gateway_v2_instance.default.id
  description      = "NACOS"
  check_type       = "1"
}
```


### Edas Usage

```hcl
resource "alibabacloudstack_api_gateway_v2_service_source" "example" {
  instance_id       = alibabacloudstack_api_gateway_v2_instance.default.id
  source_name       = var.name
  source_type       = "2"
  description       = "Microservice Space"
  
  edas_end_point_port    = 8080
  type                   = 1
  check_type             = 1
  edas_name_space_id     = "test1234"
  edas_access_key        = "root"
  edas_secret_key        = "1234"
  edas_end_point         = "127.0.0.1"
}
```

### Eureka Usage
```hcl
resource "alibabacloudstack_api_gateway_v2_service_source" "example" {
  instance_id       = alibabacloudstack_api_gateway_v2_instance.default.id
  source_name       =  var.name
  source_type       = "3"
  description       =  "Eureka"
  eureka_registry        = "https://127.0.0.1:8000"
}
```


### Database Usage

```hcl
resource "alibabacloudstack_api_gateway_v2_service_source" "example" {
  instance_id       = alibabacloudstack_api_gateway_v2_instance.default.id
  source_name       = "test5"
  source_type       = "5"
  description       = "database"
  max_connection         = 10
  max_idle_connection    = 5
  connection_idle_time   = 60
  database_type          = 0
  jdbc_url               = "jdbc:mysql://127.0.0.1:3306/testtf"
  username               = "root"
  password               = "123456"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the API Gateway instance.
* `source_name` - (Required) The name of the service source.
* `source_type` - (Required) The type of the service source. Valid values:
  * `1`: NACOS
  * `2`: Microservice Space
  * `3`: Eureka
  * `5`: Database
* `connection_idle_time` - (Optional) The connection idle time (in seconds).
* `database_type` - (Optional) The type of the database.
* `description` - (Optional) The description of the service source.
* `edas_access_key` - (Optional) The EDAS access key ID.
* `edas_end_point` - (Optional) The EDAS endpoint address.
* `edas_end_point_port` - (Optional) The EDAS endpoint port.
* `edas_name_space_id` - (Optional) The EDAS namespace ID.
* `edas_secret_key` - (Optional) The EDAS access key, sensitive information.
* `eureka_registry` - (Optional) The Eureka registry address.
* `jdbc_url` - (Optional) The JDBC connection URL.
* `max_connection` - (Optional) The maximum number of connections.
* `max_idle_connection` - (Optional) The maximum number of idle connections.
* `nacos_access_key` - (Optional) The Nacos access key ID.
* `nacos_registry` - (Optional) The Nacos registry address. 
* `nacos_secret_key` - (Optional) The Nacos access key, sensitive information.
* `password` - (Optional) The database password, sensitive information.
* `type` - (Optional) The type identifier.
* `username` - (Optional) The database username.
* `check_type` - (Optional) The check type.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. The value format is `<instance_id>:<source_id>`.
* `source_id` - The ID of the service source.
* `create_time` - The creation time of the service source.
* `update_time` - The last update time of the service source.
* `source_type_name` - The name of the service source type.

## Import

API Gateway V2 Service Source can be imported using the `instance_id` and `source_id` separated by a colon, e.g.

```
$ terraform import alibabacloudstack_api_gateway_v2_service_source.example gw-instance-123:source-456
```