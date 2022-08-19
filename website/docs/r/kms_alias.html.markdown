---
subcategory: "KMS"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_kms_alias"
sidebar_current: "docs-apsarastack-resource-kms-alias"
description: |-
  Provides a ApsaraStack KMS Alias resource.
---

# apsarastack\_kms_alias

Create an alias for the master key (CMK).



## Example Usage

Basic Usage

```
resource "apsarastack_kms_key" "key" {}

resource "apsarastack_kms_alias" "alias" {
  alias_name = "alias/test_kms_alias"
  key_id     = apsarastack_kms_key.key.id
}
```

## Argument Reference

The following arguments are supported:

* `alias_name` - (Required, ForceNew) The alias of CMK. `Encrypt`、`GenerateDataKey`、`DescribeKey` can be called using aliases. Length of characters other than prefixes: minimum length of 1 character and maximum length of 255 characters. Must contain prefix `alias/`.
* `key_id` - (Required) The id of the key.

-> **NOTE:** Each alias represents only one master key(CMK).

-> **NOTE:** Within an area of the same user, alias is not reproducible.

-> **NOTE:** UpdateAlias can be used to update the mapping relationship between alias and master key(CMK).


## Attributes Reference

* `id` - The ID of the alias.

