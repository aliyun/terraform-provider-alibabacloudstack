---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_proxy"
sidebar_current: "docs-Alibabacloudstack-polardb-proxy"
description: |-
  Provides a polardb Dbproxy resource.
---

# alibabacloudstack\_polardb\_proxy

Provides a polardb Dbproxy resource.

## 示例用法
```
variable "name" {
		default = "tf-testaccproxy867443"
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


	
variable "polardb_instance_type" {
  type      = string
  default   = ""
  sensitive = true
}

data "alibabacloudstack_polardb_instance_types" "default" {
  ids                  = var.polardb_instance_type != "" ? [var.polardb_instance_type] : null
  engine               = "MySQL"
  engine_version       = "5.7"
  sorted_by            = "CPU"
  series               = "dual_ha"
}

resource "alibabacloudstack_db_instance" "default" {
  engine               = data.alibabacloudstack_polardb_instance_types.default.instance_types.0.engine
  engine_version       = data.alibabacloudstack_polardb_instance_types.default.instance_types.0.engine_version
  instance_type        = data.alibabacloudstack_polardb_instance_types.default.instance_types.0.id
  instance_storage     = data.alibabacloudstack_polardb_instance_types.default.instance_types.0.storage_min
  instance_charge_type = "Postpaid"
  instance_name        = "${var.name}"
  vswitch_id           = "${alibabacloudstack_vpc_vswitch.default.id}"
  monitoring_period    = "60"
  storage_type         = data.alibabacloudstack_polardb_instance_types.default.instance_types.0.storage_type
}


locals {
  thrity_seconds_after = formatdate("YYYY-MM-DD'T'hh:mm:ssZ", timeadd(timestamp(), "30s"))
}

	

resource "alibabacloudstack_db_proxy" "default" {
  db_proxy_instance_num = "1"
  db_instance_id = "${alibabacloudstack_db_instance.default.id}"
}
```

## 参数参考

支持以下参数：
  * `db_instance_id` - (必填) - 数据库实例ID。
  * `db_proxy_connect_string` - (选填) - 数据库连接信息字符串。
  * `db_proxy_connect_string_port` - (选填) - 数据库连接端口。
  * `db_proxy_instance_num` - (选填) - 单节点代理实力规格CPU数。
  * `db_proxy_instance_status` - (选填) - 代理实例状态
  * `db_proxy_endpoint_aliases` - (选填) - 代理终端别名。
  * `db_proxy_endpoint_name` - (选填) - 代理终端名称。
  * `db_proxy_endpoint_type` - (选填) - 代理终端类型，取值：* **RWSplit**：读写分离代理终端。* **Custom**：自定义代理终端。
  * `db_proxy_read_write_mode` - (选填) - 代理读写分离模型。
  * `effective_specific_time` - (选填) - 指定时间生效。格式：<i>yyyy-MM-dd</i>T<i>HH:mm:ss</i>Z（UTC时间）。>**EffectiveTime**为**SpecificTime**时，该参数必传。
  * `effective_time` - (选填) - 生效时间，取值：* **Immediate**：立即生效。* **MaintainTime**：在可运维时间段内生效，请参见ModifyDBInstanceMaintainTime。* **SpecificTime**：指定时间生效。默认值：**MaintainTime**。
  * `connection_persist` - (选填) - 连接池。取值：- **0**：关闭 - **1**：开启会话级连接池 - **2**：开启事务级连接池。
  * `causal_consist_read` - (选填) - 一致性参数。取值：- **0**：始终一致 - **1**：开启会话一致性 - **2**：开启全局一致性。
  * `read_write_spliting` - (选填) - 读写分离。取值：- **0**：关闭 - **1**：开启。
  * `read_only_instance_max_delay_time` - (选填) - 只读实例流量分配的阈值（秒）。
  * `read_only_instance_weight` - (选填) - 实例的权重分配，不设置时则由系统分配。
    * `db_instance_id` - (选填) - 实例ID.
    * `weight` - (选填) - 实例权重(0-100).

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `db_proxy_connect_string` -  数据库连接信息字符串。
  * `db_proxy_connect_address` -  数据库连接端口。
  * `db_proxy_connect_string_port` - 当前数据库连接版本。
  * `db_proxy_instance_current_minor_version` - 当前数据库连接版本。
  * `db_proxy_instance_latest_minor_version` - 最后一个数据库连接版本。
  * `db_proxy_instance_status` - 数据库连接状态。
  * `db_proxy_endpoint_aliases` - 代理终端的备注信息。
  * `db_proxy_endpoint_name` - 代理终端ID。
  * `db_proxy_endpoint_type` - 代理终端类型，取值：* **RWSplit**：默认代理终端。* **Custom**：自定义代理终端。
  * `db_proxy_read_write_mode` - 代理终端模式，取值：* **ReadWrite**：读写模式。* **ReadOnly**：只读模式。
  * `db_proxy_instance_type` - 数据库代理实例类型，取值：- **common**：通用型代理- **exclusive**：独享型代理（默认值）
  * `db_proxy_service_status` - 数据库独享代理功能开关状态，取值：* **Startup**：开启。* **Shutdown**：关闭。
