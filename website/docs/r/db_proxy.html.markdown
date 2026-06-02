---
subcategory: "ApsaraDB RDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_db_proxy"
sidebar_current: "docs-Alibabacloudstack-resource-db-proxy"
description: |-
  Provides a rds Dbproxy resource.
---

# alibabacloudstack\_db\_proxy

Provides a rds Dbproxy resource.

## Example Usage
```
variable "name" {
		default = "tf-testaccdbproxy867443"
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


locals {
  thrity_seconds_after = formatdate("YYYY-MM-DD'T'hh:mm:ssZ", timeadd(timestamp(), "30s"))
}

	

resource "alibabacloudstack_db_proxy" "default" {
  db_proxy_instance_num = "1"
  db_instance_id = "${alibabacloudstack_db_instance.default.id}"
}
```

## Argument Reference

The following arguments are supported:
  * `db_instance_id` - (Required) - The ID of the RDS instance.
  * `db_proxy_connect_string` - (Optional) - The connection string of the database proxy.
  * `db_proxy_connect_string_port` - (Optional) - The port of the database proxy.
  * `db_proxy_instance_num` - (Optional) - The number of the database proxy instance.
  * `db_proxy_instance_status` - (Optional) - The status of the database proxy.
  * `db_proxy_endpoint_aliases` - (Optional) - The proxy endpoint alias.
  * `db_proxy_endpoint_name` - (Optional) - The proxy endpoint name.
  * `db_proxy_endpoint_type` - (Optional) - The proxy endpoint type.
  * `db_proxy_read_write_mode` - (Optional) - The proxy endpoint read write mode.
  * `effective_specific_time` - (Optional) - The specified time takes effect. Format: <I> yyyy-MM-dd</I> T <I> HH:mm:ss</I> Z(UTC time).> When **EffectiveTime** is set to **SpecificTime**, this parameter must be set.
  * `effective_time` - (Optional) - Effective time, value:* **Immediate**: Effective immediately.* **MaintainTime**: takes effect during the O & M period. For details, see ModifyDBInstanceMaintainTime.* **SpecificTime**: The specified time takes effect.Default value: **MaintainTime * *.
  * `connection_persist` - (Optional) - (Optional) - Connection Pool.  Value: - **0**: Disable - **1**: Enable Session Connection Pool - **2**:Enable Transaction Connection Pool.
  * `causal_consist_read` - (Optional) - Consistency Parameters.  Value: - **0**: Enable Transaction Connection Pool - **1**: Session Consistency - **2**: Global Consistency.
## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `db_proxy_connect_string` -The connection string.
  * `db_proxy_connect_address` - The connection address.
  * `db_proxy_connect_string_port` - The connection port.
  * `db_proxy_instance_current_minor_version` - The db proxy instance current minor version.
  * `db_proxy_instance_latest_minor_version` - The db proxy instance latest minor version.
  * `db_proxy_instance_status` - The db proxy instance status.
  * `db_proxy_endpoint_aliases` - Note information of the proxy 
  * `db_proxy_endpoint_name` - The ID of the proxy terminal.
  * `db_proxy_endpoint_type` - Proxy terminal type, value:* **RWSplit**: the default proxy terminal.**Custom**: Custom proxy terminal.
  * `db_proxy_read_write_mode` - Proxy terminal mode, value:* **ReadWrite**: read/write mode.* **ReadOnly**: Read-only mode.
  * `db_proxy_instance_type` - Database proxy instance type, value:-**common**: General Purpose Agent-**exclusive**: exclusive proxy (default)
  * `db_proxy_service_status` - Database exclusive proxy function switch status, value:* **Startup**: Enable.* **Shutdown**: closed.
