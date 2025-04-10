---
subcategory: "SLB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_listener"
sidebar_current: "docs-Alibabacloudstack-slb-listener"
description: |- 
  Provides a slb Listener resource.
---

# alibabacloudstack_slb_listener

Provides a slb Listener resource.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, ForceNew) ID of the load balancing instance.
* `frontend_port` - (Required, ForceNew) Port used by the Server Load Balancer instance frontend. Valid value range: [1-65535].
* `backend_port` - (Required, ForceNew) Port used by the Server Load Balancer instance backend. Valid value range: [1-65535].
* `protocol` - (Optional, ForceNew) The protocol to listen on. Valid values are [`http`, `https`, `tcp`, `udp`].
* `bandwidth` - (Required) Bandwidth peak of Listener. For the public network instance charged per traffic consumed, the Bandwidth on Listener can be set to -1, indicating the bandwidth peak is unlimited. Valid values are [-1, 1-5000] in Mbps.
* `scheduler` - (Optional) Scheduling algorithm. Valid values are `wrr`, `rr`, and `wlc`. Default to "wrr".
* `server_group_id` - (Optional) The bound server group ID.
* `master_slave_server_group_id` - (Optional) The ID of the primary/secondary server group. NOTE: You cannot set both 'server_group_id' and 'master_slave_server_group_id'.
* `sticky_session` - (Required) Whether to enable session persistence. Valid values are `on` and `off`.
* `sticky_session_type` - (Optional) Mode for handling the cookie. If `sticky_session` is "on", it is mandatory. Otherwise, it will be ignored. Valid values are `insert` and `server`.
* `cookie_timeout` - (Optional) Cookie timeout. It is mandatory when `sticky_session` is "on" and `sticky_session_type` is "insert". Otherwise, it will be ignored. Valid value range: [1-86400] in seconds.
* `cookie` - (Optional) The cookie configured on the server. It is mandatory when `sticky_session` is "on" and `sticky_session_type` is "server". Otherwise, it will be ignored. Valid value: String in line with RFC 2965, with length being 1-200. It only contains characters such as ASCII codes, English letters, and digits instead of the comma, semicolon, or spacing, and it cannot start with `$`.
* `persistence_timeout` - (Optional) Timeout of connection persistence. Valid value range: [0-3600] in seconds. Default to 0 and means closing it.
* `health_check` - (Required) Whether to enable health check. Valid values are `on` and `off`. TCP and UDP listener's HealthCheck is always on, so it will be ignored when launching TCP or UDP listeners.
* `health_check_method` - (Optional) HealthCheckMethod used for health check. HTTP and HTTPS support regions ap-northeast-1, ap-southeast-1, ap-southeast-2, ap-southeast-3, us-east-1, us-west-1, eu-central-1, ap-south-1, me-east-1, cn-huhehaote, cn-zhangjiakou, ap-southeast-5, cn-shenzhen, cn-hongkong, cn-qingdao, cn-chengdu, eu-west-1, cn-hangzhou, cn-beijing, cn-shanghai. This function does not support the TCP protocol.
* `health_check_type` - (Optional) Type of health check. Valid values are: `tcp` and `http`. Default to `tcp`.
* `health_check_domain` - (Optional) Domain name used for health check. When it used to launch TCP listener, `health_check_type` must be "http". Its length is limited to 1-80 and only characters such as letters, digits, `-`, and `.` are allowed. When it is not set or empty, Server Load Balancer uses the private network IP address of each backend server as Domain used for health check.
* `health_check_uri` - (Optional) URI used for health check. When it used to launch TCP listener, `health_check_type` must be "http". Its length is limited to 1-80 and it must start with `/`. Only characters such as letters, digits, `-`, `/`, `.`, `%`, `?`, `#`, and `&` are allowed.
* `health_check_connect_port` - (Optional) Port used for health check. Valid value range: [1-65535]. Default to "None" means the backend server port is used.
* `healthy_threshold` - (Optional) Threshold determining the result of the health check is success. It is required when `health_check` is on. Valid value range: [1-10] in seconds. Default to 3.
* `unhealthy_threshold` - (Optional) Threshold determining the result of the health check is fail. It is required when `health_check` is on. Valid value range: [1-10] in seconds. Default to 3.
* `health_check_timeout` - (Optional) Maximum timeout of each health check response. It is required when `health_check` is on. Valid value range: [1-300] in seconds. Default to 5. Note: If `health_check_timeout` < `health_check_interval`, its will be replaced by `health_check_interval`.
* `health_check_interval` - (Optional) Time interval of health checks. It is required when `health_check` is on. Valid value range: [1-50] in seconds. Default to 2.
* `health_check_http_code` - (Optional) Regular health check HTTP status code. Multiple codes are segmented by `,`. It is required when `health_check` is on. Default to `http_2xx`. Valid values are: `http_2xx`, `http_3xx`, `http_4xx`, and `http_5xx`.
* `server_certificate_id` - (Optional) SLB Server certificate ID. It is required when `protocol` is `https`.
* `ca_certificate_id` - (Optional) The ID of the certification authority (CA) certificate.
* `gzip` - (Optional) Whether to enable "Gzip Compression". If enabled, files of specific file types will be compressed, otherwise, no files will be compressed. Default to true.
* `x_forwarded_for` - (Optional) Whether to set additional HTTP Header field "X-Forwarded-For" (documented below).
  * `retrive_client_ip` - (Optional) Whether to retrieve the client ip.
  * `retrive_slb_id` - (Optional) Whether to use the XForwardedFor header to obtain the ID of the SLB instance.
  * `retrive_slb_ip` - (Optional) Whether to use the XForwardedFor_SLBIP header to obtain the public IP address of the SLB instance.
  * `retrive_slb_proto` - (Optional) Whether to use the XForwardedFor_proto header to obtain the protocol used by the listener.
