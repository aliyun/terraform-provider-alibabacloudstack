---
subcategory: "DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_domain_attachment"
sidebar_current: "docs-alibabacloudstack-resource-dns-domain-attachment"
description: |-
  Provides bind the domain name to the DNS instance resource.
---

# alibabacloudstack_dns_domain_attachment
-> **NOTE:** Alias name has: `alibabacloudstack_alidns_domainattachment`

Provides bind the domain name to the DNS instance resource.


## Example Usage

```
resource "alibabacloudstack_dns_domain_attachment" "dns" {
  instance_id     = "dns-cn-mp91lyq9xxxx"
  domain_names    = ["test111.abc", "test222.abc"]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The id of the DNS instance.
* `domain_names` - (Required) The domain names bound to the DNS instance.

## Attributes Reference

The following attributes are exported:

* `id` - This ID of this resource. The value is same as `instance_id`. 
* `domain_names` - Domain names bound to DNS instance.
* `instance_id` - The id of the DNS instance.