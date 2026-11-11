---
subcategory: "Elastic Compute Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_image_import"
description: |-
  Provides an ECS image import resource.
---

# alibabacloudstack_image_import

Import a copy of your local on-premise file to ECS, and appear as a custom replacement in the corresponding domain.

-> **NOTE:** You must upload the image file to the object storage OSS in advance.

-> **NOTE:** The region where the image is imported must be the same region as the OSS bucket where the image file is uploaded.

## Example Usage

```
resource "alibabacloudstack_image_import" "this" {
  description  = "test import image"
  architecture = "x86_64"
  image_name   = "test-import-image"
  license_type = "Auto"
  platform     = "Ubuntu"
  os_type      = "linux"
  disk_device_mapping {
    disk_image_size = 5
    oss_bucket      = "testimportimage"
    oss_object      = "root.img"
  }
}
```

## Argument Reference

The following arguments are supported:

* `architecture` - (Optional, ForceNew) Specifies the architecture of the system disk after you specify a data disk snapshot as the data source of the system disk for creating an image. Valid values: `i386` , Default is `x86_64`.
* `description` - (Optional) Description of the image. The length is 2 to 256 English or Chinese characters, and cannot begin with http: // and https: //.
* `image_name` - (Optional) The image name. The length is 2 ~ 128 English or Chinese characters. Must start with a english letter or Chinese, and cannot start with http: // and https: //. Can contain numbers, colons (:), underscores (_), or hyphens (-).
* `license_type` - (Optional, ForceNew) The type of the license used to activate the operating system after the image is imported. Default value: `Auto`. Valid values: `Auto`,`Aliyun`,`BYOL`.
* `platform` - (Optional, ForceNew) Specifies the operating system platform of the system disk after you specify a data disk snapshot as the data source of the system disk for creating an image. Valid values: `CentOS`, `RedHat`, `Fedora`, `Fedora CoreOS`, `Anolis`, `CentOS Stream`, `AlmaLinux`, `Rocky Linux`, `Ubuntu`, `Debian`, `SUSE`, `OpenSUSE`, `Aliyun`, `Kylin`, `UOS`, `CoreOS`, `Customized Linux`, `Others Linux`, `Windows Server 2003`, `Windows Server 2008`, `Windows Server 2012`, `Windows Server 2016`, `Windows Server 2019`, `Windows Server 2022`, `Windows 7`, `FreeBSD`. Default is `Ubuntu`.
* `os_type` - (Optional, ForceNew) Operating system platform type. Valid values: `windows`, Default is `linux`.
* `disk_device_mapping` - (Required, ForceNew) Description of the system with disks and snapshots under the image. 
  * `disk_image_size` - (Optional, ForceNew) Resolution size. You must ensure that the system disk space ≥ file system space. Ranges: When n = 1, the system disk: 5 ~ 500GiB, When n = 2 ~ 17, that is, data disk: 5 ~ 1000GiB, When temporary is introduced, the system automatically detects the size, which is subject to the detection result. Default is `5`.
  * `format` - (Optional, ForceNew) Image format. Value range: When the `RAW`, `VHD`, `qcow2` is imported into the image, the system automatically detects the image format, whichever comes first.
  * `oss_bucket` - (Optional, ForceNew) Save the exported OSS bucket.
  * `oss_object` - (Optional, ForceNew) The file name of your OSS Object.

-> **NOTE:** The disk_device_mapping is a list and it's first item will be used to system disk and other items are used to data disks.

## Timeouts

* `create` - (Defaults to 120 mins) Used when create the resource.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the image.
* `disk_device_mapping` - The disk device mappings of the image. The `format` of each mapping is detected and exported after the image is imported.

## Import

alibabacloudstack_image_import can be imported using the id, e.g.

```shell
$ terraform import alibabacloudstack_image_import.example <id>
```
