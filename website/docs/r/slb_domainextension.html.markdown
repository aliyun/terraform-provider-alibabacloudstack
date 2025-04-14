---
subcategory: "SLB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_domainextension"
sidebar_current: "docs-Alibabacloudstack-slb-domainextension"
description: |- 
  Provides a slb Domainextension resource.
---

# alibabacloudstack_slb_domainextension
-> **NOTE:** Alias name has: `alibabacloudstack_slb_domain_extension`

Provides a slb Domainextension resource.

## Example Usage

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
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, ForceNew) The ID of the SLB instance.
* `frontend_port` - (Optional, ForceNew) The frontend port used by the HTTPS listener of the SLB instance. Valid values: 1–65535.
* `listener_port` - (Optional, ForceNew) The listener port used by the HTTPS listener of the SLB instance. Valid values: 1–65535. 
* `domain` - (Required, ForceNew) The domain name for which the certificate is configured.
* `server_certificate_id` - (Required) The ID of the server certificate to be associated with the domain.
* `delete_protection_validation` - (Optional) Checking DeleteProtection of SLB instance before deleting. If true, this resource will not be deleted when its SLB instance enabled DeleteProtection. Default to false.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the domain extension.
* `frontend_port` - The frontend port used by the HTTPS listener of the SLB instance.
* `listener_port` - The listener port used by the HTTPS listener of the SLB instance. 
* `domain_extension_id` - The unique identifier for the domain extension.