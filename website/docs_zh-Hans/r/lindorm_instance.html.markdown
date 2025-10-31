---
subcategory: "Lindorm"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_lindorm_instance"
sidebar_current: "docs-alibabacloudstack-resource-lindorm-instance"
description: |-
  提供阿里云专有云Lindorm实例资源
---

# alibabacloudstack\_lindorm\_instance

提供一个Lindorm实例资源。

## 示例用法

### 基本用法

```terraform
variable "name" {
    default = "tf-testacc97984"
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
  vswitch_name   = "tf-testaccvpcvswitch97984"
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  vpc_id         = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block     = "172.16.0.0/24"
  enable_ipv6    = true
}

resource "alibabacloudstack_lindorm_instance" "example" {
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
```

## 参数参考

以下参数被支持：

* `zone_id` - （必需，强制新建）部署实例的可用区ID。
* `instance_alias` - （必需）实例别名。
* `cpu_brand` - （必需，强制新建）要使用的CPU品牌。
* `disk_category` - （可选，强制新建）磁盘类别。
* `engine_type` - （必需，强制新建）要使用的引擎类型。
* `instance_type` - （必需）实例规格。
* `vpc_id` - （必需，强制新建）部署实例的专有网络ID。
* `vswitch_id` - （必需，强制新建）与指定专有网络关联的交换机ID。
* `lindorm_num` - （必需）Lindorm节点数量。
* `local_disk_num` - （可选）本地磁盘数量。取值范围：1到10。默认值：1。
* `local_disk_size` - （必需，强制新建）本地磁盘大小。单位：GiB。

## 属性参考

以下属性会被导出：

* `id` - 实例ID。
* `instance_id` - 实例ID。
* `instance_status` - 实例状态。
* `create_time` - 实例创建时间。
* `instance_storage` - 实例存储容量。
* `deletion_protection` - 是否启用了删除保护。
* `disk_usage` - 磁盘使用情况。
* `enable_fs` - 是否启用了文件系统。
* `switch_l_proxy_flag` - L_Proxy开关是否启用。
* `switch_ssl_encryption_flag` - 是否启用了SSL加密。
* `ali_uid` - 阿里云账户的UID。

## 导入

Lindorm实例可以使用ID进行导入，例如：

```shell
$ terraform import alibabacloudstack_lindorm_instance.example li-12345678
```