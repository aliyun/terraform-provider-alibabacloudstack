---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_image"
sidebar_current: "docs-Alibabacloudstack-ecs-image"
description: |-  
  编排云服务器（Ecs）镜像
---

# alibabacloudstack_ecs_image
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_image`

使用Provider配置的凭证在指定的资源集下编排云服务器（Ecs）镜像。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAccEcsImageShareConfigBasic4783"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details             = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name        = "${var.name}_vsw"
  vpc_id      = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block  = "172.16.0.0/24"
  zone_id     = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_ecs_securitygroup" "default" {
  name   = "${var.name}_sg"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

resource "alibabacloudstack_security_group_rule" "default" {
  type                = "ingress"
  ip_protocol         = "tcp"
  nic_type           = "intranet"
  policy             = "accept"
  port_range         = "22/22"
  priority           = 1
  security_group_id  = "${alibabacloudstack_ecs_securitygroup.default.id}"
  cidr_ip            = "172.16.0.0/24"
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_"
  most_recent = true
  owners      = "system"
}

data "alibabacloudstack_instance_types" "all" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
}

data "alibabacloudstack_instance_types" "any_n4" {
  availability_zone     = data.alibabacloudstack_zones.default.zones[0].id
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count   = 1
  memory_size      = 1
  instance_type_family = "ecs.n4"
  sorted_by        = "Memory"
}

locals {
  default_instance_type_id = try(
    element(sort(length(data.alibabacloudstack_instance_types.default.instance_types) > 0 ? data.alibabacloudstack_instance_types.default.ids : data.alibabacloudstack_instance_types.any_n4.ids), 0),
    sort(data.alibabacloudstack_instance_types.all.ids)[0]
  )
}

resource "alibabacloudstack_ecs_instance" "default" {
  image_id              = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type         = "${local.default_instance_type_id}"
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  system_disk_size     = 20
  system_disk_name     = "test_sys_disk"
  security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
  instance_name        = "${var.name}_ecs"
  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
  zone_id              = data.alibabacloudstack_zones.default.zones.0.id
  is_outdated          = false

  lifecycle {
    ignore_changes = [
      instance_type
    ]
  }
}

resource "alibabacloudstack_ecs_image" "default" {
  instance_id = "${alibabacloudstack_ecs_instance.default.id}"
  image_name  = "${var.name}"
  description = "Custom image created by Terraform"
  tags = {
    Environment = "Test"
    Owner      = "Terraform"
  }
}

resource "alibabacloudstack_image_share_permission" "default" {
  image_id   = "${alibabacloudstack_ecs_image.default.id}"
  account_id = "123456789"
}
```

## 参数参考

支持以下参数：

* `instance_id` - (选填, 变更时重建) - 用于创建自定义镜像的实例ID。
* `snapshot_id` - (选填, 变更时重建) - 用于创建自定义镜像的快照ID。与`instance_id`和`disk_device_mapping`冲突。
* `image_name` - (选填) - 镜像名称。它必须是2到128个字符长度，以字母或中文字符开头。它可以包含数字、冒号(:)、下划线(_)或连字符(-)。默认值：null。
* `description` - (选填) - 镜像描述。它必须是2到256个字符长度，并且不能以http://或https://开头。默认值：null。
* `tags` - (选填) - 分配给资源的标签映射。最多20对标签值。
* `disk_device_mapping` - (选填, 变更时重建) - 镜像下的系统盘和快照的描述。与`snapshot_id`和`instance_id`冲突。每个`disk_device_mapping`支持以下内容：
  * `size` - (选填, 变更时重建) - 指定组合自定义镜像中磁盘的大小，单位为GiB。取值范围：5到2000。
  * `snapshot_id` - (选填, 变更时重建) - 指定用于创建组合自定义镜像的快照。
* `force` - (选填) - 是否强制删除自定义镜像。默认为`false`。
  - `true`: 强制删除自定义镜像，无论该镜像当前是否被其他实例使用。
  - `false`: 在删除镜像之前验证该镜像是否未被任何其他实例使用。

### 超时时间

* `create` - (默认为10分钟)用于创建镜像(直到达到初始“可用”状态)。
* `delete` - (默认为10分钟)用于终止镜像。

## 属性说明

除了上述所有参数外，还导出以下属性：

* `id` - 镜像的ID。
* `image_name` - 镜像名称。
* `description` - 镜像描述。
* `disk_device_mapping` - 与镜像关联的磁盘设备映射。此属性包含以下信息：
  * `size` - 磁盘的大小（以GiB为单位）。
  * `snapshot_id` - 创建镜像时使用的快照ID。
  * `device` - 磁盘设备的挂载点。
  * `format` - 磁盘的文件系统格式。
  * `import_oss_bucket` - 如果镜像是通过OSS导入的，则显示对应的OSS存储桶名称。
  * `import_oss_object` - 如果镜像是通过OSS导入的，则显示对应的OSS对象名称。