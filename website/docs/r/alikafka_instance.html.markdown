---
subcategory: "AliKafka"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_alikafka_instance"
sidebar_current: "docs-Alibabacloudstack-resource-alikafka-instance"
description: |-
  Provides a AliKafka Instance resource.
---

# alibabacloudstack_alikafka_instance

Provides an AliKafka instance resource.

## Example Usage

Basic Usage


```terraform
data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}



resource "alibabacloudstack_alikafka_instance" "default" {
  sasl = "true"
  plaintext = "true"
  spec = "Broker4C16G"
  cup_type = "Intel"
  name = "tf-testacc-alikafkainstancebasic14494"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, ForceNew) Name of your Kafka instance. The length should between 3 and 64 characters. If not set, will use instance id as instance name.
* `zone_id` - (Optional, ForceNew) The zone ID of the instance. The value can be in zone x or region id-x format. **NOTE**: When the available zone is insufficient, another availability zone may be deployed.
* `selected_zones` - (Optional, ForceNew, List) The zones among which you want to deploy the instance.
* `cpu_type` - (Required, ForceNew) The CPU type of the resource. Valid values: `intel`.
* `spec_type` - (Optional, ForceNew) The spec type of the instance.
* `replicas` - (Optional, ForceNew) Number of brokers.
* `disk_num` - (Optional, ForceNew) Number of disks per Broker.
* `vpc_id` - (Optional, ForceNew) VPC ID.
* `vswitch_id` - (Optional, ForceNew) Vswtich ID.
* `sasl` - (Optional, ForceNew) Enable SASL Access Point Type.
* `plaintext` - (Optional, ForceNew) Enable PLAINTEXT Access Point Type. PLAINTEXT type can be accessed without authentication, please pay attention to security protection.
* `message_max_bytes` - (Optional) The maximum size of messages that the server can receive. It is important that consumer and producer settings related to this attribute must be synchronized, otherwise the message published by producer is too large for consumer. The message size limit is not set as large as possible, but depends on the specific business system.
* `num_partitions` - (Optional) If the number of divided partitions is not given when creating a topic, this number will be the default number of partitions under the topic.
* `auto_create_topics_enable` - Whether to allow automatic topic creation. If true, this topic is automatically created when a message is sent or a topic that does not exist in fetch. Otherwise, you need to use the command line to create a topic.
* `num_io_threads` - (Optional) The number of I/O threads used by server to process requests; The number of threads must be at least equal to the number of hard disks.
* `queued_max_requests` - (Optional) The maximum number of requests that can be queued for I/O threads to process before network threads stop reading new requests.
* `replica_fetch_wait_max_ms` - (Optional) Replicas the maximum waiting time for communication with the leader. If it fails, it will be retried.
* `replica_lag_time_max_ms` - (Optional) If a follower does not send a fetch request within this time, the leader removes the follower from the ISR and considers the follower offline.
* `num_network_threads` - (Optional) The number of network threads used by server to process network requests
* `log_retention_bytes` - (Optional) The number of bytes saved for each partition under each topic. Note that this is the upper limit of each partition, so the number multiplied by partition is the total amount of data saved by each topic. Note that if both log.retention.hours and log.retention.bytes are set, a segment file will be deleted if any limit is exceeded. Note that this setting can be overwritten when each topic is set.
* `replica_fetch_max_bytes` - (Optional) Follower the maximum number of bytes per fetch when backing up from the leader.
* `num_replica_fetchers` - (Optional) Number of threads backing up data from leader
* `default_replication_factor` - (Optional) The default number of backups. Only topics created automatically
* `offsets_retention_minutes` - (Optional) Offsets that have exceeded this time limit will be marked as to be deleted.
* `background_threads` - (Optional) The number of threads used for background processing, such as file deletion.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Instance.
* `sasl_plaintext_endpoint` - The SSL endpoint of the instance in domain name mode.
* `sasl_ssl_endpoint` - The Simple Authentication and Security Layer (SASL) endpoint of the instance in domain name mode.
* `status` - The status of the instance.

## Timeouts

* `create` - (Defaults to 60 mins) Used when create the resource.
* `update` - (Defaults to 120 mins) Used when update the resource.
* `delete` - (Defaults to 30 mins) Used when delete the resource.

## Import

AliKafka instance can be imported using the id, e.g.

```shell
$ terraform import alibabacloudstack_alikafka_instance.instance <id>
```
