---
subcategory: "Lindorm"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_lindorm_connection"
sidebar_current: "docs-alibabacloudstack-resource-lindorm-connection"
description: |-
  提供阿里云Lindorm连接资源
---

# alibabacloudstack\_lindorm\_connection

提供Lindorm连接资源。

## 使用示例

### 基本用法

```terraform
variable "name" {
  default = "tf-testacc"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  description     = "modify_description"
  vswitch_name   = "tf-testaccvpcvswitch"
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  vpc_id         = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block     = "172.16.0.0/24"
  enable_ipv6    = true
}

resource "alibabacloudstack_lindorm_instance" "default" {
  zone_id               = "${data.alibabacloudstack_zones.default.zones.0.id}"
  instance_alias        = "${var.name}"
  cpu_brand             = "Intel"
  disk_category         = "HHD"
  engine_type           = "tsdb"
  instance_type         = "lindorm.c.xlarge"
  vpc_id                = "${alibabacloudstack_vpc_vpc.default.id}"
  vswitch_id            = "${alibabacloudstack_vpc_vswitch.default.id}"
  lindorm_num           = 2
  local_disk_num        = 2
  local_disk_size       = "400"
}

resource "alibabacloudstack_lindorm_connection" "example" {
  instance_id     = alibabacloudstack_lindorm_instance.default.id
  engine_type     = "lindorm"
  classic_network = false
}
```

## 参数参考

以下参数被支持：

* `instance_id` - (必选, 不可变) Lindorm实例的ID。
* `engine_type` - (必选, 不可变) 引擎类型。
* `classic_network` - (可选) 是否启用公网访问。默认为`false`。

## 属性参考

以下属性会被导出：

* `id` - 资源的ID，格式为`<instance_id>:<engine_type>`。
* `engine_username` - 数据库引擎的用户名。
* `engine_password` - 数据库引擎的密码。
* `net_info_list` - 网络连接信息列表。每个元素包含以下属性：
  * `port` - 端口号。
  * `connection_string` - 连接字符串。
  * `net_type` - 网络类型。有效值：`0`(公网), `2`(经典网络), `3`(VPC网络)。
  * `access_type` - 访问类型。

## 导入

Lindorm连接可以通过ID导入，例如：

```shell
$ terraform import alibabacloudstack_lindorm_connection.example li-12345678:lindorm
```