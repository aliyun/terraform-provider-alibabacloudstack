---
subcategory: "Server Load Balancer (SLB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_backendservers"
sidebar_current: "docs-Alibabacloudstack-datasource-slb-backendservers"
description: |- 
  查询负载均衡(SLB)后端服务
---

# alibabacloudstack_slb_backendservers
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_slb_backend_servers`

根据指定过滤条件列出当前凭证权限可以访问的负载均衡(SLB)后端服务列表。

## 示例用法

以下是一个完整的示例，展示了如何使用 `alibabacloudstack_slb_backendservers` 数据源来获取负载均衡实例的后端服务器信息：

```hcl
# 获取可用区信息
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

# 创建VPC
resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

# 创建VSwitch
resource "alibabacloudstack_vpc_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

# 创建安全组
resource "alibabacloudstack_ecs_securitygroup" "default" {
  name   = "${var.name}_sg"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

# 添加安全组规则
resource "alibabacloudstack_security_group_rule" "default" {
  type = "ingress"
  ip_protocol = "tcp"
  nic_type = "intranet"
  policy = "accept"
  port_range = "22/22"
  priority = 1
  security_group_id = "${alibabacloudstack_ecs_securitygroup.default.id}"
  cidr_ip = "172.16.0.0/24"
}

# 获取镜像信息
data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_"
  most_recent = true
  owners      = "system"
}

# 获取实例类型信息
data "alibabacloudstack_instance_types" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count   = 1
  memory_size      = 1
  instance_type_family = "ecs.n4"
  sorted_by        = "Memory"
}

# 创建ECS实例
resource "alibabacloudstack_ecs_instance" "default" {
  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type        = "${data.alibabacloudstack_instance_types.default.instance_types.0.id}"
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  system_disk_size     = 20
  security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
  instance_name        = "${var.name}_ecs"
  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
  zone_id             = data.alibabacloudstack_zones.default.zones.0.id
}

# 创建SLB实例
resource "alibabacloudstack_slb" "default" {
  name       = "${var.name}"
  vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
}

# 将ECS实例添加到SLB后端
resource "alibabacloudstack_slb_backend_server" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"

  backend_servers {
    server_id = "${alibabacloudstack_ecs_instance.default.id}"
    weight    = 100
  }
}

# 使用数据源获取SLB后端服务器信息
data "alibabacloudstack_slb_backend_servers" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
}

# 输出第一个后端服务器的ID
output "first_slb_backend_server_id" {
  value = "${data.alibabacloudstack_slb_backend_servers.default.backend_servers.0.id}"
}
```

## 参数参考

以下参数是支持的：

* `load_balancer_id` - (必填) 传统型负载均衡实例的 ID。
* `ids` - (可选) 作为后端服务器附加的 ECS 实例的 ID 列表。如果提供了该参数，则只返回匹配这些 ID 的后端服务器。

## 属性参考

除了上述参数外，还导出以下属性：

* `backend_servers` - 每个后端服务器的信息。每个元素包含以下字段：
  * `id` - 后端服务器的唯一标识符(ECS 实例的 ID)。
  * `weight` - 分配给后端服务器的权重，这决定了它接收的流量比例。
