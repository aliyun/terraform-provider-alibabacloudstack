---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_dedicatedhosts"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-dedicatedhosts"
description: |- 
  Provides a list of ecs dedicatedhosts owned by an alibabacloudstack account.
---

# alibabacloudstack_ecs_dedicatedhosts
-> **NOTE:** Alias name has: `alibabacloudstack_ecs_dedicated_hosts`

This data source provides a list of ECS Dedicated Hosts in an Apsara Stack Cloud account according to the specified filters.

## Example Usage

```hcl
# Declare the data source
data "alibabacloudstack_ecs_dedicated_hosts" "dedicated_hosts_ds" {
  name_regex = "tf-testAcc"
  dedicated_host_type = "ddh.g5"
  status = "Available"
}

output "first_dedicated_hosts_id" {
  value = "${data.alibabacloudstack_ecs_dedicated_hosts.dedicated_hosts_ds.hosts.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional, ForceNew) A regex string to filter results by the ECS Dedicated Host name.
* `ids` - (Optional, ForceNew) A list of ECS Dedicated Host ids.
* `dedicated_host_id` - (Optional, ForceNew) The ID of the dedicated host.
* `dedicated_host_name` - (Optional, ForceNew) The name of the dedicated host. The name must be 2 to 128 characters in length. The name must start with a letter and cannot start with `http://` or `https://`. The name can contain letters, digits, colons (`:`), underscores (`_`), and hyphens (`-`).
* `dedicated_host_type` - (Optional, ForceNew) The type of the dedicated host. You can call the [DescribeDedicatedHostTypes](https://www.alibabacloud.com/help/en/doc-detail/134240.html) operation to query the most recent list of dedicated host types.
* `operation_locks` - (Optional, ForceNew) OperationLocks.
  * `lock_reason` - (Optional, ForceNew) The reason why the dedicated host resource is locked.
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group to which the dedicated host belongs. When this parameter is specified to query resources, up to 1,000 resources that belong to the specified resource group can be displayed in the response. Resources in the default resource group are displayed in the response regardless of how this parameter is set.
* `status` - (Optional, ForceNew) The service state of the dedicated host. Valid values:
  * `Available`: The dedicated host is running normally.
  * `UnderAssessment`: The dedicated host is available but has potential risks that may cause the ECS instances on the dedicated host to fail.
  * `PermanentFailure`: The dedicated host encounters permanent failures and is unavailable.
  * `TempUnavailable`: The dedicated host is temporarily unavailable.
  * `Redeploying`: The dedicated host is being restored.
  Default value: `Available`.
* `zone_id` - (Optional, ForceNew) The zone ID of the dedicated host. You can call the [DescribeZones](https://www.alibabacloud.com/help/en/doc-detail/25610.html) operation to query the most recent zone list.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of ECS Dedicated Host ids.
* `names` - A list of ECS Dedicated Host names.
* `hosts` - A list of ECS Dedicated Hosts. Each element contains the following attributes:
  * `id` - ID of the ECS Dedicated Host.
  * `action_on_maintenance` - The policy used to migrate the instances from the dedicated host when the dedicated host fails or needs to be repaired online. Valid values:
    * `Migrate`: The instances are migrated to another physical server and restarted.
    * `Stop`: The instances are stopped. If the dedicated host cannot be repaired, the instances are migrated to another physical machine and then restarted.
  * `auto_placement` - Specifies whether to add the dedicated host to the resource pool for automatic deployment. Valid values:
    * `on`: Adds the dedicated host to the resource pool for automatic deployment.
    * `off`: Does not add the dedicated host to the resource pool for automatic deployment.
  * `auto_release_time` - The automatic release time of the dedicated host. Specify the time in the ISO 8601 standard in the `yyyy-MM-ddTHH:mm:ssZ` format. The time must be in UTC+0.
  * `capacity` - Capacity.
    * `available_local_storage` - The remaining local disk capacity. Unit: GiB.
    * `available_memory` - The remaining memory capacity, unit: GiB.
    * `available_vcpus` - The number of remaining vCPU cores.
    * `available_vgpus` - The number of available virtual GPUs.
    * `local_storage_category` - Local disk type.
    * `total_local_storage` - The total capacity of the local disk, in GiB.
    * `total_memory` - The total memory capacity, unit: GiB.
    * `total_vcpus` - The total number of vCPU cores.
    * `total_vgpus` - The total number of virtual GPUs.
  * `cores` - Cores.
  * `cpu_over_commit_ratio` - The CPU overcommit ratio. You can configure CPU overcommit ratios only for the following dedicated host types: g6s, c6s, and r6s. Valid values: 1 to 5.
  * `dedicated_host_id` - ID of the ECS Dedicated Host.
  * `dedicated_host_name` - The name of the dedicated host.
  * `dedicated_host_type` - The type of the dedicated host.
  * `description` - The description of the dedicated host.
  * `expired_time` - The expiration time of the subscription dedicated host.
  * `gpu_spec` - The GPU model.
  * `machine_id` - The machine code of the dedicated host.
  * `network_attributes` - NetworkAttributes.
    * `slb_udp_timeout` - The timeout period for a UDP session between Server Load Balancer (SLB) and the dedicated host. Unit: seconds.
    * `udp_timeout` - The timeout period for a UDP session between a user and an Apsara Stack Cloud service on the dedicated host. Unit: seconds.
  * `operation_locks` - OperationLocks.
    * `lock_reason` - The reason why the dedicated host resource is locked.
  * `payment_type` - The billing method of the dedicated host.
  * `physical_gpus` - The number of physical GPUs.
  * `resource_group_id` - The ID of the resource group to which the dedicated host belongs.
  * `sale_cycle` - The unit of the subscription billing method.
  * `sockets` - The number of physical CPUs.
  * `status` - The service status of the dedicated host.
  * `supported_custom_instance_type_families` - Custom instance type families supported by dedicated hosts.
  * `supported_instance_type_families` - SupportedInstanceTypeFamilies.
  * `supported_instance_types_list` - SupportedInstanceTypesList.
  * `tags` - Tags.
  * `zone_id` - The zone ID of the ECS Dedicated Host.