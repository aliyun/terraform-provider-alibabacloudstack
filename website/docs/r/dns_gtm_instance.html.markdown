---
subcategory: "Alibaba Cloud DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_gtm_instance"
description: |-
  Create and manage Global Traffic Manager (GTM) instances for Alibaba Cloud DNS
---

# alibabacloudstack_dns_gtm_instance

Create and manage Global Traffic Manager (GTM) instances for Alibaba Cloud DNS using credentials configured in the provider within the specified resource set.

## Example Usage

### Basic Usage

```hcl

variable "name" {
  default = "tfacc24074"
}
resource "alibabacloudstack_vpc_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}_vpc0"
}

resource "alibabacloudstack_dns_private_domain" "default" {
  name   = "${var.name}.testtf."
  remark = "Created by Terraform for DNS record test"
  vpc_ids = [
    "${alibabacloudstack_vpc_vpc.default.id}"
  ]
}

resource "alibabacloudstack_dns_gtm_instance" "default" {
  name    = var.name
  prefix  = var.name
  zone_id = alibabacloudstack_dns_private_domain.default.id
  ttl     = "300"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the scheduling instance. The name must be 1 to 128 characters in length and cannot start with `http://` or `https://`.
* `prefix` - (Required) The scheduling domain name. Used to construct the complete scheduling domain name. For example, when the prefix is "www" and the zone_name is "example.com.", the complete scheduling domain name is "www.example.com.".
* `ttl` - (Required) The global TTL (time-to-live) in seconds. Represents the duration that DNS resolution results are cached on the client side.
* `zone_id` - (Required) The domain name ID of the scheduling domain. This ID corresponds to an already created Alibaba Cloud DNS domain.

## Attributes Reference

The following attributes are exported from the API response:

* `id` - The ID of the scheduling instance.
* `create_timestamp` - The creation timestamp in seconds.
* `update_timestamp` - The update timestamp in seconds.
* `zone_name` - The domain name. Represents the domain to which the scheduling instance belongs.