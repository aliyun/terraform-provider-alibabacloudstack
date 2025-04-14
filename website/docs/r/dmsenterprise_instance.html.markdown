---
subcategory: "DMSEnterprise"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dmsenterprise_instance"
sidebar_current: "docs-Alibabacloudstack-dmsenterprise-instance"
description: |-  
  Provides a dmsenterprise Instance resource.
---

# alibabacloudstack_dms_enterprise_instance
-> **NOTE:** Alias name has: `alibabacloudstack_dmsenterprise_instance`

Provides a DMS Enterprise Instance resource.

-> **NOTE:** API users must first register in DMS.

## Example Usage

```terraform
variable "name" {
    default = "tf-testAccDmsEnterpriseInstance4641"
}

variable "password" {
}

data "alibabacloudstack_account" "current" {}

resource "alibabacloudstack_db_instance" "instance" {
    engine           = "MySQL"
    engine_version   = "5.6"
    instance_type    = "rds.mysql.t1.small"
    instance_storage = "10"
    instance_name    = "${var.name}"
    security_ips     = ["0.0.0.0/0"]
    storage_type     = "local_ssd"
}

resource "alibabacloudstack_db_account" "account" {
    instance_id = "${alibabacloudstack_db_instance.instance.id}"
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
    host              = "${alibabacloudstack_db_instance.instance.connection_string}"
    port              = 3306
    database_user     = "${alibabacloudstack_db_account.account.name}"
    database_password = "${alibabacloudstack_db_account.account.password}"
    instance_name     = "tf-testAccDmsEnterpriseInstance4641"
    dba_uid           = "${alibabacloudstack_dms_enterprise_user.default.uid}"
    safe_rule         = "自由操作"
    query_timeout     = 70
    export_timeout    = 2000
    ecs_region        = "cn-shanghai"
}
```

## Argument Reference

The following arguments are supported:

* `tid` - (Optional) The tenant ID.
* `instance_type` - (Required) Database type. Valid values: `MySQL`, `SQLServer`, `PostgreSQL`, `Oracle`, `DRDS`, `OceanBase`, `Mongo`, `Redis`.
* `instance_source` - (Required) The source of the database instance. Valid values: `PUBLIC_OWN`, `RDS`, `ECS_OWN`, `VPC_IDC`.
* `network_type` - (Required, ForceNew) Network type. Valid values: `CLASSIC`, `VPC`.
* `env_type` - (Required) Environment type. Valid values: `product` (production environment), `dev` (development environment), `pre` (pre-release environment), `test` (test environment), `sit` (SIT environment), `uat` (UAT environment), `pet` (pressure test environment), `stag` (STAG environment).
* `host` - (Required, ForceNew) Host address of the target database.
* `port` - (Required, ForceNew) Access port of the target database.
* `database_user` - (Required) Database access account.
* `database_password` - (Required) Database access password.
* `instance_name` - (Required) Instance name, to help users quickly distinguish positioning.
* `dba_uid` - (Required, ForceNew) The DBA of the instance is passed into the Alibaba Cloud UID of the DBA.
* `safe_rule` - (Required, ForceNew) The security rule of the instance is passed into the name of the security rule in the enterprise.
* `query_timeout` - (Required) Query timeout time, unit: s (seconds).
* `export_timeout` - (Required) Export timeout, unit: s (seconds).
* `ecs_instance_id` - (Optional) ECS instance ID. This value must be passed when the value of `instance_source` is `ECS_OWN`.
* `vpc_id` - (Optional) VPC ID. This value must be passed when the value of `instance_source` is `VPC_IDC`.
* `ecs_region` - (Optional) The region where the instance is located. This value must be passed when the value of `instance_source` is `RDS`, `ECS_OWN`, or `VPC_IDC`.
* `sid` - (Optional) The SID. This value must be passed when `instance_type` is `PostgreSQL` or `Oracle`.
* `data_link_name` - (Optional) Cross-database query datalink name.
* `ddl_online` - (Optional) Whether to use online services, currently only supports MySQL and PolarDB. Valid values: `0` (Not used), `1` (Native online DDL priority), `2` (DMS lock-free table structure change priority).
* `use_dsql` - (Optional) Whether to enable cross-instance query. Valid values: `0` (not open), `1` (open).
* `skip_test` - (Optional) Whether the instance ignores test connectivity. Valid values: `true`, `false`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the DMS enterprise instance and format as `<host>:<port>`.
* `dba_nick_name` - The instance DBA nickname.
* `status` - The instance status.
* `dba_id` - The DBA ID of the database instance.
* `safe_rule_id` - The safe rule ID of the database instance.
* `instance_id` - The ID of the database instance.
* `skip_test` - Whether the instance ignores test connectivity. Valid values: `true`, `false`.
* `instance_alias` - The alias of the database instance. Field 'instance_alias' has been deprecated from version 1.100.0. Use 'instance_name' instead.

## Import

DMS Enterprise can be imported using the host and port, e.g.

```bash
$ terraform import alibabacloudstack_dms_enterprise_instance.example rm-uf648hgs7874xxxx.mysql.rds.aliyuncs.com:3306
```