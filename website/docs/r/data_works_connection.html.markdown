---
subcategory: "Data Works"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_data_works_connection"
sidebar_current: "docs-alibabacloudstack-resource-data-works-connection"
description: |- Provides a AlibabacloudStack Data Works Connection resource.
---

# alibabacloudstack\_data\_works\_connection

Provides a Data Works Connection resource.

For information about Data Works Connection and how to use it,
see [What is Connection](https://help.aliyun.com/apsara/enterprise/v_3_14_0_20210519/dide/enterprise-ascm-developer-guide/CreateConnection-1-2.html?spm=a2c4g.14484438.10001.560).

## Example Usage

Basic Usage

```terraform
resource "alibabacloudstack_data_works_connection" "default" {
  project_id = "10060"
  connection_type = "rds"
  content = {
              username = "cxt_new"
              database = "cxt_test_new"
              tag = "rds"
              password = "Admin123@ascm"
              instanceName = "rm-6cq93i8k9q0045i5t"
              rdsOwnerId = "1371730998580255"
            }
  env_type = "1"
  sub_type = "mysql"
  name = "tf-testacccn-hohhot-ste3-d01dataworksconnection97435"
  description = "descriptiontf-testacccn-hohhot-ste3-d01dataworksconnection97435"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.
* `connection_type` - (Required) Type of connection.
* `content` - (Required) Details of the data source.
* `env_type` - (Required) The environment to which the data source belongs, including 0 (development environment) and 1 (production environment).
* `sub_type` - (Optional) Sub-types of strings, for scenarios where some parent types include sub-types. There are currently the following combinations:
  * Parent type: rds 
  * Sub-type: mysql, sqlserver or postgresql.
* `name` - (Required) The name of the data source.
* `description` - (Optional) Description of the connection.

## Attributes Reference

The following attributes are exported:

* `connection_id` - The resource ID of Connection. The value formats as `<connection_id>:<$.ProjectId>`.

## Import

Data Works Connection can be imported using the id, e.g.

```
$ terraform import alibabacloudstack_data_works_connection.example <connection_id>:<$.ProjectId>
```
