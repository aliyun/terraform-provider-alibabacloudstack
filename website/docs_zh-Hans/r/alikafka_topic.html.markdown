---
subcategory: "AliKafka"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_alikafka_topic"
sidebar_current: "docs-Alibabacloudstack-alikafka-topic"
description: |- 
  编排Alikafka主题资源
---

# alibabacloudstack_alikafka_topic

使用Provider配置的凭证在指定的资源集下编排Alikafka主题资源。

## 示例用法

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

## 参数说明

支持以下参数：

* `instance_id` - (必填，变更时重建) Kafka 实例的资源 ID。主题将在此实例中创建。
* `topic` - (必填，变更时重建) 主题的名称。单个实例上的两个主题不能具有相同的名称。长度不得超过 64 个字符。
* `local_topic` - (可选，变更时重建) 指示该主题是否为本地主题。默认值为 `false`。
* `compact_topic` - (可选，变更时重建) 指示该主题是否为紧凑主题。紧凑主题必须是本地主题。默认值为 `false`。
* `partition_num` - (可选) 主题的分区数。数量应在 1 到 48 之间。默认值为 `1`。
* `remark` - (必填) 主题的简要描述。长度不得超过 64 个字符。
* `tags` - (可选，v1.63.0+可用) 分配给资源的标签映射。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 资源的唯一标识符。其值被制定为 `<instance_id>:<topic>`。

### 超时时间

`timeouts` 块允许您为某些操作指定 [超时时间](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts)：

* `create` - (默认为 10 分钟)用于创建主题(直到它达到初始 `Running` 状态)。

## 导入

ALIKAFKA 主题可以使用 id 导入，例如：

```bash
$ terraform import alibabacloudstack_alikafka_topic.topic alikafka_post-cn-123455abc:topicName
```