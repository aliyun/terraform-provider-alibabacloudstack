---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_reservedinstance"
sidebar_current: "docs-Alibabacloudstack-ecs-reservedinstance"
description: |- 
  Provides a ecs Reservedinstance resource.
---

# alibabacloudstack_ecs_reservedinstance
-> **NOTE:** Alias name has: `alibabacloudstack_reserved_instance`

Provides a ecs Reservedinstance resource.

## Example Usage

```hcl
variable "name" {
    default = "tf-testaccecsreserved_instance37879"
}

resource "alibabacloudstack_ecs_reservedinstance" "default" {
  instance_type      = "ecs.t6-c4m1.large"
  instance_amount    = 1
  period_unit        = "Year"
  offering_type      = "All Upfront"
  reserved_instance_name = var.name
  description        = "ReservedInstance for testing"
  zone_id           = "cn-hangzhou-i"
  scope             = "Zone"
  period            = 1
  platform          = "Linux"
  resource_group_id = "rg-acfm5xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `instance_type` - (Required, ForceNew) The specifications of the matching pay-as-you-go instance. For example, `ecs.t6-c4m1.large`.
* `scope` - (Optional, ForceNew) Scope of the RI. Optional values: `Region`: region-level, `Zone`: zone-level. Default is `Region`.
* `zone_id` - (Optional, ForceNew) The zone ID to which the RI belongs. When `scope` is set to `Zone`, this parameter is required. For information about the zone list, see [DescribeZones](https://www.alibabacloud.com/help/doc-detail/25610.html).
* `instance_amount` - (Optional, ForceNew) Number of instances allocated to an RI. This represents the number of pay-as-you-go instances of the same specification that can be matched.
* `platform` - (Optional, ForceNew) The operating system type of the image used by the instance. Possible values:
  * `Windows`: An operating system of the Windows Server type.
  * `Linux`: Linux and Unix-like operating systems.
* `period_unit` - (Optional, ForceNew) The unit of time used to purchase reserved instance coupons. Value range:
  * International regions: `Year`
  * China regions: `Year`, `Month`
  Default value: `Month` in China regions and `Year` in international regions.
* `period` - (Optional, ForceNew) The duration of the purchase of reserved instance coupons. Value range:
  * When `PeriodUnit` is `Year`, the values range: `1`, `3`, `5`.
  * When `PeriodUnit` is `Month`, the value range is `1`.
  Default value: `1`.
* `offering_type` - (Optional, ForceNew) Payment type of the RI. Optional values:
  * `No Upfront`: No upfront payment is required.
  * `Partial Upfront`: A portion of upfront payment is required.
  * `All Upfront`: Full upfront payment is required.
* `reserved_instance_name` - (Optional) Name of the RI. The name must be a string of 2 to 128 characters in length and can contain letters, numbers, colons (`:`), underscores (`_`), and hyphens. It must start with a letter. It cannot start with `http://` or `https://`.
* `description` - (Optional) Description of the RI. 2 to 256 English or Chinese characters. It cannot start with `http://` or `https://`.
* `resource_group_id` - (Optional, ForceNew) Resource group ID.
* `reserved_instance_id` - (Optional, ForceNew) The ID of the reserved instance.
* `name` - (Optional, Deprecated) Name of the Reserved Instance.

### Removing alibabacloudstack_ecs_reservedinstance from your configuration

The alibabacloudstack_ecs_reservedinstance resource allows you to manage your Reserved Instance, but Terraform cannot destroy it. Removing this resource from your configuration will remove it from your statefile and management, but will not destroy the Reserved Instance.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - ID of the Reserved Instance.
* `instance_amount` - You can match the number of pay-as-you-go instances of the same specification.
* `platform` - The operating system type of the image used by the instance. Possible values:
  * `Windows`: An operating system of the Windows Server type.
  * `Linux`: Linux and Unix-like operating systems.
* `reserved_instance_name` - Name of the Reserved Instance.
* `resource_group_id` - Resource group ID.
* `reserved_instance_id` - The ID of the reserved instance.