---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "ApsaraStack: apsarastack_ascm_ram_policies"
sidebar_current: "docs-apsarastack-datasource-ascm-ram-policies"
description: |-
    Provides a list of ram policies to the user.
---

# apsarastack\_ascm_ram_policies

This data source provides the ram policies of the current Apsara Stack Cloud user.

## Example Usage

```
resource "apsarastack_ascm_ram_policy" "default" {
  name = "TestPolicy"
  description = "Testing"
  policy_document = "{\"Statement\":[{\"Action\":\"ecs:*\",\"Effect\":\"Allow\",\"Resource\":\"*\"}],\"Version\":\"1\"}"

}
data "apsarastack_ascm_ram_policies" "default" {
  name_regex = apsarastack_ascm_ram_policy.default.name
}
output "ram_policies" {
  value = data.apsarastack_ascm_ram_policies.default.*
}


```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of ram policy IDs.
* `name_regex` - (Optional) A regex string to filter results by ram policy name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `policies` - A list of policies. Each element contains the following attributes:
    * `id` - ID of the policy.
    * `name` - Policy name.
    * `description` - Description about the policy.
    * `ctime` -  Creation time of ram policy.
    * `cuser_id` - ID of the policy creator.
    * `region` - Name of the region where policy belongs.
    * `policy_document` - Policy Document.
     
