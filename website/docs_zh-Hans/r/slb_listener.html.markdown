---
subcategory: "SLB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_listener"
sidebar_current: "docs-Alibabacloudstack-slb-listener"
description: |- 
  编排负载均衡(SLB)监听器
---

# alibabacloudstack_slb_listener

使用Provider配置的凭证在指定的资源集编排负载均衡(SLB)监听器。

## 示例用法

```hcl
variable "name" {
    default = "tf-testaccslblistener84553"
}

resource "alibabacloudstack_slb_server_certificate" "default" {
	name = "${var.name}"
	server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgI+OuMs******XTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
	private_key = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0knDrlNdiys******ErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
}

resource "alibabacloudstack_slb_acl" "default" {
	name = "${var.name}"
	ip_version = "ipv4"
}

resource "alibabacloudstack_slb" "default" {
	name = "${var.name}"
	address_type       = "internet"
	specification        = "slb.s2.small"
}

resource "alibabacloudstack_slb_listener" "default" {
  load_balancer_id          = alibabacloudstack_slb.default.id
  frontend_port             = 80
  backend_port              = 80
  protocol                 = "http"
  bandwidth                = 10
  scheduler                = "wrr"
  sticky_session           = "off"
  health_check             = "off"

  x_forwarded_for {
    retrive_client_ip = true
    retrive_slb_id    = true
    retrive_slb_ip    = true
  }

  acl_status                = "on"
  acl_id                    = alibabacloudstack_slb_acl.default.id
  acl_type                  = "white"
}
```

## 参数参考

