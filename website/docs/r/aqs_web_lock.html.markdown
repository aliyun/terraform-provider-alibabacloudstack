---
subcategory: "Server Guard"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_aqs_web_lock"
description: |-
  Configure Web Tamper Proofing with Security Center
---

# alibabacloudstack_aqs_web_lock

Configure Web Tamper Proofing with Security Center using credentials configured in the Provider.

## Example Usage

### Basic Usage

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
  status     = "on"
  lock_configs {
    local_backup_dir    = "/usr/local/aegis/bak1"
    inclusive_file_type = "php;jsp;asp;aspx;js;cgi;html;htm;xml;shtml;shtm;jpg;gif;png;jspx"
    defence_mode        = "block"
    mode                = "whitelist"
    dir                 = "/test/tf"
  }
  lock_configs {
    dir                 = "/test2/tf"
    local_backup_dir    = "/usr/local/aegis/bak2"
    inclusive_file_type = "php;jsp;asp;aspx;js;cgi;"
    defence_mode        = "block"
    mode                = "whitelist"
  }
  lock_configs {
    exclusive_dir       = "testpath"
    defence_mode        = "audit"
    mode                = "blacklist"
    dir                 = "/test3/tf"
    local_backup_dir    = "/usr/local/aegis/bak3"
    exclusive_file_type = "log;txt;ldb"
    exclusive_file      = "aaa.txt"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instanceid` - (Required, Forces new resource) The ID of the ECS instance. Specifies the server instance to be configured with web tamper proofing protection.
* `lock_configs` - (Required) A list of protection configurations. At least one protection directory must be configured. Each configuration contains the following parameters:
  * `dir` - (Required) The path of the protection directory. Specifies the directory to be protected against web tampering, e.g., `/test/tf`.
  * `local_backup_dir` - (Required) The path of the local backup directory. Used for secure backup of the protected directory, e.g., `/usr/local/aegis/bak`. Note that the path format may differ between Linux and Windows servers; ensure the correct format is used.
  * `defence_mode` - (Required) The protection mode. Valid values:
    * `block`: Block mode, which blocks tampering attempts.
    * `audit`: Audit mode, which only records tampering attempts without blocking them.
  * `mode` - (Required) The protection directory mode. Valid values:
    * `whitelist`: Whitelist mode, which protects only the specified protection directories and file types.
    * `blacklist`: Blacklist mode, which protects all subdirectories, file types, and specified files under the protection directory that are not excluded.
* `status` - (Optional) The protection status. Valid values:
  * `off`: Disable protection (default).
  * `on`: Enable protection.

Optional parameters within `lock_configs` (configured based on the `mode` parameter):
* `inclusive_file_type` - (Optional) A list of file types to be protected (used in whitelist mode), e.g., `php;jsp;asp;aspx;js;cgi;html;htm;xml;shtml;shtm;jpg;gif;png;jspx`. Supported file types include: php, jsp, asp, aspx, js, cgi, html, htm, xml, shtml, shtm, jpg, gif, png.
* `exclusive_dir` - (Optional) A list of directories that do not require protection (used in blacklist mode), e.g., `/home/admin/test`.
* `exclusive_file_type` - (Optional) A list of file types that do not require protection (used in blacklist mode), e.g., `jpg;png`. Supported file types include: php, jsp, asp, aspx, js, cgi, html, htm, xml, shtml, shtm, jpg, gif, png.
* `exclusive_file` - (Optional) A list of files that do not require protection (used in blacklist mode), e.g., `/home/admin/tomcat/localhost.log`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID (server UUID).
* `audit_count` - The audit count, representing the number of tampering events detected in audit mode.
* `block_count` - The block count, representing the number of tampering events blocked in block mode.
* `client_status` - The client status, indicating the running status of the Security Center client.
* `defence_type` - The protection type, indicating the currently configured protection type.
* `dir_count` - The number of protected directories, representing the total number of protected directories in the current configuration.
* `intranet_ip` - The intranet IP address of the server.
* `instance_name` - The instance name, representing the name of the server instance.
* `internet_ip` - The public IP address of the server.
* `lock_configs` - A list of protection configurations, containing the following attributes:
  * `id` - The configuration ID, representing the unique identifier of the protection configuration.
* `os` - The operating system type, representing the OS identifier of the server.
* `os_name` - The operating system name, representing the name of the server's operating system.
* `service_code` - The service code, representing the code identifier of the Security Center service.
* `service_detail` - The service details, containing detailed information about the Security Center service.
* `service_status` - The service status, indicating the running status of the Security Center service.