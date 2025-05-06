---
subcategory: "Server Load Balancer (SLB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_backend_server"
sidebar_current: "docs-Alibabacloudstack-slb-backend-server"
description: |- 
  编排负载均衡(SLB)后端服务
---

# alibabacloudstack_slb_backend_server
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_slb_backendserver`

使用Provider配置的凭证在指定的资源集编排负载均衡(SLB)后端服务。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAccSlbBackendServersVpc2243280"
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
  name       = "${var.name}_vsw"
  vpc_id     = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  zone_id    = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_ecs_securitygroup" "default" {
  name   = "${var.name}_sg"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

resource "alibabacloudstack_security_group_rule" "default" {
  type                = "ingress"
  ip_protocol         = "tcp"
  nic_type            = "intranet"
  policy              = "accept"
  port_range          = "22/22"
  priority            = 1
  security_group_id   = "${alibabacloudstack_ecs_securitygroup.default.id}"
  cidr_ip             = "172.16.0.0/24"
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
  availability_zone     = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count       = 1
  memory_size          = 1
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

locals {
  default_instance_type_id = try(element(sort(length(data.alibabacloudstack_instance_types.default.instance_types) > 0 ? data.alibabacloudstack_instance_types.default.ids : data.alibabacloudstack_instance_types.any_n4.ids), 0), sort(data.alibabacloudstack_instance_types.all.ids)[0])
}

resource "alibabacloudstack_ecs_instance" "default" {
  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type        = "${local.default_instance_type_id}"
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

resource "alibabacloudstack_slb" "default" {
  name       = "${var.name}"
  vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
}

data "alibabacloudstack_instance_types" "new" {
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  eni_amount       = 2
}

resource "alibabacloudstack_network_interface" "default" {
  count           = 1
  name            = "${var.name}"
  vswitch_id      = "${alibabacloudstack_vpc_vswitch.default.id}"
  security_groups = [alibabacloudstack_ecs_securitygroup.default.id]
}

resource "alibabacloudstack_instance" "new" {
  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type        = "${data.alibabacloudstack_instance_types.new.instance_types[0].id}"
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  system_disk_size     = 40
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

resource "alibabacloudstack_network_interface_attachment" "default" {
  count               = 1
  instance_id         = "${alibabacloudstack_instance.new.id}"
  network_interface_id = "${element(alibabacloudstack_network_interface.default.*.id, count.index)}"
}

resource "alibabacloudstack_slb_backend_server" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"

  backend_servers {
    server_id = "${alibabacloudstack_ecs_instance.default.id}"
    weight    = "80"
  }

  backend_servers {
    server_id = "${alibabacloudstack_instance.new.id}"
    weight    = "100"
  }
}
```

## 参数说明

支持以下参数：

* `load_balancer_id` - (必填，强制更新)传统服务器负载均衡实例的ID。
* `backend_servers` - (可选) 要添加到SLB的后端服务器列表。每个`backend_servers`块支持以下参数：
  * `server_id` - (必填) 后端服务器的ID(ECS实例或ENI实例)。
  * `weight` - (可选) 后端服务器的权重。有效值范围为`0`到`100`。值为`0`表示该后端服务器已禁用。
* `delete_protection_validation` - (可选) 在删除此资源之前是否检查SLB实例的`DeleteProtection`属性。如果设置为`true`，当SLB实例启用了`DeleteProtection`时，资源将不会被删除。默认值为`false`。
* `force_new_validation` - (可选) 指定当某些属性发生变化时是否重新创建负载均衡器。默认值为`true`。

## 属性说明

除了上述所有参数外，还导出以下属性：

* `id` - 资源的ID，与`load_balancer_id`相同。