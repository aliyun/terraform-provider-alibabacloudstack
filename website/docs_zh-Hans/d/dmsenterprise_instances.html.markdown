---
subcategory: "DMSEnterprise"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dmsenterprise_instances"
sidebar_current: "docs-Alibabacloudstack-datasource-dmsenterprise-instances"
description: |- 
  查询企业版数据库管理实例
---

# alibabacloudstack_dmsenterprise_instances
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_dms_enterprise_instances`

根据指定过滤条件列出当前凭证权限可以访问的企业版数据库管理实例列表。

## 示例用法

```terraform
# 创建RDS实例
resource "alibabacloudstack_db_instance" "instance" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.t1.small"
  instance_storage = "10"
  instance_name    = "tf_testAccDmsEnterpriseInstancesDataSource_5978035"
  security_ips     = ["0.0.0.0/0"]
  storage_type     = "local_ssd"
}

# 创建数据库账户
resource "alibabacloudstack_db_account" "account" {
  instance_id = "${alibabacloudstack_db_instance.instance.id}"
  name        = "admin123"
  password    = "inputYourCodeHere"
  type        = "Normal"
}

# 创建企业版数据库管理实例
resource "alibabacloudstack_dms_enterprise_instance" "default" {
  dba_uid           = "1234567890"
  host              = "${alibabacloudstack_db_instance.instance.connection_string}"
  port              = "3306"
  network_type      = "CLASSIC"
  safe_rule         = "自由操作"
  tid               = "1"
  instance_type     = "mysql"
  instance_source   = "RDS"
  env_type          = "test"
  database_user     = "${alibabacloudstack_db_account.account.name}"
  database_password = "${alibabacloudstack_db_account.account.password}"
  instance_alias    = "tf_testAccDmsEnterpriseInstancesDataSource_5978035"
  query_timeout     = "70"
  export_timeout    = "2000"
  ecs_region        = ""
  ddl_online        = "0"
  use_dsql          = "0"
  data_link_name    = ""
}

# 查询企业版数据库管理实例
data "alibabacloudstack_dms_enterprise_instances" "default" {
  search_key       = "${alibabacloudstack_dms_enterprise_instance.default.host}"
  env_type         = "test"
  instance_source  = "RDS"
  instance_type    = "mysql"
  net_type         = "CLASSIC"
  status           = "NORMAL"
  name_regex       = "^tf_testAccDmsEnterpriseInstancesDataSource_.*"
  output_file      = "dms_enterprise_instances.json"
}

output "first_database_instance_id" {
  value = "${data.alibabacloudstack_dms_enterprise_instances.default.instances.0.instance_id}"
}
```

## 参数参考

以下参数是支持的：

* `name_regex` - (选填, 变更时重建） - 用于通过企业版数据库管理实例实例的别名 (`instance_alias`) 过滤结果的正则表达式字符串。
* `instance_alias_regex` - (选填, 变更时重建） - 用于通过企业版数据库管理实例的别名 (`instance_alias`) 过滤结果的正则表达式字符串。
* `env_type` - (选填, 变更时重建） - 数据库实例所属环境的类型。例如，`prod`、`test` 或 `dev`。
* `instance_source` - (选填, 变更时重建） - 数据库实例的来源。例如，`RDS`、`ECS` 或 `OnPremise`。
* `instance_type` - (选填, 变更时重建） - 数据库实例的类型。例如，`mysql`、`sqlserver` 或 `postgresql`。
* `net_type` - (选填, 变更时重建） - 数据库实例的网络类型。有效值包括 `CLASSIC` 和 `VPC`。
* `search_key` - (选填, 变更时重建） - 用于查询数据库实例的关键字。
* `status` - (选填, 变更时重建） - 按照企业版数据库管理实例的状态筛选结果。有效值包括 `NORMAL`、`UNAVAILABLE`、`UNKNOWN`、`DELETED` 和 `DISABLE`。
* `tid` - (选填, 变更时重建） - 数据管理（DMS）企业版中租户的 ID。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` -企业版数据库管理ID 列表（每个都由 `host:port` 组成）。
* `names` -企业版数据库管理名称列表。
* `instances` -企业版数据库管理实例列表。每个元素包含以下属性：
  * `id` - DMS 企业实例的唯一标识符，格式为 `<host>:<port>`。
  * `data_link_name` - 数据库实例的数据链接名称。
  * `database_password` - 数据库实例的登录密码。
  * `database_user` - 数据库实例的登录用户名。
  * `dba_id` - 数据库实例的数据库管理员（DBA）ID。
  * `dba_nick_name` - DBA 的昵称。
  * `ddl_online` - 表示是否为数据库实例启用了在线数据描述语言（DDL）服务。
  * `ecs_instance_id` - 数据库实例所属的弹性计算服务（ECS）实例的 ID。
  * `ecs_region` - 数据库实例所在的区域。
  * `env_type` - 数据库实例所属环境的类型。
  * `export_timeout` - 导出数据库实例的超时时间。
  * `host` - 数据库实例的端点。
  * `instance_alias` - 数据库实例的别名。
  * `instance_name` - 键 `instance_alias` 的别名。
  * `instance_id` - 数据库实例的 ID。
  * `instance_source` - 数据库实例的来源。
  * `instance_type` - 数据库实例的类型。
  * `port` - 数据库实例的连接端口。
  * `query_timeout` - 查询数据库实例的超时时间。
  * `safe_rule_id` - 数据库实例的安全规则 ID。
  * `sid` - 数据库实例的系统 ID（SID）。
  * `status` - 数据库实例的状态。
  * `use_dsql` - 表示是否为数据库实例启用了跨数据库查询。
  * `vpc_id` - 数据库实例所属虚拟私有云（VPC）的 ID。