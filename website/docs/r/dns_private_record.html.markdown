---
subcategory: "Alibaba Cloud DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_private_record"
description: |-
  Manages a private DNS record in Alibaba Cloud.
---

# alibabacloudstack_dns_private_record

This resource manages a private DNS record in Alibaba Cloud, used for creating and managing private records.

## Example Usage

### Basic Usage

```hcl
variable "name" {
  default = "tf-testacc58082"
}

resource "alibabacloudstack_dns_private_domain" "default" {
  name   = "${var.name}.testtf."
  remark = "Created by Terraform for DNS private record test"
}

resource "alibabacloudstack_dns_private_record" "default" {
  name         = var.name
  type         = "A"
  ttl          = 300
  lba_strategy = "ALL_RR"
  line_ids = [
    "default"
  ]
  rdatas {
    value = "192.168.1.1"
  }
  rdatas {
    value = "127.0.0.1"
  }

  zone_id = alibabacloudstack_dns_private_domain.default.id
}
```

## Argument Reference

The following arguments are supported:

* `line_ids` - (Required) The list of resolution lines. Example: `["default"]`.
* `lba_strategy` - (Required) The load balancing strategy. Valid values:
  * `ALL_RR`: Returns all addresses (non-weighted).
  * `RATIO`: Returns addresses by weight (weighted).
* `name` - (Required) The record name. Example: `"test"`.
* `rdatas` - (Required) The set of record values. Each record value contains the following attributes:
  * `value` - (Required) The value. Example: `"192.168.1.1"`.
  * `lba_weight` - (Optional) The weight. Only valid when `lba_strategy` is `RATIO`. Example: `100`.
* `ttl` - (Required) The cache time of the record (in seconds). Example: `300`.
* `type` - (Required) The record type. Valid values:
  * `A`: Points the domain name to an IPv4 address.
  * `AAAA`: Points the domain name to an IPv6 address.
  * `CNAME`: Points the domain name to another domain name.
  * `MX`: Points the domain name to a mail server address.
  * `TXT`: Text length limit is 255, usually used for SPF records (anti-spam).
  * `PTR`: Records the domain name for reverse DNS resolution of IP addresses.
  * `SRV`: Records servers that provide specific services.
  * `NAPTR`: Records domain name authority pointers, usually used for ENUM.
  * `CAA`: Sets certificate authority authorization information for the domain name.
  * `NS`: Points the domain name to a specified DNS server for resolution.
* `zone_id` - (Required, Forces new resource when changed) The domain ID. Example: `"f27ffcc8-02a4-4ca1-9ca4-6bffc4e82034"`.
* `remark` - (Optional) The remark information.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID, in the format `{ZoneId}:{Id}`.

## Import

DNS Private Record can be imported using the ZoneId and RecordId in the format `{ZoneId}:{Id}`, e.g.

```
$ terraform import alibabacloudstack_dns_private_record.example f27ffcc8-02a4-4ca1-9ca4-6bffc4e82034:record-12345
```