---
subcategory: "SLB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_domainextensions"
sidebar_current: "docs-Alibabacloudstack-datasource-slb-domainextensions"
description: |- 
  查询负载均衡(SLB)域名扩展
---

# alibabacloudstack_slb_domainextensions
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_slb_domain_extensions`

根据指定过滤条件列出当前凭证权限可以访问的负载均衡(SLB)域名扩展列表。

## 示例用法

```hcl
variable "name" {
    default = "tf-testAccCheckAlibabacloudStackSlbsDataSourceBasic-18810"
}

resource "alibabacloudstack_slb" "instance" {
    name                = "${var.name}"
    internet_charge_type = "PayByTraffic"
    address_type       = "internet"
    specification      = "slb.s2.small"
}

resource "alibabacloudstack_slb_server_certificate" "foo" {
    name               = "${var.name}"
    server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDdjCCAl4CCQCcm*******XgthAiFFjl1S9ZgdA6Zc=\n-----END CERTIFICATE-----"
    private_key        = "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQ******7l3xC00BL7Z+SAJyI4QKA\n-----END RSA PRIVATE KEY-----"
}

resource "alibabacloudstack_slb_listener" "https" {
    load_balancer_id          = "${alibabacloudstack_slb.instance.id}"
    backend_port              = 80
    frontend_port             = 443
    protocol                  = "https"
    sticky_session            = "on"
    sticky_session_type       = "insert"
    cookie                    = "testslblistenercookie"
    cookie_timeout            = 86400
    health_check              = "on"
    health_check_uri          = "/cons"
    health_check_connect_port = 20
    healthy_threshold         = 8
    unhealthy_threshold       = 8
    health_check_timeout      = 8
    health_check_interval     = 5
    health_check_http_code    = "http_2xx,http_3xx"
    bandwidth                 = 10
    ssl_certificate_id        = "${alibabacloudstack_slb_server_certificate.foo.id}"
}

resource "alibabacloudstack_slb_domain_extension" "default" {
    load_balancer_id      = "${alibabacloudstack_slb.instance.id}"
    frontend_port         = "${alibabacloudstack_slb_listener.https.frontend_port}"
    domain                = "www.test.com"
    server_certificate_id = "${alibabacloudstack_slb_server_certificate.foo.id}"
}

data "alibabacloudstack_slb_domain_extensions" "default" {
    load_balancer_id = "${alibabacloudstack_slb_domain_extension.default.load_balancer_id}"
    frontend_port    = "${alibabacloudstack_slb_domain_extension.default.frontend_port}"
}

output "slb_domain_extension_ids" {
    value = data.alibabacloudstack_slb_domain_extensions.default.extensions[*].id
}
```

## 参数说明

以下参数是支持的：

* `ids` - (可选) SLB域名扩展的ID列表。查询结果将基于此参数进行过滤。
* `load_balancer_id` - (必填) 服务器负载均衡器实例的ID。这用于筛选与指定SLB实例关联的域名扩展。
* `frontend_port` - (必填) SLB实例使用的HTTPS监听器的前端端口。有效值范围为1到65535。这用于筛选指定监听器端口的域名扩展。

## 属性说明

除了上述参数外，还导出以下属性：

* `extensions` - SLB域名扩展的列表。每个元素包含以下属性：
  * `id` - 域名扩展的唯一ID。
  * `domain` - 与域名扩展关联的域名。
  * `server_certificate_id` - 域名使用的服务器证书ID。此证书通常用于HTTPS流量的SSL/TLS加密服务。