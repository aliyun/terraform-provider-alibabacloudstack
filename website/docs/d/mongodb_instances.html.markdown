---
subcategory: "MongoDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mongodb_instances"
sidebar_current: "docs-Alibabacloudstack-datasource-mongodb-instances"
description: |- 
  Provides a list of MongoDB instances owned by an Alibabacloudstack account.
---

# alibabacloudstack_mongodb_instances

This data source provides a list of MongoDB instances in an Alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_mongodb_instances" "mongo" {
  name_regex        = "dds-.+\\d+"
  instance_type     = "replicate"
  instance_class    = "dds.mongo.mid"
  availability_zone = "eu-central-1a"
}

output "mongodb_instance_ids" {
  value = data.alibabacloudstack_mongodb_instances.mongo.ids
}

output "mongodb_instance_names" {
  value = data.alibabacloudstack_mongodb_instances.mongo.names
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to apply to the instance name. This allows filtering instances based on their names using regular expressions.
* `ids` - (Optional, Available in v1.53.0+) The list of MongoDB instance IDs. Use this parameter to filter results by specific instance IDs.
* `instance_type` - (Optional) Type of the instance to be queried. If set to `sharding`, the sharded cluster instances are listed. If set to `replicate`, replica set instances are listed. Default value is `replicate`.
* `instance_class` - (Optional) Sizing of the instance to be queried. This corresponds to the performance class of the MongoDB instance.
* `availability_zone` - (Optional) Instance availability zone. Use this parameter to filter results by a specific availability zone.
* `tags` - (Optional, Available in v1.66.0+) A mapping of tags to assign to the resource. Use this parameter to filter results by tags.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of names of all the matched MongoDB instances.
* `instances` - A list of MongoDB instances. Each element contains the following attributes:
  * `id` - The ID of the MongoDB instance.
  * `name` - The name of the MongoDB instance.
  * `charge_type` - Billing method. Value options are `PostPaid` for Pay-As-You-Go and `PrePaid` for yearly or monthly subscription.
  * `instance_type` - Instance type. Optional values are `sharding` for sharded clusters or `replicate` for replica sets.
  * `region_id` - Region ID the instance belongs to.
  * `creation_time` - Creation time of the instance in RFC3339 format.
  * `expiration_time` - Expiration time in RFC3339 format. Pay-As-You-Go instances do not expire.
  * `status` - Status of the instance.
  * `replication` - Replication factor corresponding to the number of nodes. Optional values are `1` for single node and `3` for three-node replica sets.
  * `engine` - Database engine type. Supported option is `MongoDB`.
  * `engine_version` - Database engine version.
  * `network_type` - Network type. Options include classic network or VPC.
  * `lock_mode` - Lock status of the instance.
  * `instance_class` - Sizing of the MongoDB instance.
  * `storage` - Storage size in GB.
  * `mongos` - An array composed of Mongos nodes. Each element contains:
    * `node_id` - Mongos instance ID.
    * `description` - Mongos instance description.
    * `class` - Mongos instance specification.
  * `shards` - An array composed of shards. Each element contains:
    * `node_id` - Shard instance ID.
    * `description` - Shard instance description.
    * `class` - Shard instance specification.
    * `storage` - Shard disk size in GB.
  * `availability_zone` - Instance availability zone.
  * `tags` - A mapping of tags assigned to the resource.