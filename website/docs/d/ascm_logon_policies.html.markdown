---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_logon_policies"
sidebar_current: "docs-alibabacloudstack-ascm-logon-policies"
description: |-
  Provides a list of Logon Policies.
---
# alibabacloudstack\_ascm\_logon\_policies

Provides a list of Logon Policies.

Basic Usage

```
resource "alibabacloudstack_ascm_logon_policy" "default" {
  name="Test_login_policy"
  description="testing policy"
  rule="ALLOW"
}
output "login" {
  value = alibabacloudstack_ascm_logon_policy.default.id
}
data "alibabacloudstack_ascm_logon_policies" "default"{
  name = alibabacloudstack_ascm_logon_policy.default.name
}
output "policies" {
  value = data.alibabacloudstack_ascm_logon_policies.default.*
}
```
## Argument Reference

The following arguments are supported:

* `ids` - (Optional) The ids of the Logon Policies.
* `name` - (Optional) The name of the Logon Policy.
* `name_regex` - (Optional) A regex string to filter Logon Policies by name.
* `description` - (Optional) The Logon Policies description.
* `rule` - (Optional) The Rule for the Logon Policies.
* `ip_ranges` - (Optional) The IP Ranges for the Logon Policies.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported:

* `name` - The name of the Logon Policies. 
* `policies` - The list of the Logon Policies. Each element contains the following attributes:
    * `id` - The id of the Logon Policy.
    * `name` - The name of the Logon Policy.
    * `description` - The description of the Logon Policy.
    * `rule` - The rule of the Logon Policy.
    * `ip_range` - The ip range of the Logon Policy.
    * `end_time` - The end time of the Logon Policy.
    * `start_time` - The start time of the Logon Policy.
    * `login_policy_id` - The login policy id of the Logon Policy.

