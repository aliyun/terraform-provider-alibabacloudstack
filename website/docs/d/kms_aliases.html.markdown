---
subcategory: "KMS"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_kms_aliases"
sidebar_current: "docs-apsarastack-datasource-kms-aliases"
description: |-
    Provides a list of available KMS Aliases.
---

# apsarastack\_kms\_aliases

This data source provides a list of KMS aliases in an Apsarastack Cloud account according to the specified filters.
 

## Example Usage

```
# Declare the data source
data "apsarastack_kms_aliases" "kms_aliases" {  
  name_regex = "alias/tf-testKmsAlias_123"
}

output "first_key_id" {
  value = "${data.apsarastack_kms_keys.kms_keys_ds.keys.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of KMS aliases IDs. The value is same as KMS alias_name.
* `name_regex` - (Optional) A regex string to filter the results by the KMS alias name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` -  A list of kms aliases IDs. The value is same as KMS alias_name. 
* `names` -  A list of KMS alias name.
* `aliases` - A list of KMS User alias. Each element contains the following attributes:
  * `id` - ID of the alias. The value is same as KMS alias_name.
  * `key_id` - ID of the key.
  * `alias_name` - The unique identifier of the alias.

