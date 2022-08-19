---
subcategory: "DNS"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_dns_domains"
sidebar_current: "docs-apsarastack-datasource-dns-domains"
description: |-
    Provides a list of domains available to the user.
---

# apsarastack\_dns\_domains

This data source provides a list of DNS Domains in an Apsarastack Cloud account according to the specified filters.

## Example Usage

```
resource "apsarastack_dns_domain" "default" {
  domain_name = "domaintest."
  remark = "testing Domain"
}
data "apsarastack_dns_domains" "default"{
  domain_name_regex         = apsarastack_dns_domain.default.domain_name
}
output "domains" {
  value = data.apsarastack_dns_domains.default.*
}
```

## Argument Reference

The following arguments are supported:

* `domain_name_regex` - (Optional) A regex string to filter results by the domain name. 
* `ids` (Optional) - A list of domain IDs.
* `resource_group_id` - (Optional, ForceNew) The ID of resource group which the dns belongs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of domain IDs.
* `names` - A list of domain names.
* `domains` - A list of domains. Each element contains the following attributes:
  * `domain_id` - ID of the domain.
  * `domain_name` - Name of the domain.
  * `dns_servers` - DNS list of the domain in the analysis system.
  * `resource_group_id` - The ID of resource group which the dns belongs.
