---
subcategory: "DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_records"
sidebar_current: "docs-alibabacloudstack-datasource-dns-records"
description: |-
    Provides a list of records available to the dns.
---

# alibabacloudstack\_dns\_records

This data source provides a list of DNS Domain Records in an AlibabacloudStack Cloud account according to the specified filters.

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

data "alibabacloudstack_dns_records" "default"{
  domain_id         = alibabacloudstack_dns_record.default.domain_id
  host_record_regex = alibabacloudstack_dns_record.default.host_record
}
output "records" {
  value = data.alibabacloudstack_dns_records.default.*
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Required) The domain Id associated to the records.
* `host_record_regex` - (Optional) Host record regex. 
* `value_regex` - (Optional) Host record value regex. 
* `type` - (Optional) Record type. Valid items are `A`, `NS`, `MX`, `TXT`, `CNAME`, `SRV`, `AAAA`, `REDIRECT_URL`, `FORWORD_URL` .
* `ids` - (Optional) A list of record IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of record IDs. 
* `records` - A list of records. Each element contains the following attributes:
  * `record_id` - ID of the record.
  * `domain_id` - ID of the domain the record belongs to.
  * `host_record` - Host record of the domain.
  * `type` - Type of the record.
  * `ttl` - TTL of the record.
  * `description` - Description of the record.
  * `rr_set` - RrSet for the record.
