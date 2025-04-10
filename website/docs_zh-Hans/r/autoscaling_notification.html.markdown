---
subcategory: "AutoScaling"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_autoscaling_notification"
sidebar_current: "docs-Alibabacloudstack-autoscaling-notification"
description: |- 
  编排弹性伸缩的消息通知
---

# alibabacloudstack_autoscaling_notification
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_ess_notification`

使用Provider配置的凭证在指定的资源集下编排弹性伸缩的消息通知。

## 示例用法

```hcl
variable "name" {
    default = "tf-testAccAutoscalingNotification-%d"
}

data "alibabacloudstack_regions" "default" {
    current = true
}

data "alibabacloudstack_account" "default" {
}

data "alibabacloudstack_zones" "default" {
    available_disk_category     = "cloud_efficiency"
    available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
    name       = "${var.name}"
    cidr_block = "172.16.0.0/16"
}
    
resource "alibabacloudstack_vswitch" "default" {
    vpc_id            = "${alibabacloudstack_vpc.default.id}"
    cidr_block        = "172.16.0.0/24"
    availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
    name              = "${var.name}"
}

resource "alibabacloudstack_autoscaling_scaling_group" "default" {
    min_size = 1
    max_size = 1
    scaling_group_name = "${var.name}"
    removal_policies = ["OldestInstance", "NewestInstance"]
    vswitch_ids = ["${alibabacloudstack_vswitch.default.id}"]
}

resource "alibabacloudstack_mns_queue" "default"{
    name="${var.name}"
}

resource "alibabacloudstack_autoscaling_notification" "default" {
    scaling_group_id = "${alibabacloudstack_autoscaling_scaling_group.default.id}"
    notification_types = [
        "AUTOSCALING:SCALE_OUT_SUCCESS",
        "AUTOSCALING:SCALE_OUT_ERROR"
    ]
    notification_arn = "acs:ess:${data.alibabacloudstack_regions.default.regions.0.id}:${data.alibabacloudstack_account.default.id}:queue/${alibabacloudstack_mns_queue.default.name}"
}
```

## 参数说明

支持以下参数：

* `scaling_group_id` - (必填，变更时重建) 自动伸缩组的ID。更改此参数将强制创建新资源。
* `notification_arn` - (必填，变更时重建) 通知对象的阿里云资源名称(ARN)。`notification_arn` 的格式为 `acs:ess:{region}:{account-id}:{resource-relative-id}`。其中：
  * `{region}` 是区域标识符。
  * `{account-id}` 是用户的阿里云账户ID。
  * `{resource-relative-id}` 是通知目标的相对ID，有效值包括：
    * `cloudmonitor`: 用于云监控通知。
    * `queue/{queue-name}`: 用于消息队列(MNS)基于队列的通知。
    * `topic/{topic-name}`: 用于消息队列(MNS)基于主题的通知。
* `notification_types` - (必填) 一种或多种自动伸缩事件和资源变更通知类型。支持的通知类型包括：
  * `AUTOSCALING:SCALE_OUT_SUCCESS`
  * `AUTOSCALING:SCALE_IN_SUCCESS`
  * `AUTOSCALING:SCALE_OUT_ERROR`
  * `AUTOSCALING:SCALE_IN_ERROR`
  * `AUTOSCALING:SCALE_REJECT`
  * `AUTOSCALING:SCALE_OUT_START`
  * `AUTOSCALING:SCALE_IN_START`
  * `AUTOSCALING:SCHEDULE_TASK_EXPIRING`

## 属性说明

除了上述所有参数外，还导出以下属性：

* `id` - 通知资源的唯一标识符，由 `scaling_group_id` 和 `notification_arn` 组成，格式为 `<scaling_group_id>:<notification_arn>`。