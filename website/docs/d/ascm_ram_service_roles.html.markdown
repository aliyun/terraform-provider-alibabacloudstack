---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_ram_service_roles"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-ram-service-roles"
description: |-
    Provides a list of RAM Service roles to the user.
---

# alibabacloudstack_ascm_ram_service_roles

This data source provides the ram roles of the current Apsara Stack Cloud user.

## Example Usage

```
data "alibabacloudstack_ascm_ram_service_roles" "role" {
  product = "ECS"
}
output "role" {
  value = data.alibabacloudstack_ascm_ram_service_roles.role.*
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of ram roles IDs.
* `product` - (Optional) A regex string to filter results by their product. valid values - "ECS".
* `description` - (Optional) Description about the ram role.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `roles` - A list of roles. Each element contains the following attributes:
    * `id` - ID of the role.
    * `name` - role name.
    * `description` - Description about the role.
    * `role_type` - types of role.
    * `product` - types of role.
    * `organization_name` - Name of an Organization.
    * `aliyun_user_id` - Aliyun User Id.