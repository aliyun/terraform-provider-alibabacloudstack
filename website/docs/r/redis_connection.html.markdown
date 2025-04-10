---
subcategory: "Redis"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_redis_connection"
sidebar_current: "docs-Alibabacloudstack-redis-connection"
description: |- 
  Provides a redis Connection resource.
---

# alibabacloudstack_redis_connection
-> **NOTE:** Alias name has: `alibabacloudstack_kvstore_connection`

Provides a redis Connection resource.

## Example Usage

Basic Usage:

```terraform
variable "name" {
    default = "tf-testaccredisconnection21482"
}

resource "alibabacloudstack_redis_connection" "default" {
  connection_string_prefix = "testprefix"
  instance_id              = "r-8vb6ces3yk5huhxoek"
  port                     = "6379"
}
```

## Argument Reference

The following arguments are supported:

* `connection_string_prefix` - (Required) The prefix of the connection string. The prefix can be 8 to 64 characters in length, and can contain lowercase letters and digits. It must start with a lowercase letter.
* `instance_id` - (Required, ForceNew) The ID of the Redis instance. Once set, this value cannot be changed.
* `port` - (Required) The service port number of the Redis instance. The valid range is from `1024` to `65535`.
* `connection_string` - (Optional) The connection string of the Redis instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Redis instance.
* `connection_string` - The connection string of the Redis instance.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when creating the Redis connection (until it reaches the initial `Normal` status).
* `update` - (Defaults to 2 mins) Used when updating the Redis connection (until it reaches the initial `Normal` status).
* `delete` - (Defaults to 2 mins) Used when deleting the Redis connection (until it reaches the initial `Normal` status).

## Import

Redis connection can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_redis_connection.example r-abc12345678
```