---
subcategory: "SLB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_rule"
sidebar_current: "docs-Alibabacloudstack-slb-rule"
description: |- 
  编排负载均衡(SLB)规则
---

# alibabacloudstack_slb_rule

使用Provider配置的凭证在指定的资源集编排负载均衡(SLB)规则。

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
  default = "tf-testAccSlbRuleBasic"
}

resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
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

resource "alibabacloudstack_instance" "default" {
  image_id = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type = "${local.default_instance_type_id}"
  security_groups = "${alibabacloudstack_security_group.default.*.id}"
  internet_max_bandwidth_out = "10"
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
  system_disk_category = "cloud_sperf"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
  instance_name = "${var.name}"
}

resource "alibabacloudstack_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_slb_listener" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  backend_port = 22
  frontend_port = 22
  protocol = "http"
  bandwidth = 5
  health_check_connect_port = "20"
  health_check = "on"
  sticky_session = "off"
}

resource "alibabacloudstack_slb_server_group" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  servers {
    server_ids = "${alibabacloudstack_instance.default.*.id}"
    port = 80
    weight = 100
  }
}

resource "alibabacloudstack_slb_rule" "default" {
  sticky_session_type = "server"
  frontend_port = "${alibabacloudstack_slb_listener.default.frontend_port}"
  health_check_domain = "test"
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  health_check = "on"
  health_check_uri = "/test"
  cookie_timeout = "100"
  health_check_http_code = "http_2xx"
  health_check_interval = "10"
  listener_sync = "on"
  cookie = "23ffsa"
  healthy_threshold = "3"
  health_check_connect_port = "80"
  domain = "*.aliyun.com"
  name = "${var.name}"
  unhealthy_threshold = "3"
  scheduler = "rr"
  sticky_session = "on"
  health_check_timeout = "10"
  url = "/image"
  server_group_id = "${alibabacloudstack_slb_server_group.default.id}"
}
```

## 参数参考

支持以下参数：

* `load_balancer_id` - (必填, 变更时重建) 负载均衡实例ID。
* `frontend_port` - (必填, 变更时重建) 监听器的前端端口，用于启动新的转发规则。有效范围：[1-65535]。
* `name` - (选填, 变更时重建) 转发规则的名称。
* `rule_name` - (选填, 变更时重建) 转发规则的名称。
* `listener_sync` - (选填) 转发规则是否从监听上继承健康检查、会话保持和调度算法配置。取值：**on | off**。  
  * **off**：不继承监听配置，转发规则自定义健康检查及会话保持配置。
  * **on**(默认值)：继承监听配置。
* `scheduler` - (选填) 调度算法。取值：**wrr**(默认值)、**wlc**、**rr**。此参数在 `listener_sync` 设置为 `off` 时是必填且生效的。
* `domain` - (选填, 变更时重建) 转发规则的域名。它可以包含字母a-z，数字0-9，连字符(`-`)，点(`.`)，和通配符。支持以下两种格式的域名：
  * 标准域名：`www.test.com`
  * 通配符域名：`*.test.com`。通配符(`*`)必须是格式中的第一个字符(`*.`)。
* `url` - (选填, 变更时重建) 转发规则的URL。它必须是2-80个字符长度。仅允许字母a-z，数字0-9，以及字符如 `'-'`, `'/'`, `'?'`, `'%'`, `'#'`, 和 `'&'`。URL必须以字符 `'/'` 开头，但不能仅为 `'/'`。
* `server_group_id` - (必填, 变更时重建) 将要转发的虚拟服务器组ID。
* `cookie` - (选填) 在服务器上配置的Cookie。当 `sticky_session` 是 `"on"` 且 `sticky_session_type` 是 `"server"` 时是必填项。否则将被忽略。有效值：符合RFC 2965的字符串，长度为1-200。它只能包含ASCII码、英文字母和数字，而不是逗号、分号或空格，并且不能以`$`开头。
* `cookie_timeout` - (选填) Cookie超时时间。当 `sticky_session` 是 `"on"` 且 `sticky_session_type` 是 `"insert"` 时是必填项。否则将被忽略。有效范围：[1-86400]秒。
* `health_check` - (选填) 是否发起健康检查。有效值：`on` 和 `off`。TCP和UDP监听器的健康检查始终处于开启状态，因此在启动TCP或UDP监听器时将被忽略。此参数在 `listener_sync` 设置为 `off` 时是必填且生效的。
* `health_check_http_code` - (选填) 表示健康检查成功的HTTP状态码。多个HTTP状态码用逗号(`,`)分隔。默认值：`http_2xx`。有效值：`http_2xx`, `http_3xx`, `http_4xx`, 和 `http_5xx`。此参数在 `health_check` 设置为 `on` 时是必填的。
* `health_check_interval` - (选填) 两次连续健康检查之间的时间间隔。此参数在 `health_check` 设置为 `on` 时是必填的。有效范围：[1-50]秒。默认值：`2`。
* `health_check_domain` - (选填) 健康检查使用的域名。当未设置或为空时，负载均衡器使用每个后端服务器的私网IP地址作为健康检查使用的域名。此参数在 `health_check` 设置为 `on` 时是必填的。有效范围：[1-80]字符。只允许字母、数字、`'-'` 和 `'.'`。
* `health_check_uri` - (选填) 健康检查使用的URI。此参数在 `health_check` 设置为 `on` 时是必填的。有效范围：[1-80]字符。它必须以`'/'`开头。只允许字母、数字、`'-'`, `'/'`, `'.'`, `'%'`, `'?'`, `'#'`, 和 `'&'`。
* `health_check_connect_port` - (选填) 健康检查使用的端口。有效范围：[1-65535]。默认值：`None`表示使用后端服务器端口。
* `health_check_timeout` - (选填) 每次健康检查响应的最大超时时间。此参数在 `health_check` 设置为 `on` 时是必填的。有效范围：[1-300]秒。默认值：`5`。注意：如果 `health_check_timeout` < `health_check_interval`，其将被替换为 `health_check_interval`。
* `healthy_threshold` - (选填) 确定健康检查结果为成功时的阈值。此参数在 `health_check` 设置为 `on` 时是必填的。有效范围：[1-10]秒。默认值：`3`。
* `unhealthy_threshold` - (选填) 确定健康检查结果为失败时的阈值。此参数在 `health_check` 设置为 `on` 时是必填的。有效范围：[1-10]秒。默认值：`3`。
* `sticky_session` - (选填) 是否启用会话持久性。有效值：`on` 和 `off`。默认值：`off`。此参数在 `listener_sync` 设置为 `off` 时是必填且生效的。
* `sticky_session_type` - (选填) 处理Cookie的模式。如果 `sticky_session` 是 `"on"`，它是必填项。否则将被忽略。有效值：`insert` 和 `server`。`insert` 表示从负载均衡器插入；`server` 表示负载均衡器从后端服务器学习。
* `delete_protection_validation` - (选填) 删除前检查SLB实例的删除保护。如果为true，则当其SLB实例启用了DeleteProtection时，此资源不会被删除。默认值：`false`。

## 属性参考

除了上述参数外，还导出了以下属性：

* `name` - 转发规则的名称。
* `rule_name` - 转发规则的名称。
* `health_check_connect_port` - 后端服务器的健康检查端口。有效范围：[1-65535]。若为空且 `HealthCheck` 为 `on`，表明默认使用监听后端端口配置。
