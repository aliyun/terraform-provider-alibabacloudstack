---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_disk"
sidebar_current: "docs-alibabacloudstack-resource-disk"
description: |- 
  Provides a ECS Disk resource.
---

# alibabacloudstack_disk
-> **NOTE:** Alias name has: `alibabacloudstack_ecs_disk`

Provides a ECS disk resource.

-> **NOTE:** One of `size` or `snapshot_id` is required when specifying an ECS disk. If both are specified, `size` must be greater than the size of the snapshot represented by `snapshot_id`. Currently, `alibabacloudstack_disk` does not support resizing disks.

## Example Usage

```hcl
# Create a new ECS disk with specific configurations.
resource "alibabacloudstack_disk" "example" {
  zone_id             = "${data.alibabacloudstack_zones.default.zones.0.id}"
  size                = 50
  category           = "cloud_efficiency"
  name               = "TerraformTestDisk"
  description        = "This is a test disk created by Terraform."
  
  tags = {
    Environment = "Test"
    CreatedBy   = "Terraform"
  }

  encrypted          = true
  kms_key_id        = "your-kms-key-id"
  delete_auto_snapshot = true
  delete_with_instance = false
  enable_auto_snapshot = true
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Optional, ForceNew) The ID of the zone in which to create the pay-as-you-go disk. If you do not specify `instance_id`, you must specify `zone_id`. You cannot specify both `zone_id` and `instance_id` in a request.
* `name` - (Optional) Name of the ECS disk. This name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-", ".", "_", and must not begin or end with a hyphen, and must not begin with `http://` or `https://`. Default value is null.
* `disk_name` - (Optional) The name of the disk. The name must be 2 to 128 characters in length and can contain Unicode characters under the Decimal Number category and the categories whose names contain Letter. The name can also contain colons (:), underscores (\_), periods (.), and hyphens (-).
* `description` - (Optional) Description of the disk. This description can have a string of 2 to 256 characters, It cannot begin with `http://` or `https://`. Default value is null.
* `category` - (Optional, ForceNew) Category of the disk. Valid values:
  * `all`: all disk categories
  * `cloud`: basic disk
  * `cloud_efficiency`: ultra disk
  * `cloud_ssd`: standard SSD
  * `cloud_essd`: Enterprise SSD (ESSD)
  * `cloud_auto`: ESSD AutoPL disk
  * `local_ssd_pro`: I/O-intensive local disk
  * `local_hdd_pro`: throughput-intensive local disk
  * `cloud_essd_entry`: ESSD Entry disk
  * `elastic_ephemeral_disk_standard`: standard elastic ephemeral disk
  * `elastic_ephemeral_disk_premium`: premium elastic ephemeral disk
  * `ephemeral`: retired local disk
  * `ephemeral_ssd`: retired local SSD

Default value: `all`.

* `size` - (Optional) The new disk capacity. Unit: GiB. Valid values depend on the disk category:
  * System disk:
    * Basic disk (`cloud`): 20 to 500.
    * ESSD (`cloud_essd`): Valid values vary based on the performance level of the ESSD.
      * PL0 ESSD: 1 to 2048.
      * PL1 ESSD: 20 to 2048.
      * PL2 ESSD: 461 to 2048.
      * PL3 ESSD: 1261 to 2048.
    * ESSD AutoPL disk (`cloud_auto`): 1 to 2048.
    * Other disk categories: 20 to 2048.
  * Data disk:
    * Ultra disk (`cloud_efficiency`): 20 to 32768.
    * Standard SSD (`cloud_ssd`): 20 to 32768.
    * ESSD (`cloud_essd`): Valid values vary based on the performance level of the ESSD.
      * PL0 ESSD: 1 to 32768.
      * PL1 ESSD: 20 to 32768.
      * PL2 ESSD: 461 to 32768.
      * PL3 ESSD: 1261 to 32768.
    * Basic disk (`cloud`): 5 to 2000.
    * ESSD AutoPL disk (`cloud_auto`): 1 to 32768.
    * Standard elastic ephemeral disk (`elastic_ephemeral_disk_standard`): 64 to 8192.
    * Premium elastic ephemeral disk (`elastic_ephemeral_disk_premium`): 64 to 8192.

> The new disk capacity must be larger than the original disk capacity.

* `snapshot_id` - (Optional) The ID of the snapshot to use to create the disk. Snapshots that were created on or before July 15, 2013 cannot be used to create disks. The following limits apply to `SnapshotId` and `Size`:
  * If the size of the snapshot specified by `SnapshotId` is larger than the value of `Size`, the size of the created disk is equal to the specified snapshot size.
  * If the size of the snapshot specified by `SnapshotId` is smaller than the value of `Size`, the size of the created disk is equal to the value of `Size`.
  * You cannot create elastic ephemeral disks from snapshots.

* `kms_key_id` - (Optional) The ID of the Key Management Service (KMS) key that is used by the cloud disk.
* `encrypted` - (Optional, ForceNew) Specifies whether to query only encrypted cloud disks.
  * `true`: queries only encrypted cloud disks.
  * `false`: does not query encrypted cloud disks.

Default value: `false`.

* `encrypt_algorithm` - (Optional, ForceNew) The encryption algorithm used for the disk.
* `delete_auto_snapshot` - (Optional) Specifies whether to delete the automatic snapshots of the disk when the disk is released. This parameter is empty by default, which indicates that the current value remains unchanged.
* `delete_with_instance` - (Optional) Specifies whether to release the disk along with its associated instance. This parameter is empty by default, which indicates that the current value remains unchanged. An error is returned if you set `DeleteWithInstance` to `false` in one of the following cases:
  * The disk is a local disk.
  * The disk is a basic disk and is not removable. If the `Portable` attribute of a disk is set to `false`, the disk is not removable.
* `enable_auto_snapshot` - (Optional) Specifies whether to enable the automatic snapshot policy feature for the cloud disk. Valid values:
  * `true`
  * `false`

This parameter is empty by default, which indicates that the current value remains unchanged. By default, the automatic snapshot policy feature is enabled for cloud disks. You only need to associate an automatic snapshot policy with a cloud disk before you can use the policy.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `availability_zone` - The availability zone of the disk.
* `zone_id` - The ID of the zone in which the disk was created.
* `name` - The name of the disk.
* `disk_name` - The name of the disk.
* `status` - The lifecycle status of the EBS device. For more information, see [Disk status](https://www.alibabacloud.com/help/en/doc-detail/25689.html). Valid values:
  * `In_use`: The EBS device is in use.
  * `Available`: The EBS device can be attached.
  * `Attaching`: The EBS device is being attached.
  * `Detaching`: The EBS device is being detached.
  * `Creating`: The EBS device is being created.
  * `ReIniting`: The EBS device is being initialized.
* `auto_snapshot_policy_id` - The ID of the automatic snapshot policy that is applied to the cloud disk.
* `enable_automated_snapshot_policy` - Specifies whether an automatic snapshot policy is applied to the cloud disk.
  * `true`: An automatic snapshot policy is applied to the cloud disk.
  * `false`: No automatic snapshot policy is applied to the cloud disk.

Default value: `false`.