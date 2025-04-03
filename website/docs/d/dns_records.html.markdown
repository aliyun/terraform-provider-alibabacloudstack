---
subcategory: "DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_records"
sidebar_current: "docs-alibabacloudstack-datasource-dns-records"
description: |-
    Provides a list of records available to the dns.
---

# alibabacloudstack_dns_records

This data source provides a list of DNS Domain Records in an AlibabacloudStack Cloud account according to the specified filters.

## Example Usage

```
resource "alibabacloudstack_dns_domain" "default" {
  domain_name = "domaintest."
  remark = "testing Domain"
}

# Create a new Domain record
resource "alibabacloudstack_dns_record" "default" {
  zone_id   = alibabacloudstack_dns_domain.default.domain_id
  name = "testing_record"
  type        = "A"
  remark = "testing Record"
  ttl         = 300
  lba_strategy = "ALL_RR"
  rr_set      = ["192.168.2.4","192.168.2.7","10.0.0.4"]
}

data "alibabacloudstack_dns_records" "default"{
 zone_id         = alibabacloudstack_dns_record.default.zone_id
 name = alibabacloudstack_dns_record.default.name
}
output "records" {
  value = data.alibabacloudstack_dns_records.default.*
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required) The domain Id associated to the records.
* `host_record_regex` - (Optional, ForceNew) Host record regex.
* `type` - (Optional) Record type. Valid items are `A`, `NS`, `MX`, `TXT`, `CNAME`, `SRV`, `AAAA`, `REDIRECT_URL`, `FORWORD_URL` .
* `ids` - (Optional) A list of record IDs.
* `name` - (Optional,) Name of the DNS record.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of record IDs. 
* `records` - A list of records. Each element contains the following attributes:
  * `record_id` - ID of the record.
  * `zone_id` - ID of the domain the record belongs to.
  * `name` - Host record of the domain.
  * `type` - Type of the record.
  * `ttl` - TTL of the record.
  * `remark` - Description of the record.
  * `rr_set` - RrSet for the record.