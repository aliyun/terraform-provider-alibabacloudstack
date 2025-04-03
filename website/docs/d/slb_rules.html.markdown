---
subcategory: "SLB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_rules"
sidebar_current: "docs-Alibabacloudstack-datasource-slb-rules"
description: |- 
  Provides a list of slb rules owned by an alibabacloudstack account.
---

# alibabacloudstack_slb_rules

This data source provides a list of SLB listener rules in an AlibabaCloudStack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_slb_rules" "sample_ds" {
  load_balancer_id = "${alibabacloudstack_slb.sample_slb.id}"
  frontend_port    = 80
}

output "first_slb_rule_id" {
  value = "${data.alibabacloudstack_slb_rules.sample_ds.slb_rules.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required) The ID of the Server Load Balancer instance with the listener rules.
* `frontend_port` - (Required) The port number of the SLB listener.
* `ids` - (Optional) A list of rule IDs to filter results.
* `name_regex` - (Optional) A regex string to filter results by rule name.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of SLB listener rule IDs.
* `names` - A list of SLB listener rule names.
* `slb_rules` - A list of SLB listener rules. Each element contains the following attributes:
  * `id` - The ID of the rule.
  * `name` - The name of the rule.
  * `domain` - The domain name in the HTTP request where the rule applies (e.g., `"*.aliyun.com"`).
  * `url` - The path in the HTTP request where the rule applies (e.g., `"/image"`).
  * `server_group_id` - The ID of the linked VServer group.
