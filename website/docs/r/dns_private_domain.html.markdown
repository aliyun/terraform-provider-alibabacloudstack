---
subcategory: "Cloud DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_private_domain"
sidebar_current: "docs-Alibabacloudstack-dns-private_domain"
description: |-
  Manages a private DNS domain in Alibaba Cloud.
---

# alibabacloudstack_dns_private_domain

This resource manages a private DNS domain in Alibaba Cloud, used for creating and managing private domains.

## Example Usage

### Basic Usage

```hcl
variable "name" {
  default = "tfacc77356.test."
}

resource "alibabacloudstack_vpc_vpc" "default0" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}_vpc0"
}

resource "alibabacloudstack_vpc_vpc" "default1" {
  cidr_block = "192.168.0.0/16"
  vpc_name   = "${var.name}_vpc1"
}

resource "alibabacloudstack_dns_private_domain" "default" {
  name = var.name
  vpc_ids = [
    alibabacloudstack_vpc_vpc.default0.id,
    alibabacloudstack_vpc_vpc.default1.id,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource when changed) The name of the private domain. Must end with a dot, e.g., "example.com.".

* `remark` - (Optional) The remark for the domain, used to describe the purpose of the private domain.

* `vpc_ids` - (Optional) A list of VPC IDs to associate with. Binds the private domain to the specified VPC networks, allowing ECS instances within the VPC to access the domain via the internal network.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the private domain.

* `caller_uid` - The system parameter indicating the user ID who created the private domain.

* `create_timestamp` - The creation timestamp in seconds.

* `record_count` - The total number of DNS record sets.

* `update_timestamp` - The last update timestamp in seconds.