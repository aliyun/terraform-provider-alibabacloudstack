---
subcategory: "Server Load Balancer (SLB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_domain_extensions"
sidebar_current: "docs-alibabacloudstack-resource-slb-domain-extensions"
description: |-
  Provides a Load Banlancer domain extension Resource and add it to one Listener.
---

# alibabacloudstack\_slb\_domain_extensions

This data source provides the domain extensions associated with a server load balancer listener.

## Example Usage
```
data "alibabacloudstack_slb_domain_extensions" "foo" {
  ids               = ["fake-de-id"]
  load_balancer_id  = "fake-lb-id"
  frontend_port     = "fake-port"
}

output "slb_domain_extension" {
  value = "${data.alibabacloudstack_slb_domain_extensions.foo.extensions.0.id}"
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) IDs of the SLB domain extensions.
* `load_balancer_id` - (Required) The ID of the SLB instance.
* `frontend_port` - (Required) The frontend port used by the HTTPS listener of the SLB instance. Valid values: 1–65535.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `extensions` - A list of SLB domain extension. Each element contains the following attributes:
    * `id` - The ID of the domain extension.
    * `domain` - The domain name.
    * `server_certificate_id` - The ID of the certificate used by the domain name.
