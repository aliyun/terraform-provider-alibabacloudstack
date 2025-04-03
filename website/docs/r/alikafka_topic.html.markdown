---
subcategory: "AliKafka"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_alikafka_topic"
sidebar_current: "docs-Alibabacloudstack-alikafka-topic"
description: |- 
  Provides a alikafka Topic resource.
---

# alibabacloudstack_alikafka_topic

Provides a ALIKAFKA topic resource.



-> **NOTE:**  Only the following regions support create alikafka topic.
[`cn-hangzhou`,`cn-beijing`,`cn-shenzhen`,`cn-shanghai`,`cn-qingdao`,`cn-hongkong`,`cn-huhehaote`,`cn-zhangjiakou`,`cn-chengdu`,`cn-heyuan`,`ap-southeast-1`,`ap-southeast-3`,`ap-southeast-5`,`ap-south-1`,`ap-northeast-1`,`eu-central-1`,`eu-west-1`,`us-west-1`,`us-east-1`]

## Example Usage

```hcl
variable "name" {
  default = "tf-testacc-alikafkatopicbasic12916"
}

resource "alibabacloudstack_vpc" "default" {
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  zone_id           = data.alibabacloudstack_zones.default.zones[0].id
}

resource "alibabacloudstack_alikafka_instance" "default" {
  name        = "tf-testacc-alikafkainstance"
  topic_quota = "50"
  disk_type   = "1"
  disk_size   = "500"
  deploy_type = "5"
  io_max      = "20"
  vswitch_id  = alibabacloudstack_vswitch.default.id
}

resource "alibabacloudstack_alikafka_topic" "default" {
  remark        = "alibabacloudstack_alikafka_topic_remark"
  instance_id   = alibabacloudstack_alikafka_instance.default.id
  topic         = var.name
  local_topic   = true
  compact_topic = false
  partition_num = 12
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) Resource id of your Kafka resource, the topic will be created in this instance.
* `topic` - (Required, ForceNew) Name of the topic. Two topics on a single instance cannot have the same name. The length cannot exceed 64 characters.
* `local_topic` - (Optional, ForceNew) Indicates whether the topic is a local topic or not. Default value is `false`.
* `compact_topic` - (Optional, ForceNew) Indicates whether the topic is a compact topic or not. Compact topic must be a local topic. Default value is `false`.
* `partition_num` - (Optional) The number of partitions of the topic. The number should be between 1 and 48. Default value is `1`.
* `remark` - (Required) A concise description of the topic. The length cannot exceed 64 characters.
* `tags` - (Optional, Available in v1.63.0+) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier of the resource. The value is formulated as `<instance_id>:<topic>`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the topic (until it reaches the initial `Running` status).

## Import

ALIKAFKA TOPIC can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_alikafka_topic.topic alikafka_post-cn-123455abc:topicName
```