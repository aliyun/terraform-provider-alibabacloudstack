---
subcategory: "Alibaba Cloud DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_private_domains"
sidebar_current: "docs-Alibabacloudstack-datasource-dns-private-domains"
description: |-
  Provides a list of private domains for DNS that can be used by an Alibaba Cloud account.
---

# alibabacloudstack_dns_private_domains

This data source provides private domains for DNS that can be accessed by an Alibaba Cloud account within the region configured in the provider. It is used to retrieve information of created private domains, including domain ID, name, associated VPCs, and other details.

## Example Usage

```hcl
variable "name" {
  default = "tf-testacc4077889550832453245."
}

resource "alibabacloudstack_vpc_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}_vpc"
}

resource "alibabacloudstack_dns_private_domain" "default" {
  name    = var.name
  remark  = var.name
  vpc_ids = [alibabacloudstack_vpc_vpc.default.id]
}

data "alibabacloudstack_dns_private_domains" "default" {
  name_regex = alibabacloudstack_dns_private_domain.default.name
}
```

## Argument Reference

The following arguments are supported:

* `id` (Optional): The domain ID for exact matching of a single domain.
* `ids` (Optional): A list of domain IDs for filtering multiple specific domains.
* `name` (Optional): The domain name for exact matching.
* `name_regex` (Optional): A regular expression for fuzzy matching of domain names.
* `vpc_id` (Optional): The VPC ID for querying domains associated with a specific VPC.

## Attributes Reference

The following attributes are exported:

* `id` (String): The data source ID, generated from the hash value of the query result.
* `domains` (List): A list of queried private domains. Each element contains the following attributes:
  * `id` (String): The domain ID.
  * `name` (String): The domain name.
  * `caller_uid` (String): The caller UID.
  * `create_timestamp` (Integer): The domain creation timestamp in seconds.
  * `gtm_instance_count` (Integer): The total number of associated Global Traffic Manager instances.
  * `record_count` (Integer): The total number of resolution record sets under the domain.
  * `region_and_vpcs` (List): A list of region and VPC association information. Each element contains:
    * `region_id` (String): The region ID.
    * `vpcs` (List): A list of associated VPCs. Each element contains:
      * `id` (String): The VPC ID.
      * `name` (String): The VPC name.
  * `remark` (String): The domain remark information.
  * `update_timestamp` (Integer): The domain last modification timestamp in seconds.
* `ids` (List): A list of all queried domain IDs.