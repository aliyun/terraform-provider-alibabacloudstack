---
subcategory: "Bare-Metal Management Service (BMS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bms_keypairs"
description: |-
    Query key pair information of Alibaba Cloud Bare Metal Server (BMS)
---

# alibabacloudstack_bms_keypairs

Query key pair information of Alibaba Cloud Bare Metal Server (BMS).

## Example Usage

```hcl
variable "name" {
  default = "tf-bms-keypair12432"
}

resource "alibabacloudstack_bms_keypair" "default" {
  name = var.name
}

data "alibabacloudstack_bms_keypairs" "default" {
  name_regex = alibabacloudstack_bms_keypair.default.name
}
```

## Argument Reference

The following arguments are supported:

* `ids` (List) - A list of key pair names used to filter results. If specified, only key pairs with names in the list will be returned.
* `name_regex` (String) - A regular expression for the name to filter results. If specified, only key pairs with names matching the regex will be returned.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` (String) - The ID of the key pair, which is the same as the name. In BMS, the key pair name is unique and thus used as the ID.
* `create_time` (String) - The creation time of the key pair.
* `key_pair_fingerprint` (String) - The fingerprint information of the key pair.
* `name` (String) - The name of the key pair.
* `public_key` (String) - The public key content of the key pair.
* `shared` (Integer) - The sharing status of the key pair: 0 indicates not shared, 1 indicates shared.
* `update_time` (String) - The last update time of the key pair.