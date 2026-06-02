---
subcategory: "ApsaraDB RDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_db_readonly_instance"
sidebar_current: "docs-Alibabacloudstack-resource-db-readonly-instance"
description: |-
  Provides an RDS readonly instance resource.
---

# alibabacloudstack_db_readonly_instance
Provides an RDS readonly instance resource.


## Argument Reference

The following arguments are supported:

* `engine_version` - (Required, ForceNew) Database version. Value options can refer to the latest docs [CreateDBInstance](https://www.alibabacloud.com/help/doc-detail/26228.htm) `EngineVersion`.
* `db_instance_storage_type` - (Required, ForceNew) The storage type of the instance. Valid values: `local_ssd`, `cloud_ssd`, `cloud_essd`, `cloud_essd2`, `cloud_essd3`, `cloud_pperf`, `cloud_sperf`.
* `master_db_instance_id` - (Optional, ForceNew, Deprecated) ID of the master instance. Use `master_instance_id` instead.
* `master_instance_id` - (Optional, ForceNew) ID of the master instance.
* `instance_type` - (Optional, Deprecated) DB Instance type. For details, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/26312.htm). Use `db_instance_class` instead. At least one of `instance_type` or `db_instance_class` must be specified.
* `db_instance_class` - (Optional) DB Instance class. At least one of `instance_type` or `db_instance_class` must be specified.
* `instance_storage` - (Optional, Deprecated) User-defined DB instance storage space. For details, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/26312.htm). Use `db_instance_storage` instead. At least one of `instance_storage` or `db_instance_storage` must be specified.
* `db_instance_storage` - (Optional) DB instance storage. At least one of `instance_storage` or `db_instance_storage` must be specified.
* `instance_name` - (Optional, Deprecated) The name of DB instance. It is a string of 2 to 256 characters. Use `db_instance_description` instead.
* `db_instance_description` - (Optional) The description of the DB instance. It is a string of 2 to 256 characters.
* `zone_id` - (Optional, ForceNew) The Zone to launch the DB instance.
* `vswitch_id` - (Optional, ForceNew) The virtual switch ID to launch DB instances in one VPC.
* `parameters` - (Optional) Set of parameters needs to be set after DB instance was launched. Available parameters can refer to the latest docs [View database parameter templates](https://www.alibabacloud.com/help/doc-detail/26284.htm).
    * `name` - (Required) The parameter name.
    * `value` - (Required) The parameter value.
* `tags` - (Optional) A mapping of tags to assign to the resource.
    - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
    - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.
* `force_restart` - (Optional) Whether to force restart the instance when parameters are changed. Default to `false`.

-> **NOTE:** Because of data backup and migration, change DB instance type and storage would cost 15~20 minutes. Please make full preparation before changing them.

## Attributes Reference

The following attributes are exported:

* `id` - The RDS instance ID.
* `engine` - Database type.
* `engine_version` - Database engine version.
* `port` - RDS database connection port.
* `connection_string` - RDS database connection string.
* `zone_id` - The zone ID of the instance.
* `vswitch_id` - The VSwitch ID of the instance.
* `master_instance_id` - The ID of the master instance.
* `db_instance_class` - The instance class.
* `db_instance_storage` - The instance storage.
* `db_instance_storage_type` - The storage type of the instance.
* `db_instance_description` - The description of the DB instance.