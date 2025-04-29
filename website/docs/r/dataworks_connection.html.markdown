---
subcategory: "DataWorks"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_data_works_connection"
sidebar_current: "docs-Alibabacloudstack-data-works-connection"
description: |- 
  Provides a data_works Connection resource.
---

# alibabacloudstack_data_works_connection

Provides a dataworks Connection resource.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testaccdataworksconnection44523"
}

variable "password" {
}

resource "alibabacloudstack_data_works_connection" "default" {
  project_id     = "10060"
  connection_type = "rds"
  content = {
    username      = "cxt_new"
    database     = "cxt_test_new"
    tag          = "rds"
    password     = var.password
    instanceName = "rm-6cq93i8k9q0045i5t"
    rdsOwnerId   = "1371730998580255"
  }
  env_type       = "1"
  sub_type       = "mysql"
  name           = var.name
  description    = "Description for ${var.name}"
}
```

## Argument Reference

The following arguments are supported:

* `connection_id` - (ForceNew) The ID of the connection. This is automatically generated and cannot be modified after creation.
* `project_id` - (Required) The ID of the project where the connection will be created.
* `connection_type` - (Required) Type of connection. Currently, it supports `rds`.
* `content` - (Required) Details of the data source. It is a map containing the following keys:
  * `username` - (Required) The username used to connect to the data source.
  * `database` - (Required) The name of the database to connect to.
  * `tag` - (Required) The tag associated with the connection. For `rds`, this is typically set to `rds`.
  * `password` - (Required) The password used to connect to the data source.
  * `instanceName` - (Required) The name of the RDS instance.
  * `rdsOwnerId` - (Required) The owner ID of the RDS instance.
* `env_type` - (Required) The environment to which the data source belongs. Valid values are:
  * `0` - Development environment.
  * `1` - Production environment.
* `sub_type` - (Optional) Sub-types of strings, for scenarios where some parent types include sub-types. For `rds`, valid sub-types are:
  * `mysql`
  * `sqlserver`
  * `postgresql`
* `name` - (Required, ForceNew) The name of the connection. This value must be unique within the project.
* `description` - (Optional) Description of the connection. This provides additional information about the connection.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `connection_id` - The unique identifier of the connection. This value is in the format `<connection_id>:<$.ProjectId>`.