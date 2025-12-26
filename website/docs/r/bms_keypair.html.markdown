---
subcategory: "BMS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bms_keypair"
sidebar_current: "docs-Alibabacloudstack-bms-keypair"
description: |-
  Manage Bare Metal Server key pairs
---

# alibabacloudstack_bms_keypair

Manages Bare Metal Server key pair resources for creating, reading, and deleting key pairs.

## Example Usage

### Basic Usage

```hcl

variable "name" {
  default = "tf-bms-keypair1641"
}

variable "public_key" {
  default = "ssh-rsa ********* == root@vm010017040011"
}


resource "alibabacloudstack_bms_keypair" "default" {
  name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForcesNew) The name of the key pair. It can be up to 1-128 characters in length and can contain letters, digits, underscores (_), and hyphens (-).

* `public_key` - (Optional, ForcesNew) The public key content. If not provided, the system will automatically generate a key pair. The public key format should be in SSH public key format, such as "ssh-rsa AAAAB3NzaC1yc2EAAA...".

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the key pair, which is the same as the name.
* `key_pair_fingerprint` - The fingerprint of the key pair, in the format "[bit length] [algorithm]:[hash value] [comment]", for example "4096 SHA256:+MiS793F6Nxn/ygAGtquHw2e5grziG/AA+0ESwvF3VM root@vm010017040011 (RSA)".
* `private_key` - The private key content (sensitive information), which is only returned during creation and will be empty for subsequent reads. The format is PEM format RSA private key.
* `public_key` - The public key content, in SSH public key format. If a public key was provided during creation, the provided public key is returned; if the system generated the key pair, the generated public key is returned.