---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_arms_prometheus_alert_rule"
sidebar_current: "docs-alibabacloudstack-resource-arms-prometheus-alert-rule"
description: |-
  Provides a Alibabacloudstack Application Real-Time Monitoring Service (ARMS) Prometheus Alert Rule resource.
---

# alibabacloudstack_arms_prometheus_alert_rule
-> **NOTE:** Alias name has: `alibabacloudstack_arms_prometheusalertrule`

Provides a Application Real-Time Monitoring Service (ARMS) Prometheus Alert Rule resource.

For information about Application Real-Time Monitoring Service (ARMS) Prometheus Alert Rule and how to use it, see [What is Prometheus Alert Rule](https://www.alibabacloud.com/help/en/doc-detail/212056.htm).

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alibabacloudstack_vpc.default.id
  zone_id      = data.alibabacloudstack_zones.default.zones.0.id
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone    = alibabacloudstack_vswitch.default.zone_id
  cpu_core_count       = 2
  memory_size          = 4
  kubernetes_node_role = "Worker"
  instance_type_family = "ecs.sn1ne"
}

resource "alibabacloudstack_cs_managed_kubernetes" "default" {
  name_prefix          = var.name
  cluster_spec         = "ack.pro.small"
  worker_vswitch_ids   = [alibabacloudstack_vswitch.default.id]
  new_nat_gateway      = true
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = true
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alibabacloudstack_key_pair" "default" {
  key_pair_name = "${var.name}-${random_integer.default.result}"
}

resource "alibabacloudstack_cs_kubernetes_node_pool" "default" {
  name                 = "desired_size"
  cluster_id           = alibabacloudstack_cs_managed_kubernetes.default.id
  vswitch_ids          = [alibabacloudstack_vswitch.default.id]
  instance_types       = [data.alibabacloudstack_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alibabacloudstack_key_pair.default.key_pair_name
  desired_size         = 2
}

resource "alibabacloudstack_arms_prometheus" "default" {
  cluster_type        = "aliyun-cs"
  grafana_instance_id = "free"
  cluster_id          = alibabacloudstack_cs_kubernetes_node_pool.default.cluster_id
}

resource "alibabacloudstack_arms_prometheus_alert_rule" "example" {
  cluster_id                 = alibabacloudstack_cs_managed_kubernetes.default.id
  duration                   = 1
  expression                 = "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10"
  message                    = "node available memory is less than 10%"
  prometheus_alert_rule_name = var.name
  notify_type                = "DISPATCH_RULE"
}
```

## Argument Reference

The following arguments are supported:

* `annotations` - (Optional) The annotations of the alert rule. See [`annotations`](#annotations) below.
  * `name` - (Optional) The name of the annotation.
  * `value` - (Optional) The value of the annotation.
* `cluster_id` - (Required, ForceNew) The ID of the cluster.
* `dispatch_rule_id` - (Optional) The ID of the notification policy. This parameter is required when the `notify_type` parameter is set to `DISPATCH_RULE`.
* `duration` - (Required, ForceNew) The duration of the alert.
* `expression` - (Required, ForceNew) The alert rule expression that follows the PromQL syntax.
* `labels` - (Optional) The labels of the resource. See [`labels`](#labels) below.
  * `name` - (Optional) The name of the label.
  * `value` - (Optional) The value of the label.
* `message` - (Required, ForceNew) The message of the alert notification.
* `notify_type` - (Optional) The method of sending the alert notification. Valid values: `ALERT_MANAGER`, `DISPATCH_RULE`.
* `prometheus_alert_rule_name` - (Required, ForceNew) The name of the resource.
* `status` - (Optional, ForceNew) The status of the resource.
* `type` - (Optional, ForceNew) The type of the alert rule.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Prometheus Alert Rule. The value formats as `<cluster_id>:<prometheus_alert_rule_id>`.
* `prometheus_alert_rule_id` - The first ID of the resource.
* `status` -  The status of the resource. Valid values: `0`, `1`.