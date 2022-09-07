---
subcategory: "DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_record"
sidebar_current: "docs-alibabacloudstack-resource-dns-record"
description: |-
  Provides a DNS Record resource.
---

# alibabacloudstack\_dns\_record

Provides a DNS Record resource.

## Example Usage

```
resource "alibabacloudstack_dns_domain" "default" {
  domain_name = "domaintest."
  remark = "testing Domain"
}

# Create a new Domain record
resource "alibabacloudstack_dns_record" "default" {
  domain_id   = alibabacloudstack_dns_domain.default.id
  host_record = "testing_record"
  type        = "A"
  description = "testing Record"
  ttl         = 300
  rr_set      = ["192.168.2.4","192.168.2.7","10.0.0.4"]
}

output "record" {
  value = alibabacloudstack_dns_record.default.*
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Required) ID of the Dns Domain where this record belongs.
* `host_record` - (Required) Host record for the domain record. This host_record can have at most 253 characters, and each part split with "." can have at most 63 characters, and must contain only alphanumeric characters or hyphens, such as "-",".","*","@",  and must not begin or end with "-".
* `type` - (Required) The type of domain record. Valid values are `A`,`NS`,`MX`,`TXT`,`CNAME`,`SRV`,`AAAA`,`CAA`, `REDIRECT_URL` and `FORWORD_URL`.
* `rr_set` - (Optional) The value of domain record, When the `type` is `MX`,`NS`,`CNAME`,`SRV`, the server will treat the `value` as a fully qualified domain name, so it's no need to add a `.` at the end.
* `ttl` - (Optional) The effective time of domain record. Its scope depends on the edition of the cloud resolution. Free is `[600, 86400]`, Basic is `[120, 86400]`, Standard is `[60, 86400]`, Ultimate is `[10, 86400]`, Exclusive is `[1, 86400]`. Default value is `300`.
* `description` - (Optional) The effective time of domain record. Its scope depends on the edition of the cloud resolution. Free is `[600, 86400]`, Basic is `[120, 86400]`, Standard is `[60, 86400]`, Ultimate is `[10, 86400]`, Exclusive is `[1, 86400]`. Default value is `300`.

## Attributes Reference

The following attributes are exported:

* `record_id` - ID of the Dns Record.
* `type` - The record type.
* `host_record` - The host record of record.
* `rr_set` - The record value.
* `ttl` - The record effective time.
* `domain_id` - ID of the Dns Domain where this record belongs.
