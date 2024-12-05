---
subcategory: "Table Store (OTS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ots_tables"
sidebar_current: "docs-alibabacloudstack-datasource-ots-tables"
description: |-
    Provides a list of ots tables to the user.
---

# alibabacloudstack\_ots\_tables

This data source provides the ots tables of the current Alibaba Cloud user.


## Example Usage

``` terraform
data "alibabacloudstack_ots_tables" "tables_ds" {
  instance_name = "sample-instance"
  name_regex    = "sample-table"
  output_file   = "tables.txt"
}

output "first_table_id" {
  value = "${data.alibabacloudstack_ots_tables.tables_ds.tables.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - The name of OTS instance.
* `ids` - (Optional) A list of table IDs.
* `name_regex` - (Optional) A regex string to filter results by table name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of table IDs.
* `names` - A list of table names.
* `tables` - A list of tables. Each element contains the following attributes:
  * `id` - ID of the table. The value is `<instance_name>:<table_name>`.
  * `instance_name` - The OTS instance name.
  * `table_name` - The table name of the OTS which could not be changed.
  * `primary_key` - The property of `TableMeta` which indicates the structure information of a table.
    * `name` - Name of the property.
    * `type` - Type of the property. Availably value is {1 for int, 2 for string, 3 for binary}
  * `time_to_live` - The retention time of data stored in this table.
  * `max_version` - The maximum number of versions stored in this table.
	
