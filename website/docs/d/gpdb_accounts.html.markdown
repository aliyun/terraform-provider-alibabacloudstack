---
subcategory: "GraphDatabase(GPDB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_gpdb_accounts"
sidebar_current: "docs-Alibabacloudstack-datasource-gpdb-accounts"
description: |- 
  Provides a list of gpdb accounts owned by an alibabacloudstack account.
---

# alibabacloudstack_gpdb_accounts

This data source provides a list of gpdb accounts in an alibabacloudstack account according to the specified filters.

## Example Usage

Basic Usage

```terraform
data "alibabacloudstack_gpdb_accounts" "ids" {
  db_instance_id = "example_value"
  ids            = ["my-Account-1", "my-Account-2"]
}

output "gpdb_account_id_1" {
  value = data.alibabacloudstack_gpdb_accounts.ids.accounts.0.id
}

data "alibabacloudstack_gpdb_accounts" "nameRegex" {
  db_instance_id = "example_value"
  name_regex     = "^my-Account"
}

output "gpdb_account_id_2" {
  value = data.alibabacloudstack_gpdb_accounts.nameRegex.accounts.0.id
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, ForceNew) The ID of the GPDB instance.
* `ids` - (Optional, ForceNew) A list of Account IDs. Its element value is the same as Account Name.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Account name.
* `status` - (Optional, ForceNew) The status of the account. Valid values: `Active`, `Creating`, and `Deleting`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Account names.
* `accounts` - A list of Gpdb Accounts. Each element contains the following attributes:
  * `account_description` - The description of the account.
  * `id` - The ID of the Account. Its value is the same as the Account Name.
  * `account_name` - The name of the account.
  * `db_instance_id` - The ID of the GPDB instance.
  * `status` - The status of the account. Valid values: `Active`, `Creating`, and `Deleting`.