* `established_timeout` - (Optional) Timeout of tcp listener established connection idle timeout. Valid value range: [10-900] in seconds. Default to 900.
* `acl_status` - (Optional) Specifies whether to enable access control. Default value: off. Valid values: on, off.
* `acl_type` - (Optional) The type of the network ACL. Valid values: black, white. NOTE: If `acl_status` is set to on, `acl_type` is required. Otherwise, it will be ignored.
* `acl_id` - (Optional) The ID of the network ACL that is associated with the listener. NOTE: If `acl_status` is set to "on", `acl_id` is required. Otherwise, it will be ignored.
* `listener_forward` - (Optional, ForceNew) Whether to enable http redirect to https. Valid values are `on` and `off`. Default to `off`.
* `enable_http2` - (Optional) Specifies whether to enable HTTP/2. Default value: on. Valid values: "on", "off".
* `forward_port` - (Optional, ForceNew) The port that HTTP redirects to HTTPS.
* `tls_cipher_policy` - (Optional) The Transport Layer Security (TLS) security policy. Default value: tls_cipher_policy_1_0. Valid values: tls_cipher_policy_1_0, tls_cipher_policy_1_1, tls_cipher_policy_1_2, tls_cipher_policy_1_2_strict.
* `delete_protection_validation` - (Optional) Checking DeleteProtection of SLB instance before deleting. If true, this resource will not be deleted when its SLB instance enabled DeleteProtection. Default to false.
* `logs_download_attributes` - (Optional) Attributes related to logs download configuration. 
  * `log_project` - (Required) Name of the log project.
  * `log_store` - (Required) Name of the log store.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the load balancer listener. Its format as `<load_balancer_id>:<protocol>:<frontend_port>`.
* `description` - TLSCipherPolicy parameters.
* `health_check_method` - HealthCheckMethod used for health check.
* `health_check_connect_port` - Port used for health check.
* `ssl_certificate_id` - SSL certificate ID.
* `x_forwarded_for` - Whether to set additional HTTP Header field "X-Forwarded-For".
* `listener_forward` - Whether to enable HTTP redirect to HTTPS.