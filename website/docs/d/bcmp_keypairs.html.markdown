---
subcategory: "Bare Metal Computing Platform (BMCP)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bcmp_keypairs"
sidebar_current: "docs-Alibabacloudstack-datasource-bcmp-keypairs"
description: |- 
  Provides a list of BCMP keypairs owned by an AlibabaCloudStack account.
---

# alibabacloudstack_bcmp_keypairs

This data source provides a list of BCMP key pairs in an AlibabaCloudStack account according to the specified filters.

## Example Usage

```hcl
# Declare the resource
resource "alibabacloudstack_bcmp_keypair" "default" {
  key_pair_name = "exampleKeyPair"
}

# Retrieve all key pairs matching the name_regex
data "alibabacloudstack_bcmp_keypairs" "default" {
  name_regex = "${alibabacloudstack_bcmp_keypair.default.key_pair_name}"
}

output "key_pairs" {
  value = data.alibabacloudstack_bcmp_keypairs.default.key_pairs
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter the resulting key pairs by their names.
* `ids` - (Optional) A list of key pair IDs. If provided, only the key pairs with these IDs will be returned.
* `finger_print` - (Optional) The fingerprint of the key pair.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of key pair names.
* `key_pairs` - A list of key pairs. Each element contains the following attributes:
  * `id` - ID of the key pair.
  * `key_pair_name` - Name of the key pair.
  * `finger_print` - The fingerprint of the key pair.
  * `region` - The region of the key pair.
  * `create_time` - The creation time of the key pair.
  * `update_time` - The update time of the key pair.
