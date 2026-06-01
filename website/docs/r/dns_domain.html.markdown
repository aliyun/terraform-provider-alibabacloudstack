---
subcategory: "Universal DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_domain"
sidebar_current: "docs-Alibabacloudstack-resource-dns-domain"
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
  domain_name = "tfacc-test."
  remark      = "testing Domain"
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, ForceNew) Name of the domain. Modifying this parameter will force a new resource to be created.
* `remark` - (Optional) Remarks information for your domain name.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. Format: `domain_name:domain_id`.
* `domain_id` - The domain ID.
* `domain_name` - The name of the domain.

## Import

DNS Domain can be imported using the `domain_name` and `domain_id` separated by a colon, e.g.

```
$ terraform import alibabacloudstack_dns_domain.example example.com:12345678
```
