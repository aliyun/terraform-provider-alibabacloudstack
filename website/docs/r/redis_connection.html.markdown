---
subcategory: "ApsaraDB for Redis"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_redis_connection"
sidebar_current: "docs-Alibabacloudstack-redis-connection"
description: |- 
  Provides a Redis Connection resource.
---

# alibabacloudstack_redis_connection

Provides a Redis Connection resource.

> **Note:** This resource can also be referred to by the following aliases:
> - `alibabacloudstack_kvstore_connection`

## Example Usage

Basic Usage:

```terraform
variable "name" {
    default = "tf-testaccredisconnection"
}

resource "alibabacloudstack_redis_connection" "default" {
  connection_string_prefix = var.name
  instance_id              = local.instance_id
  port                     = "6379"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the Redis instance. Modifying this parameter will force the resource to be recreated.
* `connection_string_prefix` - (Required) The prefix of the connection string. The prefix can be 8 to 64 characters in length and can contain lowercase letters and digits. It must start with a lowercase letter.
* `port` - (Required) The service port number of the Redis instance. Valid values: `1024` to `65535`.
* `connection_string` - (Optional, Computed) The connection string of the Redis instance. This attribute is returned by the API and cannot be manually set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Redis instance, same as `instance_id`.
* `connection_string` - The complete connection string of the Redis instance.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when creating the Redis connection (until it reaches the initial `Normal` status).
* `update` - (Defaults to 2 mins) Used when updating the Redis connection (until it reaches the initial `Normal` status).
* `delete` - (Defaults to 2 mins) Used when deleting the Redis connection (until it reaches the initial `Normal` status).

## Import

Redis connection can be imported using the instance ID, e.g.

```bash
$ terraform import alibabacloudstack_redis_connection.example r-abc12345678
```
