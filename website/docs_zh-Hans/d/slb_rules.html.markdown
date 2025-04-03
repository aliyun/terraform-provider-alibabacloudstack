---
subcategory: "SLB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_rules"
sidebar_current: "docs-Alibabacloudstack-datasource-slb-rules"
description: |- 
  查询负载均衡(SLB)规则
---

# alibabacloudstack_slb_rules

根据指定过滤条件列出当前凭证权限可以访问的负载均衡(SLB)规则列表。

## 示例用法

```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

data "alibabacloudstack_instance_types" "all" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
}

data "alibabacloudstack_instance_types" "any_n4" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count       = 1
  memory_size          = 1
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

locals {
	default_instance_type_id = try(element(sort(length(data.alibabacloudstack_instance_types.default.instance_types) > 0 ? data.alibabacloudstack_instance_types.default.ids : data.alibabacloudstack_instance_types.any_n4.ids), 0), sort(data.alibabacloudstack_instance_types.all.ids)[0])
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_"
  most_recent = true
  owners      = "system"
}

variable "name" {
	default = "tf-testaccslbrulesdatasourcebasic"
}

resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
  name = "${var.name}"
}

resource "alibabacloudstack_security_group" "default" {
	name = "${var.name}"
	vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_slb_listener" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  backend_port = 80
  frontend_port = 80
  protocol = "http"
  sticky_session = "on"
  sticky_session_type = "insert"
  cookie = "${var.name}"
  cookie_timeout = 86400
  health_check = "on"
  health_check_uri = "/cons"
  health_check_connect_port = 20
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  bandwidth = 10
  x_forwarded_for  {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
}

resource "alibabacloudstack_instance" "default" {
  image_id = "${data.alibabacloudstack_images.default.images.0.id}"
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
  instance_type = "${local.default_instance_type_id}"
  system_disk_category = "cloud_efficiency"
  security_groups = ["${alibabacloudstack_security_group.default.id}"]
  instance_name = "${var.name}"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_slb_server_group" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  servers {
      server_ids = ["${alibabacloudstack_instance.default.id}"]
      port = 80
      weight = 100
    }
}

resource "alibabacloudstack_slb_rule" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  frontend_port = "${alibabacloudstack_slb_listener.default.frontend_port}"
  name = "${var.name}"
  domain = "*.aliyun.com"
  url = "/image"
  server_group_id = "${alibabacloudstack_slb_server_group.default.id}"
}

data "alibabacloudstack_slb_rules" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  frontend_port = 80
}

output "first_slb_rule_id" {
  value = "${data.alibabacloudstack_slb_rules.default.slb_rules.0.id}"
}
```

## 参数参考

以下参数是支持的：

* `load_balancer_id` - (必填) 负载均衡实例ID。
* `frontend_port` - (必填) SLB监听器的前端端口号。
* `ids` - (可选) 用于过滤结果的规则ID列表。
* `name_regex` - (可选，变更时重建) 用于按规则名称过滤结果的正则表达式字符串。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - SLB监听器规则ID列表。
* `names` - SLB监听器规则名称列表。
* `slb_rules` - SLB监听器规则列表。每个元素包含以下属性：
  * `id` - 规则的ID。
  * `name` - 规则的名称。
  * `domain` - 规则适用的HTTP请求中的域名(例如，`"*.aliyun.com"`)。
  * `url` - 规则适用的HTTP请求中的路径(例如，`"/image"`)。
  * `server_group_id` - 链接的VServer组的ID。