---
subcategory: "云数据库 RDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_db_proxies"
sidebar_current: "docs-Alibabacloudstack-datasource-db-proxies"
description: |-
  提供阿里云账号下拥有的rds dbproxies列表�?
---

# alibabacloudstack\_db\_proxies

此数据源提供根据指定过滤条件列出的阿里云账号下的rds dbproxies资源列表�?

## 示例用法
```
variable "name" {
  default = "tf-testAccRdsDbProxiesDataSource-8982054"
}

resource "alibabacloudstack_db_instance" "instance" {
	engine               = "MySQL"
	engine_version       = "5.7"
	instance_type        = "rds.mysql.s2.large"
	instance_storage     = "5"
	instance_name 		   = "${var.name}"
	storage_type         = "local_ssd"
}

resource "alibabacloudstack_db_proxy" "default" {
	db_instance_id = "${alibabacloudstack_db_instance.instance.id}"
	db_proxy_instance_num = "1"
	instance_network_type = "Classic"
}

 

data "alibabacloudstack_db_proxies" "default" {
  db_instance_id = "${alibabacloudstack_db_instance.instance.id}"
}
```

## 参数参�?
以下参数是支持的�?
  * `db_instance_id` - (必填) - 数据库实例ID�?

## Attributes Reference
除了上述参数外，还导出以下属性：
  * `db_proxies` - 数据库代理列表，包含以下属�?
    * `db_instance_id` - 数据库实例ID�?
    * `db_proxy_connect_string_items` - 数据库连接信�?
    * `db_proxy_instance_current_minor_version` - 代理信息当前版本
    * `db_proxy_instance_latest_minor_version` - 代理信息最后版�?
    * `db_proxy_instance_name` - 代理实例名称
    * `db_proxy_instance_num` - 代理实例规格CPU数据
    * `db_proxy_instance_status` - 代理实例状�?
    * `db_proxy_endpoint_items` - 代理终端信息列表�?
    * `db_proxy_instance_type` - 数据库代理实例类型，取值：- **common**：通用型代�? **exclusive**：独享型代理（默认值）
    * `db_proxy_service_status` - 数据库独享代理功能开关状态，取值：* **Startup**：开启�? **Shutdown**：关闭�?
    * `instance_network_type` - 网络类型
    * `persistent_connection_status` - 是否开启连接保持。取值：- **Enabled**：开启连接保�? **Disabled**：关闭连接保�? - 仅RDS MySQL支持此参数�? - 修改连接保持状态时�?*ConfigDBProxyService**取值为**Modify**�?
