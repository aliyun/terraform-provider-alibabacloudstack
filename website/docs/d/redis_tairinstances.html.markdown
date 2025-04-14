---
subcategory: "Redis"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_redis_tairinstances"
sidebar_current: "docs-Alibabacloudstack-datasource-redis-tairinstances"
description: |- 
  Provides a list of redis tairinstances owned by an Alibabacloudstack account.

---

# alibabacloudstack_redis_tairinstances
-> **NOTE:** Alias name has: `alibabacloudstack_kvstore_instances`

This data source provides a list of Redis Tair instances in an Alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_redis_tairinstances" "default" {
    name_regex = "checkalibabacloudstacktairinstancesdatasource"
    status      = "Running"
    instance_type = "tair_rdb"
}

output "first_instance_name" {
    value = data.alibabacloudstack_redis_tairinstances.default.instances.0.name
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to apply to the Tair instance name.
* `ids` - (Optional) A list of Tair instance IDs.
* `status` - (Optional) The status of the resource. Valid values include: `Creating`, `Running`, `Restarting`, `ChangingConfig`, `FlushingData`, `Deleting`, `NetworkChanging`, `Abnormal`.
* `instance_type` - (Optional) The storage medium of the instance. Valid values: `tair_rdb`, `tair_scm`, `tair_essd`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Tair instance names.
* `ids` - A list of Tair instance IDs.
* `instances` - A list of Tair instances. Each element contains the following attributes:
  * `id` - The ID of the Tair instance.
  * `name` - The name of the Tair instance.
  * `charge_type` - Billing method. Value options: `PostPaid` for Pay-As-You-Go and `PrePaid` for subscription.
  * `region_id` - Region ID the instance belongs to.
  * `create_time` - The time when the instance was created. The time follows the ISO 8601 standard in the `yyyy-MM-ddTHH:mm:ssZ` format. The time is displayed in UTC.
  * `expire_time` - Expiration time. Pay-As-You-Go instances do not have an expiration time.
  * `status` - The status of the resource. Valid values include: `Creating`, `Running`, `Restarting`, `ChangingConfig`, `FlushingData`, `Deleting`, `NetworkChanging`, `Abnormal`.
  * `instance_type` - The storage medium of the instance. Valid values: `tair_rdb`, `tair_scm`, `tair_essd`.
  * `instance_class` - The instance type of the instance. For more information, see [Instance types](https://www.alibabacloud.com/help/en/apsaradb-for-redis/latest/instance-types).
  * `availability_zone` - The availability zone where the instance resides.
  * `vpc_id` - The ID of the virtual private cloud (VPC).
  * `vswitch_id` - The ID of the VSwitch.
  * `private_ip` - The private IP address of the instance.
  * `port` - The Tair service port. Valid values: 1024 to 65535. Default value: 6379.
  * `user_name` - The username of the instance.
  * `capacity` - The storage capacity of the instance. Unit: MB.
  * `bandwidth` - The bandwidth of the instance. Unit: Mbit/s.
  * `connections` - The connection quantity limit of the instance. Unit: count.
  * `connection_domain` - The internal endpoint of the instance.
  * `tags` - A mapping of tags to assign to the resource.