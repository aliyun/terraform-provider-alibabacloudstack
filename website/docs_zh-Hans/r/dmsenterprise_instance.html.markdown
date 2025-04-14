---
subcategory: "DMSEnterprise"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dmsenterprise_instance"
sidebar_current: "docs-Alibabacloudstack-dmsenterprise-instance"
description: |- 
  编排企业版数据库管理实例
---

# alibabacloudstack_dms_enterprise_instance
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_dmsenterprise_instance`

使用Provider配置的凭证在指定的资源集下编排企业版数据库管理实例。

## 示例用法

```hcl
variable "name" {
    default = "tf-testAccDmsEnterpriseInstance4641"
}

variable "password" {
    default = "inputYourCodeHere"
}

data "alibabacloudstack_account" "current" {}

resource "alibabacloudstack_db_instance" "instance" {
    engine           = "MySQL"
    engine_version   = "5.6"
    instance_type    = "rds.mysql.t1.small"
    instance_storage = "10"
    instance_name    = var.name
    security_ips     = ["0.0.0.0/0"]
    storage_type     = "local_ssd"
}

resource "alibabacloudstack_db_account" "account" {
    instance_id = alibabacloudstack_db_instance.instance.id
    name        = "tftest123"
    password    = var.password
    type        = "Normal"
}

resource "alibabacloudstack_ascm_user" "user" {
    cellphone_number = "13900000000"
    email           = "test@gmail.com"
    display_name    = "C2C-DELTA"
    organization_id = 33
    mobile_nation_code = "91"
    login_name      = "User_Dms_${var.name}"
    login_policy_id = 1
}

resource "alibabacloudstack_dms_enterprise_user" "default" {
    uid         = alibabacloudstack_ascm_user.user.user_id
    user_name   = alibabacloudstack_ascm_user.user.login_name
    mobile      = "15910799999"
    role_names  = ["ADMIN"]
}

resource "alibabacloudstack_dms_enterprise_instance" "default" {
    tid                = "1"
    instance_type      = "MySQL"
    instance_source   = "RDS"
    network_type      = "CLASSIC"
    env_type          = "test"
    host              = alibabacloudstack_db_instance.instance.connection_string
    port              = 3306
    database_user     = alibabacloudstack_db_account.account.name
    database_password = alibabacloudstack_db_account.account.password
    instance_alias    = "tf-testAccDmsEnterpriseInstance4641"
    dba_uid           = alibabacloudstack_dms_enterprise_user.default.uid
    safe_rule         = "自由操作"
    query_timeout     = 70
    export_timeout    = 2000
    ecs_region        = "cn-shanghai"
    use_dsql          = 0
    ddl_online        = 0
}
```

## 参数说明

支持以下参数：

* `tid` - (可选) 租户ID。
* `instance_type` - (必填) 数据库类型。有效值：`MySQL`, `SQLServer`, `PostgreSQL`, `Oracle`, `DRDS`, `OceanBase`, `Mongo`, `Redis`。
* `instance_source` - (必填) 数据库实例来源。有效值：`PUBLIC_OWN`, `RDS`, `ECS_OWN`, `VPC_IDC`。
* `network_type` - (必填, 变更时重建) 网络类型。有效值：`CLASSIC`, `VPC`。
* `env_type` - (必填) 环境类型。有效值：`product`（生产环境），`dev`（开发环境），`pre`（预发布环境），`test`（测试环境），`sit`（SIT环境），`uat`（UAT环境），`pet`（压力测试环境），`stag`（STAG环境）。
* `host` - (必填, 变更时重建) 目标数据库主机地址。
* `port` - (必填, 变更时重建) 目标数据库访问端口。
* `database_user` - (必填) 数据库访问账号。
* `database_password` - (必填) 数据库访问密码。
* `instance_name` - (必填) 实例名称，帮助用户快速区分定位。
* `dba_uid` - (必填, 变更时重建) 实例DBA的UID，必须为已注册用户的UID。
* `safe_rule` - (必填, 变更时重建) 实例的安全规则名称。
* `query_timeout` - (必填) 查询超时时间，单位：秒。
* `export_timeout` - (必填) 导出超时时间，单位：秒。
* `ecs_instance_id` - (可选) ECS实例ID。当`instance_source`为`ECS_OWN`时，必须传递此值。
* `vpc_id` - (可选) VPC ID。当`instance_source`为`VPC_IDC`时，必须传递此值。
* `ecs_region` - (可选) 实例所属区域。当`instance_source`为`RDS`, `ECS_OWN`, 或`VPC_IDC`时，必须传递此值。
* `sid` - (可选) SID。当`instance_type`为`PostgreSQL`或`Oracle`时，必须传递此值。
* `data_link_name` - (可选) 跨数据库查询datalink名称。
* `ddl_online` - (可选) 是否使用在线服务，目前仅支持MySQL和PolarDB。有效值：`0`（不使用），`1`（原生在线DDL优先），`2`（DMS无锁表结构变更优先）。
* `use_dsql` - (可选) 是否开启跨实例查询。有效值：`0`（未开启），`1`（开启）。
* `skip_test` - (可选) 是否忽略实例连通性测试。有效值：`true`, `false`。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - DMS企业实例的ID，格式为 `<host>:<port>`。
* `dba_nick_name` - 实例DBA的昵称。
* `status` - 实例状态。
* `dba_id` - 数据库实例的DBA ID。
* `safe_rule_id` - 数据库实例的安全规则ID。
* `instance_id` - 数据库实例的ID。
* `skip_test` - 是否忽略实例连通性测试。有效值：`true`, `false`。
* `instance_alias` - 数据库实例的别名。字段`instance_alias`从版本1.100.0起已被废弃，建议使用`instance_name`替代。