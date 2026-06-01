---
subcategory: "Alibaba Cloud DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_universal_dns_domain"
sidebar_current: "docs-Alibabacloudstack-resource-universal-dns-domain"
description: |-
  Cross-cloud DNS domain resource
---

# alibabacloudstack_universal_dns_domain

Creates and manages cross-cloud DNS domain resources.

## Example Usage

### Basic Usage

```hcl
variable "name" {
  default = "tf-testacc35847"
}

resource "alibabacloudstack_universal_dns_domain" "default" {
  name   = "tf-testacc35847.example."
  remark = "Created by Terraform"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The cross-cloud shared domain name (must end with ".").
* `remark` - (Optional) The remark for the domain.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the domain.
* `create_timestamp` - The creation timestamp in seconds.
* `record_count` - The total number of DNS record sets.
* `update_timestamp` - The update timestamp in seconds.