---
subcategory: "NAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_filesystem"
sidebar_current: "docs-Alibabacloudstack-nas-filesystem"
description: |- 
  Provides a nas Filesystem resource.
---

# alibabacloudstack_nas_filesystem
-> **NOTE:** Alias name has: `alibabacloudstack_nas_file_system`

Provides a nas Filesystem resource.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testAccAlibabacloudStackNasFileSystem97828"
}

data "alibabacloudstack_nas_protocols" "example" {
  type = "Capacity"
}

data "alibabacloudstack_nas_zones" "default" {}

resource "alibabacloudstack_nas_file_system" "default" {
  zone_id        = data.alibabacloudstack_nas_zones.default.zones.0.zone_id
  description    = "tf-testAccAlibabacloudStackNasFileSystem97828"
  protocol_type  = data.alibabacloudstack_nas_protocols.example.protocols.0
  storage_type   = "Capacity"
  file_system_type = "standard"
  encrypt_type   = "1"
}
```

Advanced Usage with `extreme` File System Type

```terraform
resource "alibabacloudstack_nas_file_system" "foo" {
  file_system_type = "extreme"
  protocol_type    = "NFS"
  zone_id          = "cn-hangzhou-f"
  storage_type     = "standard"
  description      = "tf-testAccNasConfig"
  capacity         = "100"
  encrypt_type     = "2"
  kms_key_id       = "your-kms-key-id"
}
```

## Argument Reference

The following arguments are supported:

* `storage_type` - (Required, ForceNew) The storage type.
  * When `file_system_type = standard`, the values are:
    * `Performance`
    * `Capacity`
    * `Premium`
  * When `file_system_type = extreme`, the values are:
    * `standard`
    * `advance`
  * When `file_system_type = cpfs`, the values are:
    * `advance_100` (100MB/s/TiB baseline)
    * `advance_200` (200MB/s/TiB baseline)

* `protocol_type` - (Required, ForceNew) File transfer protocol type.
  * When `file_system_type = standard`, the values are:
    * `NFS`
    * `SMB`
  * When `file_system_type = extreme`, the value is:
    * `NFS`
  * When `file_system_type = cpfs`, the value is:
    * `cpfs`

* `description` - (Optional) File system description. Restrictions:
  * 2~128 English or Chinese characters in length.
  * Must start with upper and lower case letters or Chinese, and cannot start with 'http:// ' or 'https://'.
  * Can contain numbers, colons (:), underscores (_), or dashes (-).

* `encrypt_type` - (Optional, ForceNew) Whether the file system is encrypted. Use the KMS service hosting key to encrypt and store the file system disk data. When reading and writing encrypted data, there is no need to decrypt it.
  * Valid values:
    * `0` (default): not encrypted.
    * `1`: NAS managed key. NAS managed keys are supported when `file_system_type = standard` or `file_system_type = extreme`.
    * `2`: User management key. You can manage keys only when `file_system_type = extreme`.

* `file_system_type` - (Optional, ForceNew) File system type.
  * Valid values:
    * `standard` (default): Universal NAS
    * `extreme`: Extreme NAS
    * `cpfs`: File Storage CPFS

* `capacity` - (Optional) The capacity of the file system. Required when `file_system_type = extreme`. Unit: GiB. **Note**: The minimum value is 100.

* `zone_id` - (Optional, ForceNew) The zone ID. The usable area refers to the physical area where power and network are independent of each other in the same region.
  * When `file_system_type = standard`, this parameter is optional. By default, a zone that meets the conditions is randomly selected based on the `protocol_type` and `storage_type` configurations.
  * This parameter is required when `file_system_type = extreme` or `file_system_type = cpfs`.

* `kms_key_id` - (Optional) The id of the KMS key. Required when `encrypt_type = 2`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the File System.
* `capacity` - The capacity of the file system.
* `zone_id` - The zone ID. The usable area refers to the physical area where power and network are independent of each other in the same region.
  * When `file_system_type = standard`, this parameter is optional. By default, a zone that meets the conditions is randomly selected based on the `protocol_type` and `storage_type` configurations.
  * This parameter is required when `file_system_type = extreme` or `file_system_type = cpfs`.
* `kms_key_id` - The id of the KMS key used for encryption when `encrypt_type = 2`.

## Import

Nas File System can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_nas_file_system.foo 1337849c59
```