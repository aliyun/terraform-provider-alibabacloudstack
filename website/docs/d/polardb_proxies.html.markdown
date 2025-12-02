---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_proxies"
sidebar_current: "docs-Alibabacloudstack-alibabacloudstack-polardb-proxies"
description: |-
  Provides a list of rds dbproxies owned by an alibabacloudstack account.
---

# alibabacloudstack_polardb_proxies

This data source provides a list of rds dbproxies in an alibabacloudstack account according to the specified filters.

## Example Usage
```
	variable "name" {
		default = "%v"
	}

	resource "alibabacloudstack_polardb_dbinstance" "default" {
	instance_storage = "5"
	instance_name = "${var.name}"
	storage_type = "local_ssd"
	engine = "MySQL"
	engine_version = "5.7"
	instance_type = "rds.mysql.t1.small"
	}

	resource "alibabacloudstack_polardb_proxy" "default" {
		db_instance_id = "${alibabacloudstack_polardb_dbinstance.default.id}"
		polardb_proxy_instance_num = "1"
	}

  data "alibabacloudstack_polardb_proxies" "default" {
    db_instance_id = "${alibabacloudstack_polardb_proxy.default.db_instance_id}"
  }
```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - The IDs of the rds dbproxies.
  * `db_instance_id` - (Required) - The ID of the RDS instance.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `db_proxies` - The list of the database proxy.
    * `db_proxy_connect_string_items` - The list of the database proxy connection string.
    * `db_proxy_instance_current_minor_version` - The current minor version of the database proxy instance.
    * `db_proxy_instance_latest_minor_version` - The latest minor version of the database proxy instance.
    * `db_proxy_instance_name` - The name of the database proxy instance.
    * `db_proxy_instance_num` - The number of the database proxy instance.
    * `db_proxy_instance_status` - The status of the database proxy instance.
    * `db_proxy_endpoint_items` - List of proxy terminal information.
    * `db_proxy_instance_type` - Database proxy instance type, value:-**common**: General Purpose Agent-**exclusive**: exclusive proxy (default)
    * `db_proxy_service_status` - Database exclusive proxy function switch status, value:* **Startup**: Enable.* **Shutdown**: closed.
    * `persistent_connection_status` - Whether to turn on the connection hold. Value:-**Enabled**: open connection hold-**Disabled**: Turn off connection hold>-only RDS MySQL supports this parameter.>-The value of **ConfigDBProxyService** is **Modify** when the connection retention status is modified * *.
