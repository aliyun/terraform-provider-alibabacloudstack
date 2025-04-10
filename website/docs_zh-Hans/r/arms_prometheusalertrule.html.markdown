---
subcategory: "应用实时监控服务 (ARMS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_arms_prometheus_alert_rule"
sidebar_current: "docs-alibabacloudstack-resource-arms-prometheus-alert-rule"
description: |-
  编排应用实时监控服务 (ARMS) Prometheus 告警规则
---

# alibabacloudstack_arms_prometheus_alert_rule
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_arms_prometheusalertrule`

使用Provider配置的凭证在指定的资源集下编排应用实时监控服务 (ARMS) Prometheus 告警规则资源。

有关 ARMS Prometheus 告警规则及其使用方法的信息，请参阅 [什么是 Prometheus 告警规则](https://www.alibabacloud.com/help/zh/doc-detail/212056.htm)。

## 示例用法

### 基础用法

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

## 参数说明

以下参数是支持的：

* `annotations` - (可选) 告警规则的注解。详见 [`annotations`](#annotations) 下面。
  * `name` - (可选) 注解的名称。
  * `value` - (可选) 注解的值。
* `cluster_id` - (必填, 变更时重建) 集群的 ID。
* `dispatch_rule_id` - (可选) 通知策略的 ID。当 `notify_type` 参数设置为 `DISPATCH_RULE` 时，此参数是必填的。
* `duration` - (必填, 变更时重建) 告警的持续时间（单位：分钟）。
* `expression` - (必填, 变更时重建) 符合 PromQL 语法的告警规则表达式。
* `labels` - (可选) 资源的标签。详见 [`labels`](#labels) 下面。
  * `name` - (可选) 标签的名称。
  * `value` - (可选) 标签的值。
* `message` - (必填, 变更时重建) 告警通知的消息内容。
* `notify_type` - (可选) 发送告警通知的方法。有效值：`ALERT_MANAGER`, `DISPATCH_RULE`。
* `prometheus_alert_rule_name` - (必填, 变更时重建) 告警规则的名称。
* `status` - (可选, 变更时重建) 告警规则的状态。有效值：`0`（禁用），`1`（启用）。
* `type` - (可选, 变更时重建) 告警规则的类型。


## 属性说明

以下属性被导出：

* `id` - Prometheus 告警规则的资源 ID。格式为 `<cluster_id>:<prometheus_alert_rule_id>`。
* `prometheus_alert_rule_id` - 告警规则的唯一标识符。
* `status` - 告警规则的状态。有效值：`0`（禁用），`1`（启用）。