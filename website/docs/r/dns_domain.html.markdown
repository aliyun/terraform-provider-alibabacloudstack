---
subcategory: "DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_domain"
sidebar_current: "docs-alibabacloudstack-resource-dns-domain"
description: |-
  Provides a DNS domain resource.
---

# alibabacloudstack_dns_domain

Provides a DNS domain resource.

-> **NOTE:** The domain name which you want to add must be already registered and had not added by another account. Every domain name can only exist in a unique group.

## Example Usage

```
# Add a new Domain.
resource "alibabacloudstack_dns_domain" "default" {
  domain_name     = "starmove."
  remark   =  "testing Domain"
}
output "dns" {
  value = alibabacloudstack_dns_domain.default.*
}
```
## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, ForceNew) Name of the domain. This name without suffix can have a string of 1 to 63 characters(domain name subject, excluding suffix), must contain only alphanumeric characters or "-", and must not begin or end with "-", and "-" must not in the 3th and 4th character positions at the same time. Suffix `.sh` and `.tel` are not supported.
* `group_id` - (Optional) Id of the group in which the domain will add. If not supplied, then use default group.
* `resource_group_id` - (Optional, ForceNew) The Id of resource group which the dns domain belongs.
* `lang` - (Optional) User language.
* `remark` - (Optional) Remarks information for your domain name.
* `domain_name` - (Required) Name of the domain. 

## Attributes Reference

The following attributes are exported:

* `id` - This ID of this resource. The value is set to `domain_name`.
* `domain_id` - The domain ID.
* `dns_servers` - A list of the dns server name.
* `domain_name` - The name of the domain. 