支持以下参数：

  * `load_balancer_id` - (必填, 变更时重建) 负载均衡实例的ID。
  * `frontend_port` - (必填, 变更时重建) Server Load Balancer 实例前端使用的端口。有效值范围：[1-65535]。
  * `backend_port` - (必填, 变更时重建) Server Load Balancer 实例后端使用的端口。有效值范围：[1-65535]。
  * `protocol` - (选填, 变更时重建) 监听协议。有效值为[`http`, `https`, `tcp`, `udp`]。
  * `bandwidth` - (必填) 监听器的带宽峰值。对于按流量计费的公网实例，监听器的带宽可以设置为 -1，表示带宽峰值不限制。有效值范围是 [-1, 1-5000] Mbps。
  * `scheduler` - (选填) 调度算法。有效值为 `wrr`, `rr`, 和 `wlc`。默认值为 "wrr"。
  * `server_group_id` - (选填) 绑定的服务器组ID。
  * `master_slave_server_group_id` - (选填) 主备服务器组的ID。注意：不能同时设置 `server_group_id` 和 `master_slave_server_group_id`。
  * `sticky_session` - (必填) 是否启用会话保持。有效值为 `on` 和 `off`。
  * `sticky_session_type` - (选填) 处理Cookie的模式。如果 `sticky_session` 是 "on"，这是必填项。否则将被忽略。有效值为 `insert` 和 `server`。
  * `cookie_timeout` - (选填) Cookie超时时间。当 `sticky_session` 是 "on" 且 `sticky_session_type` 是 "insert" 时，这是必填项。否则将被忽略。有效值范围：[1-86400] 秒。
  * `cookie` - (选填) 在服务器上配置的Cookie。当 `sticky_session` 是 "on" 且 `sticky_session_type` 是 "server" 时，这是必填项。否则将被忽略。有效值：符合RFC 2965的字符串，长度为1-200。只包含ASCII码、英文字母和数字，而不是逗号、分号或空格，并且不能以 `$` 开头。
  * `persistence_timeout` - (选填) 连接持久化超时时间。有效值范围：[0-3600] 秒。默认值为 0 表示关闭。
  * `health_check` - (必填) 是否启用健康检查。有效值为 `on` 和 `off`。TCP和UDP监听器的健康检查始终处于开启状态，因此在启动TCP或UDP监听器时将被忽略。
  * `health_check_method` - (选填) 用于健康检查的HealthCheckMethod。HTTP和HTTPS支持区域 ap-northeast-1, ap-southeast-1, ap-southeast-2, ap-southeast-3, us-east-1, us-west-1, eu-central-1, ap-south-1, me-east-1, cn-huhehaote, cn-zhangjiakou, ap-southeast-5, cn-shenzhen, cn-hongkong, cn-qingdao, cn-chengdu, eu-west-1, cn-hangzhou, cn-beijing, cn-shanghai。此功能不支持TCP协议。
  * `health_check_type` - (选填) 健康检查类型。有效值为：`tcp` 和 `http`。默认值为 `tcp`。
  * `health_check_domain` - (选填) 用于健康检查的域名。当用于启动TCP监听器时，`health_check_type` 必须为 "http"。其长度限制为1-80个字符，只允许字母、数字、`-` 和 `.`。如果不设置或为空，则Server Load Balancer使用每个后端服务器的私网IP作为用于健康检查的Domain。
  * `health_check_uri` - (选填) 用于健康检查的URI。当用于启动TCP监听器时，`health_check_type` 必须为 "http"。其长度限制为1-80个字符，必须以 `/` 开头，只允许字母、数字、`-`、`/`、`.`、`%`、`?`、`#` 和 `&`。
  * `health_check_connect_port` - (选填) 用于健康检查的端口。有效值范围：[1-65535]。默认值为 "None"，表示使用后端服务器端口。
  * `healthy_threshold` - (选填) 确定健康检查结果成功的阈值。当 `health_check` 为 on 时，这是必填项。有效值范围：[1-10] 秒。默认值为 3。
  * `unhealthy_threshold` - (选填) 确定健康检查结果失败的阈值。当 `health_check` 为 on 时，这是必填项。有效值范围：[1-10] 秒。默认值为 3。
  * `health_check_timeout` - (选填) 每次健康检查响应的最大超时时间。当 `health_check` 为 on 时，这是必填项。有效值范围：[1-300] 秒。默认值为 5。注意：如果 `health_check_timeout` < `health_check_interval`，则将其替换为 `health_check_interval`。
  * `health_check_interval` - (选填) 健康检查的时间间隔。当 `health_check` 为 on 时，这是必填项。有效值范围：[1-50] 秒。默认值为 2。
  * `health_check_http_code` - (选填) 常规健康检查HTTP状态码。多个代码用 `,` 分隔。当 `health_check` 为 on 时，这是必填项。默认值为 `http_2xx`。有效值为：`http_2xx`, `http_3xx`, `http_4xx`, 和 `http_5xx`。
  * `server_certificate_id` - (选填) SLB服务器证书ID。当 `protocol` 为 `https` 时，这是必填项。
  * `ca_certificate_id` - (选填) 认证机构(CA)证书的ID。
  * `gzip` - (选填) 是否启用"Gzip压缩"。如果启用，特定类型的文件将被压缩，否则没有任何文件会被压缩。默认值为 true。
  * `x_forwarded_for` - (选填) 是否设置额外的HTTP Header字段 "X-Forwarded-For"(如下所示)。
    * `retrive_client_ip` - (选填) 是否检索客户端IP。
    * `retrive_slb_id` - (选填) 是否使用XForwardedFor头部获取SLB实例的ID。
    * `retrive_slb_ip` - (选填) 是否使用XForwardedFor_SLBIP头部获取SLB实例的公网IP地址。
    * `retrive_slb_proto` - (选填) 是否使用XForwardedFor_proto头部获取监听器使用的协议。
  * `established_timeout` - (选填) TCP监听器建立连接的空闲超时时间。有效值范围：[10-900] 秒。默认值为 900。
  * `acl_status` - (选填) 是否启用访问控制。默认值：off。有效值：on, off。
  * `acl_type` - (选填) 网络ACL的类型。有效值：black, white。注意：如果 `acl_status` 设置为 on，`acl_type` 是必填项。否则将被忽略。
  * `acl_id` - (选填) 与监听器关联的网络ACL的ID。注意：如果 `acl_status` 设置为 "on"，`acl_id` 是必填项。否则将被忽略。
  * `listener_forward` - (选填, 变更时重建) 是否启用HTTP重定向到HTTPS。有效值为 `on` 和 `off`。默认值为 `off`。
  * `enable_http2` - (选填) 是否启用HTTP/2。默认值：on。有效值："on", "off"。
  * `forward_port` - (选填, 变更时重建) HTTP重定向到HTTPS的端口。
  * `tls_cipher_policy` - (选填) 传输层安全(TLS)加密策略。默认值：tls_cipher_policy_1_0。有效值：tls_cipher_policy_1_0, tls_cipher_policy_1_1, tls_cipher_policy_1_2, tls_cipher_policy_1_2_strict。
  * `delete_protection_validation` - (选填) 删除前检查SLB实例的DeleteProtection。如果为true，当SLB实例启用了DeleteProtection时，此资源不会被删除。默认值为 false。
  * `logs_download_attributes` - (可选) 用于定义 SLS 日志的映射。
    * `log_project`: (必填) SLS 日志服务器项目名称。
    * `log_store`: (必填) SLS 日志服务器日志存储名称。

## 属性参考

除了上述所有参数外，还导出了以下属性：

  * `id` - 负载均衡监听器的ID。格式为 `<load_balancer_id>:<protocol>:<frontend_port>`。
  * `description` - TLSCipherPolicy参数。
  * `health_check_method` - 健康检查使用的HealthCheckMethod。
  * `health_check_connect_port` - 健康检查使用的端口。
  * `ssl_certificate_id` - SSL证书ID。
  * `x_forwarded_for` - 是否设置额外的HTTP Header字段 "X-Forwarded-For"。
  * `listener_forward` - 是否启用HTTP重定向到HTTPS。