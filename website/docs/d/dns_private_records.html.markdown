---
subcategory: "Alibaba Cloud DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_private_records"
description: |-
  Query Alibaba Cloud DNS private domain resolution records
---

# alibabacloudstack_dns_private_records

> Data source for querying DNS private domain resolution records

## Example Usage

```hcl
  
variable "name" {
  default = "tf-testacc56057"
}

resource "alibabacloudstack_dns_private_domain" "default" {
  name   = "${var.name}.testtf."
  remark = "Created by Terraform for DNS record test"
}

resource "alibabacloudstack_dns_private_record" "default" {
  zone_id      = alibabacloudstack_dns_private_domain.default.id
  name         = var.name
  type         = "A"
  ttl          = 300
  lba_strategy = "ALL_RR"
  line_ids     = ["default"]
  rdatas {
    value = "192.168.1.1"
  }
  rdatas {
    value = "127.0.0.1"
  }
}

data "alibabacloudstack_dns_private_domains" "default" {
  zone_id    = alibabacloudstack_dns_private_record.default.zone_id
  name_regex = "tf-testacc[0-9]+"
}
```

## Argument Reference

The following arguments support filtering query results:

* `zone_id` (Required): Domain ID, used to specify the domain zone to query.

* `name_regex` (Optional): Regular expression for the resolution record name, used to filter results by name.

* `ids` (Optional): List of resolution record IDs to query, in the format "zone_id:id".

## Attributes Reference

The following attributes are exported:

* `id`: Unique identifier of the data source, in the format "zone_id1:id1,zone_id2:id2,..."

* `records`: List of queried resolution records. Each record contains the following attributes:
  * `id`: Resolution record ID.
  * `create_timestamp`: Creation timestamp in seconds.
  * `line_ids`: List of resolution lines.
  * `lba_strategy`: Load balancing strategy. Valid values:
    - ALL_RR: Return all addresses (non-weighted).
    - RATIO: Return addresses by weight (weighted).
  * `name`: Record name.
  * `rdatas`: List of record values. Each value contains the following attributes:
    * `lba_weight`: Weight.
    * `value`: Value.
  * `remark`: Remark.
  * `ttl`: Cache time of the record.
  * `type`: Record type. Valid values:
    - A: Point the domain to an IPv4 address.
    - AAAA: Point the domain to an IPv6 address.
    - CNAME: Point the domain to another domain name.
    - MX: Point the domain to a mail server address.
    - TXT: Text length limited to 255, usually used for SPF records (anti-spam).
    - PTR: Record the domain name for reverse DNS lookup of an IP address.
    - SRV: Record the server providing a specific service.
    - NAPTR: Record the domain name authority pointer, usually used for ENUM.
    - CAA: Set certificate authority authorization information for the domain.
    - NS: Point the domain to a specified DNS server for resolution.
  * `update_timestamp`: Modification timestamp in seconds.
  * `zone_id`: Domain ID.