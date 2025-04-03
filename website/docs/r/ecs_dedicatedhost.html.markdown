---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_dedicated_host"
sidebar_current: "docs-alibabacloudstack-resource-ecs-dedicated-host"
description: |- 
  Provides a Alibabacloudstack ecs dedicated host resource.
---

# alibabacloudstack_ecs_dedicated_host
-> **NOTE:** Alias name has: `alibabacloudstack_ecs_dedicatedhost`

This resource is used to create and manage an ECS Dedicated Host in Alibaba Cloud. A Dedicated Host allows you to deploy your instances on a physical server that is fully dedicated to you, providing better isolation and control over hardware resources.

## Example Usage

### Basic Usage

```hcl
resource "alibabacloudstack_ecs_dedicated_host" "default" {
  dedicated_host_type = "ddh.g5"
  payment_type        = "PostPaid"
  tags = {
    Create = "Terraform"
    For    = "DDH"
  }
  description         = "From_Terraform"
  dedicated_host_name = "dedicated_host_name"
}
```

### Prepaid Dedicated Host

```hcl
resource "alibabacloudstack_ecs_dedicated_host" "prepaid" {
  dedicated_host_type = "ddh.g5"
  payment_type        = "PrePaid"
  auto_renew          = true
  auto_renew_period   = 1
  sale_cycle          = "Month"
  expired_time        = 12
  tags = {
    Create = "Terraform"
    For    = "DDH"
  }
  description         = "Prepaid_DDH"
  dedicated_host_name = "prepaid_dedicated_host"
}
```

### Custom CPU Overcommit Ratio

```hcl
resource "alibabacloudstack_ecs_dedicated_host" "custom_cpu" {
  dedicated_host_type = "ddh.c6s"
  payment_type        = "PostPaid"
  cpu_over_commit_ratio = 4
  tags = {
    Create = "Terraform"
    For    = "DDH"
  }
  description         = "Custom_CPU_Ratio"
  dedicated_host_name = "custom_cpu_ratio_host"
}
```

## Argument Reference

The following arguments are supported:

* `action_on_maintenance` - (Optional) The policy used to migrate the instances from the dedicated host when the dedicated host fails or needs to be repaired online. Valid values:
  * `Migrate`: Instances are migrated to another physical server and restarted.
  * `Stop`: Instances are stopped. If the dedicated host cannot be repaired, instances are migrated to another physical machine and then restarted.
  Default value: Depends on the type of disk attached (`Migrate` for cloud disks, `Stop` for local disks).

* `auto_placement` - (Optional) Specifies whether to add the dedicated host to the resource pool for automatic deployment. If you do not specify the `DedicatedHostId` parameter when creating an instance, Alibaba Cloud automatically selects a dedicated host from the resource pool. Valid values:
  * `on`: Adds the dedicated host to the resource pool for automatic deployment.
  * `off`: Does not add the dedicated host to the resource pool for automatic deployment.
  Default value: `on`.

* `auto_release_time` - (Optional) The automatic release time of the dedicated host. Specify the time in the ISO 8601 standard in the `yyyy-MM-ddTHH:mm:ssZ` format. The time must be in UTC+0.

* `auto_renew` - (Optional) Specifies whether to automatically renew the subscription dedicated host. Default value: `false`.

* `auto_renew_period` - (Optional) The auto-renewal period of the dedicated host. Unit: months. Valid values: `1`, `2`, `3`, `6`, and `12`. This parameter takes effect and is required only when the `AutoRenew` parameter is set to `true`.

* `cpu_over_commit_ratio` - (Optional) The CPU overcommit ratio. You can configure this only for the following dedicated host types: `g6s`, `c6s`, and `r6s`. Valid values: `1` to `5`.

* `dedicated_host_cluster_id` - (Optional) The ID of the dedicated host cluster to which the dedicated host belongs.

* `dedicated_host_name` - (Optional) The name of the dedicated host. The name must be 2 to 128 characters in length. It must start with a letter but cannot start with `http://` or `https://`. It can contain letters, digits, colons (`:`), underscores (`_`), and hyphens (`-`).

* `dedicated_host_type` - (Required, ForceNew) The type of the dedicated host. You can call the [DescribeDedicatedHostTypes](https://www.alibabacloud.com/help/doc-detail/134240.htm) operation to obtain the most recent list of dedicated host types.

* `description` - (Optional) The description of the dedicated host. The description must be 2 to 256 characters in length and cannot start with `http://` or `https://`.

* `dry_run` - (Optional) Specifies whether to only validate the request. Default value: `false`.

* `expired_time` - (Optional) The subscription period of the dedicated host. This parameter takes effect and is required only when the `PaymentType` parameter is set to `PrePaid`.

* `min_quantity` - (Optional) The minimum number of dedicated hosts to create. Valid values: `1` to `100`.

* `network_attributes` - (Optional) Network attributes for the dedicated host. Contains the following attributes:
  * `slb_udp_timeout` - (Optional) The timeout period for a UDP session between Server Load Balancer (SLB) and the dedicated host. Unit: seconds. Valid values: `15` to `310`.
  * `udp_timeout` - (Optional) The timeout period for a UDP session between a user and an Alibaba Cloud service on the dedicated host. Unit: seconds. Valid values: `15` to `310`.

* `payment_type` - (Optional) The billing method of the dedicated host. Valid values:
  * `PrePaid`
  * `PostPaid`
  Default value: `PostPaid`.

* `resource_group_id` - (Optional) The ID of the resource group to which the dedicated host belongs.

* `sale_cycle` - (Optional) The unit of the subscription period of the dedicated host. Valid values:
  * `Month`
  * `Year`
  Default value: `Month`.

* `zone_id` - (Optional, ForceNew) The zone ID of the dedicated host. This parameter is empty by default. If you do not specify this parameter, the system automatically selects a zone.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the dedicated host.
* `status` - The service state of the dedicated host. Valid values:
  * `Available`: The dedicated host is running normally.
  * `UnderAssessment`: The dedicated host is available but has potential risks that may cause the ECS instances on the dedicated host to fail.
  * `PermanentFailure`: The dedicated host encounters permanent failures and is unavailable.
  * `TempUnavailable`: The dedicated host is temporarily unavailable.
  * `Redeploying`: The dedicated host is being restored.
  Default value: `Available`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 11 mins) Used when creating the dedicated host.
* `delete` - (Defaults to 1 min) Used when deleting the dedicated host.
* `update` - (Defaults to 11 mins) Used when updating the dedicated host.

## Import

ECS dedicated host can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_ecs_dedicated_host.default dh-2zedmxxxx
```