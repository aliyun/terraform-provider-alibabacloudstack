---
subcategory: "Apsara Stack Cloud Management"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_account_ascm_custom_role_arn"
description: |-
    Queries the RAM role ARN that Apsara Stack Cloud Management materializes for a role of the account the provider is authenticated with.
---

# alibabacloudstack_account_ascm_custom_role_arn

This data source queries the RAM role which Apsara Stack Cloud Management (ASCM) materializes for one of the roles
held by the account the provider is authenticated with.

~> **NOTE:** An ASCM role itself (`roleType = ROLETYPE_ASCM`, which is what `alibabacloudstack_ascm_custom_role`
and `alibabacloudstack_ascm_ram_role` create) is a control plane role and cannot be assumed by STS. For every
`organization` + `ASCM role` + `ASCM user` binding, ASCM maintains a separate role in the RAM/STS backend, and only
the ARN of that role is accepted by `AssumeRole` or by resource arguments such as the encryption `arn` of a log
store. The materialized role name follows the `ascm-role-<organization-id>-<ascm-role-id>-<ascm-user-id>` convention.

~> **NOTE:** The organization and the ASCM user are not arguments. ASCM resolves both of them from the identity of
the caller and refuses to disclose the role of another user, so this data source always answers the question
"which RAM role does the account I log in with hold for this ASCM role". To query a different identity, configure an
aliased provider with the credentials of that user.

## Example Usage

```
data "alibabacloudstack_account_ascm_custom_role_arn" "operation_admin" {
  ascm_role_id = 1
}

output "role_arn" {
  value = data.alibabacloudstack_account_ascm_custom_role_arn.operation_admin.arn
}
```

## Argument Reference

The following arguments are supported:

* `ascm_role_id` - (Required) The ID of the ASCM role on the control plane side. The authenticated account must hold this role, otherwise the read fails with `ascm.auth.EntityNotExist.User.Role`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ARN of the materialized RAM role.
* `arn` - The ARN of the RAM role materialized by ASCM, in the format `acs:ram::<account-id>:role/<role-name>`. This is the value to pass to STS `AssumeRole`.
* `role_name` - The name of the materialized RAM role.
* `ram_role_id` - The ID of the materialized RAM role on the RAM side. It is different from `ascm_role_id`, which is the ID on the ASCM control plane side.
* `assume_role_policy_document` - The trust policy document of the materialized RAM role, in JSON format. Its `Principal` entries decide which account or user is allowed to assume the role.
