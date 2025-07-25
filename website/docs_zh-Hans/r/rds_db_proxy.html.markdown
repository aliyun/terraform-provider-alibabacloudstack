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

## 示例用法
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

## 参数参考

支持以下参数：
  * `db_instance_id` - (必填) - 数据库实例ID。
  * `db_proxy_connect_string` - (选填) - 数据库连接信息字符串。
  * `db_proxy_connect_string_port` - (选填) - 数据库连接端口。
  * `db_proxy_instance_num` - (选填) - 单节点代理实力规格CPU数。
  * `effective_specific_time` - (选填) - 指定时间生效。格式：<i>yyyy-MM-dd</i>T<i>HH:mm:ss</i>Z（UTC时间）。>**EffectiveTime**为**SpecificTime**时，该参数必传。
  * `effective_time` - (选填) - 生效时间，取值：* **Immediate**：立即生效。* **MaintainTime**：在可运维时间段内生效，请参见ModifyDBInstanceMaintainTime。* **SpecificTime**：指定时间生效。默认值：**MaintainTime**。
  * `instance_network_type` - (选填) - 网络类型
  * `connection_persist` - (选填) - 连接池。取值：- **0**：关闭 - **1**：开启会话级连接池 - **2**：开启事务级连接池。
  * `causal_consist_read` - (选填) - 一致性参数。取值：- **0**：始终一致 - **1**：开启会话一致性 - **2**：开启全局一致性。

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `db_proxy_connect_string` - 数据库连接信息字符串。
  * `db_proxy_connect_string_port` - 数据库连接端口。
  * `db_proxy_instance_current_minor_version` - 当前数据库连接版本。
  * `db_proxy_instance_latest_minor_version` - 最后一个数据库连接版本。
  * `db_proxy_instance_status` - 数据库连接状态。
  * `db_proxy_endpoint_aliases` - 代理终端的备注信息。
  * `db_proxy_endpoint_name` - 代理终端ID。
  * `db_proxy_endpoint_type` - 代理终端类型，取值：* **RWSplit**：默认代理终端。* **Custom**：自定义代理终端。
  * `db_proxy_read_write_mode` - 代理终端模式，取值：* **ReadWrite**：读写模式。* **ReadOnly**：只读模式。
  * `db_proxy_instance_type` - 数据库代理实例类型，取值：- **common**：通用型代理- **exclusive**：独享型代理（默认值）
  * `db_proxy_service_status` - 数据库独享代理功能开关状态，取值：* **Startup**：开启。* **Shutdown**：关闭。
