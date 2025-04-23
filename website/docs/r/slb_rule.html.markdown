---
subcategory: "Server Load Balancer (SLB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_rule"
sidebar_current: "docs-Alibabacloudstack-slb-rule"
description: |- 
  Provides a slb Rule resource.
---

# alibabacloudstack_slb_rule

Provides a slb Rule resource.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, ForceNew) Load balancing instance ID.
* `frontend_port` - (Required, ForceNew) The listener frontend port which is used to launch the new forwarding rule. Valid range: [1-65535].
* `name` - (Optional, ForceNew) Name of the forwarding rule.
* `rule_name` - (Optional, ForceNew) Name of the forwarding rule.
* `listener_sync` - (Optional) Indicates whether the forwarding rule inherits the health check, session persistence, and scheduling algorithm configuration from listening. Default value: `on`.
* `scheduler` - (Optional) Scheduling algorithm. Valid values: `wrr`, `rr`, `wlc`. Default value: `wrr`. This parameter is required and takes effect only when `listener_sync` is set to `off`.
* `domain` - (Optional, ForceNew) Domain name of the forwarding rule. It can contain letters a-z, numbers 0-9, hyphens (-), and periods (.), and wildcard characters. The following two domain name formats are supported:
  * Standard domain name: `www.test.com`
  * Wildcard domain name: `*.test.com`. Wildcard (`*`) must be the first character in the format of (`*.`).
* `url` - (Optional, ForceNew) URL of the forwarding rule. It must be 2-80 characters in length. Only letters a-z, numbers 0-9, and characters such as `'-'`, `'/'`, `'?'`, `'%'`, `'#'`, and `'&'` are allowed. URLs must start with the character `'/'`, but cannot be `'/'` alone.
* `server_group_id` - (Required, ForceNew) ID of a virtual server group that will be forwarded.
* `cookie` - (Optional) The cookie configured on the server. It is mandatory when `sticky_session` is `"on"` and `sticky_session_type` is `"server"`. Otherwise, it will be ignored. Valid value: String in line with RFC 2965, with length being 1-200. It only contains characters such as ASCII codes, English letters, and digits instead of the comma, semicolon, or spacing, and it cannot start with `$`.
* `cookie_timeout` - (Optional) Cookie timeout. It is mandatory when `sticky_session` is `"on"` and `sticky_session_type` is `"insert"`. Otherwise, it will be ignored. Valid value range: [1-86400] in seconds.
* `health_check` - (Optional) Whether to initiate a health check. Valid values: `on` and `off`. TCP and UDP listener's HealthCheck is always on, so it will be ignored when launching TCP or UDP listener. This parameter is required and takes effect only when `listener_sync` is set to `off`.
* `health_check_http_code` - (Optional) The HTTP status code that indicates a successful health check. Separate multiple HTTP status codes with commas (`,`). Default value: `http_2xx`. Valid values: `http_2xx`, `http_3xx`, `http_4xx`, and `http_5xx`. This parameter is required when `health_check` is set to `on`.
* `health_check_interval` - (Optional) The time interval between two consecutive health checks. This parameter is required when `health_check` is set to `on`. Valid value range: [1-50] in seconds. Default value: `2`.
* `health_check_domain` - (Optional) Domain name used for health check. When it is not set or empty, Server Load Balancer uses the private network IP address of each backend server as Domain used for health check. This parameter is required when `health_check` is set to `on`. Valid value range: [1-80] characters. Only letters, digits, `'-'`, and `'.'` are allowed.
* `health_check_uri` - (Optional) URI used for health check. This parameter is required when `health_check` is set to `on`. Valid value range: [1-80] characters. It must start with `'/'`. Only letters, digits, `'-'`, `'/'`, `'.'`, `'%'`, `'?'`, `'#'`, and `'&'` are allowed.
* `health_check_connect_port` - (Optional) Port used for health check. Valid value range: [1-65535]. Default value: `None` means the backend server port is used.
* `health_check_timeout` - (Optional) Maximum timeout of each health check response. This parameter is required when `health_check` is set to `on`. Valid value range: [1-300] in seconds. Default value: `5`. Note: If `health_check_timeout` < `health_check_interval`, its will be replaced by `health_check_interval`.
* `healthy_threshold` - (Optional) Threshold determining the result of the health check is success. This parameter is required when `health_check` is set to `on`. Valid value range: [1-10] in seconds. Default value: `3`.
* `unhealthy_threshold` - (Optional) Threshold determining the result of the health check is fail. This parameter is required when `health_check` is set to `on`. Valid value range: [1-10] in seconds. Default value: `3`.
* `sticky_session` - (Optional) Whether to enable session persistence. Valid values: `on` and `off`. Default value: `off`. This parameter is required and takes effect only when `listener_sync` is set to `off`.
* `sticky_session_type` - (Optional) Mode for handling the cookie. If `sticky_session` is `"on"`, it is mandatory. Otherwise, it will be ignored. Valid values: `insert` and `server`. `insert` means it is inserted from Server Load Balancer; `server` means the Server Load Balancer learns from the backend server.
* `delete_protection_validation` - (Optional) Checking DeleteProtection of SLB instance before deleting. If true, this resource will not be deleted when its SLB instance enabled DeleteProtection. Default value: `false`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `name` - Name of the forwarding rule.
* `rule_name` - Name of the forwarding rule.
* `health_check_connect_port` - Health check the port of the backend server.