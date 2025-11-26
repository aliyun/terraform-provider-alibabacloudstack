---
subcategory: "Universal DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_universal_dns_domains"
sidebar_current: "docs-Alibabacloudstack-datasource-universal-dns-domains"
description: |-
  Queries the list of cross-cloud DNS domains.
---

# alibabacloudstack_universal_dns_domains

Queries the list of cross-cloud DNS domains provided by Alibaba Cloud.

## Example Usage

```hcl
variable "name" {
  default = "tf-testacc6879360584559214283"
}

resource "alibabacloudstack_universal_dns_domain" "default" {
  name   = "${var.name}.example."
  remark = "Created by Terraform"
}

data "alibabacloudstack_universal_dns_domains" "default" {
  name_regex = alibabacloudstack_universal_dns_domain.default.name
}
```

## Argument Reference

The following arguments support filtering query results:

* `ids` (Optional): A list of domain IDs used for filtering. Only domains with IDs in this list will be returned.

* `name_regex` (Optional): A regular expression used to filter by domain name. Only domains whose names match this regular expression will be returned.

## Attributes Reference

The following attributes are exported:

* `id` (String): The unique identifier of the data source, based on the hash value of the returned list of domain IDs.

* `domains` (List): A list of cross-cloud DNS domains matching the conditions. Each domain contains the following attributes:
  * `id` (String): The ID of the cross-cloud DNS domain.
  * `create_timestamp` (Integer): The timestamp (in seconds) when the domain was created.
  * `name` (String): The name of the cross-cloud DNS domain.
  * `record_count` (Integer): The total number of DNS records in the domain.
  * `remark` (String): The remark or description of the domain.
  * `update_timestamp` (Integer): The timestamp (in seconds) when the domain was last updated.

* `ids` (List): A list of domain IDs for the found domains.

* `names` (List): A list of domain names for the found domains.