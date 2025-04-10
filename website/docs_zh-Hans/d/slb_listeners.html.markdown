---
subcategory: "SLB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_listeners"
sidebar_current: "docs-Alibabacloudstack-datasource-slb-listeners"
description: |- 
  查询负载均衡(SLB)监听器
---

# alibabacloudstack_slb_listeners

根据指定过滤条件列出当前凭证权限可以访问的负载均衡(SLB)监听器列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAccCheckSlbListenersDataSourceHttp"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_slb" "default" {
  name = "${var.name}"
}

resource "alibabacloudstack_slb_listener" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  backend_port     = 80
  frontend_port    = 80
  protocol         = "http"
  sticky_session   = "on"
  sticky_session_type = "insert"
  cookie           = "${var.name}"
  cookie_timeout   = 86400
  health_check     = "on"
  health_check_uri = "/cons"
  health_check_connect_port = 20
  healthy_threshold = 8
  unhealthy_threshold = 8
  health_check_timeout = 8
  health_check_interval = 5
  health_check_http_code = "http_2xx,http_3xx"
  bandwidth        = 10
  x_forwarded_for {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  description      = "${var.name}"
}

data "alibabacloudstack_slb_listeners" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  protocol         = "http"
  frontend_port    = 80
  description_regex = "^${var.name}.*"
}

output "first_slb_listener_protocol" {
  value = "${data.alibabacloudstack_slb_listeners.default.slb_listeners.0.protocol}"
}
```

## 参数说明

以下参数是支持的：

* `load_balancer_id` - (必填) 负载均衡实例的ID。
* `protocol` - (可选) 按指定协议过滤监听器。有效值：`http`，`https`，`tcp` 和 `udp`。
* `frontend_port` - (可选) 按指定前端端口过滤监听器。
* `description_regex` - (可选) 用于通过SLB监听器描述筛选结果的正则表达式字符串。

## 属性说明

除了上述参数外，还导出以下属性：

* `slb_listeners` - SLB监听器的列表。每个元素包含以下属性：
  * `frontend_port` - 用于接收传入流量并将其分发到后端服务器的前端端口。
  * `backend_port` - 打开在后端服务器上以接收请求的端口。
  * `protocol` - 监听器协议。可能的值：`http`，`https`，`tcp` 和 `udp`。
  * `status` - 监听器状态。
  * `security_status` - 安全状态。仅当协议为`https`时可用。
  * `bandwidth` - 峰值带宽。如果值设置为-1，则监听器不受带宽限制。
  * `scheduler` - 用于分配流量的算法。可能的值：`wrr`(加权轮询)，`wlc`(加权最少连接)，和 `rr`(轮询)。
  * `server_group_id` - 链接的VServer组的ID。
  * `server_certificate_id` - 服务器证书的ID。
  * `master_slave_server_group_id` - 主备服务器组的ID。
  * `persistence_timeout` - TCP连接的超时值(以秒为单位)。如果值为0，则禁用会话保持功能。仅当协议为`tcp`时可用。
  * `established_timeout` - 第4层TCP监听器的连接超时时间(以秒为单位)。仅当协议为`tcp`时可用。
  * `sticky_session` - 指示是否启用了会话保持。如果启用，则来自同一客户端的所有会话请求都发送到同一后端服务器。可能的值是`on`和`off`。仅当协议为`http`或`https`时可用。
  * `sticky_session_type` - 处理Cookie的方法。可能的值是`insert`(将Cookie添加到响应中)和`server`(由后端服务器设置的Cookie)。仅当协议为`http`或`https`且`sticky_session`为`on`时可用。
  * `cookie_timeout` - Cookie的超时时间(以秒为单位)。仅当`sticky_session_type`为`insert`时可用。
  * `cookie` - 后端服务器配置的Cookie。仅当`sticky_session_type`为`server`时可用。
  * `health_check` - 指示是否启用了健康检查。可能的值是`on`和`off`。
  * `health_check_type` - 健康检查方法。可能的值是`tcp`和`http`。仅当协议为`tcp`时可用。
  * `health_check_domain` - 用于健康检查的域名。SLB向后端服务器发送HTTP头请求，域在请求验证主机字段时有用。仅当协议为`http`，`https`或`tcp`(在这种情况下`health_check_type`必须为`http`)时可用。
  * `health_check_uri` - 用于健康检查的URI。仅当协议为`http`，`https`或`tcp`(在这种情况下`health_check_type`必须为`http`)时可用。
  * `health_check_connect_port` - 用于健康检查的端口。
  * `health_check_connect_timeout` - 等待健康检查响应的时间量(以秒为单位)。
  * `healthy_threshold` - 对同一ECS实例执行的健康检查连续成功的次数(从失败到成功)。
  * `unhealthy_threshold` - 对同一ECS实例执行的健康检查连续失败的次数(从成功到失败)。
  * `health_check_timeout` - 等待健康检查响应的时间量(以秒为单位)。如果ECS实例在指定的超时期间内未发送响应，则健康检查失败。仅当协议为`http`或`https`时可用。
  * `health_check_interval` - 两次连续健康检查之间的时间间隔。
  * `health_check_http_code` - 表示健康检查正常的HTTP状态码。它可以包含几个逗号分隔的值，例如"http_2xx,http_3xx"。仅当协议为`http`，`https`或`tcp`(在这种情况下`health_check_type`必须为`http`)时可用。
  * `gzip` - 指示是否启用了Gzip压缩。可能的值是`on`和`off`。仅当协议为`http`或`https`时可用。
  * `x_forwarded_for` - 指示是否添加了HTTP头字段"X-Forwarded-For";它允许后端服务器知道用户的IP地址。可能的值是`on`和`off`。仅当协议为`http`或`https`时可用。
  * `x_forwarded_for_slb_ip` - 指示是否添加了HTTP头字段"X-Forwarded-For_SLBIP";它允许后端服务器知道SLB IP地址。可能的值是`on`和`off`。仅当协议为`http`或`https`时可用。
  * `x_forwarded_for_slb_id` - 指示是否添加了HTTP头字段"X-Forwarded-For_SLBID";它允许后端服务器知道SLB ID。可能的值是`on`和`off`。仅当协议为`http`或`https`时可用。
  * `x_forwarded_for_slb_proto` - 指示是否添加了HTTP头字段"X-Forwarded-For_proto";它允许后端服务器知道用户的协议。可能的值是`on`和`off`。仅当协议为`http`或`https`时可用。
  * `description` - SLB监听器的描述。