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

## 参数参考

支持以下参数：
  * `data_link_name` - (选填) 实例datalink名称。
  * `database_password` - (必填) 实例登录账号密码。
  * `database_user` - (必填) 实例登录用户名。
  * `dba_id` - (选填) 实例dba id，与ListUsers中的id对应。
  * `dba_uid` - (必填, 变更时重建) 实例的DBA UID，必须为已注册用户的UID。
  * `ddl_online` - (选填) 实例onlineddl配置。有效值：`0`(不使用)，`1`(原生在线DDL优先)，`2`(DMS无锁表结构变更优先)。
  * `ecs_instance_id` - (选填) 实例EcsInstanceId。当`instance_source`为`ECS_OWN`时，必须传递此值。
  * `ecs_region` - (选填) 实例所属Region。当`instance_source`为`RDS`, `ECS_OWN`, 或`VPC_IDC`时，必须传递此值。
  * `env_type` - (必填) 实例所属环境。有效值：`product`(生产环境)，`dev`(开发环境)，`pre`(预发布环境)，`test`(测试环境)，`sit`(SIT环境)，`uat`(UAT环境)，`pet`(压力测试环境)，`stag`(STAG环境)。
  * `export_timeout` - (必填) 实例导出超时时间，单位：秒。
  * `host` - (必填, 变更时重建) 实例连接地址。
  * `instance_id` - (选填) 实例ID。
  * `instance_alias` - (选填) 实例别名，帮助用户快速区分定位。
  * `instance_source` - (必填) 实例来源。有效值：`PUBLIC_OWN`, `RDS`, `ECS_OWN`, `VPC_IDC`。
  * `instance_type` - (必填) 实例DB类型。有效值：`MySQL`, `SQLServer`, `PostgreSQL`, `Oracle`, `DRDS`, `OceanBase`, `Mongo`, `Redis`。
  * `network_type` - (必填, 变更时重建) 网络类型。有效值：`CLASSIC`, `VPC`。
  * `port` - (必填, 变更时重建) 实例连接端口。
  * `query_timeout` - (必填) 实例查询超时时间，单位：秒。
  * `safe_rule` - (必填, 变更时重建) 实例的安全规则名称。
  * `safe_rule_id` - (选填) 实例所对应安全规则id。
  * `sid` - (选填) 实例sid。当`instance_type`为`PostgreSQL`或`Oracle`时，必须传递此值。
  * `skip_test` - (选填) 是否忽略测试连通性。有效值：`true`, `false`。
  * `tid` - (选填) 租户ID。
  * `use_dsql` - (选填) 是否开启跨库查询。有效值：`0`(未开启)，`1`(开启)。
  * `vpc_id` - (选填) 实例所属VPC ID。当`instance_source`为`VPC_IDC`时，必须传递此值。

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `dba_id` - 实例dba id，与ListUsers中的id对应。
  * `dba_nick_name` - 实例dba昵称。
  * `ecs_instance_id` - 实例EcsInstanceId。
  * `instance_id` - 实例ID。
  * `instance_name` - 实例别名。
  * `instance_alias` - 实例别名。
  * `safe_rule_id` - 实例所对应安全规则id。
  * `status` - 实例状态。