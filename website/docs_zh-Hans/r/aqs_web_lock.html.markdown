---
subcategory: "安骑士"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_aqs_web_lock"
sidebar_current: "docs-Alibabacloudstack-resource-aqs-web-lock"
description: |-
  安骑士防篡改设置
---

# alibabacloudstack_aqs_web_lock

使用Provider配置的凭证在指定的资源集配置安骑士网页防篡改保护。

## 示例用法

### 基础用法

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

## 参数说明

支持以下参数：

* `instanceid` - (必填, 变更时重建) 云服务器ECS实例ID。用于指定需要配置防篡改保护的服务器实例。
* `lock_configs` - (必填) 防护配置列表，至少需要配置一个防护目录。每个配置包含以下参数：
  * `dir` - (必填) 防护目录路径。指定需要进行网页防篡改保护的目录，例如`/test/tf`。
  * `local_backup_dir` - (必填) 本地备份目录路径。用于对防护目录进行安全备份的路径，例如`/usr/local/aegis/bak`。Linux服务器和Windows服务器的路径格式可能不同，请确保输入正确的格式。
  * `defence_mode` - (必填) 防护模式。取值：
    * `block`：拦截模式，对篡改行为进行拦截。
    * `audit`：告警模式，仅记录篡改行为但不拦截。
  * `mode` - (必填) 防护目录模式。取值：
    * `whitelist`：白名单模式，仅对指定的防护目录和文件类型进行保护。
    * `blacklist`：黑名单模式，对防护目录下所有未排除的子目录、文件类型和指定文件进行保护。
* `status` - (可选) 防护状态。取值：
  * `off`：关闭防护（默认值）。
  * `on`：开启防护。

`lock_configs`内部可选参数（根据`mode`参数选择配置）：
* `inclusive_file_type` - (可选) 需要防护的文件类型列表（白名单模式下使用），例如`php;jsp;asp;aspx;js;cgi;html;htm;xml;shtml;shtm;jpg;gif;png;jspx`。支持的文件类型包括：php、jsp、asp、aspx、js、cgi、html、htm、xml、shtml、shtm、jpg、gif、png。
* `exclusive_dir` - (可选) 无需防护的目录列表（黑名单模式下使用），例如`/home/admin/test`。
* `exclusive_file_type` - (可选) 无需防护的文件类型列表（黑名单模式下使用），例如`jpg;png`。支持的文件类型包括：php、jsp、asp、aspx、js、cgi、html、htm、xml、shtml、shtm、jpg、gif、png。
* `exclusive_file` - (可选) 无需防护的文件列表（黑名单模式下使用），例如`/home/admin/tomcat/localhost.log`。

## 属性说明

以下属性会从API中导出：

* `id` - 资源ID（服务器UUID）。
* `audit_count` - 审计计数，表示告警模式下检测到的篡改事件数量。
* `block_count` - 拦截计数，表示拦截模式下阻止的篡改事件数量。
* `client_status` - 客户端状态，表示安骑士客户端的运行状态。
* `defence_type` - 防护类型，表示当前配置的防护类型。
* `dir_count` - 防护目录数量，表示当前配置的防护目录总数。
* `intranet_ip` - 内网IP地址，表示服务器的内网IP。
* `instance_name` - 实例名称，表示服务器实例的名称。
* `internet_ip` - 公网IP地址，表示服务器的公网IP。
* `lock_configs` - 防护配置列表，包含以下属性：
  * `id` - 配置ID，表示该防护配置的唯一标识。
* `os` - 操作系统类型，表示服务器的操作系统标识。
* `os_name` - 操作系统名称，表示服务器的操作系统名称。
* `service_code` - 服务代码，表示安骑士服务的代码标识。
* `service_detail` - 服务详情，包含安骑士服务的详细信息。
* `service_status` - 服务状态，表示安骑士服务的运行状态。