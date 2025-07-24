---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_readonly_instance"
sidebar_current: "docs-alibabacloudstack-resource-polardb-readonly-instance"
description: |-
  Provides an PolarDB readonly instance resource.
---

# alibabacloudstack_polardb_readonly_instance
Provides an PolarDB readonly instance resource.


## Argument Reference

The following arguments are supported:

* `engine_version` - (Required, ForceNew) Database version. 
* `master_db_instance_id` - (Required) ID of the master instance.
* `instance_type` - (Required) DB Instance type. 
* `instance_storage` - (Required) User-defined DB instance storage space. 
* `instance_name` - (Optional) The name of DB instance. It a string of 2 to 256 characters.
* `parameters` - (Optional) Set of parameters needs to be set after DB instance was launched.
* `zone_id` - (Optional, ForceNew) The Zone to launch the DB instance.
* `vswitch_id` - (Optional, ForceNew) The virtual switch ID to launch DB instances in one VPC.
* `db_instance_storage_type` - (Required) The storage type of the instance. Valid values:
    local_ssd: specifies to use local SSDs. This value is recommended.
    cloud_ssd: specifies to use standard SSDs.
    cloud_essd: specifies to use enhanced SSDs (ESSDs).
    cloud_essd2: specifies to use enhanced SSDs (ESSDs).
    cloud_essd3: specifies to use enhanced SSDs (ESSDs).
* `db_instance_class` - (Optional) DB Instance class .
* `db_instance_storage` - (Optional) DB Instance storage .
* `master_instance_id` - (Optional, ForceNew) ID of the master instance .
* `db_instance_description` - (Optional) The description of the DB instance.
-> **NOTE:** Because of data backup and migration, change DB instance type and storage would cost 15~20 minutes. Please make full preparation before changing them.
* `parameters` - (Optional) Set of parameters needs to be set after DB instance was launched. Available parameters can refer to the latest docs View database parameter templates.
  * `name` - (Required) The parameter name.
  * `value` - (Required) The parameter value.

## Attributes Reference

The following attributes are exported:

* `id` - The PolarDB instance ID.
* `engine` - Database type.
* `port` - PolarDB database connection port.
* `connection_string` - PolarDB database connection string.
* `db_instance_description` - The description of the DB instance .
* `engine` - Database engine type .