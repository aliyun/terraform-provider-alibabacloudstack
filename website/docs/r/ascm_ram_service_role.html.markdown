---
subcategory: "Apsara Stack Cloud Management"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_ram_service_role"
description: |-
  Provides Ascm service ram role.
---

# alibabacloudstack_ascm_ram_service_role

Provides Ascm service ram role.

~> **Note:** This resource can also be referred to by the following aliases:
- `alibabacloudstack_ascm_service_ram_role`

## Example Usage

```
resource "alibabacloudstack_ascm_organization" "default" {
  name = "Tf-testingresource-org2"
  parent_id = "1"
} 
 resource "alibabacloudstack_ascm_ram_service_role" "default" {
  organization_id = "${alibabacloudstack_ascm_organization.default.id}"
  product_name = "ECS"
}
```

## Argument Reference

The following arguments are supported:

* `product_name` - (Required, ForceNew) The name of the product.
* `organization_id` - (Required, ForceNew) The ID of the organization.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service ram role.
* `ram_roles` - The list of ram roles.
  * `id` - The ID of the ram role.
  * `product_name` - The name of the product.
  * `organization_id` - The ID of the organization.
  * `arn` - The ARN of the ram role.
  * `role_id` - The ID of the role.
  * `role_name` - The name of the role.
  * `role_type` - The type of the role.
  * `region` - The region of the ram role.
  * `description` - The description of the ram role.
  * `aliyun_user_id` - The Aliyun user ID.
  * `assume_role_policy_document` - The assume role policy document.
  * `organization_name` - The name of the organization.
  * `policies` - The list of policies attached to the ram role.
    * `policy_id` - The ID of the policy.
    * `region` - The region of the policy.
    * `policy_name` - The name of the policy.
    * `description` - The description of the policy.
    * `policy_document` - The policy document.
    * `policy_type` - The type of the policy.
    * `default_version` - The default version of the policy.
    * `aliyun_user_id` - The Aliyun user ID.
    * `ram_group_id` - The RAM group ID.
    * `ascm_ram_policy_id` - The ASCM RAM policy ID.
    * `attach_date` - The attach date of the policy.
    * `resource_set_id` - The resource set ID.
    * `privilege_id` - The privilege ID.
    * `ram_role_id` - The RAM role ID.

## Import

ASCM Service RAM Role can be imported using the organization_id and product_name in the format `<organization_id>:<product_name>`, e.g.

```
$ terraform import alibabacloudstack_ascm_ram_service_role.example 1:ECS
```