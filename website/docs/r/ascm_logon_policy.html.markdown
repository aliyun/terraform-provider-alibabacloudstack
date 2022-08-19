---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ascm_logon_policy"
sidebar_current: "docs-apsarastack_ascm_logon_policy"
description: |-
  Provides a Apsarastack Logon Policy resource.
---
# apsarastack\_ascm_logon_policy

Provides a Apsarastack Logon Policy resource.

Basic Usage

```
resource "apsarastack_ascm_logon_policy" "login" {
  name="test_foo"
  description="testing purpose"
  rule="ALLOW"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Logon Policy.

* `description` - (Optional) The Logon Policy description.

* `rule` - (Optional)  The rules of the logon policy. Valid values: Allow and Deny.


## Attributes Reference

The following attributes are exported:

* `name` - The name of the Logon Policy.
* `description` - The Logon Policy description.
* `rule` - The Rule for the Logon Policy.
* `policy_id` - The ID of the logon policy created.

