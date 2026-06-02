---
subcategory: "Server Guard"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_aqs_web_locks"
description: |-
  Query Web tamper-proof configurations of Server Guard (AQS).
---

# alibabacloudstack_aqs_web_locks

This data source queries Web tamper-proof configurations of Server Guard (AQS). Server Guard is Alibaba Cloud's security product that provides server security protection, and Web tamper-proof is one of its features to protect website files from illegal tampering.

## Example Usage

```hcl

variable "name" {
  default = "tf-testacc11469"
}

data "alibabacloudstack_zones" "default" {
  provider                    = alibabacloudstack-common
  available_resource_creation = "VSwitch"
  enable_details              = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  provider   = alibabacloudstack-common
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  provider   = alibabacloudstack-common
  name       = "${var.name}_vsw"
  vpc_id     = alibabacloudstack_vpc_vpc.default.id
  cidr_block = "172.16.0.0/24"
  zone_id    = data.alibabacloudstack_zones.default.zones.0.id
}

resource "alibabacloudstack_ecs_securitygroup" "default" {
  provider = alibabacloudstack-common
  name     = "${var.name}_sg"
  vpc_id   = alibabacloudstack_vpc_vpc.default.id
}

data "alibabacloudstack_images" "default" {
  provider    = alibabacloudstack-common
  name_regex  = "^aliyun_"
  most_recent = true
  owners      = "system"
}

data "alibabacloudstack_instance_types" "all" {
  provider          = alibabacloudstack-common
  sorted_by         = "Memory"
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
}

resource "alibabacloudstack_ecs_instance" "default" {
  provider                      = alibabacloudstack-common
  system_disk_category          = data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0
  instance_name                 = var.name
  user_data                     = "I_am_user_data"
  security_groups               = ["${alibabacloudstack_ecs_securitygroup.default.id}"]
  vswitch_id                    = alibabacloudstack_vpc_vswitch.default.id
  image_id                      = data.alibabacloudstack_images.default.images.0.id
  security_enhancement_strategy = "Active"
  instance_type                 = data.alibabacloudstack_instance_types.all.instance_types.0.id
  availability_zone             = data.alibabacloudstack_zones.default.zones[0].id
}

resource "alibabacloudstack_aqs_web_lock" "default" {
  instanceid = alibabacloudstack_ecs_instance.default.id
  lock_configs {
    dir                 = "/test/tf/"
    local_backup_dir    = "/usr/local/aegis/bak1"
    inclusive_file_type = "php;jsp;asp;aspx;js;cgi;html;htm;xml;shtml;shtm;jpg;gif;png;jspx"
    defence_mode        = "block"
    mode                = "whitelist"
  }
}

data "alibabacloudstack_aqs_web_locks" "default" {
  ids = ["${alibabacloudstack_aqs_web_lock.default.id}"]
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) List of Web Lock configuration IDs to filter specific Web Lock configurations.

* `instanceid` - (Optional) The instance ID used to filter Web Lock configurations by ECS instance ID.

## Attributes Reference

The following attributes are exported:

* `id` - The data source ID, which is a hash value generated from the filter conditions.

* `weblocks` - A list of returned Web Lock configurations. Each element contains the following attributes:
  * `id` - The ID of the Web Lock configuration.
  * `audit_count` - The audit count.
  * `block_count` - The block count.
  * `client_status` - The client status.
  * `defence_type` - The defense type.
  * `dir_count` - The directory count.
  * `intranet_ip` - The intranet IP address.
  * `internet_ip` - The internet IP address.
  * `instance_name` - The instance name.
  * `lock_configs` - A list of lock configurations. Each element contains the following attributes:
    * `id` - The Web Lock configuration ID (integer).
    * `defence_mode` - The defense mode, which can be 'block' or others.
    * `dir` - The protected directory path.
    * `exclusive_dir` - The list of excluded directories, separated by semicolons.
    * `exclusive_file` - The list of excluded files, separated by semicolons.
    * `exclusive_file_type` - The list of excluded file types, separated by semicolons.
    * `inclusive_file_type` - The list of protected file types, separated by semicolons.
    * `local_backup_dir` - The local backup directory for locked files.
    * `mode` - The operation mode, which can be 'whitelist' or others.
  * `os` - The operating system type.
  * `os_name` - The operating system name.
  * `service_code` - The service code.
  * `service_detail` - The service details.
  * `service_status` - The service status.
  * `status` - The status.
  * `uuid` - The UUID of the associated server.