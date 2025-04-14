---
subcategory: "Log Service (SLS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_logtail_attachment"
sidebar_current: "docs-alibabacloudstack-resource-logtail-attachment"
description: |-
  编排日志接️入服务
---

# alibabacloudstack_logtail_attachment

使用Provider配置的凭证在指定的资源集编排日志接入服务。
Logtail 接入服务是日志服务提供的日志采集工具。
您可以通过 Logtail 实时采集阿里云弹性计算服务(ECS)实例等服务器上的日志，在日志服务控制台上完成配置。[详情参考](https://www.alibabacloud.com/help/doc-detail/29058.htm)

该资源旨在将一个 Logtail 配置绑定到一个机器组。

-> **注意:** 一个 Logtail 配置可以绑定到多个机器组，一个机器组也可以绑定多个 Logtail 配置。

## 示例用法

### 基础用法

```
resource "alibabacloudstack_log_project" "test" {
  name        = "test-tf2"
  description = "create by terraform"
}

resource "alibabacloudstack_log_store" "test" {
  project               = alibabacloudstack_log_project.test.name
  name                  = "tf-test-logstore"
  retention_period      = 3650
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alibabacloudstack_log_machine_group" "test" {
  project       = alibabacloudstack_log_project.test.name
  name          = "tf-log-machine-group"
  topic         = "terraform"
  identify_list = ["10.0.0.1", "10.0.0.3", "10.0.0.2"]
}

resource "alibabacloudstack_logtail_config" "test" {
  project      = alibabacloudstack_log_project.test.name
  logstore     = alibabacloudstack_log_store.test.name
  input_type   = "file"
  log_sample   = "test"
  name         = "tf-log-config"
  output_type  = "LogService"
  input_detail = <<DEFINITION
  	{
		"logPath": "/logPath",
		"filePattern": "access.log",
		"logType": "json_log",
		"topicFormat": "default",
		"discardUnmatch": false,
		"enableRawLog": true,
		"fileEncoding": "gbk",
		"maxDepth": 10
	}
	
DEFINITION

}

resource "alibabacloudstack_logtail_attachment" "test" {
  project             = alibabacloudstack_log_project.test.name
  logtail_config_name = alibabacloudstack_logtail_config.test.name
  machine_group_name  = alibabacloudstack_log_machine_group.test.name
}
```

## 参数说明

支持以下参数：

* `project` - (必填，变更时重建) 日志存储所属的项目名称。
* `logtail_config_name` - (必填，变更时重建) Logtail 配置名称，在同一项目中必须唯一。
* `machine_group_name` - (必填，变更时重建) 机器组名称，在同一项目中必须唯一。
* `force_new_property` - (可选，变更时重建) 这是由 AI 添加的一个额外属性。

## 属性说明

导出以下属性：

* `id` - Logtail 到机器组的 ID。其格式为 `<project>:<logtail_config_name>:<machine_group_name>`。
* `computed_property` - (计算属性) 这是由 AI 添加的一个额外计算属性。

## 导入

Logtail 到机器组可以使用 ID 导入，例如：

```bash
$ terraform import alibabacloudstack_logtail_to_machine_group.example tf-log:tf-log-config:tf-log-machine-group
```