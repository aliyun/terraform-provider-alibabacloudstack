---
subcategory: "SLB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_domainextension"
sidebar_current: "docs-Alibabacloudstack-slb-domainextension"
description: |- 
  编排负载均衡(SLB)域名扩展列表
---

# alibabacloudstack_slb_domainextension
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_slb_domain_extension`

使用Provider配置的凭证在指定的资源集编排负载均衡(SLB)域名扩展列表。

## 示例用法

```hcl
variable "name" {
    default = "tf-testaccslbdomain_extension97638"
}

resource "alibabacloudstack_slb" "instance" {
  name                 = var.name
  internet_charge_type = "PayByTraffic"
  internet             = "true"
}

resource "alibabacloudstack_slb_server_certificate" "foo" {
  name               = "${var.name}-certificate"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq+gAwIBAgI+OuMs******XTtI90EAxEG/bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
  private_key        = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDO0knDrlNdiys******ErVpjsckAaOW/JDG5PCSwkaMxk=\n-----END RSA PRIVATE KEY-----"
}

resource "alibabacloudstack_slb_listener" "https" {
  load_balancer_id          = alibabacloudstack_slb.instance.id
  backend_port              = 80
  frontend_port            = 443
  protocol                 = "https"
  sticky_session           = "on"
  sticky_session_type      = "insert"
  cookie                   = "testslblistenercookie"
  cookie_timeout           = 86400
  health_check             = "on"
  health_check_uri         = "/cons"
  health_check_connect_port = 20
  healthy_threshold        = 8
  unhealthy_threshold      = 8
  health_check_timeout     = 8
  health_check_interval    = 5
  health_check_http_code   = "http_2xx,http_3xx"
  bandwidth                = 10
  ssl_certificate_id       = alibabacloudstack_slb_server_certificate.foo.id
}

resource "alibabacloudstack_slb_domainextension" "default" {
  load_balancer_id      = alibabacloudstack_slb.instance.id
  frontend_port         = alibabacloudstack_slb_listener.https.frontend_port
  domain                = "www.test.com"
  server_certificate_id = alibabacloudstack_slb_server_certificate.foo.id
  delete_protection_validation = false
}
```

## 参数参考

支持以下参数：

* `load_balancer_id` - (必填, 变更时重建) SLB实例的ID。
* `frontend_port` - (选填, 变更时重建) SLB实例使用的HTTPS监听器的前端端口。有效值：1–65535。
* `domain` - (必填, 变更时重建) 为该证书配置的域名。
* `server_certificate_id` - (必填) 要与域关联的服务器证书的ID。
* `delete_protection_validation` - (选填) 在删除前检查SLB实例的DeleteProtection。如果为true，当其SLB实例启用DeleteProtection时，此资源将不会被删除。默认为false。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - 域扩展的ID。
* `frontend_port` - SLB实例使用的HTTPS监听器的前端端口。
* `listener_port` - 监听器的前端端口，与`frontend_port`相同。
* `domain_extension_id` - 域扩展的唯一标识符。