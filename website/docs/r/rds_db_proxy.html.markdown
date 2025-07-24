---
subcategory: "RDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_rds_dbproxy"
sidebar_current: "docs-Alibabacloudstack-rds-dbproxy"
description: |-
  Provides a rds Dbproxy resource.
---

# alibabacloudstack\_rds\_dbproxy

Provides a rds Dbproxy resource.

## Example Usage
```
variable "name" {
		default = "tf-testAccdbproxy-714466"
	}

resource "alibabacloudstack_db_instance" "instance" {
  engine               = "MySQL"
    engine_version       = "5.7"
    instance_type        = "rds.mysql.s2.large"
    instance_storage     = "5"
    instance_name 		 = "${var.name}"
    storage_type         = "local_ssd"
}
	
resource "alibabacloudstack_db_proxy" "default" {
  db_proxy_instance_num = "1"
  instance_network_type = "Classic"
  db_instance_id = "${alibabacloudstack_db_instance.instance.id}"
}
```

## Argument Reference

The following arguments are supported:
  * `db_instance_id` - (Required) - The ID of the RDS instance.
  * `db_proxy_connect_string` - (Optional) - The connection string of the database proxy.
  * `db_proxy_connect_string_port` - (Optional) - The port of the database proxy.
  * `db_proxy_instance_num` - (Optional) - The number of the database proxy instance.
  * `db_proxy_instance_status` - (Optional) - The status of the database proxy.
  * `effective_specific_time` - (Optional) - The specified time takes effect. Format: <I> yyyy-MM-dd</I> T <I> HH:mm:ss</I> Z(UTC time).> When **EffectiveTime** is set to **SpecificTime**, this parameter must be set.
  * `effective_time` - (Optional) - Effective time, value:* **Immediate**: Effective immediately.* **MaintainTime**: takes effect during the O & M period. For details, see ModifyDBInstanceMaintainTime.* **SpecificTime**: The specified time takes effect.Default value: **MaintainTime * *.
  * `instance_network_type` - (Optional) - Network type
  * `connection_persist` - (Optional) - Connection Pool.  Value: - **0**: Disable - **1**: Enable Session Connection Pool - **2**：Enable Transaction Connection Pool.
  * `causal_consist_read` - (Optional) - Consistency Parameters.  Value: - **0**: Enable Transaction Connection Pool - **1**: Session Consistency - **2**: Global Consistency.
## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `db_proxy_connect_string` - The connection string.
  * `db_proxy_connect_string_port` - The connection port.
  * `db_proxy_connect_string_items` - The connection information.
    * `db_proxy_connect_string` - The connection string.
    * `db_proxy_connect_string_port` - The connection port.
    * `db_proxy_connect_string_net_type` -The connection type.
    * `db_proxy_connect_string_net_work_type` - The connection network type.
    * `db_proxy_endpoint_id` - The connection ID.
  * `db_proxy_instance_current_minor_version` - The db proxy instance current minor version.
  * `db_proxy_instance_latest_minor_version` - The db proxy instance latest minor version.
  * `db_proxy_instance_status` - The db proxy instance status.
  * `db_proxy_endpoint_items` - List of proxy terminal information.
    * `db_proxy_endpoint_aliases` - Note information of the proxy terminal.
    * `db_proxy_endpoint_name` - The ID of the proxy terminal.
    * `db_proxy_endpoint_type` - Proxy terminal type, value:* **RWSplit**: the default proxy terminal.**Custom**: Custom proxy terminal.
    * `db_proxy_read_write_mode` - Proxy terminal mode, value:* **ReadWrite**: read/write mode.* **ReadOnly**: Read-only mode.
  * `db_proxy_instance_type` - Database proxy instance type, value:-**common**: General Purpose Agent-**exclusive**: exclusive proxy (default)
  * `db_proxy_service_status` - Database exclusive proxy function switch status, value:* **Startup**: Enable.* **Shutdown**: closed.
