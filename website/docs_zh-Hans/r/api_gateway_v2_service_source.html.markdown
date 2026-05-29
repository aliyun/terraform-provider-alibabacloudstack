---
subcategory: "API 网关（API Gateway）V2 版"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_service_source"
sidebar_current: "docs-Alibabacloudstack-api_gateway-api_gateway_v2_service_source"
description: |-
  创建和管理API网关v2版本的服务来源
---

# alibabacloudstack_api_gateway_v2_service_source

使用Provider配置的凭证在指定的API网关实例中创建和管理服务来源。

API网关v2版本的服务来源用于定义API的后端服务地址，支持NACOS、微服务空间、Eureka和Database等多种类型。

## 示例用法

### 基础用法

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

### 微服务空间 用法

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

### Eureka 用法
```hcl
resource "alibabacloudstack_api_gateway_v2_service_source" "example" {
  instance_id       = alibabacloudstack_api_gateway_v2_instance.default.id
  source_name       =  var.name
  source_type       = "3"
  description       =  "Eureka"
  eureka_registry        = "https://127.0.0.1:8000"
}
```


### Database 用法

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

## 参数说明

支持以下参数：

* `instance_id` - (必填, 变更时重建) API网关实例ID。
* `source_name` - (必填) 服务来源名称。
* `source_type` - (必填) 服务来源类型。有效值：
  * `1`：NACOS
  * `2`：微服务空间
  * `3`：Eureka
  * `5`：Database
* `connection_idle_time` - (可选) 连接空闲时间（秒）。
* `database_type` - (可选) 数据库类型。
* `description` - (可选) 服务来源描述。
* `edas_access_key` - (可选) EDAS访问密钥ID。
* `edas_end_point` - (可选) EDAS接入点地址。
* `edas_end_point_port` - (可选) EDAS接入点端口。
* `edas_name_space_id` - (可选) EDAS命名空间ID。
* `edas_secret_key` - (可选) EDAS访问密钥，敏感信息。
* `eureka_registry` - (可选) Eureka注册中心地址。
* `jdbc_url` - (可选) JDBC连接URL。
* `max_connection` - (可选) 最大连接数。
* `max_idle_connection` - (可选) 最大空闲连接数。
* `nacos_access_key` - (可选) Nacos访问密钥ID。
* `nacos_registry` - (可选) Nacos注册中心地址。
* `nacos_secret_key` - (可选) Nacos访问密钥，敏感信息。
* `password` - (可选) 数据库密码，敏感信息。
* `type` - (可选) 类型标识。
* `username` - (可选) 数据库用户名。
* `check_type` - (可选) 检查类型。

## 属性说明

导出以下属性：

* `id` - 服务来源ID，格式为`{instance_id}:{source_id}`。
* `create_time` - 服务来源创建时间。
* `source_id` - 服务来源ID。
* `source_type_name` - 服务来源类型名称。
* `update_time` - 服务来源最后更新时间。

## Import

API网关v2服务来源可以使用 `instance_id` 和 `source_id` 导入，两者用冒号分隔，例如：

```
$ terraform import alibabacloudstack_api_gateway_v2_service_source.example gw-instance-123:source-456
```