---
subcategory: "Universal DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_gtm_instances"
sidebar_current: "docs-Alibabacloudstack-datasource-dns-gtm-instances"
description: |-
  Provides a list of DNS GTM (Global Traffic Management) instances.
---

# alibabacloudstack_dns_gtm_instances

This data source provides a list of DNS GTM (Global Traffic Management) instances that can be accessed by an Alibaba Cloud account.

## Example Usage

```hcl
variable "name" {
  default = "tfacc60948"
}

resource "alibabacloudstack_vpc_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}_vpc0"
}

resource "alibabacloudstack_dns_private_domain" "default" {
  name   = "${var.name}.testtf."
  remark = "Created by Terraform for DNS GTM instance test"
  vpc_ids = [
    "${alibabacloudstack_vpc_vpc.default.id}"
  ]
}

resource "alibabacloudstack_dns_gtm_instance" "default" {
  name    = var.name
  prefix  = var.name
  zone_id = alibabacloudstack_dns_private_domain.default.id
  ttl     = 300
}

data "alibabacloudstack_dns_gtm_instances" "default" {
  name_regex = alibabacloudstack_dns_gtm_instance.default.name
}
```

## Argument Reference

The following arguments are supported:

* `ids` (List, Optional): Filters results by one or more instance IDs. Instance IDs are in the format like `886c6195-36d7-4d97-8bb5-d7f4228ec55c`. Multiple IDs in the list represent a logical OR relationship.

* `name_regex` (String, Optional): Filters results by regular expression matching on instance names. For example, setting it to `^test.*` will match all instance names starting with "test".

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` (String): The data source ID, which is a hash value calculated from the list of instance IDs.

* `instances` (List): A list of DNS GTM instances in the query results. Each element contains the following attributes:
  * `create_timestamp` (Integer): The timestamp (in seconds) when the instance was created.
  * `id` (String): The unique identifier of the DNS GTM instance.
  * `name` (String): The name of the DNS GTM instance.
  * `prefix` (String): The prefix of the scheduling domain name.
  * `ttl` (Integer): The global TTL (Time To Live) value in seconds.
  * `update_timestamp` (Integer): The timestamp (in seconds) of the last update to the instance.
  * `zone_id` (String): The domain name ID of the affiliated scheduling domain.
  * `zone_name` (String): The domain name.