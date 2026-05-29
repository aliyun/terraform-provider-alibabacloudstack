---
subcategory: "Relational Database Service(RDS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_db_proxies"
sidebar_current: "docs-Alibabacloudstack-datasource-rds-dbproxies"
description: |-
  Provides a list of rds dbproxies owned by an alibabacloudstack account.
---

# alibabacloudstack\_db\_proxies

This data source provides a list of rds dbproxies in an alibabacloudstack account according to the specified filters.

## Example Usage
```
variable "name" {
  default = "tf-testAccRdsDbProxiesDataSource-8446363"
}


data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}


resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.1.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}




variable "rds_instance_type" {
  type      = string
  default   = ""
  sensitive = true
}

data "alibabacloudstack_rds_instance_types" "default" {
  ids                  = var.rds_instance_type != "" ? [var.rds_instance_type] : null
  engine               = "MySQL"
  engine_version       = "5.7"
  sorted_by            = "CPU"
  series               = "dual_ha"
}

resource "alibabacloudstack_db_instance" "default" {
  engine               = data.alibabacloudstack_rds_instance_types.default.instance_types.0.engine
  engine_version       = data.alibabacloudstack_rds_instance_types.default.instance_types.0.engine_version
  instance_type        = data.alibabacloudstack_rds_instance_types.default.instance_types.0.id
  instance_storage     = data.alibabacloudstack_rds_instance_types.default.instance_types.0.storage_min
  instance_charge_type = "Postpaid"
  instance_name        = "${var.name}"
  vswitch_id           = "${alibabacloudstack_vpc_vswitch.default.id}"
  monitoring_period    = "60"
  storage_type         = data.alibabacloudstack_rds_instance_types.default.instance_types.0.storage_type
}


resource "alibabacloudstack_db_proxy" "default" {
	db_instance_id = "${alibabacloudstack_db_instance.default.id}"
	db_proxy_instance_num = "1"
}

 

data "alibabacloudstack_db_proxies" "default" {
  db_instance_id = "${alibabacloudstack_db_proxy.default.db_instance_id}"
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
