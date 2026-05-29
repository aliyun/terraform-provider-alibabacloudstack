---
subcategory: "防暴力破解安全服务"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_aqs_web_locks"
sidebar_current: "docs-Alibabacloudstack-datasource-aqs-web-locks"
description: |-
  查询安骑士（AQS）Web防篡改配置
---

# alibabacloudstack_aqs_web_locks

查询安骑士（AQS）Web防篡改配置。安骑士是阿里云的安全产品，提供服务器安全防护，Web防篡改是其中的一项功能，用于保护网站文件不被非法篡改。

## 示例用法

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

## 参数说明

以下参数支持作为过滤条件：

* `ids` (列表)：Web Lock配置ID列表，用于过滤特定的Web Lock配置。

* `instanceid` (字符串)：实例ID，用于通过ECS实例ID过滤Web Lock配置。

## 属性说明

以下属性被导出：

* `id` (字符串)：数据源ID，由过滤条件生成的哈希值。

* `weblocks` (列表)：返回的Web Lock配置列表，每个元素包含以下属性：
  * `id` (字符串)：Web Lock配置的ID。
  * `audit_count` (整数)：审计计数。
  * `block_count` (整数)：阻断计数。
  * `client_status` (字符串)：客户端状态。
  * `defence_type` (字符串)：防御类型。
  * `dir_count` (整数)：目录计数。
  * `intranet_ip` (字符串)：内网IP地址。
  * `internet_ip` (字符串)：公网IP地址。
  * `instance_name` (字符串)：实例名称。
  * `lock_configs` (列表)：锁定配置列表，每个元素包含以下属性：
    * `id` (整数)：Web Lock配置ID。
    * `defence_mode` (字符串)：防御模式，可以是'block'或其他。
    * `dir` (字符串)：受保护的目录路径。
    * `exclusive_dir` (字符串)：排除的目录列表，分号分隔。
    * `exclusive_file` (字符串)：排除的文件列表，分号分隔。
    * `exclusive_file_type` (字符串)：排除的文件类型列表，分号分隔。
    * `inclusive_file_type` (字符串)：受保护的文件类型列表，分号分隔。
    * `local_backup_dir` (字符串)：锁定文件的本地备份目录。
    * `mode` (字符串)：操作模式，可以是'whitelist'或其他。
  * `os` (字符串)：操作系统类型。
  * `os_name` (字符串)：操作系统名称。
  * `service_code` (字符串)：服务代码。
  * `service_detail` (字符串)：服务详情。
  * `service_status` (字符串)：服务状态。
  * `status` (字符串)：状态。
  * `uuid` (字符串)：关联服务器的UUID。